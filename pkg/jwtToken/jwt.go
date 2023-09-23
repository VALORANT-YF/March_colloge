package jwtToken

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("夏天") //密钥

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt 包自带的jwt.StandardClaims 只包含了官方字段 想使用user_id字段,所以要自定义结构体
type MyClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	//创建一个自己声明的数据
	c := MyClaims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //设置过期时间
			Issuer:    "bluebell",                                 //签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret 签名并获得完整的编码后的字符串 token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析Token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret , nil
	})
	if err != nil {
		return nil , err
	}
	if token.Valid { //校验token
		return mc , nil
	}
	return nil , errors.New("invalid token")
}