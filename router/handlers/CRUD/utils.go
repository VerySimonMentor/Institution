package CRUD

import (
	"Institution/logs"
	"Institution/mysql"
	"Institution/redis"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CountryInstanceForm struct {
	CountryId int   `json:"countryId"`
	ListIndex int64 `json:"listIndex"`
}

type SchoolInstanceForm struct {
	CountryListIndex int64 `json:"countryListIndex"`
	SchoolId         int   `json:"schoolId"`
	SchoolListIndex  int64 `json:"schoolListIndex"`
}

type ItemInstanceForm struct {
	CountryListIndex int64 `json:"countryListIndex"`
	SchoolListIndex  int64 `json:"schoolListIndex"`
	ItemListIndex    int64 `json:"itemListIndex"`
}

func pageRange(page, pageNum, countryNum int) (int, int) {
	start := (page - 1) * pageNum
	end := start + pageNum
	if start > countryNum {
		start = countryNum
	}
	if end > countryNum {
		end = countryNum
	}
	return start, end
}

func getCountryInRedis(ctx *gin.Context) []Country {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()

	if !checkCountryInRedis(ctx) {
		countryListSQL := make([]mysql.CountrySQL, 0)
		if err := mysqlClient.Model(&mysql.CountrySQL{}).Find(&countryListSQL).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
			logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
			return nil
		}
		if len(countryListSQL) == 0 {
			return make([]Country, 0)
		}
		countryList := make([]Country, len(countryListSQL))
		for i, country := range countryListSQL {
			countryAndSchool := make([]int, 0)
			province := make([]Province, 0)
			json.Unmarshal(country.CountryAndSchool, &countryAndSchool)
			json.Unmarshal(country.Province, &province)
			countryList[i] = Country{
				CountryId:        country.CountryId,
				CountryEngName:   country.CountryEngName,
				CountryChiName:   country.CountryChiName,
				CountryAndSchool: countryAndSchool,
				Province:         province,
			}
		}
		countryJSON := make([][]byte, len(countryList))
		for i, country := range countryList {
			var err error
			countryJSON[i] = make([]byte, 0)
			countryJSON[i], err = json.Marshal(country)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
			if err = redisClient.RPush(context.Background(), "country", countryJSON[i]).Err(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
		}
	}

	countryString, err := redisClient.LRange(context.Background(), "country", 0, -1).Result()
	countryList := make([]Country, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return nil
	} else {
		for _, country := range countryString {
			var countryStruct Country
			if err := json.Unmarshal([]byte(country), &countryStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
				return nil
			}
			countryList = append(countryList, countryStruct)
		}
	}

	return countryList
}

func checkCountryInRedis(ctx *gin.Context) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), "country").Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return false
	}
	if keytype != "list" {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis类型错误"})
		// logs.GetInstance().Logger.Errorf("ShowCountryHandler error %s", err)
		return false
	}

	return true
}

func getSchoolInRedis(ctx *gin.Context, schoolKey string, countryAndSchool []int) []School {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()

	if !checkSchoolInRedis(ctx, schoolKey) {
		schoolListSQL := make([]mysql.SchoolSQL, len(countryAndSchool))
		for i, schoolId := range countryAndSchool {
			school := mysql.SchoolSQL{}
			if err := mysqlClient.Model(&mysql.SchoolSQL{}).Where("schoolId = ?", schoolId).Find(&school).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
				logs.GetInstance().Logger.Errorf("get school in redis error %s", err)
				return nil
			}
			schoolListSQL[i] = school
		}
		if len(schoolListSQL) == 0 {
			return make([]School, 0)
		}

		schoolList := make([]School, len(schoolListSQL))
		for i, school := range schoolListSQL {
			schoolAndItem := make([]int, 0)
			json.Unmarshal(school.SchoolAndItem, &schoolAndItem)
			schoolList[i] = School{
				SchoolId:           school.SchoolId,
				SchoolEngName:      school.SchoolEngName,
				SchoolChiName:      school.SchoolChiName,
				SchoolAbbreviation: school.SchoolAbbreviation,
				SchoolType:         school.SchoolType,
				Province:           school.Province,
				OfficialWebLink:    school.OfficialWebLink,
				SchoolRemark:       school.SchoolRemark,
				SchoolAndItem:      schoolAndItem,
			}
		}
		schoolJSON := make([][]byte, len(schoolList))
		for i, school := range schoolList {
			var err error
			schoolJSON[i] = make([]byte, 0)
			schoolJSON[i], err = json.Marshal(school)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
				return nil
			}
			if err = redisClient.RPush(context.Background(), schoolKey, schoolJSON[i]).Err(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
				logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
				return nil
			}
		}
	}

	schoolString, err := redisClient.LRange(context.Background(), schoolKey, 0, -1).Result()
	schoolList := make([]School, len(countryAndSchool))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return nil
	} else {
		for i, school := range schoolString {
			var schoolStruct School
			if err := json.Unmarshal([]byte(school), &schoolStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
				return nil
			}
			schoolList[i] = schoolStruct
		}
	}

	return schoolList
}

func checkSchoolInRedis(ctx *gin.Context, schoolKey string) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), schoolKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return false
	}
	if keytype != "list" {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis类型错误"})
		// logs.GetInstance().Logger.Errorf("ShowSchoolHandler error %s", err)
		return false
	}

	return true
}

func getItemInRedis(ctx *gin.Context, itemKey string, schoolAndItem []int) []Item {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()

	if !checkItemInRedis(ctx, itemKey) {
		itemListSQL := make([]mysql.ItemSQL, len(schoolAndItem))
		for i, itemId := range schoolAndItem {
			item := mysql.ItemSQL{}
			if err := mysqlClient.Model(&mysql.ItemSQL{}).Where("itemId = ?", itemId).Find(&item).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
				logs.GetInstance().Logger.Errorf("get item in redis error %s", err)
				return nil
			}
			itemListSQL[i] = item
		}
		if len(itemListSQL) == 0 {
			return make([]Item, 0)
		}

		itemList := make([]Item, len(itemListSQL))
		for i, item := range itemListSQL {
			var levelRate []Level
			json.Unmarshal(item.LevelRate, &levelRate)
			itemList[i] = Item{
				ItemId:           item.ItemId,
				ItemName:         item.ItemName,
				LevelDescription: item.LevelDescription,
				LevelRate:        levelRate,
				ItemRemark:       item.ItemRemark,
			}
		}
		itemJSON := make([][]byte, len(itemList))
		for i, item := range itemList {
			var err error
			itemJSON[i] = make([]byte, 0)
			itemJSON[i], err = json.Marshal(item)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowItemHandler error %s", err)
				return nil
			}
			if err = redisClient.RPush(context.Background(), itemKey, itemJSON[i]).Err(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
				logs.GetInstance().Logger.Errorf("ShowItemHandler error %s", err)
				return nil
			}
		}
	}

	itemString, err := redisClient.LRange(context.Background(), itemKey, 0, -1).Result()
	itemList := make([]Item, len(schoolAndItem))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowItemHandler error %s", err)
		return nil
	} else {
		for i, item := range itemString {
			var itemStruct Item
			if err := json.Unmarshal([]byte(item), &itemStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowItemHandler error %s", err)
				return nil
			}
			itemList[i] = itemStruct
		}
	}

	return itemList
}

func checkItemInRedis(ctx *gin.Context, itemKey string) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), itemKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowItemHandler error %s", err)
		return false
	}
	if keytype != "list" {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis类型错误"})
		// logs.GetInstance().Logger.Errorf("ShowItemHandler error %s", err)
		return false
	}

	return true
}

func getSystemInRedis(ctx *gin.Context) System {
	redisClient := redis.GetClient()

	if !checkSystemInRedis(ctx) {
		mysqlClient := mysql.GetClient()
		systemSQL := mysql.SystemSQL{}
		var count int64

		mysqlClient.Model(&mysql.SystemSQL{}).Count(&count)
		if count == 0 {
			return System{
				MaxUserLevel:   0,
				SchoolTyepList: make([]SchoolType, 0),
			}
		}

		if err := mysqlClient.Model(&mysql.SystemSQL{}).First(&systemSQL).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
			logs.GetInstance().Logger.Errorf("ShowSystemHandler error %s", err)
			return System{}
		}
		var schoolTyepList []SchoolType
		json.Unmarshal(systemSQL.SchoolTyepList, &schoolTyepList)
		system := System{
			MaxUserLevel:   systemSQL.MaxUserLevel,
			SchoolTyepList: schoolTyepList,
		}
		systemJSON, _ := json.Marshal(system)
		if err := redisClient.Set(context.Background(), "system", systemJSON, 0).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
			logs.GetInstance().Logger.Errorf("ShowSystemHandler error %s", err)
			return System{}
		}
	}

	systemString, err := redisClient.Get(context.Background(), "system").Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSystemHandler error %s", err)
		return System{}
	}
	var system System
	if err := json.Unmarshal([]byte(systemString), &system); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
		logs.GetInstance().Logger.Errorf("ShowSystemHandler error %s", err)
		return System{}
	}

	return system
}

func checkSystemInRedis(ctx *gin.Context) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), "system").Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowSystemHandler error %s", err)
		return false
	}
	if keytype != "string" {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis类型错误"})
		// logs.GetInstance().Logger.Errorf("ShowSystemHandler error %s", err)
		return false
	}

	return true
}

func getUserInRedis(ctx *gin.Context) []User {
	redisClient := redis.GetClient()
	mysqlClient := mysql.GetClient()

	if !checkUserInRedis(ctx) {
		userListSQL := make([]mysql.UserSQL, 0)
		if err := mysqlClient.Model(&mysql.UserSQL{}).Find(&userListSQL).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "查询失败"})
			logs.GetInstance().Logger.Errorf("ShowUserHandler error %s", err)
			return nil
		}
		if len(userListSQL) == 0 {
			return make([]User, 0)
		}
		userList := make([]User, len(userListSQL))
		for i, user := range userListSQL {
			userList[i] = User{
				UserId:       user.UserId,
				UserAccount:  user.UserAccount,
				UserPassWord: user.UserPassWord,
				UserEmail:    user.UserEmail,
				UserNumber:   user.UserNumber,
				UserLevel:    user.UserLevel,
				StudentCount: user.StudentCount,
			}
		}
		userJSON := make([][]byte, len(userList))
		for i, user := range userList {
			var err error
			userJSON[i] = make([]byte, 0)
			userJSON[i], err = json.Marshal(user)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowUserHandler error %s", err)
				return nil
			}
			if err = redisClient.RPush(context.Background(), "user", userJSON[i]).Err(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis存储失败"})
				logs.GetInstance().Logger.Errorf("ShowUserHandler error %s", err)
				return nil
			}
		}
	}

	userString, err := redisClient.LRange(context.Background(), "user", 0, -1).Result()
	userList := make([]User, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowUserHandler error %s", err)
		return nil
	} else {
		for _, user := range userString {
			var userStruct User
			if err := json.Unmarshal([]byte(user), &userStruct); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": "json转换失败"})
				logs.GetInstance().Logger.Errorf("ShowUserHandler error %s", err)
				return nil
			}
			userList = append(userList, userStruct)
		}
	}

	return userList
}

func checkUserInRedis(ctx *gin.Context) bool {
	redisClient := redis.GetClient()
	keytype, err := redisClient.Type(context.Background(), "user").Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "redis查询失败"})
		logs.GetInstance().Logger.Errorf("ShowUserHandler error %s", err)
		return false
	}
	if keytype != "list" {
		return false
	}

	return true
}
