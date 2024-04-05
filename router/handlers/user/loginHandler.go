package user

import (
	"Institution/config"
	"Institution/logs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(ctx *gin.Context, adminConfig *config.AdminConfig) {
	var loginForm LoginForm
	if err := ctx.ShouldBindJSON(&loginForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("LoginHandler error %s", err)
		return
	}

	if loginForm.Username != adminConfig.Username || loginForm.Password != adminConfig.Password {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "用户名或密码错误"})
		logs.GetInstance().Logger.Info("login failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功"})
}
