package system

import (
	"context"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type initApi struct{}

const initOrderApi = system.InitOrderSystem + 1

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ApiGroup: "系统用户", Method: "DELETE", Path: "/user/deleteUser", Description: "删除用户"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/admin_register", Description: "用户注册"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/getUserList", Description: "获取用户列表"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setUserInfo", Description: "设置用户信息"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfInfo", Description: "设置自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "GET", Path: "/user/getUserInfo", Description: "获取自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthorities", Description: "设置权限组"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/changePassword", Description: "修改密码（建议选择)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthority", Description: "修改用户角色(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/resetPassword", Description: "重置用户密码"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfSetting", Description: "用户界面配置"},

		{ApiGroup: "api", Method: "POST", Path: "/api/createApi", Description: "创建api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/deleteApi", Description: "删除Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/updateApi", Description: "更新Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiList", Description: "获取api列表"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "获取所有api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiById", Description: "获取api详细信息"},
		{ApiGroup: "api", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "批量删除api"},
		{ApiGroup: "api", Method: "GET", Path: "/api/syncApi", Description: "获取待同步API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/getApiGroups", Description: "获取路由组"},
		{ApiGroup: "api", Method: "POST", Path: "/api/enterSyncApi", Description: "确认同步API"},
		{ApiGroup: "api", Method: "POST", Path: "/api/ignoreApi", Description: "忽略API"},

		{ApiGroup: "角色", Method: "POST", Path: "/authority/copyAuthority", Description: "拷贝角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/createAuthority", Description: "创建角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/deleteAuthority", Description: "删除角色"},
		{ApiGroup: "角色", Method: "PUT", Path: "/authority/updateAuthority", Description: "更新角色信息"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/getAuthorityList", Description: "获取角色列表"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/setDataAuthority", Description: "设置角色资源权限"},

		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "更改角色api权限"},
		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "获取权限列表"},

		{ApiGroup: "菜单", Method: "POST", Path: "/menu/addBaseMenu", Description: "新增菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenu", Description: "获取菜单树(必选)"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/deleteBaseMenu", Description: "删除菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/updateBaseMenu", Description: "更新菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuById", Description: "根据id获取菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuList", Description: "分页获取基础menu列表"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuTree", Description: "获取用户动态路由"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuAuthority", Description: "获取指定角色menu"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/addMenuAuthority", Description: "增加menu和角色关联关系"},

		{ApiGroup: "分片上传", Method: "GET", Path: "/fileUploadAndDownload/findFile", Description: "寻找目标文件（秒传）"},
		{ApiGroup: "分片上传", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinue", Description: "断点续传"},
		{ApiGroup: "分片上传", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinueFinish", Description: "断点续传完成"},
		{ApiGroup: "分片上传", Method: "POST", Path: "/fileUploadAndDownload/removeChunk", Description: "上传完成移除文件"},

		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "文件上传（建议选择）"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/deleteFile", Description: "删除文件"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/editFileName", Description: "文件名或者备注编辑"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/getFileList", Description: "获取上传文件列表"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/importURL", Description: "导入URL"},

		{ApiGroup: "系统服务", Method: "POST", Path: "/system/getServerInfo", Description: "获取服务器信息"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/system/getSystemConfig", Description: "获取配置文件内容"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/system/setSystemConfig", Description: "设置配置文件内容"},

		{ApiGroup: "客户", Method: "PUT", Path: "/customer/customer", Description: "更新客户"},
		{ApiGroup: "客户", Method: "POST", Path: "/customer/customer", Description: "创建客户"},
		{ApiGroup: "客户", Method: "DELETE", Path: "/customer/customer", Description: "删除客户"},
		{ApiGroup: "客户", Method: "GET", Path: "/customer/customer", Description: "获取单一客户"},
		{ApiGroup: "客户", Method: "GET", Path: "/customer/customerList", Description: "获取客户列表"},

		{ApiGroup: "代码生成器", Method: "GET", Path: "/autoCode/getDB", Description: "获取所有数据库"},
		{ApiGroup: "代码生成器", Method: "GET", Path: "/autoCode/getTables", Description: "获取数据库表"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/createTemp", Description: "自动化代码"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/preview", Description: "预览自动化代码"},
		{ApiGroup: "代码生成器", Method: "GET", Path: "/autoCode/getColumn", Description: "获取所选table的所有字段"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/installPlugin", Description: "安装插件"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/pubPlug", Description: "打包插件"},

		{ApiGroup: "模板配置", Method: "POST", Path: "/autoCode/createPackage", Description: "配置模板"},
		{ApiGroup: "模板配置", Method: "GET", Path: "/autoCode/getTemplates", Description: "获取模板文件"},
		{ApiGroup: "模板配置", Method: "POST", Path: "/autoCode/getPackage", Description: "获取所有模板"},
		{ApiGroup: "模板配置", Method: "POST", Path: "/autoCode/delPackage", Description: "删除模板"},

		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/getMeta", Description: "获取meta信息"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/rollback", Description: "回滚自动生成代码"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/getSysHistory", Description: "查询回滚记录"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/delSysHistory", Description: "删除回滚记录"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/addFunc", Description: "增加模板方法"},

		{ApiGroup: "系统字典详情", Method: "PUT", Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Description: "更新字典内容"},
		{ApiGroup: "系统字典详情", Method: "POST", Path: "/sysDictionaryDetail/createSysDictionaryDetail", Description: "新增字典内容"},
		{ApiGroup: "系统字典详情", Method: "DELETE", Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Description: "删除字典内容"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/sysDictionaryDetail/findSysDictionaryDetail", Description: "根据ID获取字典内容"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Description: "获取字典内容列表"},

		{ApiGroup: "系统字典", Method: "POST", Path: "/sysDictionary/createSysDictionary", Description: "新增字典"},
		{ApiGroup: "系统字典", Method: "DELETE", Path: "/sysDictionary/deleteSysDictionary", Description: "删除字典"},
		{ApiGroup: "系统字典", Method: "PUT", Path: "/sysDictionary/updateSysDictionary", Description: "更新字典"},
		{ApiGroup: "系统字典", Method: "GET", Path: "/sysDictionary/findSysDictionary", Description: "根据ID获取字典（建议选择）"},
		{ApiGroup: "系统字典", Method: "GET", Path: "/sysDictionary/getSysDictionaryList", Description: "获取字典列表"},

		{ApiGroup: "操作记录", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "新增操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "根据ID获取操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "获取操作记录列表"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "删除操作记录"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "批量删除操作历史"},

		{ApiGroup: "断点续传(插件版)", Method: "POST", Path: "/simpleUploader/upload", Description: "插件版分片上传"},
		{ApiGroup: "断点续传(插件版)", Method: "GET", Path: "/simpleUploader/checkFileMd5", Description: "文件完整度验证"},
		{ApiGroup: "断点续传(插件版)", Method: "GET", Path: "/simpleUploader/mergeFileMd5", Description: "上传完成合并文件"},

		{ApiGroup: "email", Method: "POST", Path: "/email/emailTest", Description: "发送测试邮件"},
		{ApiGroup: "email", Method: "POST", Path: "/email/sendEmail", Description: "发送邮件"},

		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/setAuthorityBtn", Description: "设置按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/getAuthorityBtn", Description: "获取已有按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/canRemoveAuthorityBtn", Description: "删除按钮"},

		{ApiGroup: "导出模板", Method: "POST", Path: "/sysExportTemplate/createSysExportTemplate", Description: "新增导出模板"},
		{ApiGroup: "导出模板", Method: "DELETE", Path: "/sysExportTemplate/deleteSysExportTemplate", Description: "删除导出模板"},
		{ApiGroup: "导出模板", Method: "DELETE", Path: "/sysExportTemplate/deleteSysExportTemplateByIds", Description: "批量删除导出模板"},
		{ApiGroup: "导出模板", Method: "PUT", Path: "/sysExportTemplate/updateSysExportTemplate", Description: "更新导出模板"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/sysExportTemplate/findSysExportTemplate", Description: "根据ID获取导出模板"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/sysExportTemplate/getSysExportTemplateList", Description: "获取导出模板列表"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/sysExportTemplate/exportExcel", Description: "导出Excel"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/sysExportTemplate/exportTemplate", Description: "下载模板"},
		{ApiGroup: "导出模板", Method: "POST", Path: "/sysExportTemplate/importExcel", Description: "导入Excel"},

		{ApiGroup: "公告", Method: "POST", Path: "/info/createInfo", Description: "新建公告"},
		{ApiGroup: "公告", Method: "DELETE", Path: "/info/deleteInfo", Description: "删除公告"},
		{ApiGroup: "公告", Method: "DELETE", Path: "/info/deleteInfoByIds", Description: "批量删除公告"},
		{ApiGroup: "公告", Method: "PUT", Path: "/info/updateInfo", Description: "更新公告"},
		{ApiGroup: "公告", Method: "GET", Path: "/info/findInfo", Description: "根据ID获取公告"},
		{ApiGroup: "公告", Method: "GET", Path: "/info/getInfoList", Description: "获取公告列表"},

		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/createClientUser", Description: "新增客户端用户"},
		{ApiGroup: "客户端用户", Method: "DELETE", Path: "/clientUser/deleteClientUser", Description: "删除客户端用户"},
		{ApiGroup: "客户端用户", Method: "DELETE", Path: "/clientUser/deleteClientUserByIds", Description: "批量删除客户端用户"},
		{ApiGroup: "客户端用户", Method: "PUT", Path: "/clientUser/updateClientUser", Description: "更新客户端用户"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/findClientUser", Description: "根据ID获取客户端用户"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getClientUserList", Description: "获取客户端用户列表"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getUserInfo", Description: "获取自身信息"},
		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/login", Description: "客户端登录"},
		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/register", Description: "注册客户端用户"},

		{ApiGroup: "轮播图", Method: "POST", Path: "/banner/createBanner", Description: "新增轮播图"},
		{ApiGroup: "轮播图", Method: "DELETE", Path: "/banner/deleteBanner", Description: "删除轮播图"},
		{ApiGroup: "轮播图", Method: "DELETE", Path: "/banner/deleteBannerByIds", Description: "批量删除轮播图"},
		{ApiGroup: "轮播图", Method: "PUT", Path: "/banner/updateBanner", Description: "更新轮播图"},
		{ApiGroup: "轮播图", Method: "GET", Path: "/banner/findBanner", Description: "根据ID获取轮播图"},
		{ApiGroup: "轮播图", Method: "GET", Path: "/banner/getBannerList", Description: "获取轮播图列表"},

		{ApiGroup: "商品分类", Method: "POST", Path: "/category/createCategory", Description: "新增商品分类"},
		{ApiGroup: "商品分类", Method: "DELETE", Path: "/category/deleteCategory", Description: "删除商品分类"},
		{ApiGroup: "商品分类", Method: "DELETE", Path: "/category/deleteCategoryByIds", Description: "批量删除商品分类"},
		{ApiGroup: "商品分类", Method: "PUT", Path: "/category/updateCategory", Description: "更新商品分类"},
		{ApiGroup: "商品分类", Method: "GET", Path: "/category/findCategory", Description: "根据ID获取商品分类"},
		{ApiGroup: "商品分类", Method: "GET", Path: "/category/getCategoryList", Description: "获取商品分类列表"},

		{ApiGroup: "商品", Method: "POST", Path: "/good/createGood", Description: "新增商品"},
		{ApiGroup: "商品", Method: "DELETE", Path: "/good/deleteGood", Description: "删除商品"},
		{ApiGroup: "商品", Method: "DELETE", Path: "/good/deleteGoodByIds", Description: "批量删除商品"},
		{ApiGroup: "商品", Method: "PUT", Path: "/good/updateGood", Description: "更新商品"},
		{ApiGroup: "商品", Method: "GET", Path: "/good/findGood", Description: "根据ID获取商品"},
		{ApiGroup: "商品", Method: "GET", Path: "/good/getGoodList", Description: "获取商品列表"},

		{ApiGroup: "sku", Method: "POST", Path: "/sku/createSku", Description: "新增sku"},
		{ApiGroup: "sku", Method: "DELETE", Path: "/sku/deleteSku", Description: "删除sku"},
		{ApiGroup: "sku", Method: "DELETE", Path: "/sku/deleteSkuByIds", Description: "批量删除sku"},
		{ApiGroup: "sku", Method: "PUT", Path: "/sku/updateSku", Description: "更新sku"},
		{ApiGroup: "sku", Method: "GET", Path: "/sku/findSku", Description: "根据ID获取sku"},
		{ApiGroup: "sku", Method: "GET", Path: "/sku/getSkuList", Description: "获取sku列表"},

		{ApiGroup: "购物车", Method: "POST", Path: "/cart/createCart", Description: "新增购物车"},
		{ApiGroup: "购物车", Method: "DELETE", Path: "/cart/deleteCart", Description: "删除购物车"},
		{ApiGroup: "购物车", Method: "DELETE", Path: "/cart/deleteCartByIds", Description: "批量删除购物车"},
		{ApiGroup: "购物车", Method: "PUT", Path: "/cart/updateCart", Description: "更新购物车"},
		{ApiGroup: "购物车", Method: "GET", Path: "/cart/findCart", Description: "根据ID获取购物车"},
		{ApiGroup: "购物车", Method: "GET", Path: "/cart/getCartList", Description: "获取购物车列表"},
		{ApiGroup: "购物车", Method: "POST", Path: "/cart/cutCart", Description: "删除购物车"},
		{ApiGroup: "购物车", Method: "POST", Path: "/cart/addCart", Description: "添加购物车"},
		{ApiGroup: "购物车", Method: "GET", Path: "/cart/getSelfCart", Description: "获取自身购物车"},
		{ApiGroup: "购物车", Method: "GET", Path: "/cart/clearCart", Description: "全部删除购物车"},

		{ApiGroup: "订单", Method: "POST", Path: "/order/createOrder", Description: "新增订单"},
		{ApiGroup: "订单", Method: "DELETE", Path: "/order/deleteOrder", Description: "删除订单"},
		{ApiGroup: "订单", Method: "DELETE", Path: "/order/deleteOrderByIds", Description: "批量删除订单"},
		{ApiGroup: "订单", Method: "PUT", Path: "/order/updateOrder", Description: "更新订单"},
		{ApiGroup: "订单", Method: "GET", Path: "/order/findOrder", Description: "根据ID获取订单"},
		{ApiGroup: "订单", Method: "GET", Path: "/order/getOrderList", Description: "获取订单列表"},
		{ApiGroup: "订单", Method: "POST", Path: "/order/placeOrder", Description: "直接下单"},
		{ApiGroup: "订单", Method: "POST", Path: "/order/placeOrderByCart", Description: "购物车下单"},
		{ApiGroup: "订单", Method: "GET", Path: "/order/selfOrderList", Description: "我的订单列表"},
		{ApiGroup: "订单", Method: "GET", Path: "/order/selfOrder", Description: "我的订单详情"},
		{ApiGroup: "订单", Method: "POST", Path: "/order/updateOrderStatus", Description: "变更订单状态"},
		{ApiGroup: "订单", Method: "GET", Path: "/order/selfOrderComment", Description: "获取单商品评价"},
		{ApiGroup: "订单", Method: "GET", Path: "/order/checkRouters", Description: "获取快递路径"},
		{ApiGroup: "订单", Method: "POST", Path: "/order/updateOrder", Description: "更新订单"},

		{ApiGroup: "城市管理", Method: "GET", Path: "/geo/getGeos", Description: "获取城市列表"},
		{ApiGroup: "城市管理", Method: "GET", Path: "/geo/getGeo", Description: "获取单一城市"},
		{ApiGroup: "城市管理", Method: "PUT", Path: "/geo/editGeo", Description: "修改城市"},
		{ApiGroup: "城市管理", Method: "POST", Path: "/geo/createGeo", Description: "创建城市"},
		{ApiGroup: "城市管理", Method: "DELETE", Path: "/geo/deleteGeo", Description: "删除城市"},

		{ApiGroup: "用户地址", Method: "POST", Path: "/address/createAddress", Description: "新增用户地址"},
		{ApiGroup: "用户地址", Method: "DELETE", Path: "/address/deleteAddress", Description: "删除用户地址"},
		{ApiGroup: "用户地址", Method: "DELETE", Path: "/address/deleteAddressByIds", Description: "批量删除用户地址"},
		{ApiGroup: "用户地址", Method: "PUT", Path: "/address/updateAddress", Description: "更新用户地址"},
		{ApiGroup: "用户地址", Method: "GET", Path: "/address/findAddress", Description: "根据ID获取用户地址"},
		{ApiGroup: "用户地址", Method: "GET", Path: "/address/getAddressList", Description: "获取用户地址列表"},
		{ApiGroup: "用户地址", Method: "GET", Path: "/address/getDefaultAddress", Description: "获取默认用户地址"},

		{ApiGroup: "收藏", Method: "POST", Path: "/collect/createCollect", Description: "新增收藏"},
		{ApiGroup: "收藏", Method: "DELETE", Path: "/collect/deleteCollect", Description: "删除收藏"},
		{ApiGroup: "收藏", Method: "DELETE", Path: "/collect/deleteCollectByIds", Description: "批量删除收藏"},
		{ApiGroup: "收藏", Method: "PUT", Path: "/collect/updateCollect", Description: "更新收藏"},
		{ApiGroup: "收藏", Method: "GET", Path: "/collect/findCollect", Description: "根据ID获取收藏"},
		{ApiGroup: "收藏", Method: "GET", Path: "/collect/getCollectList", Description: "获取收藏列表"},

		{ApiGroup: "支付", Method: "GET", Path: "/wxpay/getOrderById", Description: "获取订单"},
		{ApiGroup: "支付", Method: "POST", Path: "/wxpay/getPayParams", Description: "获取支付参数"},
		{ApiGroup: "支付", Method: "POST", Path: "/wxpay/payAction", Description: "支付完成"},

		{ApiGroup: "用户评论", Method: "POST", Path: "/comment/createComment", Description: "新增用户评论"},
		{ApiGroup: "用户评论", Method: "DELETE", Path: "/comment/deleteComment", Description: "删除用户评论"},
		{ApiGroup: "用户评论", Method: "DELETE", Path: "/comment/deleteCommentByIds", Description: "批量删除用户评论"},
		{ApiGroup: "用户评论", Method: "PUT", Path: "/comment/updateComment", Description: "更新用户评论"},
		{ApiGroup: "用户评论", Method: "GET", Path: "/comment/findComment", Description: "根据ID获取用户评论"},
		{ApiGroup: "用户评论", Method: "GET", Path: "/comment/getComment", Description: "后台获取用户评论"},
		{ApiGroup: "用户评论", Method: "GET", Path: "/comment/getCommentList", Description: "获取用户评论列表"},

		{ApiGroup: "参数管理", Method: "POST", Path: "/sysParams/createSysParams", Description: "新建参数"},
		{ApiGroup: "参数管理", Method: "DELETE", Path: "/sysParams/deleteSysParams", Description: "删除参数"},
		{ApiGroup: "参数管理", Method: "DELETE", Path: "/sysParams/deleteSysParamsByIds", Description: "批量删除参数"},
		{ApiGroup: "参数管理", Method: "PUT", Path: "/sysParams/updateSysParams", Description: "更新参数"},
		{ApiGroup: "参数管理", Method: "GET", Path: "/sysParams/findSysParams", Description: "根据ID获取参数"},
		{ApiGroup: "参数管理", Method: "GET", Path: "/sysParams/getSysParamsList", Description: "获取参数列表"},
		{ApiGroup: "参数管理", Method: "GET", Path: "/sysParams/getSysParam", Description: "获取参数列表"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/authorityBtn/canRemoveAuthorityBtn", "POST").
		First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
