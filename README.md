## go抽奖，发红包，抢红包
### 技术栈：gin+gorm+mysql+redis
入口：main.go
/ping：抽奖接口，使用redis分布式锁，redis string， hash， set缓存
/deliverRedPacket：发红包接口，使用互斥锁，sync.Map, channel, 分布式锁来保证共享变量并发安全
/fetchRedPacket：抢红包接口，使用互斥锁，sync.Map, channel, 分布式锁来保证共享变量并发安全
