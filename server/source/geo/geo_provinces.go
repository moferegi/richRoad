package geo

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/model"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//
//geo_areas
//geo_provinces

type initProvince struct{}

const initOrderProvince = initOrderCity + 1

// auto run
func init() {
	system.RegisterInit(initOrderProvince, &initProvince{})
}

func (i initProvince) InitializerName() string {
	return model.Province{}.TableName()
}

func (i *initProvince) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&model.Province{})
}

func (i *initProvince) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&model.Province{})
}

func (i *initProvince) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []model.Province{

		{Name: "北京市", Level: 0, Code: "11", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "天津市", Level: 0, Code: "12", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "河北省", Level: 0, Code: "13", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "山西省", Level: 0, Code: "14", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "内蒙古自治区", Level: 0, Code: "15", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "辽宁省", Level: 0, Code: "21", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "吉林省", Level: 0, Code: "22", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "黑龙江省", Level: 0, Code: "23", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "上海市", Level: 0, Code: "31", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "江苏省", Level: 0, Code: "32", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "浙江省", Level: 0, Code: "33", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "安徽省", Level: 0, Code: "34", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "福建省", Level: 0, Code: "35", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "江西省", Level: 0, Code: "36", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "山东省", Level: 0, Code: "37", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "河南省", Level: 0, Code: "41", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "湖北省", Level: 0, Code: "42", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "湖南省", Level: 0, Code: "43", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "广东省", Level: 0, Code: "44", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "广西壮族自治区", Level: 0, Code: "45", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "海南省", Level: 0, Code: "46", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "重庆市", Level: 0, Code: "50", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "四川省", Level: 0, Code: "51", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "贵州省", Level: 0, Code: "52", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "云南省", Level: 0, Code: "53", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "西藏自治区", Level: 0, Code: "54", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "陕西省", Level: 0, Code: "61", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "甘肃省", Level: 0, Code: "62", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "青海省", Level: 0, Code: "63", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "宁夏回族自治区", Level: 0, Code: "64", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
		{Name: "新疆维吾尔自治区", Level: 0, Code: "65", Geocode: "", Latitude: "", Longitude: "", Sort: 0, ParentCode: "0"},
	}
	if err := db.Create(&entities).Error; err != nil {

		return ctx, errors.Wrap(err, "geo_cities表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initProvince) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.First(&model.Province{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
