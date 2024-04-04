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

func pageRange(page, pageNum, countryNum int) (int, int) {
	start := (page - 1) * pageNum
	end := start + pageNum
	if start > countryNum {
		start = countryNum
	}
	if end > countryNum {
		end = countryNum
	}
	return start, end
}

func checkCountryInRedis(ctx *gin.Context) []Country {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()
	countryList := make([]Country, 0)
	countryString, err := redisClient.LRange(context.Background(), "country", 0, -1).Result()
	if redis.CheckNil(err) {
		if err = mysqlClient.Find(&countryList).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "查询失败"})
			logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
			return nil
		}
		countryJSON, err := json.Marshal(countryList)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "json转换失败"})
			logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
			return nil
		}
		if err = redisClient.Set(context.Background(), "country", countryJSON, 0).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "redis存储失败"})
			logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
			return nil
		}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return nil
	} else {
		for _, country := range countryString {
			var countryStruct Country
			if err := json.Unmarshal([]byte(country), &countryStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
			countryList = append(countryList, countryStruct)
		}
	}

	return countryList
}
