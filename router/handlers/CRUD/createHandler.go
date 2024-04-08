package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCountryHandler(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}
	mysqlClient := mysql.GetClient()
	redisClient := redis.GetClient()
	countrySQL := mysql.CountrySQL{
		CountryEngName:   "default",
		CountryChiName:   "默认",
		CountryAndSchool: []byte("{}"),
		Province:         []byte("{}"),
	}

	if err := mysqlClient.Create(&countrySQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	if err := mysqlClient.Last(&countrySQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}
	// if !checkCountryInRedis(ctx) {

	// }
	countryByte, _ := json.Marshal(Country{
		CountryId:        countrySQL.CountryId,
		CountryEngName:   countrySQL.CountryEngName,
		CountryChiName:   countrySQL.CountryChiName,
		CountryAndSchool: make(map[int]struct{}),
		Province:         make(map[string]struct{}),
	})
	redisClient.RPush(context.Background(), "country", countryByte)
	countryList := getCountryInRedis(ctx)
	countryNum := len(countryList)
	var totalPage int
	if countryNum%pageNum == 0 {
		totalPage = countryNum / pageNum
	} else {
		totalPage = countryNum/pageNum + 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"countryId": countrySQL.CountryId,
		"totalPage": totalPage})
}
