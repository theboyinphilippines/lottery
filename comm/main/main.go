package main

import (
	"log"
	"lottery/comm"
)

func main() {
	timeEnd := comm.NowUnix() + 30*86400*12*4 //从现在加4年时间
	log.Println("timeEnd：", timeEnd)
}
