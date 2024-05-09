package wx

import (
	"Institution/config"
	"Institution/logs"
	"Institution/mysql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FastLoginHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	code := ctx.Query("code")
	phoneNumber := GetPhoneNumber(code, wxConfig)
	if phoneNumber == "" {
		logs.GetInstance().Logger.Errorf("get phone number error")
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	mysqlClient := mysql.GetClient()
	var user mysql.UserSQL
	result := mysqlClient.Where("userNumber = ?", phoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			user.UserNumber = phoneNumber
			user.UserLevel = 0
			if err := mysqlClient.Create(&user).Error; err != nil {
				logs.GetInstance().Logger.Errorf("create user error %s", err)
				ctx.JSON(http.StatusBadRequest, gin.H{})
			}
		} else {
			logs.GetInstance().Logger.Errorf("get user info error %s", result.Error)
			ctx.JSON(http.StatusBadRequest, gin.H{})
		}
	}

	var loginState int
	if user.UserLevel == 0 {
		loginState = 1
	} else {
		loginState = 2
	}
	ctx.JSON(http.StatusOK, gin.H{
		"loginState": loginState,
	})
}

type LoginForm struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func LoginHandler(ctx *gin.Context) {
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

	var loginState int
	if user.UserLevel == 0 {
		loginState = 1
	} else {
		loginState = 2
	}
	ctx.JSON(http.StatusOK, gin.H{
		"loginState": loginState,
	})
}
