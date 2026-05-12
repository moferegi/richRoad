package i18n

import "github.com/gin-gonic/gin"

// 语言代码常量
const (
	LangZh = "zh"
	LangEn = "en"
	LangMn = "mn"
)

// 消息字典：key -> { lang -> text }
var messages = map[string]map[string]string{
	// ========== 通用 CRUD ==========
	"ok":                              {LangZh: "操作成功", LangEn: "Success", LangMn: "Амжилттай"},
	"fail":                            {LangZh: "操作失败", LangEn: "Operation failed", LangMn: "Амжилтгүй"},
	"success":                         {LangZh: "成功", LangEn: "Success", LangMn: "Амжилттай"},
	"createSuccess":                   {LangZh: "创建成功", LangEn: "Created successfully", LangMn: "Амжилттай үүсгэлээ"},
	"createFail":                      {LangZh: "创建失败", LangEn: "Create failed", LangMn: "Үүсгэж чадсангүй"},
	"deleteSuccess":                   {LangZh: "删除成功", LangEn: "Deleted successfully", LangMn: "Амжилттай устгалаа"},
	"deleteFail":                      {LangZh: "删除失败", LangEn: "Delete failed", LangMn: "Устгаж чадсангүй"},
	"batchDeleteSuccess":              {LangZh: "批量删除成功", LangEn: "Batch deleted", LangMn: "Бөөнөөр устгалаа"},
	"batchDeleteFail":                 {LangZh: "批量删除失败", LangEn: "Batch delete failed", LangMn: "Бөөнөөр устгаж чадсангүй"},
	"updateSuccess":                   {LangZh: "更新成功", LangEn: "Updated successfully", LangMn: "Амжилттай шинэчлэлээ"},
	"updateFail":                      {LangZh: "更新失败", LangEn: "Update failed", LangMn: "Шинэчлэж чадсангүй"},
	"querySuccess":                    {LangZh: "查询成功", LangEn: "Query successful", LangMn: "Хайлт амжилттай"},
	"queryFail":                       {LangZh: "查询失败", LangEn: "Query failed", LangMn: "Хайлт амжилтгүй"},
	"getSuccess":                      {LangZh: "获取成功", LangEn: "Retrieved successfully", LangMn: "Амжилттай авлаа"},
	"getFail":                         {LangZh: "获取失败", LangEn: "Retrieve failed", LangMn: "Авч чадсангүй"},
	"setSuccess":                      {LangZh: "设置成功", LangEn: "Set successfully", LangMn: "Амжилттай тохируулалаа"},
	"setFail":                         {LangZh: "设置失败", LangEn: "Set failed", LangMn: "Тохируулж чадсангүй"},
	"invalidParams":                   {LangZh: "参数错误", LangEn: "Invalid parameters", LangMn: "Параметр буруу"},
	"invalidID":                       {LangZh: "ID参数错误", LangEn: "Invalid ID parameter", LangMn: "ID параметр буруу"},
	"invalidIDs":                      {LangZh: "IDs参数错误", LangEn: "Invalid IDs parameter", LangMn: "IDs параметр буруу"},
	"invalidTaskNo":                   {LangZh: "taskNo参数错误", LangEn: "Invalid task number parameter", LangMn: "taskNo параметр буруу"},
	"invalidOrder":                    {LangZh: "订单参数错误", LangEn: "Invalid order parameter", LangMn: "Захиалгын параметр буруу"},
	"userIDRequired":                  {LangZh: "用户ID不能为空", LangEn: "User ID is required", LangMn: "Хэрэглэгчийн ID хоосон байж болохгүй"},
	"visitorIDRequired":               {LangZh: "visitorId不能为空", LangEn: "visitorId is required", LangMn: "visitorId хоосон байж болохгүй"},
	"dateRequired":                    {LangZh: "date不能为空", LangEn: "date is required", LangMn: "date хоосон байж болохгүй"},
	"configGroupRequired":             {LangZh: "configGroup不能为空", LangEn: "configGroup is required", LangMn: "configGroup хоосон байж болохгүй"},
	"configKeyRequired":               {LangZh: "configKey不能为空", LangEn: "configKey is required", LangMn: "configKey хоосон байж болохгүй"},
	"configNotPublic":                 {LangZh: "该配置不对外开放", LangEn: "This config is not public", LangMn: "Энэ тохиргоо олон нийтэд нээлттэй биш"},
	"noPermission":                    {LangZh: "无权限操作", LangEn: "Permission denied", LangMn: "Эрх хүрэлцэхгүй"},
	"adjustSuccess":                   {LangZh: "调整成功", LangEn: "Adjusted successfully", LangMn: "Амжилттай тохирууллаа"},
	"adjustAmountMustPositive":        {LangZh: "调整数量必须大于0", LangEn: "Adjustment amount must be greater than 0", LangMn: "Тохируулгын хэмжээ 0-оос их байх ёстой"},
	"adjustReasonRequired":            {LangZh: "调整原因不能为空", LangEn: "Adjustment reason is required", LangMn: "Тохируулгын шалтгаан хоосон байж болохгүй"},
	"pointChangeTypeInvalid":          {LangZh: "增减类型必须是 increase 或 decrease", LangEn: "Change type must be increase or decrease", LangMn: "Өөрчлөлтийн төрөл нь increase эсвэл decrease байх ёстой"},
	"pointRecordRequiredFields":       {LangZh: "用户ID、增减类型和积分变化不能为空", LangEn: "User ID, change type and point change are required", LangMn: "Хэрэглэгчийн ID, өөрчлөлтийн төрөл, онооны өөрчлөлт хоосон байж болохгүй"},
	"assetTypeMustPointOrTryonPoint":  {LangZh: "资产类型必须是 point 或 tryon_point", LangEn: "Asset type must be point or tryon_point", LangMn: "Хөрөнгийн төрөл нь point эсвэл tryon_point байх ёстой"},
	"pointInsufficient":               {LangZh: "积分不足，无法扣除", LangEn: "Insufficient points, cannot deduct", LangMn: "Оноо хүрэлцэхгүй тул суутгах боломжгүй"},
	"tryonPointInsufficient":          {LangZh: "试衣币不足，无法扣除", LangEn: "Insufficient try-on points, cannot deduct", LangMn: "Өмсгөл туршилтын оноо хүрэлцэхгүй тул суутгах боломжгүй"},
	"submitSuccess":                   {LangZh: "提交成功", LangEn: "Submitted successfully", LangMn: "Амжилттай илгээлээ"},
	"cancelSuccess":                   {LangZh: "取消成功", LangEn: "Cancelled successfully", LangMn: "Амжилттай цуцаллаа"},
	"confirmSuccess":                  {LangZh: "确认成功", LangEn: "Confirmed successfully", LangMn: "Амжилттай баталгаажууллаа"},
	"tryonTaskDone":                   {LangZh: "试衣任务处理完成", LangEn: "Try-on task completed", LangMn: "Өмсгөл туршилтын даалгавар дууслаа"},
	"tryonTaskProcessing":             {LangZh: "试衣任务已提交，处理中", LangEn: "Try-on task submitted and processing", LangMn: "Өмсгөл туршилтын даалгавар илгээгдэж, боловсруулж байна"},
	"tryonBeautifyDone":               {LangZh: "智能美肤处理完成", LangEn: "Smart skin beautify completed", LangMn: "Ухаалаг арьс сайжруулалт дууслаа"},
	"tryonBeautifyProcessing":         {LangZh: "智能美肤处理中", LangEn: "Smart skin beautify is processing", LangMn: "Ухаалаг арьс сайжруулалт боловсруулж байна"},
	"tryonBeautifyUnsupported":        {LangZh: "当前模型不支持智能美肤", LangEn: "Current model does not support smart skin beautify", LangMn: "Одоогийн загвар ухаалаг арьс сайжруулалтыг дэмжихгүй"},
	"tryonBeautifyAlreadyUsed":        {LangZh: "每个试衣任务仅可使用一次智能美肤", LangEn: "Smart skin beautify can be used only once per try-on task", LangMn: "Өмсгөлийн даалгавар бүрт ухаалаг сайжруулалтыг нэг л удаа ашиглана"},
	"tryonTaskNotReadyForBeautify":    {LangZh: "试衣任务未完成，暂不可进行智能美肤", LangEn: "Try-on task is not completed, beautify is unavailable", LangMn: "Өмсгөлийн даалгавар дуусаагүй тул сайжруулалт хийх боломжгүй"},
	"requestReusedHistory":            {LangZh: "请求已受理，返回历史任务", LangEn: "Request reused an existing task", LangMn: "Хүсэлтийг хүлээн авч, өмнөх даалгаврыг буцаалаа"},
	"orderNotFound":                   {LangZh: "订单不存在", LangEn: "Order not found", LangMn: "Захиалга олдсонгүй"},
	"payMethodRequired":               {LangZh: "支付方式不能为空", LangEn: "Payment method is required", LangMn: "Төлбөрийн арга хоосон байж болохгүй"},
	"loginRequired":                   {LangZh: "请先登录", LangEn: "Please login first", LangMn: "Эхлээд нэвтэрнэ үү"},
	"tryonRechargeAmountFormatError":  {LangZh: "充值金额格式不正确", LangEn: "Invalid recharge amount format", LangMn: "Цэнэглэх дүнгийн формат буруу"},
	"tryonRechargeOrderNoGenFail":     {LangZh: "订单号生成失败，请稍后重试", LangEn: "Failed to generate order number, please try again later", LangMn: "Захиалгын дугаар үүсгэж чадсангүй, дараа дахин оролдоно уу"},
	"tryonRechargePlanNotConfigured":  {LangZh: "充值套餐未配置", LangEn: "Recharge plans are not configured", LangMn: "Цэнэглэх багцууд тохируулагдаагүй байна"},
	"tryonRechargePlanConfigInvalid":  {LangZh: "充值套餐配置格式错误", LangEn: "Recharge plan configuration format is invalid", LangMn: "Цэнэглэх багцын тохиргооны формат буруу"},
	"tryonRechargePlanChanged":        {LangZh: "充值套餐已变更，请刷新页面后重试", LangEn: "Recharge plans changed, please refresh and try again", LangMn: "Цэнэглэх багц өөрчлөгдсөн тул хуудсыг шинэчлээд дахин оролдоно уу"},
	"tryonRechargePointsMustPositive": {LangZh: "充值点数必须大于0", LangEn: "Recharge points must be greater than 0", LangMn: "Цэнэглэх оноо 0-оос их байх ёстой"},
	"tryonRechargeAmountMustPositive": {LangZh: "充值金额必须大于0", LangEn: "Recharge amount must be greater than 0", LangMn: "Цэнэглэх дүн 0-оос их байх ёстой"},
	"tryonRechargeOrderStateInvalidForSubmitPayment": {LangZh: "当前订单状态不支持提交付款确认", LangEn: "Current order status does not support payment submission", LangMn: "Одоогийн захиалгын төлөв төлбөр илгээхийг дэмжихгүй"},
	"tryonRechargeOrderPayMethodNotQrcode":           {LangZh: "仅扫码支付订单可提交付款确认", LangEn: "Only QR-code payment orders can submit payment confirmation", LangMn: "Зөвхөн QR төлбөртэй захиалга төлбөрийн баталгаажуулалт илгээж болно"},
	"tryonRechargeOrderNotEditable":                  {LangZh: "订单不存在或当前状态不支持修改", LangEn: "Order not found or current status cannot be edited", LangMn: "Захиалга олдсонгүй эсвэл одоогийн төлөвт өөрчлөх боломжгүй"},
	"tryonRechargeOrderNotCancelable":                {LangZh: "订单不存在或当前状态不支持取消", LangEn: "Order not found or current status cannot be cancelled", LangMn: "Захиалга олдсонгүй эсвэл одоогийн төлөвт цуцлах боломжгүй"},
	"tryonRechargeOrderStateInvalidForConfirm":       {LangZh: "当前订单状态不支持确认支付", LangEn: "Current order status does not support payment confirmation", LangMn: "Одоогийн захиалгын төлөв төлбөр баталгаажуулахыг дэмжихгүй"},
	"tryonTaskNotFound":                              {LangZh: "试衣任务不存在", LangEn: "Try-on task not found", LangMn: "Өмсгөл туршилтын даалгавар олдсонгүй"},
	"tryonRefundFailedContactAdmin":                  {LangZh: "试衣失败且退币异常，请联系管理员", LangEn: "Try-on failed and refund failed, please contact administrator", LangMn: "Өмсгөл туршилт амжилтгүй болж, буцаалт ч амжилтгүй болсон тул админтай холбогдоно уу"},
	"tryonRefinerRefundFailedContactAdmin":           {LangZh: "试衣精修失败且退币异常，请联系管理员", LangEn: "Refiner failed and refund failed, please contact administrator", LangMn: "Зураг сайжруулалт амжилтгүй болж, буцаалт ч амжилтгүй болсон тул админтай холбогдоно уу"},
	"tryonBeautifyRefundFailedContactAdmin":          {LangZh: "智能美肤失败且退币异常，请联系管理员", LangEn: "Beautify failed and refund failed, please contact administrator", LangMn: "Ухаалаг сайжруулалт амжилтгүй болж, буцаалт ч амжилтгүй болсон тул админтай холбогдоно уу"},
	"tryonTaskStatusConflictUpdateAsync":             {LangZh: "任务状态已变更，无法更新异步任务", LangEn: "Task status changed and async task cannot be updated", LangMn: "Даалгаврын төлөв өөрчлөгдсөн тул асинк даалгаврыг шинэчлэх боломжгүй"},
	"tryonTaskStatusConflictUpdateRefiner":           {LangZh: "任务状态已变更，无法更新精修任务", LangEn: "Task status changed and refiner task cannot be updated", LangMn: "Даалгаврын төлөв өөрчлөгдсөн тул сайжруулалтын даалгаврыг шинэчлэх боломжгүй"},
	"tryonTaskStatusConflictComplete":                {LangZh: "任务状态已变更，无法完成", LangEn: "Task status changed and cannot be completed", LangMn: "Даалгаврын төлөв өөрчлөгдсөн тул дуусгах боломжгүй"},
	"tryonTaskStatusConflictRefund":                  {LangZh: "任务状态已变更，无法执行退币", LangEn: "Task status changed and refund cannot be executed", LangMn: "Даалгаврын төлөв өөрчлөгдсөн тул буцаалт гүйцэтгэх боломжгүй"},
	"tryonBeautifyStatusConflict":                    {LangZh: "智能美肤状态已变更，请刷新后重试", LangEn: "Beautify status changed, please refresh and retry", LangMn: "Сайжруулалтын төлөв өөрчлөгдсөн тул шинэчлээд дахин оролдоно уу"},
	"tryonTimeRangeInvalid":                          {LangZh: "开始时间不能晚于结束时间", LangEn: "Start time cannot be later than end time", LangMn: "Эхлэх хугацаа дуусах хугацаанаас хойш байж болохгүй"},
	"invalidUserID":                                  {LangZh: "用户ID无效", LangEn: "Invalid user ID", LangMn: "Хэрэглэгчийн ID буруу"},
	"tryonModelConfigInvalid":                        {LangZh: "试衣模型配置格式错误", LangEn: "Try-on model configuration format is invalid", LangMn: "Өмсгөл турших загварын тохиргооны формат буруу"},
	"tryonModelUnavailable":                          {LangZh: "所选模型不可用或未启用", LangEn: "Selected model is unavailable or disabled", LangMn: "Сонгосон загвар боломжгүй эсвэл идэвхгүй"},
	"tryonModelNameRequired":                         {LangZh: "模特名称不能为空", LangEn: "Model name is required", LangMn: "Загварын нэр хоосон байж болохгүй"},
	"tryonModelNotFoundOrNoPermission":               {LangZh: "模特不存在或无权限", LangEn: "Model not found or no permission", LangMn: "Загвар олдсонгүй эсвэл эрх хүрэлцэхгүй"},
	"tryonClothNameRequired":                         {LangZh: "衣橱名称不能为空", LangEn: "Wardrobe item name is required", LangMn: "Шкафын нэр хоосон байж болохгүй"},
	"tryonClothCategoryInvalid":                      {LangZh: "衣橱分类不正确", LangEn: "Invalid wardrobe category", LangMn: "Шкафын ангилал буруу"},
	"tryonClothNotFoundOrNoPermission":               {LangZh: "衣橱记录不存在或无权限", LangEn: "Wardrobe item not found or no permission", LangMn: "Шкафын бичлэг олдсонгүй эсвэл эрх хүрэлцэхгүй"},
	"tryonRefinerProcessFail":                        {LangZh: "图片精修处理失败", LangEn: "Image refiner processing failed", LangMn: "Зураг сайжруулалтын боловсруулалт амжилтгүй"},
	"tryonRefinerStatusUnknown":                      {LangZh: "图片精修返回未知状态", LangEn: "Image refiner returned unknown status", LangMn: "Зураг сайжруулагч тодорхойгүй төлөв буцаалаа"},
	"tryonModelProcessFail":                          {LangZh: "试衣模型处理失败", LangEn: "Try-on model processing failed", LangMn: "Өмсгөл турших загварын боловсруулалт амжилтгүй"},
	"tryonModelStatusUnknown":                        {LangZh: "试衣模型返回未知状态", LangEn: "Try-on model returned unknown status", LangMn: "Өмсгөл турших загвар тодорхойгүй төлөв буцаалаа"},
	"tryonModelServiceNotConfigured":                 {LangZh: "试衣模型服务未配置", LangEn: "Try-on model service is not configured", LangMn: "Өмсгөл турших загварын үйлчилгээ тохируулагдаагүй"},
	"tryonModelResponseEmpty":                        {LangZh: "试衣模型返回为空", LangEn: "Try-on model response is empty", LangMn: "Өмсгөл турших загвар хоосон хариу өглөө"},
	"tryonModelResultEmpty":                          {LangZh: "试衣模型返回结果为空", LangEn: "Try-on model result is empty", LangMn: "Өмсгөл турших загварын үр дүн хоосон байна"},
	"modelImageRequired":                             {LangZh: "模特图不能为空", LangEn: "Model image is required", LangMn: "Загварын зураг хоосон байж болохгүй"},
	"clothesImageRequired":                           {LangZh: "服饰图不能为空", LangEn: "Garment image is required", LangMn: "Хувцасны зураг хоосон байж болохгүй"},
	"gradioModelURLNotConfigured":                    {LangZh: "Gradio试衣模型地址未配置", LangEn: "Gradio try-on model URL is not configured", LangMn: "Gradio өмсгөл турших загварын хаяг тохируулагдаагүй"},
	"gradioModelURLInvalid":                          {LangZh: "Gradio试衣模型地址格式错误", LangEn: "Gradio try-on model URL format is invalid", LangMn: "Gradio өмсгөл турших загварын URL формат буруу"},
	"aliyunTryonApiKeyMissing":                       {LangZh: "阿里试衣API Key未配置", LangEn: "Aliyun try-on API key is not configured", LangMn: "Aliyun өмсгөл турших API key тохируулагдаагүй"},
	"aliyunTryonResponseEmpty":                       {LangZh: "阿里试衣返回为空", LangEn: "Aliyun try-on response is empty", LangMn: "Aliyun өмсгөл туршилтын хариу хоосон байна"},
	"aliyunTryonTaskIDMissing":                       {LangZh: "阿里试衣未返回task_id", LangEn: "Aliyun try-on did not return task_id", LangMn: "Aliyun өмсгөл туршилт task_id буцаасангүй"},
	"refinerInputImageRequired":                      {LangZh: "精修输入图不能为空", LangEn: "Refiner input image is required", LangMn: "Сайжруулалтын оролтын зураг хоосон байж болохгүй"},
	"aliyunRefinerApiKeyMissing":                     {LangZh: "阿里图片精修API Key未配置", LangEn: "Aliyun image refiner API key is not configured", LangMn: "Aliyun зураг сайжруулагч API key тохируулагдаагүй"},
	"aliyunRefinerResponseEmpty":                     {LangZh: "阿里图片精修返回为空", LangEn: "Aliyun image refiner response is empty", LangMn: "Aliyun зураг сайжруулагч хоосон хариу өглөө"},
	"aliyunRefinerTaskIDMissing":                     {LangZh: "阿里图片精修未返回task_id", LangEn: "Aliyun image refiner did not return task_id", LangMn: "Aliyun зураг сайжруулагч task_id буцаасангүй"},
	"aliyunBeautifyAccessKeyMissing":                 {LangZh: "阿里智能美肤AccessKey未配置", LangEn: "Aliyun beautify AccessKey is not configured", LangMn: "Aliyun сайжруулалтын AccessKey тохируулагдаагүй"},
	"aliyunParsingApiKeyMissing":                     {LangZh: "阿里图片分割API Key未配置", LangEn: "Aliyun image parsing API key is not configured", LangMn: "Aliyun зураг ялгалтын API key тохируулагдаагүй"},
	"aliyunParsingResponseEmpty":                     {LangZh: "阿里图片分割返回为空", LangEn: "Aliyun image parsing response is empty", LangMn: "Aliyun зураг ялгалтын хариу хоосон байна"},
	"aliyunTaskIDRequired":                           {LangZh: "task_id不能为空", LangEn: "task_id is required", LangMn: "task_id хоосон байж болохгүй"},
	"aliyunTaskQueryResponseEmpty":                   {LangZh: "查询阿里试衣任务返回为空", LangEn: "Aliyun try-on task query response is empty", LangMn: "Aliyun өмсгөл туршилтын даалгаврын асуулгын хариу хоосон байна"},
	"mediaURLRequired":                               {LangZh: "图片地址不能为空", LangEn: "Image URL is required", LangMn: "Зургийн URL хоосон байж болохгүй"},
	"mediaURLDataURINotAllowed":                      {LangZh: "图片地址不能是 data URI，请先上传并使用公网 URL", LangEn: "Image URL cannot be data URI; upload first and use a public URL", LangMn: "Зургийн URL нь data URI байж болохгүй; эхлээд байршуулж, олон нийтэд нээлттэй URL ашиглана уу"},
	"mediaURLPublicURLRequired":                      {LangZh: "图片地址不是公网 URL，请在系统参数配置 tryon_media_public_base_url 后重试", LangEn: "Image URL is not public; configure tryon_media_public_base_url and retry", LangMn: "Зургийн URL олон нийтэд нээлттэй биш; tryon_media_public_base_url тохируулаад дахин оролдоно уу"},
	"mediaURLInvalid":                                {LangZh: "图片地址格式错误", LangEn: "Invalid image URL format", LangMn: "Зургийн URL формат буруу"},
	"mediaURLPrivateHostNotAccessible":               {LangZh: "当前图片地址为 localhost/内网地址，模型服务无法访问；请配置 tryon_media_public_base_url 为公网域名", LangEn: "Image URL points to localhost/private host; configure tryon_media_public_base_url to a public domain", LangMn: "Зургийн URL localhost/дотоод хост руу зааж байна; tryon_media_public_base_url-ийг олон нийтийн домэйнээр тохируулна уу"},
	"publicBaseURLRequired":                          {LangZh: "公网地址前缀不能为空", LangEn: "Public base URL prefix is required", LangMn: "Олон нийтийн суурь URL префикс хоосон байж болохгүй"},
	"publicBaseURLInvalid":                           {LangZh: "tryon_media_public_base_url 配置格式错误，请填写如 https://back.example.com", LangEn: "Invalid tryon_media_public_base_url format, e.g. https://back.example.com", LangMn: "tryon_media_public_base_url формат буруу, жишээ нь https://back.example.com"},
	"tryonResultURLRequired":                         {LangZh: "结果图片地址为空", LangEn: "Result image URL is empty", LangMn: "Үр дүнгийн зургийн URL хоосон байна"},
	"tryonResultDownloadFailed":                      {LangZh: "下载结果图失败", LangEn: "Failed to download result image", LangMn: "Үр дүнгийн зургийг татаж чадсангүй"},
	"tryonResultContentEmpty":                        {LangZh: "结果图片内容为空", LangEn: "Result image content is empty", LangMn: "Үр дүнгийн зургийн агуулга хоосон байна"},
	"tryonResultTooLarge":                            {LangZh: "结果图片超过大小限制", LangEn: "Result image exceeds size limit", LangMn: "Үр дүнгийн зураг хэмжээний хязгаараас хэтэрсэн"},
	"tryonBeautifyEndpointInvalid":                   {LangZh: "智能美肤服务地址格式错误", LangEn: "Beautify service endpoint is invalid", LangMn: "Сайжруулалтын үйлчилгээний хаягийн формат буруу"},
	"tryonBeautifyRequestFailed":                     {LangZh: "智能美肤请求失败", LangEn: "Beautify request failed", LangMn: "Сайжруулалтын хүсэлт амжилтгүй"},
	"tryonBeautifyResponseInvalid":                   {LangZh: "智能美肤响应格式错误", LangEn: "Beautify response format is invalid", LangMn: "Сайжруулалтын хариуны формат буруу"},
	"tryonBeautifyResultEmpty":                       {LangZh: "智能美肤未返回结果图", LangEn: "Beautify did not return result image", LangMn: "Сайжруулалт үр дүнгийн зураг буцаасангүй"},
	"reason_tryonBeautifyDeduct":                     {LangZh: "智能美肤扣费", LangEn: "Beautify deduction", LangMn: "Сайжруулалтын суутгал"},
	"reason_tryonBeautifyRefund":                     {LangZh: "智能美肤退款", LangEn: "Beautify refund", LangMn: "Сайжруулалтын буцаалт"},
	"tryonGuestInitTxMissing":                        {LangZh: "试衣赠币事务上下文缺失", LangEn: "Try-on bonus transaction context is missing", LangMn: "Өмсгөл бонусын гүйлгээний контекст дутуу байна"},
	"tryonTaskFailed":                                {LangZh: "试衣任务失败", LangEn: "Try-on task failed", LangMn: "Өмсгөл туршилтын даалгавар амжилтгүй"},
	"gradioTryonRequestFailed":                       {LangZh: "Gradio试衣请求失败", LangEn: "Gradio try-on request failed", LangMn: "Gradio өмсгөл туршилтын хүсэлт амжилтгүй"},
	"gradioTryonCreateRespParseFail":                 {LangZh: "Gradio试衣创建响应解析失败", LangEn: "Failed to parse Gradio try-on create response", LangMn: "Gradio өмсгөл туршилтын үүсгэх хариуг задлаж чадсангүй"},
	"gradioEventIDMissing":                           {LangZh: "Gradio试衣未返回event_id", LangEn: "Gradio try-on did not return event_id", LangMn: "Gradio өмсгөл туршилт event_id буцаасангүй"},
	"gradioResultEmpty":                              {LangZh: "Gradio试衣返回结果为空", LangEn: "Gradio try-on result is empty", LangMn: "Gradio өмсгөл туршилтын үр дүн хоосон байна"},
	"gradioTryonQueryFailed":                         {LangZh: "Gradio试衣结果查询失败", LangEn: "Failed to query Gradio try-on result", LangMn: "Gradio өмсгөл туршилтын үр дүнг асууж чадсангүй"},
	"gradioTryonTaskFailed":                          {LangZh: "Gradio试衣任务失败", LangEn: "Gradio try-on task failed", LangMn: "Gradio өмсгөл туршилтын даалгавар амжилтгүй"},
	"gradioTryonErrorWithoutDetails":                 {LangZh: "Gradio任务失败，但未返回可解析的错误详情", LangEn: "Gradio task failed without parsable error details", LangMn: "Gradio даалгавар амжилтгүй болсон ч задлах боломжтой алдааны дэлгэрэнгүй ирсэнгүй"},
	"aliyunParsingNoUsableGarmentArea":               {LangZh: "图片分割未识别到可用服饰区域", LangEn: "Image parsing did not detect usable garment area", LangMn: "Зургийн ялгалт ашиглах боломжтой хувцасны хэсгийг илрүүлсэнгүй"},
	"aliyunTaskSuccessNoImage":                       {LangZh: "试衣任务成功但未返回图片", LangEn: "Task succeeded but no image returned", LangMn: "Даалгавар амжилттай болсон ч зураг буцаасангүй"},
	"modelAliyunRetouchNeedStandardOSS":              {LangZh: "智能美肤图片地址不符合阿里云要求：请使用上海区标准 OSS 域名（不支持 CDN/自定义域名）", LangEn: "Beautify image URL is not accepted by Aliyun: use a standard OSS domain in cn-shanghai (CDN/custom domains are not supported)", LangMn: "Aliyun сайжруулалтын зургийн URL шаардлага хангахгүй: cn-shanghai стандарт OSS домэйн ашиглана уу (CDN/захиалгат домэйн дэмжихгүй)"},
	"modelCannotDownloadResource":                    {LangZh: "模型服务无法下载图片资源，请确认 sourceImage/templateImage 为公网可访问 URL（不要使用 localhost 或内网地址）", LangEn: "Model service cannot download image resources, ensure sourceImage/templateImage are publicly accessible URLs", LangMn: "Загварын үйлчилгээ зургийн нөөцийг татаж чадсангүй, sourceImage/templateImage-ийг олон нийтэд нээлттэй URL болгож шалгана уу"},
	"modelDownloadRefused":                           {LangZh: "模型服务无法下载图片资源（URL被拒绝），请更换可直连的图片域名后重试", LangEn: "Model service cannot download image resources (URL refused), switch to a directly accessible image domain and retry", LangMn: "Загварын үйлчилгээ зургийн нөөцийг татаж чадсангүй (URL татгалзсан), шууд хандах боломжтой домэйнээр солиод дахин оролдоно уу"},
	"modelImageURLInvalidOrUnreachable":              {LangZh: "图片地址无效或模型服务不可达，请更换可直连的图片域名后重试", LangEn: "Image URL is invalid or model service is unreachable, switch to a directly accessible image domain and retry", LangMn: "Зургийн URL буруу эсвэл загварын үйлчилгээ хүрэх боломжгүй байна, шууд хандах домэйнээр солиод дахин оролдоно уу"},
	"modelImageRiskCheckFailed":                      {LangZh: "图片风控或资源检测失败，请更换图片后重试", LangEn: "Image risk control or resource inspection failed, please change image and retry", LangMn: "Зургийн эрсдэлийн хяналт эсвэл нөөц шалгалт амжилтгүй, зургийг сольж дахин оролдоно уу"},

	// ========== 认证/登录 ==========
	"loginSuccess":           {LangZh: "登录成功", LangEn: "Login successful", LangMn: "Амжилттай нэвтэрлээ"},
	"loginFail":              {LangZh: "用户名或密码错误", LangEn: "Invalid username or password", LangMn: "Хэрэглэгчийн нэр эсвэл нууц үг буруу"},
	"phoneLoginFail":         {LangZh: "手机号不存在或密码错误", LangEn: "Phone number not found or wrong password", LangMn: "Утасны дугаар эсвэл нууц үг буруу"},
	"captchaError":           {LangZh: "验证码错误", LangEn: "Invalid captcha", LangMn: "Баталгаажуулах код буруу"},
	"codeEmpty":              {LangZh: "code不能为空", LangEn: "Code is required", LangMn: "Код хоосон байж болохгүй"},
	"openidFail":             {LangZh: "获取openid失败", LangEn: "Failed to get OpenID", LangMn: "OpenID авч чадсангүй"},
	"tokenFail":              {LangZh: "获取token失败", LangEn: "Failed to get token", LangMn: "Token авч чадсангүй"},
	"loginStatusFail":        {LangZh: "设置登录状态失败", LangEn: "Failed to set login status", LangMn: "Нэвтрэлтийн төлөв тохируулж чадсангүй"},
	"jwtBlacklistFail":       {LangZh: "jwt作废失败", LangEn: "Failed to revoke JWT", LangMn: "JWT хүчингүй болгож чадсангүй"},
	"passwordMismatch":       {LangZh: "两次输入的密码不一致", LangEn: "Passwords do not match", LangMn: "Нууц үг таарахгүй байна"},
	"usernameOrPhoneExists":  {LangZh: "用户名或手机号码已存在", LangEn: "Username or phone already exists", LangMn: "Хэрэглэгчийн нэр эсвэл утасны дугаар аль хэдийн бүртгэгдсэн байна"},
	"phoneAlreadyRegistered": {LangZh: "该手机号已注册", LangEn: "Phone number already registered", LangMn: "Энэ утасны дугаар аль хэдийн бүртгэгдсэн байна"},
	"signInSuccess":          {LangZh: "签到成功", LangEn: "Sign-in successful", LangMn: "Амжилттай тэмдэглэлээ"},
	"signInAlreadyToday":     {LangZh: "今日已签到", LangEn: "Already signed in today", LangMn: "Өнөөдөр аль хэдийн тэмдэглэсэн байна"},
	"cannotModify":           {LangZh: "无法修改", LangEn: "Cannot modify", LangMn: "Өөрчлөх боломжгүй"},
	"seriesNotFound":         {LangZh: "剧集不存在", LangEn: "Series not found", LangMn: "Цуврал олдсонгүй"},
	"loginLocked":            {LangZh: "登录失败次数过多，请稍后再试", LangEn: "Too many failed attempts, please try again later", LangMn: "Нэвтрэлт хэт олон удаа амжилтгүй, дараа дахин оролдоно уу"},
	"accountBanned":          {LangZh: "账号已被封禁", LangEn: "Account has been banned", LangMn: "Бүртгэл хориглогдсон"},
	"registerIPLimit":        {LangZh: "该IP注册次数已达上限", LangEn: "Registration limit reached for this IP", LangMn: "Энэ IP-ээс бүртгэх хязгаарт хүрсэн"},
	"areaCodeInvalid":        {LangZh: "区号无效或未启用", LangEn: "Invalid or disabled area code", LangMn: "Бүсийн код буруу эсвэл идэвхгүй"},
	"phoneFormatError":       {LangZh: "手机号格式不正确", LangEn: "Invalid phone number format", LangMn: "Утасны дугаарын формат буруу"},
	"changeSuccess":          {LangZh: "修改成功", LangEn: "Changed successfully", LangMn: "Амжилттай өөрчлөлөө"},
	"phoneRequired":          {LangZh: "手机号不能为空", LangEn: "Phone number is required", LangMn: "Утасны дугаар хоосон байж болохгүй"},
	"passwordError":          {LangZh: "密码错误", LangEn: "Incorrect password", LangMn: "Нууц үг буруу"},
	"userNotExist":           {LangZh: "用户不存在", LangEn: "User not found", LangMn: "Хэрэглэгч олдсонгүй"},

	// ========== JWT 中间件 ==========
	"notLogin":     {LangZh: "未登录或非法访问，请登录", LangEn: "Please login to continue", LangMn: "Нэвтэрнэ үү"},
	"tokenInvalid": {LangZh: "您的帐户异地登陆或令牌失效", LangEn: "Your session is invalid, please login again", LangMn: "Таны токен хүчингүй болсон"},
	"tokenExpired": {LangZh: "登录已过期，请重新登录", LangEn: "Session expired, please login again", LangMn: "Нэвтрэлт дууссан, дахин нэвтэрнэ үү"},
}

// contextKey 用于 gin.Context 存储语言
const langKey = "accept-language"

// GetLang 从 gin.Context 获取当前语言
func GetLang(c *gin.Context) string {
	if lang, exists := c.Get(langKey); exists {
		if s, ok := lang.(string); ok && s != "" {
			return s
		}
	}
	return LangZh
}

// SetLang 将语言存入 gin.Context
func SetLang(c *gin.Context, lang string) {
	c.Set(langKey, lang)
}

// T 根据消息 key 和当前请求语言返回翻译文本
func T(c *gin.Context, key string) string {
	lang := GetLang(c)
	if m, ok := messages[key]; ok {
		if s, ok := m[lang]; ok {
			return s
		}
		// 降级中文
		if s, ok := m[LangZh]; ok {
			return s
		}
	}
	return key
}

// TWithSuffix 返回翻译文本 + 后缀（如错误详情）
func TWithSuffix(c *gin.Context, key string, suffix string) string {
	msg := T(c, key)
	if suffix != "" {
		return msg + ": " + suffix
	}
	return msg
}
