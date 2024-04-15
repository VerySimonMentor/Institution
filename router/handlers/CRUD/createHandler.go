package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"encoding/json"
	"fmt"
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
		CountryAndSchool: make([]int, 0),
		Province:         make([]Province, 0),
	})
	redisClient.RPush(context.Background(), "country", countryByte)
	countryNum := int(redisClient.LLen(context.Background(), "country").Val())
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

type CreateSchoolForm struct {
	CountryListIndex int64 `json:"countryListIndex"`
	PageNum          int   `json:"pageNum"`
}

func CreateSchoolHandler(ctx *gin.Context) {
	var createSchoolForm CreateSchoolForm
	if err := ctx.ShouldBindJSON(&createSchoolForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("CreateSchoolHandler error %s", err)
		return
	}
	if createSchoolForm.CountryListIndex < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		return
	}
	mysqlClient := mysql.GetClient()
	redisClient := redis.GetClient()
	schoolSQL := mysql.SchoolSQL{
		SchoolEngName:      "default",
		SchoolChiName:      "默认",
		SchoolAbbreviation: "默认",
		SchoolType:         0,
		Province:           "",
		OfficialWebLink:    "",
		SchoolRemark:       "",
		SchoolAndItem:      []byte("{}"),
	}

	if err := mysqlClient.Create(&schoolSQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateSchoolHandler error %s", err)
		return
	}
	if err := mysqlClient.Last(&schoolSQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateCountryHandler error %s", err)
		return
	}

	schoolByte, _ := json.Marshal(School{
		SchoolId:           schoolSQL.SchoolId,
		SchoolEngName:      schoolSQL.SchoolEngName,
		SchoolChiName:      schoolSQL.SchoolChiName,
		SchoolAbbreviation: schoolSQL.SchoolAbbreviation,
		SchoolType:         schoolSQL.SchoolType,
		Province:           schoolSQL.Province,
		OfficialWebLink:    schoolSQL.OfficialWebLink,
		SchoolRemark:       schoolSQL.SchoolRemark,
		SchoolAndItem:      make([]int, 0),
	})
	countryString, err := redisClient.LIndex(context.Background(), "country", createSchoolForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("CreateSchoolHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("CreateSchoolHandler error %s", err)
		return
	}
	logs.GetInstance().Logger.Info(country)
	country.CountryAndSchool = append(country.CountryAndSchool, schoolSQL.SchoolId)
	countryByte, _ := json.Marshal(country)
	redisClient.LSet(context.Background(), "country", createSchoolForm.CountryListIndex, countryByte)

	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	redisClient.RPush(context.Background(), schoolKey, schoolByte)
	schoolNum := int(redisClient.LLen(context.Background(), schoolKey).Val())
	var totalPage int
	if schoolNum%createSchoolForm.PageNum == 0 {
		totalPage = schoolNum / createSchoolForm.PageNum
	} else {
		totalPage = schoolNum/createSchoolForm.PageNum + 1
	}

	go func(country Country) {
		countryAndSchool, _ := json.Marshal(country.CountryAndSchool)
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.CountrySQL{}).Where("countryId = ?", country.CountryId).Update("countryAndSchool", countryAndSchool).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("CreateSchoolHandler error %s", err)
		}
	}(country)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"schoolId":  schoolSQL.SchoolId,
		"totalPage": totalPage,
	})
}
