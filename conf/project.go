package conf

import "time"

const (
	IpLimitMax   = 10000
	UserPrizeMax = 6000
)

const GtypeVirtual = 0   // 虚拟币
const GtypeCodeSame = 1  // 虚拟券，相同的码
const GtypeCodeDiff = 2  // 虚拟券，不同的码
const GtypeGiftSmall = 3 // 实物小奖
const GtypeGiftLarge = 4 // 实物大奖

var SysTimeLocation, _ = time.LoadLocation("Asia/Chongqing")
