package system

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbin = initOrderApiIgnore + 1

type initCasbin struct{}

// auto run
func init() {
	system.RegisterInit(initOrderCasbin, &initCasbin{})
}

func (i *initCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&adapter.CasbinRule{})
}

func (i *initCasbin) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i *initCasbin) InitializerName() string {
	var entity adapter.CasbinRule
	return entity.TableName()
}

func (i *initCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []adapter.CasbinRule{
		{Ptype: "p", V0: "888", V1: "/user/admin_register", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/createApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/getApiById", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/updateApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/getAllApis", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/deleteApisByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/api/syncApi", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/getApiGroups", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/enterSyncApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/ignoreApi", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/authority/copyAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/updateAuthority", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/authority/createAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/deleteAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/getAuthorityList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authority/setDataAuthority", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/menu/getMenu", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/getMenuList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/addBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/getBaseMenuTree", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/addMenuAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/getMenuAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/deleteBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/updateBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/menu/getBaseMenuById", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/user/setUserInfo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/user/setSelfInfo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/deleteUser", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/user/changePassword", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/setUserAuthority", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/setUserAuthorities", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/resetPassword", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/user/setSelfSetting", V2: "PUT"},

		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/findFile", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/breakpointContinueFinish", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/breakpointContinue", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/removeChunk", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/upload", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/deleteFile", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/editFileName", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/getFileList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/fileUploadAndDownload/importURL", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/casbin/updateCasbin", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/casbin/getPolicyPathByAuthorityId", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/jwt/jsonInBlacklist", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/system/getSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/system/setSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/system/getServerInfo", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/customer/customer", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/customer/customerList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/autoCode/getDB", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/autoCode/getMeta", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/preview", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/getTables", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/autoCode/getColumn", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/autoCode/rollback", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/createTemp", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/delSysHistory", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/getSysHistory", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/createPackage", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/getTemplates", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/autoCode/getPackage", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/delPackage", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/createPlug", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/installPlugin", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/pubPlug", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/addFunc", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/mcp", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/mcpTest", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/autoCode/mcpList", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/sysDictionaryDetail/findSysDictionaryDetail", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysDictionaryDetail/updateSysDictionaryDetail", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysDictionaryDetail/createSysDictionaryDetail", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysDictionaryDetail/getSysDictionaryDetailList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysDictionaryDetail/deleteSysDictionaryDetail", V2: "DELETE"},

		{Ptype: "p", V0: "888", V1: "/sysDictionary/findSysDictionary", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysDictionary/updateSysDictionary", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysDictionary/getSysDictionaryList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysDictionary/createSysDictionary", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysDictionary/deleteSysDictionary", V2: "DELETE"},

		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/findSysOperationRecord", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/updateSysOperationRecord", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/createSysOperationRecord", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/getSysOperationRecordList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/deleteSysOperationRecord", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysOperationRecord/deleteSysOperationRecordByIds", V2: "DELETE"},

		{Ptype: "p", V0: "888", V1: "/email/emailTest", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/email/sendEmail", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/simpleUploader/upload", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/simpleUploader/checkFileMd5", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/simpleUploader/mergeFileMd5", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/authorityBtn/setAuthorityBtn", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authorityBtn/getAuthorityBtn", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/authorityBtn/canRemoveAuthorityBtn", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/createSysExportTemplate", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/deleteSysExportTemplate", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/deleteSysExportTemplateByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/updateSysExportTemplate", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/findSysExportTemplate", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/getSysExportTemplateList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/exportExcel", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/exportTemplate", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysExportTemplate/importExcel", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/info/createInfo", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/info/deleteInfo", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/info/deleteInfoByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/info/updateInfo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/info/findInfo", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/info/getInfoList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/clientUser/createClientUser", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/clientUser/deleteClientUser", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/clientUser/deleteClientUserByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/clientUser/updateClientUser", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/clientUser/findClientUser", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/clientUser/getClientUserList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/clientUser/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/clientUser/login", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/clientUser/register", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/banner/createBanner", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/banner/deleteBanner", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/banner/deleteBannerByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/banner/updateBanner", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/banner/findBanner", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/banner/getBannerList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/category/createCategory", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/category/deleteCategory", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/category/deleteCategoryByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/category/updateCategory", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/category/findCategory", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/category/getCategoryList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/good/createGood", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/good/deleteGood", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/good/deleteGoodByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/good/updateGood", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/good/findGood", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/good/getGoodList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/sku/createSku", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sku/deleteSku", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sku/deleteSkuByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sku/updateSku", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sku/findSku", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sku/getSkuList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/cart/createCart", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/cart/deleteCart", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/cart/deleteCartByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/cart/updateCart", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/cart/findCart", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/cart/getCartList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/cart/cutCart", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/cart/addCart", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/cart/getSelfCart", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/cart/clearCart", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/order/createOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/order/deleteOrder", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/order/deleteOrderByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/order/updateOrder", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/order/findOrder", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/order/getOrderList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/order/placeOrder", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/order/placeOrderByCart", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/order/selfOrderList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/order/selfOrder", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/order/updateOrderStatus", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/order/selfOrderComment", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/order/checkRouters", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/order/updateOrder", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/geo/getGeos", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/geo/getGeo", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/geo/editGeo", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/geo/createGeo", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/geo/deleteGeo", V2: "DELETE"},

		{Ptype: "p", V0: "888", V1: "/address/createAddress", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/address/deleteAddress", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/address/deleteAddressByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/address/updateAddress", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/address/findAddress", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/address/getAddressList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/address/getDefaultAddress", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/collect/createCollect", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/collect/deleteCollect", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/collect/deleteCollectByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/collect/updateCollect", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/collect/findCollect", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/collect/getCollectList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/wxpay/getOrderById", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/wxpay/getPayParams", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/wxpay/payAction", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/comment/createComment", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/comment/deleteComment", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/comment/deleteCommentByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/comment/updateComment", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/comment/findComment", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/comment/getComment", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/comment/getCommentList", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/sysParams/createSysParams", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/sysParams/deleteSysParams", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysParams/deleteSysParamsByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/sysParams/updateSysParams", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/sysParams/findSysParams", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysParams/getSysParamsList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/sysParams/getSysParam", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/attachmentCategory/getCategoryList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/attachmentCategory/addCategory", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/attachmentCategory/deleteCategory", V2: "POST"},

		// 标签管理权限
		{Ptype: "p", V0: "888", V1: "/tag/createTag", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/tag/deleteTag", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/tag/deleteTagByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/tag/updateTag", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/tag/findTag", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/tag/getTagList", V2: "GET"},

		// 优惠券管理权限
		{Ptype: "p", V0: "888", V1: "/cou/adminIssueCouponToAll", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/cou/claimCouponByUser", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/Cou/createCoupon", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/cou/createCouponOrderUser", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/Cou/deleteCoupon", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/Cou/deleteCouponByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/cou/deleteCouponOrderUser", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/cou/deleteCouponOrderUserByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/Cou/findCoupon", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/cou/findCouponOrderUser", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/cou/getAllClaimCoupon", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/Cou/getCouponList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/cou/getCouponOrderUserList", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/Cou/updateCoupon", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/cou/updateCouponOrderUser", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/order/changeOrderCoupon", V2: "POST"},

		// 促销信息管理权限
		{Ptype: "p", V0: "888", V1: "/promo/createPromotion", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/promo/deletePromotion", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/promo/deletePromotionByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/promo/updatePromotion", V2: "PUT"},
		{Ptype: "p", V0: "888", V1: "/promo/findPromotion", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/promo/getPromotionList", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/user/admin_register", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/createApi", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/getApiById", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/updateApi", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/getAllApis", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/createAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/deleteAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/getAuthorityList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/authority/setDataAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/getMenu", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/getMenuList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/addBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/getBaseMenuTree", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/addMenuAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/getMenuAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/deleteBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/updateBaseMenu", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/menu/getBaseMenuById", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/user/changePassword", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/user/setUserAuthority", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/upload", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/getFileList", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/deleteFile", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/editFileName", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/fileUploadAndDownload/importURL", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/casbin/updateCasbin", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/casbin/getPolicyPathByAuthorityId", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/jwt/jsonInBlacklist", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/system/getSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/system/setSystemConfig", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "PUT"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "DELETE"},
		{Ptype: "p", V0: "8881", V1: "/customer/customer", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/customer/customerList", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/user/getUserInfo", V2: "GET"},

		// 8881角色的促销信息管理权限
		{Ptype: "p", V0: "8881", V1: "/promo/createPromotion", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/promo/deletePromotion", V2: "DELETE"},
		{Ptype: "p", V0: "8881", V1: "/promo/deletePromotionByIds", V2: "DELETE"},
		{Ptype: "p", V0: "8881", V1: "/promo/updatePromotion", V2: "PUT"},
		{Ptype: "p", V0: "8881", V1: "/promo/findPromotion", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/promo/getPromotionList", V2: "GET"},

		{Ptype: "p", V0: "8080", V1: "/address/createAddress", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/address/deleteAddress", V2: "DELETE"},
		{Ptype: "p", V0: "8080", V1: "/address/deleteAddressByIds", V2: "DELETE"},
		{Ptype: "p", V0: "8080", V1: "/address/findAddress", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/address/getAddressList", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/address/getDefaultAddress", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/address/updateAddress", V2: "PUT"},
		{Ptype: "p", V0: "8080", V1: "/cart/addCart", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/cart/clearCart", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/cart/cutCart", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/cart/getSelfCart", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/clientUser/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/collect/createCollect", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/collect/deleteCollect", V2: "DELETE"},
		{Ptype: "p", V0: "8080", V1: "/collect/deleteCollectByIds", V2: "DELETE"},
		{Ptype: "p", V0: "8080", V1: "/collect/findCollect", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/collect/getCollectList", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/collect/updateCollect", V2: "PUT"},
		{Ptype: "p", V0: "8080", V1: "/comment/createComment", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/comment/deleteComment", V2: "DELETE"},
		{Ptype: "p", V0: "8080", V1: "/comment/deleteCommentByIds", V2: "DELETE"},
		{Ptype: "p", V0: "8080", V1: "/comment/findComment", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/comment/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/comment/updateComment", V2: "PUT"},
		{Ptype: "p", V0: "8080", V1: "/fileUploadAndDownload/upload", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/order/checkRouters", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/order/placeOrder", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/order/placeOrderByCart", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/order/selfOrder", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/order/selfOrderComment", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/order/selfOrderList", V2: "GET"},
		{Ptype: "p", V0: "8080", V1: "/order/updateOrderStatus", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/order/updateOrder", V2: "POST"},

		// 8080角色的优惠券权限
		{Ptype: "p", V0: "8080", V1: "/cou/claimCouponByUser", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/cou/getAllClaimCoupon", V2: "POST"},
		{Ptype: "p", V0: "8080", V1: "/order/changeOrderCoupon", V2: "POST"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, "Casbin 表 ("+i.InitializerName()+") 数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initCasbin) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/user/getUserInfo", V2: "GET"}).
		First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
