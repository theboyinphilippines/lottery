package service

import (
	"context"
	"encoding/json"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/global"
	"lottery/internal/dao"
	models "lottery/model"
	"strconv"
	"strings"
)

// 用useCache bool来说明是否要用缓存
// 更新数据库数据时，首先要更新缓存，才能保持数据一致性
type GiftService interface {
	GetAll(useCache bool) []models.LtGift
	CountAll() int64
	Get(id int, useCache bool) *models.LtGift
	Delete(id int) error
	Update(data *models.LtGift) error
	Create(data *models.LtGift) error
	GetAllUse(useCache bool) []models.LtGift
	IncrLeftNum(id, num int) error
	DecrLeftNum(id, num int) error
	IsPrize(prizeCode int, limitBlack bool) *models.ObjGiftPrize
}

type giftService struct {
	giftDao *dao.GiftDao
}

func NewGiftService() GiftService {
	return &giftService{
		giftDao: dao.NewGiftDao(global.DB),
	}
}

func (g *giftService) IsPrize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift models.ObjGiftPrize
	//先从奖品表中拿到所有的奖品列表，用prizeCode去匹配
	giftList := g.GetAllUse(true)
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

func (g *giftService) GetAll(useCache bool) []models.LtGift {
	//从数据库中读
	if !useCache {
		return g.giftDao.GetAll()
	}
	//从缓存中读
	gifts := g.getAllByCache()
	// 如果缓存中读到空，从数据库中取，并放到缓存中
	if len(gifts) < 1 {
		gifts = g.giftDao.GetAll()
		g.setAllByCache(gifts)
	}
	return gifts
}

func (g *giftService) CountAll() int64 {
	//return g.giftDao.CountAll()
	//从缓存里面读（GetAll缓存中没有就在数据库中读）
	gifts := g.GetAll(true)
	return int64(len(gifts))
}

func (g *giftService) Get(id int, useCache bool) *models.LtGift {
	if !useCache {
		return g.giftDao.Get(id)
	}
	//GetAll缓存中没有就在数据库中读
	gifts := g.GetAll(true)
	for _, gift := range gifts {
		if gift.Id == id {
			return &gift
		}
	}
	return nil
}

// 软删除，实际是更新状态
func (g *giftService) Delete(id int) error {
	//更新数据库数据时，首先要更新缓存，才能保持数据一致性
	data := models.LtGift{Id: id}
	//更新缓存，实际就是直接清空缓存
	g.updateByCache(&data)
	return g.giftDao.Delete(id)
}

func (g *giftService) Update(data *models.LtGift) error {
	//更新缓存，实际就是直接清空缓存
	g.updateByCache(data)
	return g.giftDao.Update(data)
}

func (g *giftService) Create(data *models.LtGift) error {
	g.updateByCache(data)
	return g.giftDao.Create(data)
}

func (g *giftService) GetAllUse(useCache bool) []models.LtGift {
	if !useCache {
		return g.giftDao.GetAllUse()
	} else {
		gifts := g.GetAll(true)
		datalist := make([]models.LtGift, 0)
		for _, gift := range gifts {
			now := comm.NowUnix()
			//筛选出符合条件的奖品
			if gift.Id > 0 && gift.PrizeNum >= 0 && gift.SysStatus == 0 &&
				gift.TimeBegin <= now && gift.TimeEnd >= now {
				datalist = append(datalist, gift)
			}
		}
		return datalist
	}
}

func (g *giftService) IncrLeftNum(id, num int) error {
	return g.giftDao.IncrLeftNum(id, num)
}

func (g *giftService) DecrLeftNum(id, num int) error {
	return g.giftDao.DecrLeftNum(id, num)
}

func (g *giftService) getAllByCache() []models.LtGift {
	key := "allgift"
	result, err := global.Rdb.Do(context.Background(), "GET", key).Result()
	if err != nil {
		log.Println("gift_service.getAllByCache GET key=", key, ", error=", err)
		return nil
	}
	str := comm.GetString(result, "")
	if str == "" {
		return nil
	}
	//redis中缓存数据反序列化给dataList（本来可以发序列化给[]models.LtGift，但是model中有些字段json设置为-）
	dataList := []map[string]interface{}{}
	err = json.Unmarshal([]byte(str), &dataList)
	if err != nil {
		log.Println("gift_service.json.Unmarshal str error=", err)
		return nil
	}
	gifts := make([]models.LtGift, len(dataList))
	for i := 0; i < len(dataList); i++ {
		data := dataList[i]
		id := comm.GetInt64FromMap(data, "Id", 0)
		if id < 0 {
			gifts[i] = models.LtGift{}
		} else {
			gift := models.LtGift{
				Id:           int(id),
				Title:        comm.GetStringFromMap(data, "Title", ""),
				PrizeNum:     int(comm.GetInt64FromMap(data, "PrizeNum", 0)),
				LeftNum:      int(comm.GetInt64FromMap(data, "LeftNum", 0)),
				PrizeCode:    comm.GetStringFromMap(data, "PrizeCode", ""),
				PrizeTime:    int(comm.GetInt64FromMap(data, "PrizeTime", 0)),
				Img:          comm.GetStringFromMap(data, "Img", ""),
				Displayorder: int(comm.GetInt64FromMap(data, "Displayorder", 0)),
				Gtype:        int(comm.GetInt64FromMap(data, "Gtype", 0)),
				Gdata:        comm.GetStringFromMap(data, "Gdata", ""),
				TimeBegin:    int(comm.GetInt64FromMap(data, "TimeBegin", 0)),
				TimeEnd:      int(comm.GetInt64FromMap(data, "TimeEnd", 0)),
				//PrizeData:    comm.GetStringFromMap(data, "PrizeData", ""),
				PrizeBegin: int(comm.GetInt64FromMap(data, "PrizeBegin", 0)),
				PrizeEnd:   int(comm.GetInt64FromMap(data, "PrizeEnd", 0)),
				SysStatus:  int(comm.GetInt64FromMap(data, "SysStatus", 0)),
				SysCreated: int(comm.GetInt64FromMap(data, "SysCreated", 0)),
				SysUpdated: int(comm.GetInt64FromMap(data, "SysUpdated", 0)),
				SysIp:      comm.GetStringFromMap(data, "SysIp", ""),
			}
			gifts[i] = gift
		}
	}
	return gifts
}

func (g *giftService) setAllByCache(gifts []models.LtGift) {
	strValue := ""
	if len(gifts) > 0 {
		dataList := make([]map[string]interface{}, len(gifts))
		//格式转化为 gifts转化为dataList
		for i := 0; i < len(gifts); i++ {
			// gifts中一个一个取出来
			gift := gifts[i]
			data := make(map[string]interface{})
			data["Id"] = gift.Id
			data["Title"] = gift.Title
			data["PrizeNum"] = gift.PrizeNum
			data["LeftNum"] = gift.LeftNum
			data["PrizeCode"] = gift.PrizeCode
			data["PrizeTime"] = gift.PrizeTime
			data["Img"] = gift.Img
			data["Displayorder"] = gift.Displayorder
			data["Gtype"] = gift.Gtype
			data["Gdata"] = gift.Gdata
			data["TimeBegin"] = gift.TimeBegin
			data["TimeEnd"] = gift.TimeEnd
			//data["PrizeData"] = gift.PrizeData
			data["PrizeBegin"] = gift.PrizeBegin
			data["PrizeEnd"] = gift.PrizeEnd
			data["SysStatus"] = gift.SysStatus
			data["SysCreated"] = gift.SysCreated
			data["SysUpdated"] = gift.SysUpdated
			data["SysIp"] = gift.SysIp
			dataList[i] = data
		}
		str, err := json.Marshal(dataList)
		if err != nil {
			log.Println("gift_service.setAllByCache json.Marshal(dataList) error = ", err)
		}
		strValue = string(str)
	}
	key := "allgift"
	_, err := global.Rdb.Do(context.Background(), "SET", key, strValue).Result()
	if err != nil {
		log.Println("gift_service.setAllByCache global.Rdb.Do SET error = ", err)
	}
}

// 更新缓存，直接清空缓存数据
func (g *giftService) updateByCache(data *models.LtGift) {
	if data == nil || data.Id <= 0 {
		return
	}
	// 集群模式，redis缓存
	key := "allgift"
	_, err := global.Rdb.Do(context.Background(), "DEL", key).Result()
	if err != nil {
		log.Println("gift_service.setAllByCache global.Rdb.Do DEL key error = ", err)
	}

}

var _ GiftService = &giftService{}
