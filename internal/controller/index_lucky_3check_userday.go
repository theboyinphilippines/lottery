package controller

import (
	"fmt"
	"log"
	"lottery/conf"
	"lottery/global"
	models "lottery/model"
	"lottery/utils"
	"strconv"
	"time"
)

// 3.验证用户今日参与次数（从数据库中取，表It_userday）
func checkUserday(uid int, num int64) bool {
	y, m, d := time.Now().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	day, _ := strconv.Atoi(strDay)
	var userdayInfo models.LtUserday
	result := global.DB.Where("uid = ? and day = ?", uid, day).First(&userdayInfo)
	if result.RowsAffected == 1 && userdayInfo.Uid == uid {
		//数据库有记录 判断次数是否超出设置值
		if userdayInfo.Num > conf.UserPrizeMax {
			// redis中的次数小于数据库中的次数时，将数据库中次数重置redis中的次数
			if userdayInfo.Num > int(num) {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			return false
		} else {
			//没有超出就+1，更新数据库
			userdayInfo.Num++
			if userdayInfo.Num > int(num) {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			result = global.DB.Save(&userdayInfo)
			if result.Error != nil {
				log.Println("checkUserday.Save err: " + result.Error.Error())
			}
		}
	} else {
		//数据库没有记录,创建记录
		userdayInfo.Uid = uid
		userdayInfo.Day = day
		userdayInfo.Num = 1
		userdayInfo.SysCreated = int(time.Now().Unix())
		result = global.DB.Create(&userdayInfo)
		if result.Error != nil {
			log.Println("checkUserday.Create err: " + result.Error.Error())
		}
		utils.InitUserLuckyNum(uid, 1)
	}
	return true
}
