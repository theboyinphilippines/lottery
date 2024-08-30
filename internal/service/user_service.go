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

type UserService interface {
	GetAll(page, size int) []models.LtUser
	CountAll() int64
	Get(id int) *models.LtUser
	Update(data *models.LtUser) error
	Create(data *models.LtUser) error
	CheckBlackUser(uid int) (bool, *models.LtUser)
}

type userService struct {
	userDao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		userDao: dao.NewUserDao(global.DB),
	}
}

// 6.验证用户黑名单（黑名单不能抽实体奖） 数据库It_user
func (g *userService) CheckBlackUser(uid int) (bool, *models.LtUser) {
	userInfo := g.Get(uid)
	// 判断是否超过黑名单有效期，还在有效期，返回false
	if userInfo.Blacktime > int(time.Now().Unix()) {
		return false, userInfo
	}
	return true, nil
}
func (g *userService) GetAll(page, size int) []models.LtUser {
	return g.userDao.GetAll(page, size)
}

func (g *userService) CountAll() int64 {
	return g.userDao.CountAll()
}

func (g *userService) Get(id int) *models.LtUser {
	data := g.getByCache(id)
	if data == nil || data.Id <= 0 {
		data = g.userDao.Get(id)
		//数据库中也没读到数据，返回一个空数据
		if data == nil || data.Id <= 0 {
			data = &models.LtUser{Id: id}
		}
		g.setByCache(data)
	}
	return data
}

func (g *userService) Update(data *models.LtUser) error {
	g.updateByCache(data)
	return g.userDao.Update(data)
}

func (g *userService) Create(data *models.LtUser) error {
	return g.userDao.Create(data)
}

// 从缓存中得到信息（因为user表数据太多，无法全量保存，用部分保存，用hash保存单个用户信息）
func (g *userService) getByCache(id int) *models.LtUser {
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", id)
	dataMap := global.Rdb.HGetAll(context.Background(), key).Val()
	dataId := comm.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataId <= 0 {
		return nil
	}
	data := &models.LtUser{
		Id:         int(dataId),
		Username:   comm.GetStringFromStringMap(dataMap, "Username", ""),
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		Realname:   comm.GetStringFromStringMap(dataMap, "Realname", ""),
		Mobile:     comm.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIp:      comm.GetStringFromStringMap(dataMap, "SysIp", ""),
	}
	return data

}

// 将信息更新到缓存
func (s *userService) setByCache(data *models.LtUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	id := data.Id
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", id)
	// 数据更新到redis缓存
	params := []interface{}{}
	params = append(params, "Id", id)
	if data.Username != "" {
		params = append(params, "Username", data.Username)
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "Realname", data.Realname)
		params = append(params, "Mobile", data.Mobile)
		params = append(params, "Address", data.Address)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
		params = append(params, "SysIp", data.SysIp)
	}
	//err := global.Rdb.Do(context.Background(), "HMSET", params).Err()
	err := global.Rdb.HMSet(context.Background(), key, params...).Err()
	if err != nil {
		log.Println("user_service.setByCache HMSET params=", params, ", error=", err)
	}
}

// 数据更新了，直接清空缓存数据
func (s *userService) updateByCache(data *models.LtUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", data.Id)
	// 删除redis中的缓存
	global.Rdb.Do(context.Background(), "DEL", key)
}

var _ UserService = &userService{}
