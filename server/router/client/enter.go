package client

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	ClientUserRouter
	AddressRouter
	CollectRouter
	PointRecordRouter
}

var cprApi = api.ApiGroupApp.ClientApiGroup.PointRecordApi
