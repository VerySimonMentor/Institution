package itemprocess

import (
	"Institution/config"
	"Institution/logs"
	"Institution/router/handlers/CRUD"
	"Institution/router/handlers/wx"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCountryHandler(ctx *gin.Context, wxConfig *config.WxConfig) {
	loginTocken := ctx.Query("loginTocken")
	check, _ := wx.CheckLoginTocken(wxConfig, loginTocken)
	if !check {
		logs.GetInstance().Logger.Warn("check login tocken error")
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	countryList := CRUD.GetCountryInRedis(ctx)
	countryResp := make([]string, len(countryList))
	for i, country := range countryList {
		countryResp[i] = country.CountryChiName
	}

	ctx.JSON(http.StatusOK, gin.H{
		"countryList": countryResp,
	})
}
