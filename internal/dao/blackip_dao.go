/**
 * 抽奖系统的数据库操作
 */
package dao

import (
	"gorm.io/gorm"
	models "lottery/model"
)

type BlackipDao struct {
	db *gorm.DB
}

func NewBlackipDao(db *gorm.DB) *BlackipDao {
	return &BlackipDao{
		db: db,
	}
}

func (d *BlackipDao) Get(id int) *models.LtBlackip {
	data := models.LtBlackip{}
	result := d.db.First(&data, id)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &data
}

func (d *BlackipDao) GetAll(page, size int) []models.LtBlackip {
	offset := (page - 1) * size
	datalist := make([]models.LtBlackip, 0)
	result := d.db.Order("id desc").Offset(offset).Limit(size).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *BlackipDao) CountAll() int64 {
	var count int64
	result := d.db.Model(&models.LtBlackip{}).Count(&count)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return count
}

func (d *BlackipDao) Search(ip string) []models.LtBlackip {
	datalist := make([]models.LtBlackip, 0)
	result := d.db.Order("id desc").Where("ip=?", ip).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *BlackipDao) Update(data *models.LtBlackip) error {
	result := d.db.Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *BlackipDao) Create(data *models.LtBlackip) error {
	result := d.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 根据IP获取信息 取1条
func (d *BlackipDao) GetByIp(ip string) *models.LtBlackip {
	data := models.LtBlackip{}
	result := d.db.Order("id desc").Where("ip=?", ip).First(&data)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &data
}
