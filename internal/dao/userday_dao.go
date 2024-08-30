/**
 * 抽奖系统的数据库操作
 */
package dao

import (
	"gorm.io/gorm"
	models "lottery/model"
)

type UserdayDao struct {
	db *gorm.DB
}

func NewUserdayDao(db *gorm.DB) *UserdayDao {
	return &UserdayDao{
		db: db,
	}
}

func (d *UserdayDao) Get(id int) *models.LtUserday {
	data := models.LtUserday{}
	result := d.db.First(&data, id)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &data
}

func (d *UserdayDao) GetAll(page, size int) []models.LtUserday {
	offset := (page - 1) * size
	datalist := make([]models.LtUserday, 0)
	result := d.db.Order("id desc").Offset(offset).Limit(size).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *UserdayDao) CountAll() int64 {
	var count int64
	result := d.db.Model(&models.LtUserday{}).Count(&count)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return count
}

func (d *UserdayDao) Search(uid, day int) []models.LtUserday {
	datalist := make([]models.LtUserday, 0)
	result := d.db.Order("id desc").Where("uid=? and day=?", uid, day).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *UserdayDao) Count(uid, day int) int {
	info := &models.LtUserday{}
	result := d.db.Where("uid=? and day=?", uid, day).First(info)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return info.Num
}

//func (d *UserdayDao) Delete(id int) error {
//	data := &models.LtUserday{Id: id, SysStatus: 1}
//	_, err := d.engine.Id(data.Id).Update(data)
//	return err
//}

func (d *UserdayDao) Update(data *models.LtUserday) error {
	result := d.db.Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *UserdayDao) Create(data *models.LtUserday) error {
	result := d.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
