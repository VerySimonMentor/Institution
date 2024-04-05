package router

import (
	router "Institution/router/handlers/CRUD"

	"github.com/gin-gonic/gin"
)

// 注册路由
func RouterInit() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	ginRouter.LoadHTMLGlob("html/*")

	ginRouter.GET("/country/create", router.CreateCountryHandler)

	ginRouter.POST("/country/show", router.ShowCountryHandler)

	return ginRouter
}
