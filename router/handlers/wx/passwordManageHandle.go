package wx

import (
	"Institution/config"
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckPasswordHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	loginTocken := ctx.Query("loginTocken")
	check, phoneNumber := CheckLoginTocken(wxConfig, loginTocken)
	if !check {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	mysqlClient := mysql.GetClient()
	var user mysql.UserSQL
	mysqlClient.Where("userNumber = ?", phoneNumber).First(&user)
	var noPassWord bool
	if user.UserPassWd == "" {
		noPassWord = true
	} else {
		noPassWord = false
	}

	ctx.JSON(http.StatusOK, gin.H{
		"noPassWord":  noPassWord,
		"phoneNumber": phoneNumber,
	})
}

type InitPasswordForm struct {
	LoginTocken string `json:"loginTocken"`
	Password    string `json:"password"`
}

func InitPasswordHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	var passwordForm InitPasswordForm
	if err := ctx.ShouldBindJSON(&passwordForm); err != nil {
		logs.GetInstance().Logger.Errorf("bind json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	check, phoneNumber := CheckLoginTocken(wxConfig, passwordForm.LoginTocken)
	if !check {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	updatePassword(ctx, passwordForm.Password, phoneNumber)

	ctx.JSON(http.StatusOK, gin.H{
		"state": true,
	})
}

type NewPasswordForm struct {
	LoginTocken string `json:"loginTocken"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func NewPasswordHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	var passwordForm NewPasswordForm
	if err := ctx.ShouldBindJSON(&passwordForm); err != nil {
		logs.GetInstance().Logger.Errorf("bind json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	check, phoneNumber := CheckLoginTocken(wxConfig, passwordForm.LoginTocken)
	if !check {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	var user mysql.UserSQL
	mysqlClient := mysql.GetClient()
	mysqlClient.Where("userNumber = ?", phoneNumber).First(&user)
	if user.UserPassWd != passwordForm.OldPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"state": false,
		})
		return
	}

	updatePassword(ctx, passwordForm.NewPassword, phoneNumber)

	ctx.JSON(http.StatusOK, gin.H{
		"state": true,
	})
}

func updatePassword(ctx *gin.Context, newPassword, phoneNumber string) {
	md5Hash := md5.Sum([]byte(newPassword))
	hashedPassword := hex.EncodeToString(md5Hash[:])

	redisClient := redis.GetClient()
	allUsers, err := redisClient.LRange(context.Background(), "user", 0, -1).Result()
	if err != nil {
		logs.GetInstance().Logger.Errorf("get user list error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	for i, user := range allUsers {
		var userSQL mysql.UserSQL
		json.Unmarshal([]byte(user), &userSQL)
		if userSQL.UserNumber == phoneNumber {
			userSQL.UserPassWd = hashedPassword
			userByte, _ := json.Marshal(userSQL)
			redisClient.LSet(context.Background(), "user", int64(i), userByte)
			break
		}
	}

	go func(hashedPassword string) {
		mysqlClient := mysql.GetClient()
		mysqlClient.Model(&mysql.UserSQL{}).Where("userNumber = ?", phoneNumber).Update("userPassWd", hashedPassword)
	}(hashedPassword)
}
