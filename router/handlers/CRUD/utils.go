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

type InstanceForm struct {
	CountryId int   `json:"countryId"`
	ListIndex int64 `json:"listIndex"`
}

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

func getCountryInRedis(ctx *gin.Context) []Country {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()
	countryListSQL := make([]mysql.CountrySQL, 0)
	if !checkCountryInRedis(ctx) {
		if err := mysqlClient.Model(&mysql.CountrySQL{}).Find(&countryListSQL).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
			logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
			return nil
		}
		if len(countryListSQL) == 0 {
			return make([]Country, 0)
		}
		countryList := make([]Country, len(countryListSQL))
		for i, country := range countryListSQL {
			countryAndSchool := make(map[int]struct{})
			province := make(map[string]struct{})
			json.Unmarshal(country.CountryAndSchool, &countryAndSchool)
			json.Unmarshal(country.Province, &province)
			countryList[i] = Country{
				CountryId:        country.CountryId,
				CountryEngName:   country.CountryEngName,
				CountryChiName:   country.CountryChiName,
				CountryAndSchool: countryAndSchool,
				Province:         province,
			}
		}
		countryJSON := make([][]byte, len(countryList))
		for i, country := range countryList {
			var err error
			countryJSON[i] = make([]byte, 0)
			countryJSON[i], err = json.Marshal(country)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
			if err = redisClient.RPush(context.Background(), "country", countryJSON[i]).Err(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
		}
	}

	countryString, err := redisClient.LRange(context.Background(), "country", 0, -1).Result()
	countryList := make([]Country, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return nil
	} else {
		for _, country := range countryString {
			var countryStruct Country
			if err := json.Unmarshal([]byte(country), &countryStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
			countryList = append(countryList, countryStruct)
		}
	}

	return countryList
}

func checkCountryInRedis(ctx *gin.Context) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), "country").Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return false
	}
	if keytype != "list" {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis类型错误"})
		// logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return false
	}

	return true
}
