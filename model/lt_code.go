package models

// 优惠券的优惠编码表
type LtCode struct {
	Id         int    `gorm:"not null pk autoincr INT(10)"`
	GiftId     int    `gorm:"not null default 0 comment('奖品ID，关联lt_gift表') INT(10)"`
	Code       string `gorm:"not null default '' comment('虚拟券编码') VARCHAR(255)"`
	SysCreated int    `gorm:"not null default 0 comment('创建时间') INT(10)"`
	SysUpdated int    `gorm:"not null default 0 comment('更新时间') INT(10)"`
	SysStatus  int    `gorm:"not null default 0 comment('状态，0正常，1作废，2已发放') SMALLINT(5)"`
}
