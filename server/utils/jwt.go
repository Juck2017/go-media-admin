package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 10

var MySecret = []byte("夏天夏天悄悄过去")

// 生成token
func GenToken(username string) (string,error)  {
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer: username, // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	token,err := jwt.ParseWithClaims(tokenString,&MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret,nil
	})
	if err != nil {
		return nil,err
	}
	if claims, ok := token.Claims.(*MyClaims);ok && token.Valid {
		return claims, err
	}
	return nil, errors.New("token不可用或者已经失效")
}