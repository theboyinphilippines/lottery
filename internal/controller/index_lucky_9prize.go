package controller

import (
	"lottery/comm"
	"lottery/conf"
	"lottery/global"
	models "lottery/model"
	"strconv"
	"strings"
)

// 8.匹配奖品是否中奖，从lt_gift表中去匹配
func prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift models.ObjGiftPrize
	//先从奖品表中拿到所有的奖品列表，用prizeCode去匹配
	var giftList []models.LtGift
	now := comm.NowUnix()
	global.DB.Select("id", "title", "prize_num", "left_num", "prize_code", "prize_time", "img", "displayorder", "gtype", "gdata").
		Where("prize_num >= ? and sys_status = ? and time_begin <= ? and time_end >= ?", 0, 0, now, now).
		Order("gtype desc").
		Order("displayorder asc").
		Find(&giftList)
	for _, gift := range giftList {
		prizeCodeStr := strings.Split(gift.PrizeCode, "-")
		prizeCodeA, _ := strconv.Atoi(prizeCodeStr[0])
		prizeCodeB, _ := strconv.Atoi(prizeCodeStr[1])
		if prizeCodeA <= prizeCode && prizeCodeB >= prizeCode {
			//中奖编码满足区间,说明可以中奖
			// 不在黑名单或者抽到不是实体奖时，可以中奖
			//（如果不在黑名单，可以直接中奖；如果在黑名单同时抽到不是实体奖时，才中奖）
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				prizeGift.Id = gift.Id
				prizeGift.Title = gift.Title
				prizeGift.PrizeNum = gift.PrizeNum
				prizeGift.LeftNum = gift.LeftNum
				prizeGift.PrizeCodeA = prizeCodeA
				prizeGift.PrizeCodeB = prizeCodeB
				prizeGift.Img = gift.Img
				prizeGift.Displayorder = gift.Displayorder
				prizeGift.Gtype = gift.Gtype
				prizeGift.Gdata = gift.Gdata
				break
			}
		}
	}
	return &prizeGift
}
