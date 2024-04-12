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

type CountryInstanceForm struct {
	CountryId int   `json:"countryId"`
	ListIndex int64 `json:"listIndex"`
}

type SchoolInstanceForm struct {
	CountryListIndex int64 `json:"countryListIndex"`
	SchoolId         int   `json:"schoolId"`
	SchoolListIndex  int64 `json:"schoolListIndex"`
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

	if !checkCountryInRedis(ctx) {
		countryListSQL := make([]mysql.CountrySQL, 0)
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
			countryAndSchool := make([]int, 0)
			province := make([]Province, 0)
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

func getSchoolInRedis(ctx *gin.Context, schoolKey string, countryAndSchool []int) []School {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()

	if !checkSchoolInRedis(ctx, schoolKey) {
		schoolListSQL := make([]mysql.SchoolSQL, len(countryAndSchool))
		for i, schoolId := range countryAndSchool {
			school := mysql.SchoolSQL{}
			if err := mysqlClient.Model(&mysql.SchoolSQL{}).Where("schoolId = ?", schoolId).Find(&school).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
				logs.GetInstance().Logger.Errorf("get school in redis error %s", err)
				return nil
			}
			schoolListSQL[i] = school
		}
		if len(schoolListSQL) == 0 {
			return make([]School, 0)
		}

		schoolList := make([]School, len(schoolListSQL))
		for i, school := range schoolListSQL {
			schoolAndItem := make([]int, 0)
			json.Unmarshal(school.SchoolAndItem, &schoolAndItem)
			schoolList[i] = School{
				SchoolId:           school.SchoolId,
				SchoolEngName:      school.SchoolEngName,
				SchoolChiName:      school.SchoolChiName,
				SchoolAbbreviation: school.SchoolAbbreviation,
				SchoolType:         school.SchoolType,
				Province:           school.Province,
				OfficialWebLink:    school.OfficialWebLink,
				SchoolRemark:       school.SchoolRemark,
				SchoolAndItem:      schoolAndItem,
			}
		}
		schoolJSON := make([][]byte, len(schoolList))
		for i, school := range schoolList {
			var err error
			schoolJSON[i] = make([]byte, 0)
			schoolJSON[i], err = json.Marshal(school)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
				return nil
			}
			if err = redisClient.RPush(context.Background(), schoolKey, schoolJSON[i]).Err(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
				logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
				return nil
			}
		}
	}

	schoolString, err := redisClient.LRange(context.Background(), schoolKey, 0, -1).Result()
	schoolList := make([]School, len(countryAndSchool))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return nil
	} else {
		for i, school := range schoolString {
			var schoolStruct School
			if err := json.Unmarshal([]byte(school), &schoolStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
				return nil
			}
			schoolList[i] = schoolStruct
		}
	}

	return schoolList
}

func checkSchoolInRedis(ctx *gin.Context, schoolKey string) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), schoolKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return false
	}
	if keytype != "list" {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis类型错误"})
		// logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return false
	}

	return true
}
