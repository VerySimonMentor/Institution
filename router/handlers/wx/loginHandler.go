package wx

import (
	"Institution/config"
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"Institution/router/handlers/CRUD"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FastLoginHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	loginCode := ctx.Query("loginCode")
	sessionKey := code2Session(wxConfig, loginCode)
	if sessionKey == "" {
		logs.GetInstance().Logger.Errorf("get session key error")
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}
	md5Hash := md5.Sum([]byte(sessionKey))
	loginTocken := hex.EncodeToString(md5Hash[:])
	redisClient := redis.GetClient()

	code := ctx.Query("code")
	phoneNumber := GetPhoneNumber(code, wxConfig)
	if phoneNumber == "" {
		logs.GetInstance().Logger.Errorf("get phone number error")
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	redisClient.Set(context.Background(), loginTocken, phoneNumber, loginTockenExpireTime*time.Hour*24)

	mysqlClient := mysql.GetClient()
	var userSQL mysql.UserSQL
	result := mysqlClient.Where("userNumber = ?", phoneNumber).First(&userSQL)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			userSQL.UserNumber = phoneNumber
			userSQL.UserLevel = 0
			if err := mysqlClient.Create(&userSQL).Error; err != nil {
				logs.GetInstance().Logger.Errorf("create user error %s", err)
				ctx.JSON(http.StatusBadRequest, gin.H{})
			}
			if err := mysqlClient.Last(&userSQL).Error; err != nil {
				logs.GetInstance().Logger.Errorf("get user info error %s", err)
				ctx.JSON(http.StatusBadRequest, gin.H{})
			}

			user := CRUD.User{
				UserId:       userSQL.UserId,
				UserNumber:   userSQL.UserNumber,
				UserLevel:    userSQL.UserLevel,
				UserAccount:  userSQL.UserAccount,
				UserPassWd:   userSQL.UserPassWd,
				UserEmail:    userSQL.UserEmail,
				StudentCount: userSQL.StudentCount,
			}
			userByte, _ := json.Marshal(user)
			redisClient.RPush(context.Background(), "user", userByte)
		} else {
			logs.GetInstance().Logger.Errorf("get user info error %s", result.Error)
			ctx.JSON(http.StatusBadRequest, gin.H{})
		}
	}

	var loginState int
	if userSQL.UserLevel == 0 {
		loginState = 1
	} else {
		loginState = 2
	}
	ctx.JSON(http.StatusOK, gin.H{
		"loginState":  loginState,
		"loginTocken": loginTocken,
	})
}

type LoginForm struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	LoginCode   string `json:"loginCode"`
}

func LoginHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	var loginForm LoginForm
	if err := ctx.ShouldBindJSON(&loginForm); err != nil {
		logs.GetInstance().Logger.Errorf("bind json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	mysqlClient := mysql.GetClient()
	var user mysql.UserSQL
	result := mysqlClient.Where("userNumber = ?", loginForm.PhoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"err": "用户不存在"})
			return
		} else {
			logs.GetInstance().Logger.Errorf("get user info error %s", result.Error)
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}

	if user.UserPassWd != loginForm.Password {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "密码错误"})
		return
	}

	sessionKey := code2Session(wxConfig, loginForm.LoginCode)
	if sessionKey == "" {
		logs.GetInstance().Logger.Errorf("get session key error")
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	md5Hash := md5.Sum([]byte(sessionKey))
	loginTocken := hex.EncodeToString(md5Hash[:])
	redisClient := redis.GetClient()
	redisClient.Set(context.Background(), loginTocken, loginForm.PhoneNumber, loginTockenExpireTime*time.Hour*24)

	var loginState int
	if user.UserLevel == 0 {
		loginState = 1
	} else {
		loginState = 2
	}
	ctx.JSON(http.StatusOK, gin.H{
		"loginState":  loginState,
		"loginTocken": loginTocken,
	})
}

func CheckLoginTockenHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	loginTocken := ctx.Query("loginTocken")
	if CheckLoginTocken(wxConfig, loginTocken) {
		redisClien := redis.GetClient()
		phone, err := redisClien.Get(context.Background(), loginTocken).Result()
		if err != nil {
			logs.GetInstance().Logger.Errorf("get phone number error %s", err)
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		mysqlClient := mysql.GetClient()
		var user mysql.UserSQL
		mysqlClient.Where("userNumber = ?", phone).First(&user)
		var loginState int
		if user.UserLevel == 0 {
			loginState = 1
		} else {
			loginState = 2
		}
		ctx.JSON(http.StatusOK, gin.H{
			"state":      true,
			"loginState": loginState,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"state":      false,
			"loginState": 0,
		})
	}
}
