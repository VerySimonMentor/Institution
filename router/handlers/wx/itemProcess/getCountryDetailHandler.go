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
	json.Unmarshal([]byte(selectedProvinceMapStr), &selectedProvinceMap)
	json.Unmarshal([]byte(selectedSchoolTypeMapStr), &selectedSchoolTypeMap)

	for i := range country.CountryAndSchool {
		countryDetail := CountryDetailResp{}

		schoolStr, err := redisClient.LIndex(context.Background(), schoolKey, cast.ToInt64(i)).Result()
		if err != nil {
			logs.GetInstance().Logger.Errorf("get school error %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		var school CRUD.School
		json.Unmarshal([]byte(schoolStr), &school)
		if _, ok := selectedProvinceMap[school.Province]; !ok && len(selectedProvinceMap) != 0 {
			continue
		}
		if _, ok := selectedSchoolTypeMap[school.SchoolType]; !ok && len(selectedSchoolTypeMap) != 0 {
			continue
		}
		if searchContent != "" {
			match1, _ := regexp.MatchString(searchContent, school.SchoolChiName)
			match2, _ := regexp.MatchString(searchContent, school.SchoolEngName)
			match3, _ := regexp.MatchString(searchContent, school.SchoolAbbreviation)
			if !match1 && !match2 && !match3 {
				continue
			}
		}

		countryDetail.SchoolChiName = school.SchoolChiName
		countryDetail.SchoolEngName = school.SchoolEngName
		countryDetail.SchoolType = system.SchoolTypeList[school.SchoolType].SchoolTypeName

		countryDetail.CountryItem = make([]CountryItem, len(school.SchoolAndItem))
		for j := range school.SchoolAndItem {
			itemStr, err := redisClient.LIndex(context.Background(), fmt.Sprintf(CRUD.ItemKey, country.CountryId, school.SchoolId), cast.ToInt64(j)).Result()
			if err != nil {
				logs.GetInstance().Logger.Errorf("get item error %v", err)
				ctx.JSON(http.StatusBadRequest, gin.H{})
				return
			}

			var item CRUD.Item
			json.Unmarshal([]byte(itemStr), &item)
			countryDetail.CountryItem[j].ItemName = item.ItemName
			countryDetail.CountryItem[j].ItemRemark = item.ItemRemark
			var levelIndex int
			for i := len(item.LevelRate) - 1; i >= 0; i++ {
				if item.LevelRate[i].LevelId < user.UserLevel {
					levelIndex = i + 1
					break
				}
			}
			if item.LevelRate[levelIndex].IfNotCombine {
				countryDetail.CountryItem[j].ItemDetail = item.LevelRate[levelIndex].LevelRate
			} else {
				countryDetail.CountryItem[j].ItemDetail = fmt.Sprintf(item.LevelDescription, item.LevelRate[levelIndex].LevelRate)
			}
			countryDetailResp = append(countryDetailResp, countryDetail)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"countryDetail": countryDetailResp,
		"schoolType":    system.SchoolTypeList,
		"province":      country.Province,
	})
}
