package controller

import (
	"github.com/gin-gonic/gin"
	"lottery/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/*
改造red_packet.go文件，使用分布式锁 保证并发安全
*/
// 有并发安全问题，共享变量packetList （用3种方法处理：互斥锁，sync.map，channel，分布式锁？）
var packetList = make(map[uint32][]uint)

// 获取所有红包的总个数和总金额
func (i *IndexController) GetAllRedPacket(c *gin.Context) {
	var totalMoney uint
	var totalNum int
	for _, v := range packetList {
		totalNum++
		for _, money := range v {
			totalMoney += money
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"totalNum":   totalNum,
		"totalMoney": totalMoney,
	})

}

// 发红包（分布式锁）
func (i *IndexController) DeliverRedPacket(c *gin.Context) {
	uid := c.Query("uid") //红包id
	num := c.Query("num") //红包数目
	//id := c.Query("id")                    //用户id
	totalMoneyStr := c.Query("totalMoney") //红包总金额
	numInt, _ := strconv.Atoi(num)
	totalMoneyInt, _ := strconv.Atoi(totalMoneyStr)
	//idInt, _ := strconv.Atoi(id)
	uidInt, _ := strconv.Atoi(uid)
	uidUint := uint32(uidInt)
	totalMoneyUint := totalMoneyInt * 100 //换做分
	leftMoney := totalMoneyUint
	leftNum := numInt
	rMax := 0.55 //最大红包金额比率
	//红包数目大的时候 分配均匀一点
	if numInt >= 100 {
		rMax = 0.1
	}

	//随机发红包（发num个）
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	moneyList := make([]uint, numInt)
	for j := 0; j < numInt; j++ {
		if leftNum == 1 {
			moneyList[j] = uint(leftMoney)
			break
		}
		//当剩下的分钱等于剩余数目时，分钱不能再细分，每个给1分钱
		if leftMoney == leftNum {
			//每个给1分钱
			for k := numInt - leftNum; k < numInt; k++ {
				moneyList[k] = 1
			}
			break
		}
		// leftMoney-leftNum 给剩余的红包数，每个预留1分钱
		money := seed.Intn(int(float64(leftMoney-leftNum) * rMax))
		if money < 1 {
			money = 1
		}
		moneyList[j] = uint(money)
		leftMoney = leftMoney - money
		leftNum--
	}
	//使用分布式锁（使用同一个key占位置）
	ok := utils.LockLucky(100001)
	if ok {
		defer utils.UnLockLucky(100001)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 102,
			"msg":  "发红包失败，请稍后重试",
		})
		return
	}

	packetList[uidUint] = moneyList
	c.JSON(http.StatusOK, gin.H{
		"uid":        uidInt,
		"num":        numInt,
		"totalmoney": totalMoneyUint,
		"packetList": packetList,
	})
	return
}

// 抢红包（分布式锁）
func (ic *IndexController) FetchRedPacket(c *gin.Context) {
	id := c.Query("id") //红包id
	//uid := c.Query("uid") //用户id
	idInt, _ := strconv.Atoi(id)
	//uidInt, _ := strconv.Atoi(uid)
	idUint := uint32(idInt)

	//抢红包，随机生成index
	//使用分布式锁（使用同一个key占位置）
	ok := utils.LockLucky(100000)
	if ok {
		defer utils.UnLockLucky(100000)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 102,
			"msg":  "抢红包失败，请稍后重试",
		})
		return
	}
	MoneyList := packetList[idUint]
	if len(MoneyList) >= 1 {
		var money uint
		seed := rand.New(rand.NewSource(time.Now().UnixNano()))
		i := seed.Intn(len(MoneyList))
		money = MoneyList[i]
		MoneyList = append(MoneyList[:i], MoneyList[i+1:]...)
		packetList[idUint] = MoneyList
		c.JSON(http.StatusOK, gin.H{
			"id":         idInt,
			"money":      money,
			"packetList": packetList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "红包抢完了",
	})
	return

}
