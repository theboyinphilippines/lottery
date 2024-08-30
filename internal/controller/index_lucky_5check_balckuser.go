package controller

import (
	"lottery/global"
	models "lottery/model"
	"time"
)

// //6.验证用户黑名单（黑名单不能抽实体奖） 数据库It_user
func checkBlackUser(uid int) (bool, *models.LtUser) {
	var userInfo models.LtUser
	global.DB.First(&userInfo, uid)
	// 判断是否超过黑名单有效期，还在有效期，返回false
	if userInfo.Blacktime > int(time.Now().Unix()) {
		return false, &userInfo
	}
	return true, nil
}
