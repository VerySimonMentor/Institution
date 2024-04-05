package router

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCountryHandler(ctx *gin.Context) {
	mysqlClient := mysql.GetClient()
	redisClient := redis.GetClient()
	country := Country{
		CountryEngName:   "default",
		CountryChiName:   "默认",
		CountryAndSchool: make(map[int]struct{}),
		Province:         make(map[string]struct{}),
	}

	if err := mysqlClient.Create(&country).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	if err := mysqlClient.Last(&country).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}
	if !checkCountryInRedis(ctx) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "redis查询失败"})
		logs.GetInstance().Logger.Error("CreateCountryHandler error")
		return
	}
	redisClient.RPush(context.Background(), "country", country)
	ctx.JSON(http.StatusOK, gin.H{"msg": "创建成功",
		"countryId": country.CountryId})
}
