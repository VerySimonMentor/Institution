package CRUD

import (
	"Institution/logs"
	"Institution/redis"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageShow struct {
	Page    int `json:"page"`
	PageNum int `json:"pageNum"`
}

type CountryResp struct {
	CountryId      int    `json:"countryId"`
	CountryChiName string `json:"countryChiName"`
	CountryEngName string `json:"countryEngName"`
	SchoolNum      int    `json:"schoolNum"`
	ProvinceNum    int    `json:"provinceNum"`
}

func ShowCountryHandler(ctx *gin.Context) {
	var pageShow PageShow
	if err := ctx.ShouldBindJSON(&pageShow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return
	}

	countryList := getCountryInRedis(ctx)
	if countryList == nil {
		return
	}

	countryNum := len(countryList)
	if countryNum == 0 {
		ctx.JSON(http.StatusOK, gin.H{"results": []CountryResp{}})
		return
	}
	start, end := pageRange(pageShow.Page, pageShow.PageNum, countryNum)
	var totalPage int
	if countryNum%pageShow.PageNum == 0 {
		totalPage = countryNum / pageShow.PageNum
	} else {
		totalPage = countryNum/pageShow.PageNum + 1
	}
	// logs.GetInstance().Logger.Infof("start: %d, end: %d", start, end)
	countryResp := make([]CountryResp, end-start)
	for i := start; i < end; i++ {
		index := (pageShow.Page-1)*pageShow.PageNum + i - start
		countryResp[i-start] = CountryResp{
			CountryId:      countryList[index].CountryId,
			CountryChiName: countryList[index].CountryChiName,
			CountryEngName: countryList[index].CountryEngName,
			SchoolNum:      len(countryList[index].CountryAndSchool),
			ProvinceNum:    len(countryList[index].Province),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"results": countryResp,
		"totalPage": totalPage,
	})
}

func ShowProvinceHandler(ctx *gin.Context) {
	var provinceForm CountryInstanceForm
	if err := ctx.ShouldBindJSON(&provinceForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("ShowProvinceHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryString, err := redisClient.LIndex(ctx, "country", provinceForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowProvinceHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("ShowProvinceHandler error %s", err)
		return
	}
	if country.CountryId != provinceForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("ShowProvinceHandler error %d != %d", country.CountryId, provinceForm.CountryId)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"country": country,
	})
}

type ShowSchoolForm struct {
	CountryId      int    `json:"countryId"`
	ListIndex      int64  `json:"listIndex"`
	CountryChiName string `json:"countryChiName"`
	Page           int    `json:"page"`
	PageNum        int    `json:"pageNum"`
}

func ShowSchoolHandler(ctx *gin.Context) {
	var showSchoolForm ShowSchoolForm
	if err := ctx.ShouldBindJSON(&showSchoolForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryList := getCountryInRedis(ctx)
	if countryList == nil {
		return
	}
	if len(countryList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"results": []School{}})
		return
	}
	countryString, err := redisClient.LIndex(ctx, "country", showSchoolForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return
	}
	if country.CountryId != showSchoolForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %d != %d", country.CountryId, showSchoolForm.CountryId)
		return
	}

	schoolKey := fmt.Sprintf(SchoolKey, showSchoolForm.CountryId)
	schoolList := getSchoolInRedis(ctx, schoolKey, country.CountryAndSchool)

	start, end := pageRange(showSchoolForm.Page, showSchoolForm.PageNum, len(schoolList))
	var totalPage int
	if len(schoolList)%showSchoolForm.PageNum == 0 {
		totalPage = len(schoolList) / showSchoolForm.PageNum
	} else {
		totalPage = len(schoolList)/showSchoolForm.PageNum + 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"results":   schoolList[start:end],
		"totalPage": totalPage,
	})
}
