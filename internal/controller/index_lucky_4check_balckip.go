package controller

import (
	"lottery/global"
	models "lottery/model"
	"time"
)

// 5.验证ip黑名单 数据库It_balckip（避免同一个ip大量用户的刷奖，黑名单不能抽实体奖）
func checkBlackIp(ip string) (bool, *models.LtBlackip) {
	var blackIpInfo models.LtBlackip
	result := global.DB.Where("ip = ? ", ip).First(&blackIpInfo)
	if result.RowsAffected == 0 {
		//没有黑名单记录，返回true
		return true, nil
	}
	// 判断是否超过黑名单有效期，还在有效期，返回false
	if blackIpInfo.Blacktime > int(time.Now().Unix()) {
		return false, &blackIpInfo
	}
	return true, nil
}
