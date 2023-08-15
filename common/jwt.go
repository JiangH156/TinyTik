package common

import (
	"TinyTik/utils/logger"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserName string
	jwt.RegisteredClaims
}

// JWT过期时间
const TokenExpireDuration = time.Hour * 24

// 用于签名的字符串
var jwtScret = []byte("TinyTik")

func GenToken(UserName string) (string, error) {
	//创建一个自己的声明
	claims := UserClaims{
		UserName,
		jwt.RegisteredClaims{
			//发行者
			Issuer: "TinyTik",
			//主题
			Subject: "TinyTik_TOKEN",
			//过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			//生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			//发布时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	logger.Fatal(token)
	//使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString(jwtScret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*UserClaims, error) {
	//自定义claims使用
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtScret, nil
	})
	if err != nil {
		return nil, err
	}
	// 判断类型是否正常和token是否有效
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("权限不够")
}
