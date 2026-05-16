package service

import (
	"errors"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/model"
)

type GeoService struct{}

const maxGeoLevel = 5

func tableByLevel(level int) string {
	if level <= 0 {
		return "geo_provinces"
	}
	if level == 1 {
		return "geo_cities"
	}
	return "geo_areas"
}

func (e *GeoService) CreateGeo(geo model.Geo) (g model.Geo, err error) {
	parentLevel := geo.Level
	childLevel := parentLevel + 1
	if childLevel < 0 || childLevel > maxGeoLevel {
		return geo, errors.New("无效层级")
	}
	geo.Level = childLevel
	err = global.GVA_DB.Table(tableByLevel(childLevel)).Create(&geo).Error
	return geo, err
}

func (e *GeoService) GetGeo(geo model.Geo) (g model.Geo, err error) {
	err = global.GVA_DB.Table(tableByLevel(geo.Level)).Where("id = ?", geo.ID).First(&geo).Error
	return geo, err
}

func (e *GeoService) GetGeos(level string, code string) (g []model.Geo, err error) {
	var Geos []model.Geo
	parentLevel, convErr := strconv.Atoi(level)
	if convErr != nil {
		parentLevel = -1
	}
	childLevel := parentLevel + 1
	if childLevel < 0 || childLevel > maxGeoLevel {
		return Geos, nil
	}
	err = global.GVA_DB.Table(tableByLevel(childLevel)).
		Where("parentCode = ? AND level = ?", code, childLevel).
		Order("sort asc, code asc").
		Find(&Geos).Error
	if err != nil {
		return Geos, err
	}
	for idx := range Geos {
		if childLevel >= maxGeoLevel {
			Geos[idx].HasChildren = false
			continue
		}
		nextLevel := childLevel + 1
		var count int64
		err = global.GVA_DB.Table(tableByLevel(nextLevel)).Where("parentCode = ? AND level = ?", Geos[idx].Code, nextLevel).Count(&count).Error
		if err != nil {
			return Geos, err
		}
		Geos[idx].HasChildren = count > 0
	}
	return Geos, err
}

func (e *GeoService) EditGeo(geo model.Geo) (g model.Geo, err error) {
	err = global.GVA_DB.Table(tableByLevel(geo.Level)).Where("code = ?", geo.Code).Updates(&geo).Error
	return geo, err
}

func (e *GeoService) DeleteGeo(geo model.Geo) (err error) {
	var count int64
	if geo.Level < maxGeoLevel {
		childLevel := geo.Level + 1
		global.GVA_DB.Table(tableByLevel(childLevel)).Where("parentCode = ? AND level = ?", geo.Code, childLevel).Count(&count)
		if count > 0 {
			return errors.New("存在子区域不允许删除")
		}
	}
	err = global.GVA_DB.Table(tableByLevel(geo.Level)).Where("id = ?", geo.ID).Delete(&model.Geo{}).Error
	return err
}
