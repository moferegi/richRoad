package client

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	ClientUserRouter
	AddressRouter
	CollectRouter
	PointRecordRouter
	TryonRechargeOrderRouter
	TryonTaskRouter
	TryonModelRouter
	TryonClothRouter
	VisitorRouter
	SysConfigRouter
	LanguageRouter
	PhoneAreaCodeRouter
	SignInRouter
	ExternalLinkDomainRouter
	VideoTagRouter
	DiaryTagRouter
	GameRouter
}

var cprApi = api.ApiGroupApp.ClientApiGroup.PointRecordApi
var videoTagApi = api.ApiGroupApp.ClientApiGroup.VideoTagApi
var diaryTagApi = api.ApiGroupApp.ClientApiGroup.DiaryTagApi
