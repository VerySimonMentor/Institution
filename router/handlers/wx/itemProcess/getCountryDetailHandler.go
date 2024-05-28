package itemprocess

import (
	"Institution/config"
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"Institution/router/handlers/CRUD"
	"Institution/router/handlers/wx"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type CountryDetailResp struct {
	SchoolChiName string        `json:"schoolChiName"`
	SchoolEngName string        `json:"schoolEngName"`
	SchoolType    string        `json:"schoolType"`
	CountryItem   []CountryItem `json:"countryItems"`
}

type CountryItem struct {
	ItemName   string `json:"itemName"`
	ItemDetail string `json:"itemDetail"`
	ItemRemark string `json:"itemRemark"`
}

func GetCountryDetailHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	loginTocken := ctx.Query("loginTocken")
	countryListIndex := cast.ToInt64(ctx.Query("countryListIndex"))
	selectedProvinceMapStr := ctx.Query("selectedProvinceMap")
	selectedSchoolTypeMapStr := ctx.Query("selectedSchoolTypeMap")
	searchContent := ctx.Query("searchContent")
	check, phoneNumber := wx.CheckLoginTocken(wxConfig, loginTocken)
	if !check {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	mysqlClient := mysql.GetClient()
	var user mysql.UserSQL
	mysqlClient.Where("userNumber = ?", phoneNumber).First(&user)

	redisClient := redis.GetClient()
	system := CRUD.GetSystemInRedis(ctx)

	countryStr, err := redisClient.LIndex(context.Background(), "country", countryListIndex).Result()
	if err != nil {
		logs.GetInstance().Logger.Errorf("get country error %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var country CRUD.Country
	json.Unmarshal([]byte(countryStr), &country)
	countryDetailResp := make([]CountryDetailResp, 0)
	schoolKey := fmt.Sprintf(CRUD.SchoolKey, country.CountryId)
	selectedProvinceMap, selectedSchoolTypeMap := make(map[int]bool), make(map[int]bool)
	e := json.Unmarshal([]byte(selectedProvinceMapStr), &selectedProvinceMap)
	if e != nil {
		logs.GetInstance().Logger.Errorf("selectedProvinceMapStr %v", e)
	}
	json.Unmarshal([]byte(selectedSchoolTypeMapStr), &selectedSchoolTypeMap)
	logs.GetInstance().Logger.Infof("selectedProvinceMap %v", selectedProvinceMap)
	logs.GetInstance().Logger.Infof("selectedSchoolTypeMap %v", selectedSchoolTypeMap)

	schoolList := CRUD.GetSchoolInRedis(ctx, schoolKey, country.CountryAndSchool)
	for _, school := range schoolList {
		countryDetail := CountryDetailResp{}

		if _, ok := selectedProvinceMap[school.Province]; !ok && len(selectedProvinceMap) != 0 {
			continue
		}
		if _, ok := selectedSchoolTypeMap[school.SchoolType]; !ok && len(selectedSchoolTypeMap) != 0 {
			continue
		}
		if searchContent != "" {
			match1, err1 := regexp.MatchString(searchContent, school.SchoolChiName)
			match2, err2 := regexp.MatchString(searchContent, school.SchoolEngName)
			match3, err3 := regexp.MatchString(searchContent, school.SchoolAbbreviation)
			if err1 != nil || err2 != nil || err3 != nil {
				logs.GetInstance().Logger.Errorf("searchContent %v err1 %v err2 %v err3 %v", searchContent, err1, err2, err3)
				continue
			}
			logs.GetInstance().Logger.Infof("match1 %v match2 %v match3 %v", match1, match2, match3)
			if !match1 && !match2 && !match3 {
				continue
			}
		}

		countryDetail.SchoolChiName = school.SchoolChiName
		countryDetail.SchoolEngName = school.SchoolEngName
		if school.SchoolType > 0 {
			countryDetail.SchoolType = system.SchoolTypeList[school.SchoolType].SchoolTypeName
		}

		countryDetail.CountryItem = make([]CountryItem, len(school.SchoolAndItem))
		itemKey := fmt.Sprintf(CRUD.ItemKey, country.CountryId, school.SchoolId)
		itemList := CRUD.GetItemInRedis(ctx, itemKey, school.SchoolAndItem)
		for j, item := range itemList {
			countryDetail.CountryItem[j].ItemName = item.ItemName
			countryDetail.CountryItem[j].ItemRemark = item.ItemRemark
			var levelIndex int
			for i := len(item.LevelRate) - 1; i >= 0; i-- {
				if item.LevelRate[i].LevelId < user.UserLevel {
					if i == len(item.LevelRate)-1 {
						levelIndex = i
					} else {
						levelIndex = i + 1
					}
					break
				}
			}
			if len(item.LevelRate) > 0 && item.LevelRate[levelIndex].IfNotCombine {
				countryDetail.CountryItem[j].ItemDetail = item.LevelRate[levelIndex].LevelRate
			} else if len(item.LevelRate) > 0 {
				countryDetail.CountryItem[j].ItemDetail = fmt.Sprintf(item.LevelDescription, item.LevelRate[levelIndex].LevelRate)
			}
		}
		countryDetailResp = append(countryDetailResp, countryDetail)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"countryDetail": countryDetailResp,
		"schoolType":    system.SchoolTypeList,
		"province":      country.Province,
	})
}
