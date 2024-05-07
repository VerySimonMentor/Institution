package router

import (
	"Institution/config"
	"Institution/router/handlers"
	crud "Institution/router/handlers/CRUD"
	"Institution/router/handlers/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册路由
func RouterInit(config *config.Config, rootPath string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	ginRouter.LoadHTMLGlob(rootPath + "/html/*.html")
	ginRouter.Static("/static", rootPath+"/html/static")
	ginRouter.Static("/script", rootPath+"/html/script")

	ginRouter.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/login")
	})
	ginRouter.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})
	ginRouter.POST("/login", func(ctx *gin.Context) {
		user.LoginHandler(ctx, &config.Admin)
	})
	ginRouter.Use(CookieVerify())

	ginRouter.GET("/manage", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "manage.html", gin.H{})
	})
	ginRouter.GET("/country/create", crud.CreateCountryHandler)
	ginRouter.GET("/school/initPage", crud.InitSchoolHandler)
	ginRouter.GET("/item/getSchool", crud.InitItemHandler)
	ginRouter.GET("/user/create", crud.CreateUserHandler)
	ginRouter.GET("/system/show", crud.ShowSystemHandler)
	ginRouter.GET("/system/create", crud.CreateSystemHandler)
	ginRouter.GET("/flush", handlers.FlushRedisHandler)

	ginRouter.POST("/country/show", crud.ShowCountryHandler)
	ginRouter.POST("/changeCountry", crud.UpdateCountryHandler)
	ginRouter.POST("/country/changeProvince/show", crud.ShowProvinceHandler)
	ginRouter.POST("/country/changeProvince/save", crud.UpdateProvinceHandler)
	ginRouter.POST("/country/editSchool", crud.ShowSchoolHandler)
	ginRouter.POST("/school/create", crud.CreateSchoolHandler)
	ginRouter.POST("/school/change", crud.UpdateSchoolHandler)
	ginRouter.POST("/school/editItem", crud.ShowItemHandler)
	ginRouter.POST("/item/create", crud.CreateItemHandler)
	ginRouter.POST("/item/paste", crud.PasteItemHandler)
	ginRouter.POST("/item/change", crud.UpdateItemHandler)
	ginRouter.POST("/item/changeProportion/show", crud.ShowLevelHandler)
	ginRouter.POST("/item/changeProportion/save", crud.UpdateLevelHandler)
	ginRouter.POST("/user/show", crud.ShowUserHandler)
	ginRouter.POST("/user/changeUser", crud.UpdateUserHandler)
	ginRouter.POST("/system/change", crud.UpdateSystemHandler)

	ginRouter.DELETE("/country/delete", crud.DeleteCountryHandler)
	ginRouter.DELETE("/school/delete", crud.DeleteSchoolHandler)
	ginRouter.DELETE("/item/delete", crud.DeleteItemHandler)
	ginRouter.DELETE("/user/delete", crud.DeleteUserHandler)
	ginRouter.DELETE("/system/delete", crud.DeleteSystemHandler)

	return ginRouter
}
