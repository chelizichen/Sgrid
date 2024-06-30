package jwt

import (
	h "Sgrid/src/http"
	"Sgrid/src/pool"
	"Sgrid/src/storage/vo"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	JwtTokenKey = "____Sgrid.src.public.jwt______"
	JwtUserInfo = "userInfo"
)

var SgridJwtValidate = new(sgridJwtValidate)
var RdsToken = new(rdsToken)

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

var jwtKey = []byte(JwtTokenKey)
var UserInfo = JwtUserInfo

type sgridJwtValidate struct{}

func (s *sgridJwtValidate) GenToken(username string, expiresAt time.Time) (string, error) {
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

func (s *sgridJwtValidate) ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

func (s *sgridJwtValidate) Validate(c *gin.Context) {
	tokenString := c.Request.Header["Token"]
	if len(tokenString) == 0 {
		h.AbortWithError(c, "empty token")
		return
	}
	value, err := s.validateToken(tokenString)
	if value == "" || err != nil {
		h.AbortWithError(c, "token validate error")
	} else {
		c.Set(UserInfo, value)
		c.Next()
	}
}

func (s *sgridJwtValidate) validateToken(values []string) (val string, err error) {
	for _, v := range values {
		s := RdsToken.GetToken(v)
		if s != "" {
			return s, nil
		}
		continue
	}
	return "", nil
}

type rdsToken struct{}

func (r *rdsToken) SetToken(Key string, expireTime time.Duration, val vo.VoUser) error {
	b, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Error to Marshal", err.Error(), val)
		return err
	}
	return pool.GRDB.SetEX(pool.RDBContext, Key, b, expireTime).Err()
}

func (r *rdsToken) GetToken(Key string) string {
	s, err := pool.GRDB.Get(pool.RDBContext, Key).Result()
	if err != nil {
		fmt.Println("get  key Error", err.Error())
		return ""
	}
	return s
}
