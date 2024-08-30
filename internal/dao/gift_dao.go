package dao

import (
	"gorm.io/gorm"
	"lottery/comm"
	models "lottery/model"
)

type GiftDao struct {
	db *gorm.DB
}

func NewGiftDao(db *gorm.DB) *GiftDao {
	return &GiftDao{db: db}
}

// 通过id获取奖品信息
func (d *GiftDao) Get(id int) *models.LtGift {
	data := models.LtGift{}
	result := d.db.First(&data, id)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return &data

}

// 获取所有奖品信息
func (d *GiftDao) GetAll() []models.LtGift {
	datalist := make([]models.LtGift, 0)
	result := d.db.Order("sys_status asc").Order("displayorder asc").Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

// 统计所有奖品
func (d *GiftDao) CountAll() int64 {
	var count int64
	result := d.db.Model(&models.LtGift{}).Count(&count)
	if result.RowsAffected == 0 || result.Error != nil {
		return 0
	}
	return count
}

// 删除奖品，实际将状态修改为1
func (d *GiftDao) Delete(id int) error {
	result := d.db.Model(&models.LtGift{}).Where("id=?", id).Update("sys_status", 1)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 更新奖品
func (d *GiftDao) Update(data *models.LtGift) error {
	result := d.db.Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 创建奖品
func (d *GiftDao) Create(data *models.LtGift) error {
	result := d.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 获取到当前可以获取的奖品列表
// 有奖品限定，状态正常，时间期间内
// gtype倒序， displayorder正序
func (d *GiftDao) GetAllUse() []models.LtGift {
	now := comm.NowUnix()
	datalist := make([]models.LtGift, 0)
	result := d.db.Select("id", "title", "prize_num", "left_num", "prize_code", "prize_time", "img", "displayorder", "gtype", "gdata").
		Where("prize_num >= ? and sys_status = ? and time_begin <= ? and time_end >= ?", 0, 0, now, now).
		Order("gtype desc").
		Order("displayorder asc").
		Find(&datalist)
	if result.RowsAffected == 0 || result.Error != nil {
		return nil
	}
	return datalist
}

func (d *GiftDao) IncrLeftNum(id, num int) error {
	result := d.db.Model(&models.LtGift{}).Select("left_num").Where("id =? ", id).Update("left_num", gorm.Expr("left_num + ?", num))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *GiftDao) DecrLeftNum(id, num int) error {
	result := d.db.Model(&models.LtGift{}).Select("left_num").Where("id =? and left_num>=? ", id, num).Update("left_num", gorm.Expr("left_num - ?", num))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
