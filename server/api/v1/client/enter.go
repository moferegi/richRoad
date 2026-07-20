package client

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ClientUserApi
	AddressApi
	CollectApi
	PointRecordApi
	TryonRechargeOrderApi
	TryonTaskApi
	TryonModelApi
	TryonClothApi
	VisitorApi
	SysConfigApi
	LanguageApi
	PhoneAreaCodeApi
	SignInApi
	ExternalLinkDomainApi
	VideoTagApi
}

var (
	jwtService             = service.ServiceGroupApp.SystemServiceGroup.JwtService
	cprService             = service.ServiceGroupApp.ClientServiceGroup.PointRecordService
	sysConfigService       = service.ServiceGroupApp.ClientServiceGroup.SysConfigService
	marketingRewardService = service.ServiceGroupApp.ShopServiceGroup.MarketingRewardService
	extDomainService       = service.ServiceGroupApp.ClientServiceGroup.ExternalLinkDomainService
	videoTagService        = service.ServiceGroupApp.ClientServiceGroup.VideoTagService
	cloudStorageService    = service.ServiceGroupApp.ClientServiceGroup.CloudStorageService
)
