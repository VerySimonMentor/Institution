package CRUD

const (
	SchoolKey = "%d-school"
	ItemKey   = "%d-%d-item"

	DeleteTypeResp = "国家:%s的学校:%s正在使用该类型,请先更改再删除!"
)

type Country struct {
	CountryId        int        `json:"countryId"`
	CountryEngName   string     `json:"countryEngName"`
	CountryChiName   string     `json:"countryChiName"`
	CountryAndSchool []int      `json:"countryAndSchool"`
	Province         []Province `json:"province"`
}

type Province struct {
	EngName string `json:"engName"`
	ChiName string `json:"chiName"`
}

type School struct {
	SchoolId           int    `json:"schoolId"`
	SchoolEngName      string `json:"schoolEngName"`
	SchoolChiName      string `json:"schoolChiName"`
	SchoolAbbreviation string `json:"schoolAbbreviation"`
	SchoolType         int    `json:"schoolType"`
	Province           string `json:"province"`
	OfficialWebLink    string `json:"officialWebLink"`
	SchoolRemark       string `json:"schoolRemark"`
	SchoolAndItem      []int  `json:"schoolAndItem"`
}

type Item struct {
	ItemId          int     `json:"itemId"`
	ItemName        string  `json:"itemName"`
	LevelDescrption string  `json:"levelDescrption"`
	LevelRate       []Level `json:"levelRate"`
	ItemRemark      string  `json:"itemRemark"`
}

type Level struct {
	LevelId   int `json:"levelId"`
	LevelRate int `json:"levelRate"`
}

type User struct {
	UserId       int    `json:"userId"`
	UserAccount  string `json:"userAccount"`
	UserPassWord string `json:"userPassWord"`
	UserEmail    string `json:"userEmail"`
	UserNumber   string `json:"userNumber"`
	UserLevel    int    `json:"userLevel"`
	StudentCount int    `json:"studentCount"`
}

type System struct {
	MaxUserLevel   int          `json:"maxUserLevel"`
	SchoolTyepList []SchoolType `json:"schoolTyepList"`
}

type SchoolType struct {
	SchoolTypeId   int    `json:"schoolTypeId"`
	SchoolTypeName string `json:"schoolTypeName"`
}
