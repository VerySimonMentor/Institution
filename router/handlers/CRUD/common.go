package router

type Country struct {
	CountryId        int                 `json:"countryId"`
	CountryEngName   string              `json:"countryEngName"`
	CountryChiName   string              `json:"countryChiName"`
	CountryAndSchool map[int]struct{}    `json:"countryAndSchool"`
	Province         map[string]struct{} `json:"province"`
}
