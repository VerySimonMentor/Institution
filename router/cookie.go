package router

import (
	"Institution/config"
	"Institution/logs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CookieVerify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		adminConfig := config.GetServerConfig().Admin

		username, err := ctx.Cookie("username")
		if err != nil || username != adminConfig.Username {
			ctx.Redirect(http.StatusFound, "/login")
			ctx.Abort()
			logs.GetInstance().Logger.Warnf("user %s cookie false: %s", username, err)
			return
		}

		password, err := ctx.Cookie("password")
		if err != nil || password != adminConfig.Password {
			ctx.Redirect(http.StatusFound, "/login")
			ctx.Abort()
			logs.GetInstance().Logger.Warnf("user %s cookie false: %s", username, err)
			return
		}
		ctx.Next()
	}
}
