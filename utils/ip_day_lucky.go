package utils

import (
	"context"
	"fmt"
	"log"
	"lottery/comm"
	"lottery/global"
	"lottery/initialize"
	"math"
	"time"
)

func init() {
	//resetGroupIpList() 依赖redis，这里要先初始化redis，才能调用resetGroupIpList
	initialize.InitRedis()
	resetGroupIpList()
}

const ipFrameSize = 2

// 4.验证ip今日的参与次数 redis hash存储
func IncrIpLuckyNum(strIp string) int64 {
	//ip先转化为整数
	ip := comm.Ip4toInt(strIp)
	// 分两段保存，随机生成2个hash key，保存在2个redis hash中
	i := ip % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i)
	rs, err := global.Rdb.Do(context.Background(), "HINCRBY", key, ip, 1).Result()

	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error=", err)
		return math.MaxInt32
	} else {
		return rs.(int64)
	}
}

// 定时删除，每天凌晨0点清空ip今日的参与次数
func resetGroupIpList() {
	log.Println("ip_day_lucky.resetGroupIpList start")
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_ips_%d", i)
		// TODO 删除一个hash key会导致线程阻塞？
		global.Rdb.Do(context.Background(), "DEL", key)
	}
	log.Println("ip_day_lucky.resetGroupIpList stop")

	//第二天零点归0，设置定时器执行
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIpList)
}
