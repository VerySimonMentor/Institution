package CRUD

import (
	"Institution/logs"
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
