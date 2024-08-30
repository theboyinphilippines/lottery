package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/global"
	models "lottery/model"
	"lottery/utils"
	"net/http"
)

// 抽奖
func (i *IndexController) GetLucky(c *gin.Context) {
	//1.验证登录（jwt中间件自动验证）

	//2.用户抽奖分布式锁（锁用户，占位置，防止用户快速点击，或恶意刷进口）
	uid, _ := c.Get("userId")
	userName, _ := c.Get("userName")
	uidInt := int(comm.GetInt64(uid, 0))
	userNameStr := comm.GetString(userName, "")
	ok := utils.LockLucky(uidInt)
	if ok {
		defer utils.UnLockLucky(uidInt)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 102,
			"msg":  "正在抽奖，请稍后重试",
		})
		return
	}

	//3.验证用户今日参与次数
	//先从缓存中拿到今天的参与次数
	userDayNum := utils.IncrUserLuckyNum(uidInt)
	if userDayNum > conf.UserPrizeMax {
		c.JSON(http.StatusOK, gin.H{
			"code": 103,
			"msg":  "今日抽奖次数已用完，明天再来把",
		})
	} else {
		//从数据库中去验证用户今日次数
		ok = checkUserday(uidInt, userDayNum)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": 103,
				"msg":  "今日抽奖次数已用完，明天再来把",
			})
			return
		}
	}

	//4.验证ip今日的参与次数
	ip := comm.ClientIP(c.Request)
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		c.JSON(http.StatusOK, gin.H{
			"code": 104,
			"msg":  "相同ip参与次数太多，明天再来把",
		})
		return
	}

	//5.验证ip黑名单（避免同一个ip大量用户的刷奖，黑名单不能抽实体奖）
	//设置一个黑名单参数
	limitBlack := false
	//ok, blackIpInfo := checkBlackIp(ip)
	ok, blackIpInfo := i.ServiceBlackIp.CheckBlackIp(ip)
	if !ok {
		//在ip黑名单
		limitBlack = true
		log.Println("黑名单中的ip", ip, limitBlack)
	}

	//6.验证用户黑名单（黑名单不能抽实体奖）
	//ok, userInfo := checkBlackUser(uidInt)
	ok, userInfo := i.ServiceUser.CheckBlackUser(uidInt)
	if !ok {
		//在用户黑名单
		limitBlack = true
		log.Println("黑名单中的用户", uidInt, limitBlack)
	}

	//7.获得抽奖编码
	// 随机得到一个抽奖编码，匹配有效的奖品中奖编码区间
	// 奖品是有序的，后台可以设置顺序，避免编码区间的包含关系
	prizeCode := comm.Random(10000)

	//8.匹配奖品是否中奖

	//prizeGift := prize(prizeCode, limitBlack)

	prizeGift := i.ServiceGift.IsPrize(prizeCode, limitBlack)
	// 判断是否匹配到奖品，或奖品剩余数量不足，或者已在黑名单但抽到的是实体奖，也不能让他中奖，此时prizeGift id为初始值0（prizeGift.Id <= 0）
	if prizeGift == nil || prizeGift.Id <= 0 || prizeGift.PrizeNum < 0 ||
		(prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		c.JSON(http.StatusOK, gin.H{
			"code": 205,
			"msg":  "很遗憾，没有中奖，请下次再试",
		})
		return
	}

	//10.假如匹配到的奖是优惠券 不同编码的优惠券的发放
	//不同优惠券的发放，需要用到分布式锁
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id)
		//数据库没有不同编码的优惠券
		if code == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 208,
				"msg":  "很遗憾，没有中奖，请下次再试",
			})
			return
		}
	}

	//9.有限制奖品发放（奖品有剩余才可以发出去）
	//数据库实现原子性奖品库存数量的递减
	if prizeGift.PrizeNum > 0 {
		// 奖品剩余数量扣减
		//先从缓存里扣，看是否能扣减 （i.ServiceGift.IsPrize这步已经放到奖品池了）
		giftPoolNum := utils.GetGiftPoolNum(prizeGift.Id)
		if giftPoolNum <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": 206,
				"msg":  "很遗憾，没有中奖，请下次再试",
			})
			return
		}
		//从缓存中扣减，然后更新数据库  utils.PrizeGift包含的缓存扣减
		ok = utils.PrizeGift(prizeGift.Id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": 207,
				"msg":  "很遗憾，没有中奖，请下次再试",
			})
			return
		}
	}

	//11. 记录中奖记录
	//奖品信息可能修改，但中奖记录中不能修改
	//如果中了实物大奖，还需要把用户，ip设置为黑名单一段时间
	result := models.LtResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        uidInt,
		Username:   userNameStr,
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	//保存中奖记录
	res := global.DB.Create(&result)
	if res.RowsAffected == 0 || res.Error != nil {
		log.Println("index_lucky.Create result，err:", res.Error)
		c.JSON(http.StatusOK, gin.H{
			"code": 209,
			"msg":  "很遗憾，没有中奖，请下次再试",
		})
		return
	}
	//如果中了实物大奖，需要加入黑名单一段时间
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		PrizeLarge(ip, uidInt, blackIpInfo, userInfo)
	}
	log.Println("prizeGift：", prizeGift)

	//12. 返回抽奖结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "恭喜你中奖",
		"data": prizeGift,
	})
	return
}
