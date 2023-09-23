package jwtToken

import (
	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt 包自带的jwt.StandardClaims 只包含了官方字段 想使用user_id字段,所以要自定义结构体
type MyClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}
