package models

// 用户表
type LtUser struct {
	Id         int    `gorm:"not null pk autoincr INT(10)"`
	Username   string `gorm:"not null default '' comment('用户名') VARCHAR(50)"`
	Blacktime  int    `gorm:"not null default 0 comment('黑名单限制到期时间') INT(10)"`
	Realname   string `gorm:"not null default '' comment('联系人') VARCHAR(50)"`
	Mobile     string `gorm:"not null default '' comment('手机号') VARCHAR(50)"`
	Address    string `gorm:"not null default '' comment('联系地址') VARCHAR(255)"`
	SysCreated int    `gorm:"not null default 0 comment('创建时间') INT(10)"`
	SysUpdated int    `gorm:"not null default 0 comment('修改时间') INT(10)"`
	SysIp      string `gorm:"not null default '' comment('IP地址') VARCHAR(50)"`
}
