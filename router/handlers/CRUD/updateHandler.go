package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"encoding/json"
	"fmt"
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
	var updateCountry Country
	if err := json.Unmarshal([]byte(updateCountryString), &updateCountry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %s", err)
		return
	}
	if updateCountry.CountryId != updateCountryForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateCountryHandler error %d != %d", updateCountry.CountryId, updateCountryForm.CountryId)
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
	CountryId      int        `json:"countryId"`
	ListIndex      int64      `json:"listIndex"`
	CountryEngName string     `json:"countryEngName"`
	CountryChiName string     `json:"countryChiName"`
	Province       []Province `json:"province"`
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
	var updateCountry Country
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
	updateCountry.Province = updateCountryForm.Province
	updateCountryByte, _ := json.Marshal(updateCountry)
	err = redisClient.LSet(ctx, "country", updateCountryForm.ListIndex, updateCountryByte).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		return
	}

	countryAndSchoolByte, _ := json.Marshal(updateCountry.CountryAndSchool)
	provinceByte, _ := json.Marshal(updateCountryForm.Province)
	updateCountrySQL := mysql.CountrySQL{
		CountryId:        updateCountryForm.CountryId,
		CountryEngName:   updateCountryForm.CountryEngName,
		CountryChiName:   updateCountryForm.CountryChiName,
		CountryAndSchool: countryAndSchoolByte,
		Province:         provinceByte,
	}
	go func(updateCountry mysql.CountrySQL) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.CountrySQL{}).Where("countryId = ?", updateCountry.CountryId).Updates(&updateCountry).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateProvinceHandler error %s", err)
		}
	}(updateCountrySQL)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type UpdateSchoolForm struct {
	CountryListIndex int64       `json:"countryListIndex"`
	SchoolId         int         `json:"schoolId"`
	SchoolListIndex  int64       `json:"schoolListIndex"`
	UpdateField      string      `json:"updateField"`
	UpdateValue      interface{} `json:"updateValue"`
}

func UpdateSchoolHandler(ctx *gin.Context) {
	var updateSchoolForm UpdateSchoolForm
	if err := ctx.ShouldBindJSON(&updateSchoolForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryString, err := redisClient.LIndex(ctx, "country", updateSchoolForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		return
	}

	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	schoolString, err := redisClient.LIndex(ctx, schoolKey, updateSchoolForm.SchoolListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		return
	}
	var school School
	if err := json.Unmarshal([]byte(schoolString), &school); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		return
	}
	if school.SchoolId != updateSchoolForm.SchoolId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %d != %d", school.SchoolId, updateSchoolForm.SchoolId)
		return
	}

	updateValue := updateSchoolForm.UpdateValue.(string)
	switch updateSchoolForm.UpdateField {
	case "schoolEngName":
		school.SchoolEngName = updateValue
	case "schoolChiName":
		school.SchoolChiName = updateValue
	case "schoolAbbreviation":
		school.SchoolAbbreviation = updateValue
	case "schoolType":
		school.SchoolType = updateSchoolForm.UpdateValue.(int)
	case "province":
		school.Province = updateValue
	case "officialWebLink":
		school.OfficialWebLink = updateValue
	case "schoolRemark":
		school.SchoolRemark = updateValue
	}
	schoolByte, _ := json.Marshal(school)
	err = redisClient.LSet(ctx, schoolKey, updateSchoolForm.SchoolListIndex, schoolByte).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		return
	}

	go func(updateSchoolForm UpdateSchoolForm) {
		mysqlClient := mysql.GetClient()
		updateValue := updateSchoolForm.UpdateValue.(string)
		err := mysqlClient.Model(&mysql.SchoolSQL{}).Where("schoolId = ?", updateSchoolForm.SchoolId).UpdateColumn(updateSchoolForm.UpdateField, updateValue).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		}
	}(updateSchoolForm)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}
