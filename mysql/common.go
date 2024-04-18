package mysql

type CountrySQL struct {
	CountryId        int    `gorm:"column:countryId" json:"countryId"`
	CountryEngName   string `gorm:"column:countryEngName" json:"countryEngName"`
	CountryChiName   string `gorm:"column:countryChiName" json:"countryChiName"`
	CountryAndSchool []byte `gorm:"column:countryAndSchool" json:"countryAndSchool"`
	Province         []byte `gorm:"column:province" json:"province"`
}

func (c CountrySQL) TableName() string {
	return "country"
}

type SchoolSQL struct {
	SchoolId           int    `gorm:"column:schoolId"`
	SchoolEngName      string `gorm:"column:schoolEngName"`
	SchoolChiName      string `gorm:"column:schoolChiName"`
	SchoolAbbreviation string `gorm:"column:schoolAbbreviation"`
	SchoolType         int    `gorm:"column:schoolType"`
	Province           string `gorm:"column:province"`
	OfficialWebLink    string `gorm:"column:officialWebLink"`
	SchoolRemark       string `gorm:"column:schoolRemark"`
	SchoolAndItem      []byte `gorm:"column:schoolAndItem"`
}

func (s SchoolSQL) TableName() string {
	return "school"
}

type ItemSQL struct {
	ItemId           int    `gorm:"column:itemId"`
	ItemName         string `gorm:"column:itemName"`
	LevelDescription string `gorm:"column:levelDescription"`
	LevelRate        []byte `gorm:"column:levelRate"`
	ItemRemark       string `gorm:"column:itemRemark"`
}

func (i ItemSQL) TableName() string {
	return "item"
}

type UserSQL struct {
	UserId       int    `gorm:"column:userId"`
	UserAccount  string `gorm:"column:userAccount"`
	UserPassWd   string `gorm:"column:userPassWd"`
	UserEmail    string `gorm:"column:userEmail"`
	UserNumber   string `gorm:"column:userNumber"`
	UserLevel    int    `gorm:"column:userLevel"`
	StudentCount int    `gorm:"column:studentCount"`
}

func (u UserSQL) TableName() string {
	return "user"
}

type SystemSQL struct {
	SystemId       int    `gorm:"column:systemId"`
	MaxUserLevel   int    `gorm:"column:maxUserLevel"`
	SchoolTypeList []byte `gorm:"column:schoolTypeList"`
}

func (s SystemSQL) TableName() string {
	return "systemSetting"
}
