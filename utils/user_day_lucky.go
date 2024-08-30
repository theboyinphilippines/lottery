package utils

import (
	"context"
	"fmt"
	"log"
	"lottery/comm"
	"lottery/global"
	"math"
	"time"
)

func init() {
	resetGroupUserList()
}

const userFrameSize = 2

// 3.验证用户今日参与次数的缓存 hash
func IncrUserLuckyNum(uid int) int64 {
	// 分两段保存，随机生成2个hash key，保存在2个redis hash中
	i := uid % userFrameSize
	key := fmt.Sprintf("day_user_%d", i)
	rs, err := global.Rdb.Do(context.Background(), "HINCRBY", key, uid, 1).Result()
	if err != nil {
		log.Println("user_day_lucky.IncrUserLuckyNum redis HINCRBY error=", err)
		return math.MaxInt32
	} else {
		return rs.(int64)
	}
}

// 用于redis挂掉之后，参与次数已经落后数据库中的次数，重置redis中用户今天参与次数
func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % userFrameSize
	key := fmt.Sprintf("day_user_%d", i)
	_, err := global.Rdb.Do(context.Background(), "HSET", key, uid, num).Result()
	if err != nil {
		log.Println("user_day_lucky.InitUserLuckyNum redis HSET error=", err)
	}
}

// 定时删除，每天凌晨0点清空ip今日的参与次数
func resetGroupUserList() {
	log.Println("user_day_lucky.resetGroupUserList start")
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_user_%d", i)
		// TODO 删除一个hash key会导致线程阻塞？
		global.Rdb.Do(context.Background(), "DEL", key)
	}
	log.Println("user_day_lucky.resetGroupUserList stop")

	//第二天零点归0，设置定时器执行
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupUserList)
}
