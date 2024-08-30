package main

import (
	"github.com/gin-gonic/gin"
	"lottery/internal/controller"
	"lottery/internal/service"
	"lottery/middleware"
)

func main() {
	r := gin.Default()
	//go controller.FetchPacketListMoneyWithChannel()

	indexCtl := controller.IndexController{
		ServiceGift:    service.NewGiftService(),
		ServiceBlackIp: service.NewBlackipService(),
		ServiceUser:    service.NewUserService(),
	}
	r.POST("/login", indexCtl.Login)                      //登录
	r.GET("/getAllRedPacket", indexCtl.GetAllRedPacket)   //微博多个红包，获取所有红包
	r.GET("/deliverRedPacket", indexCtl.DeliverRedPacket) //微博多个红包，发红包
	r.GET("/fetchRedPacket", indexCtl.FetchRedPacket)     //微博多个红包，抢红包
	r.Use(middleware.JWTAuth())
	r.GET("/ping", indexCtl.GetLucky) //抽奖
	r.Run()                           // listen and serve on 0.0.0.0:8080
}
