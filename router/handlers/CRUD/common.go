package CRUD

type Country struct {
	CountryId        int              `json:"countryId"`
	CountryEngName   string           `json:"countryEngName"`
	CountryChiName   string           `json:"countryChiName"`
	CountryAndSchool map[int]struct{} `json:"countryAndSchool"`
	Province         []Province       `json:"province"`
}

type Province struct {
	EngName string `json:"engName"`
	ChiName string `json:"chiName"`
}
