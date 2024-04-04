package router

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
	var country mysql.CountrySQL
	if err := ctx.ShouldBindJSON(&country); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	mysqlClient := mysql.GetClient()
	redisClient := redis.GetClient()
	if err := mysqlClient.Create(&country).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	countryJSON, err := json.Marshal(country)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "json转换失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}
	err = redisClient.RPush(context.Background(), "country", countryJSON).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "redis存储失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
}
