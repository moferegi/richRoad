package client

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	ClientUserRouter
	AddressRouter
	CollectRouter
	PointRecordRouter
	TryonTaskRouter
	VisitorRouter
	SysConfigRouter
	LanguageRouter
	PhoneAreaCodeRouter
	SignInRouter
	ExternalLinkDomainRouter
}

var cprApi = api.ApiGroupApp.ClientApiGroup.PointRecordApi
