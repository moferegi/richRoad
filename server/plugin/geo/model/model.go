package model

type Geo struct {
	ID         uint   `gorm:"primary_key" json:"id" form:"id"`
	Name       string `gorm:"column:name" json:"name" form:"name"`
	Level      int    `gorm:"column:level" json:"level" form:"level"`
	Code       string `gorm:"column:code" json:"code" form:"code"`
	Geocode    string `gorm:"column:geocode" json:"geocode" form:"geocode"`
	Latitude   string `gorm:"column:latitude" json:"latitude" form:"latitude"`
	Longitude  string `gorm:"column:longitude" json:"longitude" form:"longitude"`
	Sort       int    `gorm:"column:sort" json:"sort" form:"id"`
	Children   []Geo  `gorm:"-" json:"children" form:"children"`
	ParentCode string `gorm:"column:parentCode" json:"parentCode" form:"parentCode"`
}

type Province struct {
	ID         int    `gorm:"primary_key" json:"id" form:"id"`
	Name       string `gorm:"column:name" json:"name" form:"name"`
	Level      int    `gorm:"column:level" json:"level" form:"level"`
	Code       string `gorm:"column:code" json:"code" form:"code"`
	Geocode    string `gorm:"column:geocode" json:"geocode" form:"geocode"`
	Latitude   string `gorm:"column:latitude" json:"latitude" form:"latitude"`
	Longitude  string `gorm:"column:longitude" json:"longitude" form:"longitude"`
	Sort       int    `gorm:"column:sort" json:"sort" form:"sort"`
	ParentCode string `gorm:"column:parentCode" json:"parentCode" form:"parentCode"`
}

func (Province) TableName() string {
	return "geo_provinces"
}

type City struct {
	ID           int    `gorm:"primary_key" json:"id" form:"id"`
	Name         string `gorm:"column:name" json:"name" form:"name"`
	Level        int    `gorm:"column:level" json:"level" form:"level"`
	Code         string `gorm:"column:code" json:"code" form:"code"`
	Geocode      string `gorm:"column:geocode" json:"geocode" form:"geocode"`
	Latitude     string `gorm:"column:latitude" json:"latitude" form:"latitude"`
	Longitude    string `gorm:"column:longitude" json:"longitude" form:"longitude"`
	Sort         int    `gorm:"column:sort" json:"sort" form:"sort"`
	ParentCode   string `gorm:"column:parentCode" json:"parentCode" form:"parentCode"`
	ProvinceCode string `gorm:"column:provinceCode" json:"provinceCode" form:"provinceCode"`
}

func (City) TableName() string {
	return "geo_cities"
}

type Area struct {
	ID           int    `gorm:"primary_key" json:"id" form:"id"`
	Name         string `gorm:"column:name" json:"name" form:"name"`
	Level        int    `gorm:"column:level" json:"level" form:"level"`
	Code         string `gorm:"column:code" json:"code" form:"code"`
	Geocode      string `gorm:"column:geocode" json:"geocode" form:"geocode"`
	Latitude     string `gorm:"column:latitude" json:"latitude" form:"latitude"`
	Longitude    string `gorm:"column:longitude" json:"longitude" form:"longitude"`
	Sort         int    `gorm:"column:sort" json:"sort" form:"sort"`
	CityCode     string `gorm:"column:cityCode" json:"cityCode" form:"cityCode"`
	ProvinceCode string `gorm:"column:provinceCode" json:"provinceCode" form:"provinceCode"`
	ParentCode   string `gorm:"column:parentCode" json:"parentCode" form:"parentCode"`
}

func (Area) TableName() string {
	return "geo_areas"
}
