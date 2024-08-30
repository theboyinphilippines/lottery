package dao

import (
	"gorm.io/gorm"
	models "lottery/model"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

// 通过id获取用户信息
func (d *UserDao) Get(id int) *models.LtUser {
	data := models.LtUser{}
	result := d.db.First(&data, id)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &data

}

// 获取所有用户信息
func (d *UserDao) GetAll(page, size int) []models.LtUser {
	offset := (page - 1) * size
	datalist := make([]models.LtUser, 0)
	result := d.db.Order("id desc").Offset(offset).Limit(size).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

// 统计所有用户
func (d *UserDao) CountAll() int64 {
	var count int64
	result := d.db.Model(&models.LtUser{}).Count(&count)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return count
}

// 更新用户
func (d *UserDao) Update(data *models.LtUser) error {
	result := d.db.Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 创建用户
func (d *UserDao) Create(data *models.LtUser) error {
	result := d.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
