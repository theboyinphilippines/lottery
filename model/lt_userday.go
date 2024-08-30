package models

// 用户今日参与次数
type LtUserday struct {
	Id         int `gorm:"not null pk autoincr INT(10)"`
	Uid        int `gorm:"not null default 0 comment('用户ID') INT(10)"`
	Day        int `gorm:"not null default 0 comment('日期，如：20180725') INT(10)"`
	Num        int `gorm:"not null default 0 comment('次数') INT(10)"`
	SysCreated int `gorm:"not null default 0 comment('创建时间') INT(10)"`
	SysUpdated int `gorm:"not null default 0 comment('修改时间') INT(10)"`
}
