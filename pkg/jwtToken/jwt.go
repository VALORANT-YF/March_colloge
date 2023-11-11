package jwtToken

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("张云飞") //密钥

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt 包自带的jwt.StandardClaims 只包含了官方字段 想使用user_id字段,所以要自定义结构体
type MyClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID string) (aToken, rToken string, err error) {
	//创建一个自己声明的数据
	c := MyClaims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //设置过期时间
			Issuer:    "college",                                  //签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	//refresh Token 不需要存自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 30).Unix(), //设置过期时间
			Issuer:    "college",                               //签发人
		},
	}).SignedString(mySecret)
	// 使用指定的secret 签名并获得完整的编码后的字符串 token
	return
}

// ParseToken 解析JWT 的Access Token
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析Token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { //校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
