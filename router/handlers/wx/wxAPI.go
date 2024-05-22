package wx

import (
	"Institution/config"
	"Institution/logs"
	"Institution/redis"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessToken(wxConfig *config.WxConfig) string {
	redisClient := redis.GetClient()
	accessToken, err := redisClient.Get(context.Background(), accessTokenKey).Result()
	if redis.CheckNil(err) {
		resp, e := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=" + wxConfig.Grant_type + "&appid=" + wxConfig.Appid + "&secret=" + wxConfig.Secret)
		if e != nil {
			logs.GetInstance().Logger.Errorf("get access token error %s", e)
			return ""
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var token AccessToken
		json.Unmarshal(body, &token)
		logs.GetInstance().Logger.Infof("token %+v", token)

		accessToken = token.AccessToken
		redisClient.Set(context.Background(), accessTokenKey, accessToken, 7200*time.Second)
	} else if err != nil {
		logs.GetInstance().Logger.Errorf("get access token error %s", err)
		return ""
	}

	return accessToken
}

type PhoneNumberResp struct {
	ErrCode   int       `json:"errcode"`
	PhoneInfo PhoneInfo `json:"phone_info"`
}

type PhoneInfo struct {
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     int    `json:"countryCode"`
	PhoneNumber     string `json:"phoneNumber"`
}

func GetPhoneNumber(code string, wxConfig *config.WxConfig) string {
	accessToken := GetAccessToken(wxConfig)

	requestJSON, _ := json.Marshal(map[string]string{
		"code": code,
	})
	resp, err := http.Post("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="+accessToken, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		logs.GetInstance().Logger.Errorf("get phone number error %s", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var phoneResp PhoneNumberResp
	json.Unmarshal(body, &phoneResp)
	if phoneResp.ErrCode != 0 {
		redisClient := redis.GetClient()
		redisClient.Del(context.Background(), accessTokenKey)
		accessToken = GetAccessToken(wxConfig)
		if accessToken == "" {
			logs.GetInstance().Logger.Errorf("get access token error")
			return ""
		}
		resp, err = http.Post("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="+accessToken, "application/json", bytes.NewBuffer(requestJSON))
		if err != nil {
			logs.GetInstance().Logger.Errorf("get phone number error %s", err)
			return ""
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		json.Unmarshal(body, &phoneResp)
		logs.GetInstance().Logger.Infof("phone number %+v", phoneResp)
	}

	if phoneResp.PhoneInfo.CountryCode == 86 {
		return phoneResp.PhoneInfo.PurePhoneNumber
	}
	return phoneResp.PhoneInfo.PhoneNumber
}

func CheckTokenHandler(ctx *gin.Context) {
	tocken := ctx.Query("token")
	redisClient := redis.GetClient()
	_, err := redisClient.Get(context.Background(), tocken).Result()
	if redis.CheckNil(err) {
		ctx.JSON(http.StatusOK, gin.H{
			"state": false,
		})
	} else if err != nil {
		logs.GetInstance().Logger.Infof("redis get state tocken err %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"state": true,
	})
}

type SessionForm struct {
	SessionKey string `json:"session_key"`
	Errcode    int    `json:"errcode"`
}

func code2Session(wxConfig *config.WxConfig, code string) string {
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + wxConfig.Appid + "&secret=" + wxConfig.Secret + "&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		logs.GetInstance().Logger.Errorf("get session key error %s", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var session SessionForm
	json.Unmarshal(body, &session)
	if session.Errcode != 0 {
		logs.GetInstance().Logger.Errorf("get session key error %v", session)
		return ""
	}

	return session.SessionKey
}

func CheckLoginTocken(wxConfig *config.WxConfig, loginTocken string) (bool, string) {
	redisClient := redis.GetClient()
	exists, err := redisClient.Exists(context.Background(), loginTocken).Result()
	if err != nil {
		logs.GetInstance().Logger.Errorf("check session key error %s", err)
		return false, ""
	}

	if exists == 0 {
		return false, ""
	}
	phoneNumber, err := redisClient.Get(context.Background(), loginTocken).Result()
	if err != nil {
		logs.GetInstance().Logger.Errorf("get phone number error %s", err)
		return false, ""
	}

	return true, phoneNumber
}
