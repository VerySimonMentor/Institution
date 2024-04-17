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

type CreateItemForm struct {
	CountryListIndex int64 `json:"countryListIndex"`
	SchoolListIndex  int64 `json:"schoolListIndex"`
	PageNum          int   `json:"pageNum"`
}

func CreateItemHandler(ctx *gin.Context) {
	var createItemForm CreateItemForm
	if err := ctx.ShouldBindJSON(&createItemForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}
	if createItemForm.CountryListIndex < 0 || createItemForm.SchoolListIndex < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		return
	}
	mysqlClient := mysql.GetClient()
	redisClient := redis.GetClient()
	itemSQL := mysql.ItemSQL{
		ItemName:         "default",
		LevelDescription: "默认",
		LevelRate:        []byte("{}"),
		ItemRemark:       "",
	}

	if err := mysqlClient.Create(&itemSQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}
	if err := mysqlClient.Last(&itemSQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}

	itemByte, _ := json.Marshal(Item{
		ItemId:           itemSQL.ItemId,
		ItemName:         itemSQL.ItemName,
		LevelDescription: itemSQL.LevelDescription,
		LevelRate:        make([]Level, 0),
		ItemRemark:       itemSQL.ItemRemark,
	})
	countryString, err := redisClient.LIndex(context.Background(), "country", createItemForm.CountryListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}
	var country Country
	if err := json.Unmarshal([]byte(countryString), &country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}
	schoolKey := fmt.Sprintf(SchoolKey, country.CountryId)
	schoolString, err := redisClient.LIndex(context.Background(), schoolKey, createItemForm.SchoolListIndex).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}
	var school School
	if err := json.Unmarshal([]byte(schoolString), &school); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		return
	}
	school.SchoolAndItem = append(school.SchoolAndItem, itemSQL.ItemId)
	schoolByte, _ := json.Marshal(school)
	redisClient.LSet(context.Background(), schoolKey, createItemForm.SchoolListIndex, schoolByte)

	itemKey := fmt.Sprintf(ItemKey, country.CountryId, school.SchoolId)
	redisClient.RPush(context.Background(), itemKey, itemByte)
	itemNum := int(redisClient.LLen(context.Background(), itemKey).Val())
	var totalPage int
	if itemNum%createItemForm.PageNum == 0 {
		totalPage = itemNum / createItemForm.PageNum
	} else {
		totalPage = itemNum/createItemForm.PageNum + 1
	}

	go func(school School) {
		schoolAndItem, _ := json.Marshal(school.SchoolAndItem)
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.SchoolSQL{}).Where("schoolId = ?", school.SchoolId).Update("schoolAndItem", schoolAndItem).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("CreateItemHandler error %s", err)
		}
	}(school)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"itemId":    itemSQL.ItemId,
		"totalPage": totalPage,
	})
}

func CreateUserHandler(ctx *gin.Context) {
	pageNumStr := ctx.Query("pageNum")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		logs.GetInstance().Logger.Errorf("CreateUserHandler error %s", err)
		return
	}
	mysqlClient := mysql.GetClient()
	userSQL := mysql.UserSQL{
		UserAccount:  "default",
		UserPassWd:   "",
		UserEmail:    "",
		UserNumber:   "00000000000",
		UserLevel:    0,
		StudentCount: 0,
	}

	if err := mysqlClient.Create(&userSQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateUserHandler error %s", err)
		return
	}
	if err := mysqlClient.Last(&userSQL).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "创建失败"})
		logs.GetInstance().Logger.Errorf("CreateUserHandler error %s", err)
		return
	}

	redisClinet := redis.GetClient()
	user := User{
		UserId:       userSQL.UserId,
		UserAccount:  userSQL.UserAccount,
		UserPassWord: userSQL.UserPassWd,
		UserEmail:    userSQL.UserEmail,
		UserNumber:   userSQL.UserNumber,
		UserLevel:    userSQL.UserLevel,
		StudentCount: userSQL.StudentCount,
	}
	userByte, _ := json.Marshal(user)
	redisClinet.RPush(context.Background(), "user", userByte)
	userNum := int(redisClinet.LLen(context.Background(), "user").Val())
	var totalPage int
	if userNum%pageNum == 0 {
		totalPage = userNum / pageNum
	} else {
		totalPage = userNum/pageNum + 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"userId":    userSQL.UserId,
		"totalPage": totalPage,
	})
}

func CreateSystemHandler(ctx *gin.Context) {
	system := getSystemInRedis(ctx)
	maxSchoolTypeId := system.SchoolTyepList[len(system.SchoolTyepList)-1].SchoolTypeId
	system.SchoolTyepList = append(system.SchoolTyepList, SchoolType{
		SchoolTypeId:   maxSchoolTypeId + 1,
		SchoolTypeName: "默认",
	})
	redisClient := redis.GetClient()
	systemByte, _ := json.Marshal(system)
	redisClient.Set(context.Background(), "system", systemByte, 0)

	go func(system System) {
		mysqlClient := mysql.GetClient()
		err := mysqlClient.Model(&mysql.SystemSQL{}).Updates(system).Error
		if err != nil {
			logs.GetInstance().Logger.Errorf("CreateSystemHandler error %s", err)
		}
	}(system)

	ctx.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
}
