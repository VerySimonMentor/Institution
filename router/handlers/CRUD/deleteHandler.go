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

type DeleteForm struct {
	CountryId int   `json:"countryId"`
	ListIndex int64 `json:"listIndex"`
}

func DeleteCountryHandler(ctx *gin.Context) {
	var deleteForm DeleteForm
	if err := ctx.ShouldBindJSON(&deleteForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	deleteCountryString, err := redisClient.LIndex(context.Background(), "country", deleteForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}
	var deleteCountry Country
	if err := json.Unmarshal([]byte(deleteCountryString), &deleteCountry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}
	if deleteCountry.CountryId != deleteForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}
	_, err = redisClient.LRem(context.Background(), "country", 0, deleteCountryString).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "redis删除失败"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}

	go func(countryId int) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Where("countryId = ?", countryId).Delete(&Country{}).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		}
	}(deleteForm.CountryId)

	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
