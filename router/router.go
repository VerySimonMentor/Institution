package router

import (
	crud "Institution/router/handlers/CRUD"

	"github.com/gin-gonic/gin"
)

// 注册路由
func RouterInit() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	ginRouter.LoadHTMLGlob("html/*")

	ginRouter.GET(("/"))
	ginRouter.GET("/country/create", crud.CreateCountryHandler)

	ginRouter.POST("/country/show", crud.ShowCountryHandler)

	ginRouter.DELETE("/country/delete", crud.DeleteCountryHandler)

	return ginRouter
}
