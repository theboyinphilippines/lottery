package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// init方法：被引用这个包时，会自动调用init方法
var (
	DB  *gorm.DB
	Rdb *redis.Client
)
