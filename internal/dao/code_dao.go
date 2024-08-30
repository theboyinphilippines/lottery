package dao

import (
	"gorm.io/gorm"
	models "lottery/model"
)

type CodeDao struct {
	db *gorm.DB
}

func NewCodeDao(db *gorm.DB) *CodeDao {
	return &CodeDao{db: db}
}

func (d *CodeDao) Get(id int) *models.LtCode {
	data := models.LtCode{}
	result := d.db.First(&data, id)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &data
}

func (d *CodeDao) GetAll(page, size int) []models.LtCode {
	offset := (page - 1) * size
	datalist := make([]models.LtCode, 0)
	result := d.db.Order("id desc").Offset(offset).Limit(size).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *CodeDao) GetAllUsingCode() []models.LtCode {
	datalist := make([]models.LtCode, 0)
	result := d.db.Order("id asc").Where("sys_status=?", 0).Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *CodeDao) CountAll() int64 {
	var count int64
	result := d.db.Model(&models.LtCode{}).Count(&count)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return count
}

func (d *CodeDao) CountByGift(giftId int) int64 {
	var count int64
	result := d.db.Model(&models.LtCode{}).Where("gift_id=?", giftId).Count(&count)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return count
}

func (d *CodeDao) Search(giftId int) []models.LtCode {
	datalist := make([]models.LtCode, 0)
	result := d.db.Where("gift_id=? ", giftId).Order("id desc").Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *CodeDao) Delete(id int) error {
	result := d.db.Model(&models.LtCode{}).Where("id=?", id).Update("sys_status", 1)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *CodeDao) Update(data *models.LtCode) error {
	result := d.db.Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *CodeDao) Create(data *models.LtCode) error {
	result := d.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 找到下一个可用的最小的优惠券
func (d *CodeDao) NextUsingCode(giftId, codeId int) *models.LtCode {
	var datalist models.LtCode
	result := d.db.Where("gift_id=? and sys_status=? and id>?", giftId, 0, codeId).
		Order("id asc").First(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &datalist
}

// 根据唯一的code来更新
func (d *CodeDao) UpdateByCode(data *models.LtCode) error {
	result := d.db.Where("code=?", data.Code).Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 根据唯一的code来更新
func (d *CodeDao) UpdateStatusByCode(data *models.LtCode) error {
	result := d.db.Model(&models.LtCode{}).Where("code=?", data.Code).Updates(map[string]interface{}{"sys_status": data.SysStatus, "sys_updated": data.SysUpdated})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
