package models

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID          uint
	UserName    string
	AuthorityId uint // 角色id
	jwt.StandardClaims
}
