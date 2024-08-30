package controller

import (
	"lottery/comm"
	"lottery/global"
	models "lottery/model"
)

func PrizeLarge(ip string, uid int, blackIpInfo *models.LtBlackip, userInfo *models.LtUser) {
	//更新用户黑名单
	now := comm.NowUnix()
	blackTime := now + 30*86400 //黑30天
	global.DB.Model(&models.LtUser{}).Where("id=?", uid).Update("blacktime", blackTime)

	//更新ip黑名单
	if blackIpInfo == nil || blackIpInfo.Id <= 0 {
		//没有ip黑名单记录就创建
		blackIpInfo = &models.LtBlackip{
			Ip:         ip,
			Blacktime:  blackTime,
			SysCreated: now,
		}
		global.DB.Create(blackIpInfo)
	} else {
		//有记录就更新
		blackIpInfo.Blacktime = blackTime
		blackIpInfo.SysUpdated = now
		global.DB.Save(blackIpInfo)
	}
}
