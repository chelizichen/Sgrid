package jwt

import (
	"Sgrid/src/configuration"
	h "Sgrid/src/http"
	"Sgrid/src/storage/vo"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

var jwtKey = []byte("____Sgrid.src.public.jwt______")
var UserInfo = "userInfo"

var GenToken = func(username string, expiresAt time.Time) (string, error) {
	fmt.Println("username", username)
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// 解析从前端获取到的token值
var ParseToken = func(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

func Validate(c *gin.Context) {
	s := c.Request.Header["Token"]
	if len(s) == 0 {
		h.AbortWithError(c, "empty token")
		return
	}
	value, err := ValidateToken(s)
	if value == "" || err != nil {
		h.AbortWithError(c, "token validate error")
	} else {
		c.Set(UserInfo, value)
		c.Next()
	}
}

func ValidateToken(values []string) (val string, err error) {
	for _, v := range values {
		s := GetToken(v)
		if s != "" {
			return s, nil
		}
		continue
	}
	return "", nil
}

func SetToken(Key string, expireTime time.Duration, val vo.VoUser) error {
	b, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Error to Marshal", err.Error(), val)
		return err
	}
	return configuration.GRDB.SetEX(configuration.RDBContext, Key, b, expireTime).Err()
}

func GetToken(Key string) string {
	s, err := configuration.GRDB.Get(configuration.RDBContext, Key).Result()
	if err != nil {
		fmt.Println("get  key Error", err.Error())
		return ""
	}
	return s
}
