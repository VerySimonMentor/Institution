package router

import (
	"github.com/gin-gonic/gin"
)

// 注册路由
func RouterInit() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	ginRouter.LoadHTMLGlob("html/*")

	return ginRouter
}
