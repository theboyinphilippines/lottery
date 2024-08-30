package models

type LtBlackip struct {
	Id         int    `gorm:"not null pk autoincr INT(10)"`
	Ip         string `gorm:"not null default '' comment('IP地址') VARCHAR(50)"`
	Blacktime  int    `gorm:"not null default 0 comment('黑名单限制到期时间') INT(10)"`
	SysCreated int    `gorm:"not null default 0 comment('创建时间') INT(10)"`
	SysUpdated int    `gorm:"not null default 0 comment('修改时间') INT(10)"`
}
