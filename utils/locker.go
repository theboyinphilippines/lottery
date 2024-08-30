package utils

import (
	"context"
	"fmt"
	"lottery/global"
	"time"
)

//2.用户抽奖分布式锁（占位置）

// 获取锁（加过期时间，防止死锁）
func LockLucky(uid int) bool {
	key := getLuckyLockKey(uid)
	ok, _ := global.Rdb.SetNX(context.Background(), key, 1, 3*time.Second).Result()
	if ok {
		return true
	} else {
		return false
	}

}

// 释放锁（删key）
func UnLockLucky(uid int) bool {
	key := getLuckyLockKey(uid)
	res, _ := global.Rdb.Del(context.Background(), key).Result()
	if res != 1 {
		return false
	} else {
		return true
	}
}

// 获取锁的redis key
func getLuckyLockKey(uid int) string {
	return fmt.Sprintf("lucky_lock_%d", uid)
}
