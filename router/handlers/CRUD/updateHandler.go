package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		err := mysqlClient.Model(&mysql.CountrySQL{}).Where("countryId = ?", updateCountryForm.CountryId).UpdateColumn(updateCountryForm.UpdateField, updateCountryForm.UpdateValue).Error
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
		err := mysqlClient.Model(&mysql.SchoolSQL{}).Where("schoolId = ?", updateSchoolForm.SchoolId).UpdateColumn(updateSchoolForm.UpdateField, updateSchoolForm.UpdateValue).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateSchoolHandler error %s", err)
		}
	}(updateSchoolForm)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type UpdateItemForm struct {
	CountryListIndex int64       `json:"countryListIndex"`
	SchoolListIndex  int64       `json:"schoolListIndex"`
	ItemListIndex    int64       `json:"itemListIndex"`
	ItemId           int         `json:"itemId"`
	UpdateField      string      `json:"updateField"`
	UpdateValue      interface{} `json:"updateValue"`
}

func UpdateItemHandler(ctx *gin.Context) {
	var updateItemForm UpdateItemForm
	if err := ctx.ShouldBindJSON(&updateItemForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryString, err := redisClient.LIndex(ctx, "country", updateItemForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}

	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	schoolString, err := redisClient.LIndex(ctx, schoolKey, updateItemForm.SchoolListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}
	var school School
	if err := json.Unmarshal([]byte(schoolString), &school); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}

	itemKey := fmt.Sprintf(ItemKey, country.CountryId, school.SchoolId)
	itemString, err := redisClient.LIndex(ctx, itemKey, updateItemForm.ItemListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}
	var item Item
	if err := json.Unmarshal([]byte(itemString), &item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}

	updateValue := updateItemForm.UpdateValue.(string)
	switch updateItemForm.UpdateField {
	case "itemName":
		item.ItemName = updateValue
	case "levelDescrption":
		item.LevelDescription = updateValue
	case "itemRemark":
		item.ItemRemark = updateValue
	}
	itemByte, _ := json.Marshal(item)
	err = redisClient.LSet(ctx, itemKey, updateItemForm.ItemListIndex, itemByte).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		return
	}

	go func(updateItemForm UpdateItemForm) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.ItemSQL{}).Where("itemId = ?", updateItemForm.ItemId).UpdateColumn(updateItemForm.UpdateField, updateItemForm.UpdateValue).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateItemHandler error %s", err)
		}
	}(updateItemForm)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type UpdateLevelForm struct {
	Item
	CountryListIndex int64 `json:"countryListIndex"`
	SchoolListIndex  int64 `json:"schoolListIndex"`
	ItemListIndex    int64 `json:"itemListIndex"`
}

func UpdateLevelHandler(ctx *gin.Context) {
	var updateLevelForm UpdateLevelForm
	if err := ctx.ShouldBindJSON(&updateLevelForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryString, err := redisClient.LIndex(ctx, "country", updateLevelForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}

	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	schoolString, err := redisClient.LIndex(ctx, schoolKey, updateLevelForm.SchoolListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}
	var school School
	if err := json.Unmarshal([]byte(schoolString), &school); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}

	itemKey := fmt.Sprintf(ItemKey, country.CountryId, school.SchoolId)
	itemString, err := redisClient.LIndex(ctx, itemKey, updateLevelForm.ItemListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}
	var item Item
	if err := json.Unmarshal([]byte(itemString), &item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}

	item.ItemName = updateLevelForm.ItemName
	item.LevelDescription = updateLevelForm.LevelDescription
	item.LevelRate = updateLevelForm.LevelRate
	item.ItemRemark = updateLevelForm.ItemRemark
	itemByte, _ := json.Marshal(item)
	err = redisClient.LSet(ctx, itemKey, updateLevelForm.ItemListIndex, itemByte).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		return
	}

	go func(item Item) {
		mysqlClient := mysql.GetClient()
		levelRate, _ := json.Marshal(item.LevelRate)
		itemSQL := mysql.ItemSQL{
			ItemId:           item.ItemId,
			ItemName:         item.ItemName,
			LevelDescription: item.LevelDescription,
			LevelRate:        levelRate,
			ItemRemark:       item.ItemRemark,
		}
		err := mysqlClient.Model(&mysql.ItemSQL{}).Where("itemId = ?", item.ItemId).Updates(itemSQL).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateLevelHandler error %s", err)
		}
	}(item)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type UpdateUserForm struct {
	ListIndex   int64       `json:"listIndex"`
	UserId      int         `json:"userId"`
	UpdateField string      `json:"updateField"`
	UpdateValue interface{} `json:"updateValue"`
}

func UpdateUserHandler(ctx *gin.Context) {
	var updateUserForm UpdateUserForm
	if err := ctx.ShouldBindJSON(&updateUserForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateUserHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	userString, err := redisClient.LIndex(ctx, "user", updateUserForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("UpdateUserHandler error %s", err)
		return
	}
	var user User
	if err := json.Unmarshal([]byte(userString), &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("UpdateUserHandler error %s", err)
		return
	}

	updateValue := updateUserForm.UpdateValue.(string)
	switch updateUserForm.UpdateField {
	case "userAccount":
		user.UserAccount = updateValue
	case "userPassWord":
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateValue), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "密码加密失败"})
			logs.GetInstance().Logger.Errorf("UpdateUserHandler error %s", err)
			return
		}
		user.UserPassWord = string(hashedPassword)
	case "userEmail":
		user.UserEmail = updateValue
	case "userNumber":
		user.UserNumber = updateValue
	case "userLevel":
		user.UserLevel = updateUserForm.UpdateValue.(int)
	case "studentCount":
		user.StudentCount = updateUserForm.UpdateValue.(int)
	}
	userByte, _ := json.Marshal(user)
	err = redisClient.LSet(ctx, "user", updateUserForm.ListIndex, userByte).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateUserHandler error %s", err)
		return
	}

	go func(updateUserForm UpdateUserForm) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.UserSQL{}).Where("userId = ?", updateUserForm.UserId).UpdateColumn(updateUserForm.UpdateField, updateUserForm.UpdateValue).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateUserHandler error %s", err)
		}
	}(updateUserForm)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

type UpdateSystemForm struct {
	ListIndex   int64       `json:"listIndex"`
	UpdateField string      `json:"updateField"`
	UpdateValue interface{} `json:"updateValue"`
}

func UpdateSystemHandler(ctx *gin.Context) {
	var updateSystemForm UpdateSystemForm
	if err := ctx.ShouldBindJSON(&updateSystemForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("UpdateSystemHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	system := getSystemInRedis(ctx)
	switch updateSystemForm.UpdateField {
	case "maxUserLevel":
		system.MaxUserLevel = updateSystemForm.UpdateValue.(int)
	case "schoolTyepList":
		system.SchoolTyepList[updateSystemForm.ListIndex].SchoolTypeName = updateSystemForm.UpdateValue.(string)
	}
	systemByte, _ := json.Marshal(system)
	err := redisClient.Set(ctx, "system", systemByte, 0).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis更新失败"})
		logs.GetInstance().Logger.Errorf("UpdateSystemHandler error %s", err)
		return
	}

	go func(system System) {
		mysqlClient := mysql.GetClient()
		schoolTypeList, _ := json.Marshal(system.SchoolTyepList)
		systemSQL := mysql.SystemSQL{
			MaxUserLevel:   system.MaxUserLevel,
			SchoolTyepList: schoolTypeList,
		}
		err := mysqlClient.Model(&mysql.SystemSQL{}).Updates(systemSQL).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("UpdateSystemHandler error %s", err)
		}
	}(system)

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}
