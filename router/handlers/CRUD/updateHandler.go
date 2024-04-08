package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateCountryForm struct {
	CountryId   int         `json:"countryId"`
	ListIndex   int64       `json:"listIndex"`
	UpdateField string      `json:"updateField"`
	UpdateValue interface{} `json:"updateValue"`
}

func UpdateCountryHandler(ctx *gin.Context) {
	var updateCountryForm UpdateCountryForm
	if err := ctx.ShouldBindJSON(&updateCountryForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	updateCountryString, err := redisClient.LIndex(ctx, "country", updateCountryForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		return
	}
	var updateCountry mysql.CountrySQL
	if err := json.Unmarshal([]byte(updateCountryString), &updateCountry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		return
	}
	if updateCountry.CountryId != updateCountryForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		return
	}

	switch updateCountryForm.UpdateField {
	case "countryEngName":
		updateCountry.CountryEngName = updateCountryForm.UpdateValue.(string)
	case "countryChiName":
		updateCountry.CountryChiName = updateCountryForm.UpdateValue.(string)
	}
	updateCountryByte, _ := json.Marshal(updateCountry)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
	// 	logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
	// 	return
	// }
	_, err = redisClient.LSet(ctx, "country", updateCountryForm.ListIndex, updateCountryByte).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		return
	}

	go func(updateCountryForm UpdateCountryForm) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.CountrySQL{}).Where("countryId = ?", updateCountryForm.CountryId).UpdateColumn(updateCountryForm.UpdateField, updateCountryForm.UpdateValue.(string)).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		}
	}(updateCountryForm)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type UpdateProvinceForm struct {
	CountryId      int                 `json:"countryId"`
	ListIndex      int64               `json:"listIndex"`
	CountryEngName string              `json:"countryEngName"`
	CountryChiName string              `json:"countryChiName"`
	Province       map[string]struct{} `json:"province"`
}

func UpdateProvinceHandler(ctx *gin.Context) {
	var updateCountryForm UpdateProvinceForm
	if err := ctx.ShouldBindJSON(&updateCountryForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	updateCountryString, err := redisClient.LIndex(ctx, "country", updateCountryForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		return
	}
	var updateCountry mysql.CountrySQL
	if err := json.Unmarshal([]byte(updateCountryString), &updateCountry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		return
	}
	if updateCountry.CountryId != updateCountryForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		return
	}
	updateCountry.CountryEngName = updateCountryForm.CountryEngName
	updateCountry.CountryChiName = updateCountryForm.CountryChiName
	updateCountry.Province, _ = json.Marshal(updateCountryForm.Province)
	updateCountryByte, _ := json.Marshal(updateCountry)
	err = redisClient.LSet(ctx, "country", updateCountryForm.ListIndex, updateCountryByte).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		return
	}

	go func(updateCountry mysql.CountrySQL) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.CountrySQL{}).Where("countryId = ?", updateCountry.CountryId).Updates(&updateCountry).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		}
	}(updateCountry)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}
