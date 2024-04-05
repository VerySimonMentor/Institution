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
	SchoolType         string `gorm:"column:schoolType"`
	Province           string `gorm:"column:province"`
	OfficialWebLink    string `gorm:"column:officialWebLink"`
	SchoolRemark       string `gorm:"column:schoolRemark"`
	SchoolAndItem      []byte `gorm:"column:schoolAndItem"`
}

func (s SchoolSQL) TableName() string {
	return "school"
}

type ItemSQL struct {
	ItemId          int    `gorm:"column:itemId"`
	ItemName        string `gorm:"column:itemName"`
	LevelDescrption string `gorm:"column:levelDescrption"`
	LevelRate       string `gorm:"column:levelRate"`
	ItemRemark      string `gorm:"column:itemRemark"`
}

func (i ItemSQL) TableName() string {
	return "item"
}

type UserSQL struct {
	UserId       int    `gorm:"column:userId"`
	UserAccount  string `gorm:"column:userAccount"`
	UserPassWord string `gorm:"column:userPassWord"`
	UserEmail    string `gorm:"column:userEmail"`
	UserNumber   string `gorm:"column:userNumber"`
	UserLevel    int    `gorm:"column:userLevel"`
	StudentCount int    `gorm:"column:studentCount"`
}

func (u UserSQL) TableName() string {
	return "user"
}
