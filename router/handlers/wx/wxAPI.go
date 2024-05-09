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
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessToken(wxConfig *config.WxConfig) string {
	redisClient := redis.GetClient()
	accessToken, err := redisClient.Get(context.Background(), tokenKey).Result()
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
		redisClient.Set(context.Background(), tokenKey, accessToken, 7200*time.Second)
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
	logs.GetInstance().Logger.Infof("phone number %+v", phoneResp)
	if phoneResp.ErrCode != 0 {
		redisClient := redis.GetClient()
		redisClient.Del(context.Background(), tokenKey)
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
