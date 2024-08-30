/**
 * 抽奖系统数据处理（包括数据库，也包括缓存等其他形式数据）
 */
package service

import (
	"context"
	"fmt"
	"log"
	"lottery/comm"
	"lottery/global"
	"lottery/internal/dao"
	models "lottery/model"
	"time"
)

type BlackipService interface {
	GetAll(page, size int) []models.LtBlackip
	CountAll() int64
	Search(ip string) []models.LtBlackip
	Get(id int) *models.LtBlackip
	//Delete(id int) error
	Update(user *models.LtBlackip) error
	Create(user *models.LtBlackip) error
	GetByIp(ip string) *models.LtBlackip
	CheckBlackIp(ip string) (bool, *models.LtBlackip)
}

type blackipService struct {
	dao *dao.BlackipDao
}

func NewBlackipService() BlackipService {
	return &blackipService{
		dao: dao.NewBlackipDao(global.DB),
	}
}

// 5.验证ip黑名单 数据库It_balckip（避免同一个ip大量用户的刷奖，黑名单不能抽实体奖）
func (s *blackipService) CheckBlackIp(ip string) (bool, *models.LtBlackip) {
	blackIpInfo := s.GetByIp(ip)
	if blackIpInfo == nil || blackIpInfo.Id <= 0 {
		//没有黑名单记录，返回true
		return true, nil
	}
	// 判断是否超过黑名单有效期，还在有效期，返回false
	if blackIpInfo.Blacktime > int(time.Now().Unix()) {
		return false, blackIpInfo
	}
	return true, nil
}

func (s *blackipService) GetAll(page, size int) []models.LtBlackip {
	return s.dao.GetAll(page, size)
}

func (s *blackipService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *blackipService) Search(ip string) []models.LtBlackip {
	return s.dao.Search(ip)
}

func (s *blackipService) Get(id int) *models.LtBlackip {
	return s.dao.Get(id)
}

//func (s *blackipService) Delete(id int) error {
//	return s.dao.Delete(id)
//}

func (s *blackipService) Update(data *models.LtBlackip) error {
	// 先更新缓存的数据
	s.updateByCache(data)
	// 再更新数据的数据
	return s.dao.Update(data)
}

func (s *blackipService) Create(data *models.LtBlackip) error {
	return s.dao.Create(data)
}

// 根据IP读取IP的黑名单数据
func (s *blackipService) GetByIp(ip string) *models.LtBlackip {
	// 先从缓存中读取数据
	data := s.getByCache(ip)
	if data == nil || data.Ip == "" {
		// 再从数据库中读取数据
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.LtBlackip{Ip: ip}
		}
		s.setByCache(data)
	}
	return data
}

func (s *blackipService) getByCache(ip string) *models.LtBlackip {
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", ip)
	dataMap := global.Rdb.HGetAll(context.Background(), key).Val()
	dataIp := comm.GetStringFromStringMap(dataMap, "Ip", "")
	if dataIp == "" {
		return nil
	}
	data := &models.LtBlackip{
		Id:         int(comm.GetInt64FromStringMap(dataMap, "Id", 0)),
		Ip:         dataIp,
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
	}
	return data
}

// 缓存放在redis hash
func (s *blackipService) setByCache(data *models.LtBlackip) {
	if data == nil || data.Ip == "" {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	// 数据更新到redis缓存
	params := []interface{}{}
	params = append(params, "Ip", data.Ip)
	if data.Id > 0 {
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
	}
	//err := global.Rdb.Do(context.Background(), "HMSET", params...).Err()
	err := global.Rdb.HMSet(context.Background(), key, params...).Err()
	if err != nil {
		log.Println("blackip_service.setByCache HMSET params=", params, ", error=", err)
	}
}

// 数据更新了，直接清空缓存数据
func (s *blackipService) updateByCache(data *models.LtBlackip) {
	if data == nil || data.Ip == "" {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	// 删除redis中的缓存
	global.Rdb.Do(context.Background(), "DEL", key)
}
