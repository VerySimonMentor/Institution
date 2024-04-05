package CRUD

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetCountryResp struct {
	CountryId      int    `json:"countryId"`
	CountryChiName string `json:"countryChiName"`
}

func GetCountryHandler(ctx *gin.Context) {
	countryList := getCountryInRedis(ctx)
	if countryList == nil {
		return
	}

	allCountries := make([]GetCountryResp, len(countryList))
	for i, country := range countryList {
		allCountries[i] = GetCountryResp{
			CountryId:      country.CountryId,
			CountryChiName: country.CountryChiName,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"result": allCountries})
}
