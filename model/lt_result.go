package models

// 中奖记录表
type LtResult struct {
	Id         int    `gorm:"not null pk autoincr INT(10)" json:"-"`
	GiftId     int    `gorm:"not null default 0 comment('奖品ID，关联lt_gift表') INT(10)" json:"gift_id"`
	GiftName   string `gorm:"not null default '' comment('奖品名称') VARCHAR(255)" json:"gift_name"`
	GiftType   int    `gorm:"not null default 0 comment('奖品类型，同lt_gift. gtype') INT(10)" json:"gift_type"`
	Uid        int    `gorm:"not null default 0 comment('用户ID') INT(10)" json:"uid"`
	Username   string `gorm:"not null default '' comment('用户名') VARCHAR(50)" json:"username"`
	PrizeCode  int    `gorm:"not null default 0 comment('抽奖编号（4位的随机数）') INT(10)" json:"-"`
	GiftData   string `gorm:"not null default '' comment('获奖信息') VARCHAR(255)" json:"-"`
	SysCreated int    `gorm:"not null default 0 comment('创建时间') INT(10)" json:"-"`
	SysIp      string `gorm:"not null default '' comment('用户抽奖的IP') VARCHAR(50)" json:"-"`
	SysStatus  int    `gorm:"not null default 0 comment('状态，0 正常，1删除，2作弊') SMALLINT(5)" json:"-"`
}
