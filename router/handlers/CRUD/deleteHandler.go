package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCountryHandler(ctx *gin.Context) {
	var deleteForm CountryInstanceForm
	if err := ctx.ShouldBindJSON(&deleteForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	deleteCountryString, err := redisClient.LIndex(context.Background(), "country", deleteForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}
	var deleteCountry Country
	if err := json.Unmarshal([]byte(deleteCountryString), &deleteCountry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}
	if deleteCountry.CountryId != deleteForm.CountryId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}
	_, err = redisClient.LRem(context.Background(), "country", 0, deleteCountryString).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis删除失败"})
		logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		return
	}

	go func(countryId int) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Where("countryId = ?", countryId).Delete(&mysql.CountrySQL{}).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleteCountryHandler error %s", err)
		}
	}(deleteForm.CountryId)

	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func DeleteSchoolHandler(ctx *gin.Context) {
	var deleteForm SchoolInstanceForm
	if err := ctx.ShouldBindJSON(&deleteForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryString, err := redisClient.LIndex(context.Background(), "country", deleteForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}
	country.CountryAndSchool = append(country.CountryAndSchool[:deleteForm.SchoolListIndex], country.CountryAndSchool[deleteForm.SchoolListIndex+1:]...)
	countryStringNew, _ := json.Marshal(country)
	_, err = redisClient.LSet(context.Background(), "country", deleteForm.CountryListIndex, countryStringNew).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis修改失败"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}

	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	deleteSchoolString, err := redisClient.LIndex(context.Background(), schoolKey, deleteForm.SchoolListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}
	var deleteSchool School
	if err := json.Unmarshal([]byte(deleteSchoolString), &deleteSchool); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}
	if deleteSchool.SchoolId != deleteForm.SchoolId {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}
	_, err = redisClient.LRem(context.Background(), schoolKey, 0, deleteSchoolString).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis删除失败"})
		logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		return
	}

	go func(schoolId int, country Country) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Where("schoolId = ?", schoolId).Delete(&mysql.SchoolSQL{}).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		}

		countryAndSchoolString, _ := json.Marshal(country.CountryAndSchool)
		err = mysqlClient.Model(&mysql.CountrySQL{}).Where("countryId = ?", country.CountryId).UpdateColumn("countryAndSchool", countryAndSchoolString).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleteSchoolHandler error %s", err)
		}
	}(deleteForm.SchoolId, country)

	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func DeleteItemHandler(ctx *gin.Context) {
	var deleteItemForm ItemInstanceForm
	if err := ctx.ShouldBindJSON(&deleteItemForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	countryString, err := redisClient.LIndex(context.Background(), "country", deleteItemForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}
	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	schoolString, err := redisClient.LIndex(context.Background(), schoolKey, deleteItemForm.SchoolListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}
	var school School
	if err := json.Unmarshal([]byte(schoolString), &school); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}
	school.SchoolAndItem = append(school.SchoolAndItem[:deleteItemForm.ItemListIndex], school.SchoolAndItem[deleteItemForm.ItemListIndex+1:]...)
	schoolStringNew, _ := json.Marshal(school)
	_, err = redisClient.LSet(context.Background(), schoolKey, deleteItemForm.SchoolListIndex, schoolStringNew).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis修改失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}

	itemKey := fmt.Sprintf(ItemKey, country.CountryId, school.SchoolId)
	deleteItemString, err := redisClient.LIndex(context.Background(), itemKey, deleteItemForm.ItemListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}
	var deleteItem Item
	if err := json.Unmarshal([]byte(deleteItemString), &deleteItem); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}

	err = redisClient.LRem(context.Background(), itemKey, 0, deleteItemString).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis删除失败"})
		logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		return
	}

	go func(itemId int, school School) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Where("itemId = ?", itemId).Delete(&mysql.ItemSQL{}).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		}

		schoolAndItemString, _ := json.Marshal(school.SchoolAndItem)
		err = mysqlClient.Model(&mysql.SchoolSQL{}).Where("schoolId = ?", school.SchoolId).UpdateColumn("schoolAndItem", schoolAndItemString).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleleItemHandler error %s", err)
		}
	}(deleteItem.ItemId, school)

	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

type DeleteForm struct {
	ListIndex int64 `json:"listIndex"`
}

func DeleteUserHandler(ctx *gin.Context) {
	var deleteUserForm DeleteForm
	if err := ctx.ShouldBindJSON(&deleteUserForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleleUserHandler error %s", err)
		return
	}

	redisClient := redis.GetClient()
	userString, err := redisClient.LIndex(context.Background(), "user", deleteUserForm.ListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("DeleteUserHandler error %s", err)
		return
	}
	var user User
	if err := json.Unmarshal([]byte(userString), &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("DeleteUserHandler error %s", err)
		return
	}

	_, err = redisClient.LRem(context.Background(), "user", 0, userString).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis删除失败"})
		logs.GetInstance().Logger.Errorf("DeleteUserHandler error %s", err)
		return
	}

	go func(userId int) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Where("userId = ?", userId).Delete(&mysql.UserSQL{}).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("DeleteUserHandler error %s", err)
		}
	}(user.UserId)

	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func DeleteSystemHandler(ctx *gin.Context) {
	var deleteSystemForm DeleteForm
	if err := ctx.ShouldBindJSON(&deleteSystemForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("DeleleSystemHandler error %s", err)
		return
	}

	usedSchool := make([]string, 0)
	listIndex := deleteSystemForm.ListIndex
	system := getSystemInRedis(ctx)
	currentSchoolTypeId := system.SchoolTyepList[listIndex].SchoolTypeId
	country := getCountryInRedis(ctx)
	for _, c := range country {
		schoolKey := fmt.Sprintf(SchoolKey, c.CountryId)
		school := getSchoolInRedis(ctx, schoolKey, c.CountryAndSchool)
		for _, s := range school {
			if s.SchoolType == currentSchoolTypeId {
				usedSchool = append(usedSchool, fmt.Sprintf(DeleteTypeResp, c.CountryChiName, s.SchoolChiName))
			}
		}
	}

	if len(usedSchool) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "该类型正在使用中",
			"results": usedSchool,
		})
		return
	}
	system.SchoolTyepList = append(system.SchoolTyepList[:listIndex], system.SchoolTyepList[listIndex+1:]...)
	redisClient := redis.GetClient()
	systemString, _ := json.Marshal(system)
	_, err := redisClient.Set(context.Background(), "system", systemString, 0).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis修改失败"})
		logs.GetInstance().Logger.Errorf("DeleteSystemHandler error %s", err)
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
			logs.GetInstance().Logger.Errorf("DeleteSystemHandler error %s", err)
		}
	}(system)

	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
