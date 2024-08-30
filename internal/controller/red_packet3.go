package controller

//
//import (
//	"github.com/gin-gonic/gin"
//	"math/rand"
//	"net/http"
//	"strconv"
//	"sync"
//	"time"
//)
//
///*
//改造red_packet.go文件，使用sync.Map，channel 保证并发安全
//*/
//// 有并发安全问题，共享变量packetList （用3种方法处理：互斥锁，sync.map，channel，分布式锁？）
//type task struct {
//	id       uint32   //红包id
//	callback chan uint  //返回金额的通道
//}
//
//var packetList sync.Map
//var chTasks = make(chan task)
//
//// 获取所有红包的总个数和总金额
//func (i *IndexController) GetAllRedPacket(c *gin.Context) {
//	var totalMoney uint
//	var totalNum int
//	packetList.Range(func(key, value any) bool {
//		totalNum++
//		v := value.([]uint)
//		for _, money := range v {
//			totalMoney += money
//		}
//		return true
//	})
//
//	c.JSON(http.StatusOK, gin.H{
//		"totalNum":   totalNum,
//		"totalMoney": totalMoney,
//	})
//
//}
//
//// 发红包（sync Map）
//func (i *IndexController) DeliverRedPacket(c *gin.Context) {
//	uid := c.Query("uid")                  //用户id
//	num := c.Query("num")                  //红包数目
//	totalMoneyStr := c.Query("totalMoney") //红包总金额
//	numInt, _ := strconv.Atoi(num)
//	totalMoneyInt, _ := strconv.Atoi(totalMoneyStr)
//	uidInt, _ := strconv.Atoi(uid)
//	uidUint := uint32(uidInt)
//	totalMoneyUint := totalMoneyInt * 100 //换做分
//	leftMoney := totalMoneyUint
//	leftNum := numInt
//	rMax := 0.55 //最大红包金额比率
//	//红包数目大的时候 分配均匀一点
//	if numInt >= 100 {
//		rMax = 0.1
//	}
//
//	//随机发红包（发num个）
//	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
//	moneyList := make([]uint, numInt)
//	for j := 0; j < numInt; j++ {
//		if leftNum == 1 {
//			moneyList[j] = uint(leftMoney)
//			break
//		}
//		//当剩下的分钱等于剩余数目时，分钱不能再细分，每个给1分钱
//		if leftMoney == leftNum {
//			//每个给1分钱
//			for k := numInt - leftNum; k < numInt; k++ {
//				moneyList[k] = 1
//			}
//			break
//		}
//		// leftMoney-leftNum 给剩余的红包数，每个预留1分钱
//		money := seed.Intn(int(float64(leftMoney-leftNum) * rMax))
//		if money < 1 {
//			money = 1
//		}
//		moneyList[j] = uint(money)
//		leftMoney = leftMoney - money
//		leftNum--
//	}
//
//	packetList.Store(uidUint, moneyList)
//	c.JSON(http.StatusOK, gin.H{
//		"uid":        uidInt,
//		"num":        numInt,
//		"totalmoney": totalMoneyUint,
//		"packetList": moneyList,
//	})
//}
//
//// 抢红包（sync Map）
//func (ic *IndexController) FetchRedPacket(c *gin.Context) {
//	id := c.Query("id") //红包id
//	idInt, _ := strconv.Atoi(id)
//	idUint := uint32(idInt)
//
//	//构造一个抢红包任务
//	callback := make(chan uint)
//	t := task{
//		id:       idUint,
//		callback: callback,
//	}
//	//发送任务
//	chTasks <- t
//	go FetchPacketListMoneyWithChannel()
//	//接收返回结果
//	money := <-callback
//	if money <= 0 {
//		c.JSON(http.StatusOK, gin.H{
//			"msg": "红包抢完了",
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"id":    idInt,
//		"money": money,
//		//"packetList": MoneyList,
//	})
//	return
//}
//
//// 用channel抢红包
//func FetchPacketListMoneyWithChannel() {
//	for {
//		t := <-chTasks
//		id := t.id
//		MoneyListAny, _ := packetList.Load(id)
//		MoneyList := MoneyListAny.([]uint)
//		if len(MoneyList) >= 1 {
//			var money uint
//			seed := rand.New(rand.NewSource(time.Now().UnixNano()))
//			i := seed.Intn(len(MoneyList))
//			money = MoneyList[i]
//			MoneyList = append(MoneyList[:i], MoneyList[i+1:]...)
//			packetList.Store(id, MoneyList)
//			//返回抽奖金额回去
//			t.callback <- money
//		} else {
//			t.callback <- 0
//		}
//	}
//}
