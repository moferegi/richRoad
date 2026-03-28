package client

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ClientUserApi
	AddressApi
	CollectApi
	PointRecordApi
	VisitorApi
	SysConfigApi
}

var (
	jwtService       = service.ServiceGroupApp.SystemServiceGroup.JwtService
	cprService       = service.ServiceGroupApp.ClientServiceGroup.PointRecordService
	sysConfigService = service.ServiceGroupApp.ClientServiceGroup.SysConfigService
)
