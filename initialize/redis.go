package initialize

import (
	"context"
	"github.com/redis/go-redis/v9"
	"lottery/global"
)

func InitRedis() {
	//默认有连接池
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	global.Rdb = rdb
}
