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
	logs.GetInstance().Logger.Info(start, end, totalPage)
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
	CountryListIndex int64 `json:"countryListIndex"`
	Page             int   `json:"page"`
	PageNum          int   `json:"pageNum"`
}

type SchoolResp struct {
	School
	ItemNum int `json:"itemNum"`
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
	countryString, err := redisClient.LIndex(ctx, "country", showSchoolForm.CountryListIndex).Result()
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

	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	logs.GetInstance().Logger.Info(schoolKey, country.CountryAndSchool)
	schoolList := getSchoolInRedis(ctx, schoolKey, country.CountryAndSchool)
	logs.GetInstance().Logger.Info(schoolList)
	if schoolList == nil {
		return
	}
	if len(schoolList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"results": []School{}})
		return
	}

	start, end := pageRange(showSchoolForm.Page, showSchoolForm.PageNum, len(schoolList))
	var totalPage int
	if len(schoolList)%showSchoolForm.PageNum == 0 {
		totalPage = len(schoolList) / showSchoolForm.PageNum
	} else {
		totalPage = len(schoolList)/showSchoolForm.PageNum + 1
	}
	schoolResp := make([]SchoolResp, end-start)
	for i := start; i < end; i++ {
		index := (showSchoolForm.Page-1)*showSchoolForm.PageNum + i - start
		schoolResp[i-start] = SchoolResp{
			School:  schoolList[index],
			ItemNum: len(schoolList[index].SchoolAndItem),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"results":    schoolResp,
		"schoolType": TypeList,
		"province":   country.Province,
		"totalPage":  totalPage,
	})
}

func InitSchoolHandler(ctx *gin.Context) {
	countryList := getCountryInRedis(ctx)
	allCountry := make([]string, 0)
	if countryList == nil {
		logs.GetInstance().Logger.Errorf("InitSchoolHandler error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
		return
	}

	for _, country := range countryList {
		allCountry = append(allCountry, country.CountryChiName)
	}

	ctx.JSON(http.StatusOK, gin.H{"results": allCountry})
}
