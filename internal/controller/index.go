package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"lottery/global"
	"lottery/internal/service"
	"lottery/middleware"
	models "lottery/model"
	"net/http"
	"time"
)

type IndexController struct {
	ServiceGift    service.GiftService
	ServiceBlackIp service.BlackipService
	ServiceUser    service.UserService
}

func (i *IndexController) Login(c *gin.Context) {
	userName := c.PostForm("userName")
	mobile := c.PostForm("mobile")
	//登录成功
	var userInfo models.LtUser
	global.DB.Where("username=? and mobile =? ", userName, mobile).First(&userInfo)
	if userInfo.Username == "haha" && userInfo.Mobile == "18754585874" {
		// 生成token
		j := middleware.NewJWT()
		token, err := j.CreateToken(models.CustomClaims{
			ID:       uint(userInfo.Id),
			UserName: userInfo.Username,
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),               //生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
				Issuer:    "shy",
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"msg": "生成token失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":        userInfo.Id,
			"user_name": userInfo.Username,
			"token":     token,
			"expire_at": (time.Now().Unix() + 60*60*24*30) * 1000, //毫秒级别
		})

	}
}
