package utils

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"lottery/comm"
	"lottery/global"
	"lottery/initialize"
	"lottery/internal/dao"
	"lottery/internal/service"
	models "lottery/model"
)

func init() {
	//从数据库中查找可用奖品数据，将奖品放到奖品池 redis hash
	initialize.InitDB()
	giftList := dao.NewGiftDao(global.DB).GetAllUse()
	params := []interface{}{}
	for _, gift := range giftList {
		params = append(params, gift.Id, gift.LeftNum)
	}
	key := "gift_pool"
	err := global.Rdb.HMSet(context.Background(), key, params...).Err()
	if err != nil {
		log.Println("gift_service.IsPrize HMSET params=", params, ", error=", err)
	}
}

func init() {
	//这里项目没做后台接口，所以从数据库中拿到可用优惠券，放到redis set中
	codeInfoList := dao.NewCodeDao(global.DB).GetAllUsingCode()
	//pipe := global.Rdb.Pipeline()
	for _, codeInfo := range codeInfoList {
		key := fmt.Sprintf("gift_code_%d", codeInfo.GiftId)
		_, err := global.Rdb.Do(context.Background(), "SADD", key, codeInfo.Code).Result()
		if err != nil {
			log.Println("prizedata.init gift_code SADD error=", err)
		}

		//pipe.SAdd(context.Background(), key, codeInfo.Code)
	}
	//_, err := pipe.Exec(context.Background())
	//if err != nil {
	//	log.Println("prizedata.init gift_code pipeline SADD error=", err)
	//}

}

// 获取当前奖品池中的奖品数量，从redis中
func GetGiftPoolNum(id int) int {
	key := "gift_pool"
	rs, err := global.Rdb.Do(context.Background(), "HGET", key, id).Result()
	if err != nil {
		log.Println("prizedata.GetGiftPoolNum error=", err)
		return 0
	}
	num := comm.GetInt64(rs, 0)
	return int(num)
}

func DecrGiftPoolNum(id int) int {
	key := "gift_pool"
	rs, err := global.Rdb.Do(context.Background(), "HGET", key, id).Result()
	if err != nil {
		log.Println("prizedata.GetGiftPoolNum error=", err)
		return 0
	}
	num := comm.GetInt64(rs, 0)
	return int(num)
}

// 设置奖品池的数量
func setGiftPool(id, num int) {
	key := "gift_pool"
	_, err := global.Rdb.Do(context.Background(), "HSET", key, id, num).Result()
	if err != nil {
		log.Println("prizedata.setGiftPool error=", err)
	}
}

// 发奖，redis缓存
func PrizeServGift(id int) bool {
	key := "gift_pool"
	rs, err := global.Rdb.Do(context.Background(), "HINCRBY", key, id, -1).Result()
	log.Println("prizedata.prizeServGift rs=", rs)

	if err != nil {
		log.Println("prizedata.prizeServGift error=", err)
		return false
	}
	num := comm.GetInt64(rs, -1)
	log.Println("prizedata.prizeServGift num=", num)
	if num >= 0 {
		return true
	} else {
		return false
	}
}

// //9.有限制奖品发放（奖品有剩余才可以发出去） 奖品剩余数量扣减
func PrizeGift(id int) bool {
	//缓存扣减奖品剩余数量
	ok := PrizeServGift(id)
	if ok {
		//缓存中扣减成功，就更新数据库
		//扣减剩余数量，这里增量操作+1，-1，mysql默认是原子操作，不存在并发安全问题导致超卖
		result := global.DB.Model(&models.LtGift{}).Select("left_num").Where("id =? and left_num>=? ", id, 1).Update("left_num", gorm.Expr("left_num - ?", 1))
		if result.RowsAffected == 0 || result.Error != nil {
			log.Println("prizedata.PrizeGift error")
			return false
		}
	}
	return ok
}

// 10.不同编码的优惠券的发放 lt_code表
// 不同优惠券的发放，需要用到分布式锁（从数据库中取code）
func PrizeCodeDiff(giftId int) string {
	// 从缓存里面取一个code
	code := DeliverCodeFromCache(giftId)
	//找到了编码，更新数据库状态
	codeInfo := models.LtCode{
		Code:       code,
		SysUpdated: comm.NowUnix(),
		SysStatus:  2,
	}
	err := dao.NewCodeDao(global.DB).UpdateStatusByCode(&codeInfo)
	if err != nil {
		log.Println("prizedata.PrizeCodeDiff error")
		return ""
	}
	return code
}

// 10.不同编码的优惠券的发放 lt_code表（老式用法）
// 不同优惠券的发放，需要用到分布式锁（从数据库中取code）
func PrizeCodeDiffFromMysql(giftId int) string {
	lockUid := 0 - giftId - 100000000
	LockLucky(lockUid)
	defer UnLockLucky(lockUid)
	// 在lt_code表中找到可用的优惠券的编码
	var codeInfo models.LtCode
	codeId := 0
	result := global.DB.Where("id>? and gift_id=? and sys_status=? ", codeId, giftId, 0).Order("id asc").First(&codeInfo)
	if result.RowsAffected == 0 || result.Error != nil {
		log.Println("prizedata.PrizeCodeDiff error")
		return ""
	}
	//找到了编码，更新状态
	codeInfo.SysStatus = 2
	codeInfo.SysUpdated = comm.NowUnix()
	global.DB.Save(&codeInfo)
	return codeInfo.Code
}

// 导入新的优惠券编码 用redis中set来做（后台导入code，保存到数据库同时，也要保存到redis中）
func ImportCacheCodes(id int, code string) bool {
	// 集群版本需要放入到redis中
	// [暂时]本机版本的就直接从数据库中处理吧
	// redis中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	_, err := global.Rdb.Do(context.Background(), "SADD", key, code).Result()
	if err != nil {
		log.Println("prizedata.RecacheCodes SADD error=", err)
		return false
	} else {
		return true
	}
}

// 重新整理优惠券的编码到缓存中（用于后台倒入新的code时，重新更新缓存）
func RecacheCodes(id int, codeService service.CodeService) (sucNum, errNum int) {
	// 集群版本需要放入到redis中
	// [暂时]本机版本的就直接从数据库中处理吧
	list := codeService.Search(id)
	if list == nil || len(list) <= 0 {
		return 0, 0
	}
	// redis中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	tmpKey := "tmp_" + key
	for _, data := range list {
		//挑出正常状态的code
		if data.SysStatus == 0 {
			code := data.Code
			_, err := global.Rdb.Do(context.Background(), "SADD", tmpKey, code).Result()
			if err != nil {
				log.Println("prizedata.RecacheCodes SADD error=", err)
				errNum++
			} else {
				sucNum++
			}
		}
	}
	//修改key的名称（new key存在时，将覆盖旧值）
	_, err := global.Rdb.Do(context.Background(), "RENAME", tmpKey, key).Result()
	if err != nil {
		log.Println("prizedata.RecacheCodes RENAME error=", err)
	}
	return sucNum, errNum
}

// 从缓存中抽取一个优惠券code
func DeliverCodeFromCache(giftId int) string {
	key := fmt.Sprintf("gift_code_%d", giftId)
	rs, err := global.Rdb.Do(context.Background(), "SPOP", key).Result()
	if err != nil {
		log.Println("prizedata.init gift_code SADD error=", err)
	}
	code := comm.GetString(rs, "")
	return code
}
