package service

import (
	handlers "Sgrid/src/http"
	"Sgrid/src/public/jwt"
	"Sgrid/src/storage"
	"Sgrid/src/storage/vo"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Registry(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	GROUP.POST("/login", login)
	GROUP.POST("/loginByCache", jwt.SgridJwtValidate.Validate, loginByCache)
	GROUP.GET("/getUserMenusByUserId", getUserMenusByUserId)
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) //加密处理
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	v, err := storage.QueryUser(username)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	compareError := bcrypt.CompareHashAndPassword(hash, []byte(v.Password))
	if compareError != nil {
		handlers.AbortWithError(c, compareError.Error())
		return
	}

	expireTime := time.Hour * 12
	token, err := jwt.SgridJwtValidate.GenToken(username, time.Now().Add(expireTime))
	if err != nil {
		fmt.Println("gen token error")
		handlers.AbortWithError(c, err.Error())
		return
	}
	rsp := vo.VoUser{
		UserName:   v.UserName,
		Password:   v.Password,
		CreateTime: v.CreateTime,
		Id:         v.Id,
		Token:      token,
	}
	setTokenError := jwt.RdsToken.SetToken(token, expireTime, rsp)
	if err != nil {
		fmt.Println("set token error", err.Error())
		handlers.AbortWithError(c, setTokenError.Error())
		return
	}

	handlers.AbortWithSucc(c, rsp)
}

func loginByCache(c *gin.Context) {
	value, exists := c.Get(jwt.UserInfo)
	if !exists {
		handlers.AbortWithError(c, "login User Error")
	} else {
		var vo *vo.VoUser
		json.Unmarshal([]byte(value.(string)), &vo)
		handlers.AbortWithSucc(c, vo)
	}
}

func getUserMenusByUserId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	rm := storage.GetUserMenusByUserId(id)
	handlers.AbortWithSucc(c, rm)
}
