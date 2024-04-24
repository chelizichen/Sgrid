package jwt

import (
	"Sgrid/src/storage/pojo"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

var jwtKey = []byte("____Sgrid@Jwt______")

var GenToken = func(user pojo.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: uint(user.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析从前端获取到的token值
var ParseToken = func(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
