package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"encoding/json"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	if err := mysqlClient.Last(&country).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}
	if !checkCountryInRedis(ctx) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Error("CreateCountryHandler error")
		return
	}
	countryByte, _ := json.Marshal(country)
	redisClient.RPush(context.Background(), "country", countryByte)
	ctx.JSON(http.StatusOK, gin.H{"msg": "创建成功",
		"countryId": country.CountryId})
}
