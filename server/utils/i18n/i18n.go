package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// 语言代码常量
const (
	LangZh = "zh"
	LangEn = "en"
	LangMn = "mn"
	LangVi = "vi"
	LangAr = "ar"
	LangJa = "ja"
	LangKo = "ko"
	LangMs = "ms"
)

// 消息字典：key -> { lang -> text }
var messages = map[string]map[string]string{
	// ========== 通用 CRUD ==========
	"ok":                                             {LangZh: "操作成功", LangEn: "Success", LangMn: "Амжилттай", LangJa: "操作に成功しました", LangKo: "작업에 성공했습니다", LangMs: "Operasi berjaya"},
	"fail":                                           {LangZh: "操作失败", LangEn: "Operation failed", LangMn: "Амжилтгүй", LangJa: "操作に失敗しました", LangKo: "작업에 실패했습니다", LangMs: "Operasi gagal"},
	"success":                                        {LangZh: "成功", LangEn: "Success", LangMn: "Амжилттай", LangJa: "成功", LangKo: "성공", LangMs: "Berjaya"},
	"createSuccess":                                  {LangZh: "创建成功", LangEn: "Created successfully", LangMn: "Амжилттай үүсгэлээ", LangJa: "作成に成功しました", LangKo: "생성에 성공했습니다", LangMs: "Berjaya dicipta"},
	"createFail":                                     {LangZh: "创建失败", LangEn: "Create failed", LangMn: "Үүсгэж чадсангүй", LangJa: "作成に失敗しました", LangKo: "생성에 실패했습니다", LangMs: "Gagal dicipta"},
	"deleteSuccess":                                  {LangZh: "删除成功", LangEn: "Deleted successfully", LangMn: "Амжилттай устгалаа", LangJa: "削除に成功しました", LangKo: "삭제에 성공했습니다", LangMs: "Berjaya dipadam"},
	"deleteFail":                                     {LangZh: "删除失败", LangEn: "Delete failed", LangMn: "Устгаж чадсангүй", LangJa: "削除に失敗しました", LangKo: "삭제에 실패했습니다", LangMs: "Gagal dipadam"},
	"batchDeleteSuccess":                             {LangZh: "批量删除成功", LangEn: "Batch deleted", LangMn: "Бөөнөөр устгалаа", LangJa: "一括削除に成功しました", LangKo: "일괄 삭제에 성공했습니다", LangMs: "Berjaya dipadam secara pukal"},
	"batchDeleteFail":                                {LangZh: "批量删除失败", LangEn: "Batch delete failed", LangMn: "Бөөнөөр устгаж чадсангүй", LangJa: "一括削除に失敗しました", LangKo: "일괄 삭제에 실패했습니다", LangMs: "Gagal padam pukal"},
	"updateSuccess":                                  {LangZh: "更新成功", LangEn: "Updated successfully", LangMn: "Амжилттай шинэчлэлээ", LangJa: "更新に成功しました", LangKo: "업데이트에 성공했습니다", LangMs: "Berjaya dikemas kini"},
	"updateFail":                                     {LangZh: "更新失败", LangEn: "Update failed", LangMn: "Шинэчлэж чадсангүй", LangJa: "更新に失敗しました", LangKo: "업데이트에 실패했습니다", LangMs: "Gagal dikemas kini"},
	"querySuccess":                                   {LangZh: "查询成功", LangEn: "Query successful", LangMn: "Хайлт амжилттай", LangJa: "検索に成功しました", LangKo: "조회에 성공했습니다", LangMs: "Pertanyaan berjaya"},
	"queryFail":                                      {LangZh: "查询失败", LangEn: "Query failed", LangMn: "Хайлт амжилтгүй", LangJa: "検索に失敗しました", LangKo: "조회에 실패했습니다", LangMs: "Pertanyaan gagal"},
	"getSuccess":                                     {LangZh: "获取成功", LangEn: "Retrieved successfully", LangMn: "Амжилттай авлаа", LangJa: "取得に成功しました", LangKo: "가져오기에 성공했습니다", LangMs: "Berjaya diperoleh"},
	"getFail":                                        {LangZh: "获取失败", LangEn: "Retrieve failed", LangMn: "Авч чадсангүй", LangJa: "取得に失敗しました", LangKo: "가져오기에 실패했습니다", LangMs: "Gagal diperoleh"},
	"setSuccess":                                     {LangZh: "设置成功", LangEn: "Set successfully", LangMn: "Амжилттай тохируулалаа", LangJa: "設定に成功しました", LangKo: "설정에 성공했습니다", LangMs: "Tetapan berjaya"},
	"setFail":                                        {LangZh: "设置失败", LangEn: "Set failed", LangMn: "Тохируулж чадсангүй", LangJa: "設定に失敗しました", LangKo: "설정에 실패했습니다", LangMs: "Tetapan gagal"},
	"invalidParams":                                  {LangZh: "参数错误", LangEn: "Invalid parameters", LangMn: "Параметр буруу", LangJa: "パラメータが不正です", LangKo: "잘못된 매개변수입니다", LangMs: "Parameter tidak sah"},
	"importSuccess":                                  {LangZh: "导入成功", LangEn: "Imported successfully", LangMn: "Амжилттай импортоллоо", LangJa: "インポートに成功しました", LangKo: "가져오기에 성공했습니다", LangMs: "Berjaya diimport"},
	"importFail":                                     {LangZh: "导入失败", LangEn: "Import failed", LangMn: "Импорт амжилтгүй", LangJa: "インポートに失敗しました", LangKo: "가져오기에 실패했습니다", LangMs: "Import gagal"},
	"fileGetFail":                                    {LangZh: "文件获取失败", LangEn: "Failed to get file", LangMn: "Файл авч чадсангүй", LangJa: "ファイルの取得に失敗しました", LangKo: "파일을 가져오지 못했습니다", LangMs: "Gagal mendapatkan fail"},
	"exportTemplateIDRequired":                       {LangZh: "模板ID不能为空", LangEn: "Template ID is required", LangMn: "Загварын ID хоосон байж болохгүй", LangJa: "テンプレートIDは必須です", LangKo: "템플릿 ID는 필수입니다", LangMs: "ID templat diperlukan"},
	"exportTokenRequired":                            {LangZh: "导出token不能为空", LangEn: "Export token is required", LangMn: "Экспортын token хоосон байж болохгүй", LangJa: "エクスポートトークンは必須です", LangKo: "내보내기 토큰은 필수입니다", LangMs: "Token eksport diperlukan"},
	"exportTokenQueryNotAllowed":                     {LangZh: "当前环境禁止通过URL传递导出token，请使用请求头", LangEn: "Export token in URL is disabled in current environment, please use request header", LangMn: "Одоогийн орчинд URL-ээр экспортын token дамжуулахыг хориглоно, header ашиглана уу", LangJa: "現在の環境ではURLでのエクスポートトークン送信は無効です。ヘッダーを使用してください", LangKo: "현재 환경에서는 URL로 내보내기 토큰 전달이 금지되어 있습니다. 헤더를 사용하세요", LangMs: "Token eksport melalui URL dinyahaktifkan dalam persekitaran semasa, sila gunakan pengepala permintaan"},
	"exportTokenInvalidOrExpired":                    {LangZh: "导出token无效或已过期", LangEn: "Export token is invalid or expired", LangMn: "Экспортын token буруу эсвэл хугацаа дууссан", LangJa: "エクスポートトークンが無効か期限切れです", LangKo: "내보내기 토큰이 유효하지 않거나 만료되었습니다", LangMs: "Token eksport tidak sah atau telah tamat tempoh"},
	"exportParamsInvalid":                            {LangZh: "解析导出参数失败", LangEn: "Failed to parse export parameters", LangMn: "Экспортын параметр задлахад алдаа гарлаа", LangJa: "エクスポートパラメータの解析に失敗しました", LangKo: "내보내기 파라미터를 파싱하지 못했습니다", LangMs: "Gagal menghuraikan parameter eksport"},
	"exportTokenTypeInvalid":                         {LangZh: "token类型错误", LangEn: "Invalid token type", LangMn: "Token төрөл буруу", LangJa: "トークンタイプが不正です", LangKo: "토큰 유형이 올바르지 않습니다", LangMs: "Jenis token tidak sah"},
	"exportTemplateNil":                              {LangZh: "导出模板不能为空", LangEn: "Export template cannot be empty", LangMn: "Экспортын загвар хоосон байж болохгүй", LangJa: "エクスポートテンプレートは必須です", LangKo: "내보내기 템플릿은 비워둘 수 없습니다", LangMs: "Templat eksport tidak boleh kosong"},
	"exportTemplateTableNameInvalid":                 {LangZh: "表名不合法", LangEn: "Invalid table name", LangMn: "Хүснэгтийн нэр буруу", LangJa: "テーブル名が不正です", LangKo: "테이블 이름이 유효하지 않습니다", LangMs: "Nama jadual tidak sah"},
	"exportTemplateLimitNegative":                    {LangZh: "limit 不能小于 0", LangEn: "Limit cannot be less than 0", LangMn: "Limit 0-ээс бага байж болохгүй", LangJa: "limit は 0 未満にできません", LangKo: "limit은 0보다 작을 수 없습니다", LangMs: "Limit tidak boleh kurang daripada 0"},
	"exportTemplateLimitExceeded":                    {LangZh: "limit 超过安全上限", LangEn: "Limit exceeds security threshold", LangMn: "Limit аюулгүй дээд хязгаараас хэтэрсэн", LangJa: "limit が安全上限を超えています", LangKo: "limit이 보안 상한을 초과했습니다", LangMs: "Limit melebihi ambang keselamatan"},
	"exportTemplateOrderInvalid":                     {LangZh: "排序字段不合法", LangEn: "Invalid order field", LangMn: "Эрэмбэлэх талбар буруу", LangJa: "並び替えフィールドが不正です", LangKo: "정렬 필드가 유효하지 않습니다", LangMs: "Medan susunan tidak sah"},
	"exportTemplateOrderDirectionInvalid":            {LangZh: "排序方向不合法", LangEn: "Invalid order direction", LangMn: "Эрэмбэлэх чиглэл буруу", LangJa: "並び替え方向が不正です", LangKo: "정렬 방향이 유효하지 않습니다", LangMs: "Arah susunan tidak sah"},
	"exportTemplateOrderFieldInvalid":                {LangZh: "排序字段不在可用字段中", LangEn: "Order field is not available", LangMn: "Эрэмбэлэх талбар боломжит талбаруудад байхгүй", LangJa: "並び替えフィールドが利用可能な項目にありません", LangKo: "정렬 필드가 사용 가능한 필드에 없습니다", LangMs: "Medan susunan tiada dalam medan yang tersedia"},
	"exportTemplateConditionFromInvalid":             {LangZh: "条件参数键不合法", LangEn: "Invalid condition parameter key", LangMn: "Нөхцлийн параметрийн түлхүүр буруу", LangJa: "条件パラメータキーが不正です", LangKo: "조건 매개변수 키가 유효하지 않습니다", LangMs: "Kunci parameter syarat tidak sah"},
	"exportTemplateConditionColumnInvalid":           {LangZh: "条件字段不合法", LangEn: "Invalid condition field", LangMn: "Нөхцлийн талбар буруу", LangJa: "条件フィールドが不正です", LangKo: "조건 필드가 유효하지 않습니다", LangMs: "Medan syarat tidak sah"},
	"exportTemplateConditionOperatorInvalid":         {LangZh: "条件操作符不合法", LangEn: "Invalid condition operator", LangMn: "Нөхцлийн оператор буруу", LangJa: "条件演算子が不正です", LangKo: "조건 연산자가 유효하지 않습니다", LangMs: "Operator syarat tidak sah"},
	"exportTemplateJoinTypeInvalid":                  {LangZh: "JOIN 类型不合法", LangEn: "Invalid JOIN type", LangMn: "JOIN төрөл буруу", LangJa: "JOIN タイプが不正です", LangKo: "JOIN 유형이 유효하지 않습니다", LangMs: "Jenis JOIN tidak sah"},
	"exportTemplateJoinTableInvalid":                 {LangZh: "JOIN 表名不合法", LangEn: "Invalid JOIN table name", LangMn: "JOIN хүснэгтийн нэр буруу", LangJa: "JOIN テーブル名が不正です", LangKo: "JOIN 테이블 이름이 유효하지 않습니다", LangMs: "Nama jadual JOIN tidak sah"},
	"exportTemplateJoinOnInvalid":                    {LangZh: "JOIN ON 条件不合法", LangEn: "Invalid JOIN ON condition", LangMn: "JOIN ON нөхцөл буруу", LangJa: "JOIN ON 条件が不正です", LangKo: "JOIN ON 조건이 유효하지 않습니다", LangMs: "Syarat JOIN ON tidak sah"},
	"exportTemplateRawSQLDisabled":                   {LangZh: "当前环境禁止自定义SQL", LangEn: "Custom SQL is disabled in current environment", LangMn: "Одоогийн орчинд custom SQL хориглосон", LangJa: "現在の環境ではカスタムSQLが無効です", LangKo: "현재 환경에서는 커스텀 SQL이 비활성화되어 있습니다", LangMs: "SQL tersuai dinyahaktifkan dalam persekitaran semasa"},
	"exportTemplateRawSQLSelectOnly":                 {LangZh: "仅允许 SELECT 语句", LangEn: "Only SELECT statements are allowed", LangMn: "Зөвхөн SELECT өгүүлбэр зөвшөөрнө", LangJa: "SELECT 文のみ許可されています", LangKo: "SELECT 문만 허용됩니다", LangMs: "Hanya pernyataan SELECT dibenarkan"},
	"exportTemplateRawSQLUnsafeFragment":             {LangZh: "SQL 包含不安全片段", LangEn: "SQL contains unsafe fragment", LangMn: "SQL аюултай хэсэг агуулж байна", LangJa: "SQL に危険な断片が含まれています", LangKo: "SQL에 안전하지 않은 조각이 포함되어 있습니다", LangMs: "SQL mengandungi serpihan tidak selamat"},
	"exportTemplateRawSQLKeywordForbidden":           {LangZh: "SQL 包含不允许的关键字", LangEn: "SQL contains forbidden keyword", LangMn: "SQL зөвшөөрөөгүй түлхүүр үг агуулж байна", LangJa: "SQL に禁止キーワードが含まれています", LangKo: "SQL에 금지된 키워드가 포함되어 있습니다", LangMs: "SQL mengandungi kata kunci terlarang"},
	"exportTemplateImportSQLDisabled":                {LangZh: "当前环境禁止自定义导入SQL", LangEn: "Custom import SQL is disabled in current environment", LangMn: "Одоогийн орчинд custom import SQL хориглосон", LangJa: "現在の環境ではカスタムインポートSQLが無効です", LangKo: "현재 환경에서는 커스텀 import SQL이 비활성화되어 있습니다", LangMs: "SQL import tersuai dinyahaktifkan dalam persekitaran semasa"},
	"exportTemplateImportSQLTypeInvalid":             {LangZh: "导入SQL仅允许 INSERT INTO 或 UPDATE", LangEn: "Import SQL only allows INSERT INTO or UPDATE", LangMn: "Import SQL нь зөвхөн INSERT INTO эсвэл UPDATE зөвшөөрнө", LangJa: "インポートSQLは INSERT INTO または UPDATE のみ許可されます", LangKo: "가져오기 SQL은 INSERT INTO 또는 UPDATE만 허용됩니다", LangMs: "SQL import hanya membenarkan INSERT INTO atau UPDATE"},
	"exportTemplateImportSQLUnsafeFragment":          {LangZh: "导入SQL包含不安全片段", LangEn: "Import SQL contains unsafe fragment", LangMn: "Import SQL аюултай хэсэг агуулж байна", LangJa: "インポートSQLに危険な断片が含まれています", LangKo: "가져오기 SQL에 안전하지 않은 조각이 포함되어 있습니다", LangMs: "SQL import mengandungi serpihan tidak selamat"},
	"exportTemplateImportSQLKeywordForbidden":        {LangZh: "导入SQL包含不允许的关键字", LangEn: "Import SQL contains forbidden keyword", LangMn: "Import SQL зөвшөөрөөгүй түлхүүр үг агуулж байна", LangJa: "インポートSQLに禁止キーワードが含まれています", LangKo: "가져오기 SQL에 금지된 키워드가 포함되어 있습니다", LangMs: "SQL import mengandungi kata kunci terlarang"},
	"exportTemplateImportSQLWhereRequired":           {LangZh: "UPDATE 导入SQL必须包含 WHERE 条件", LangEn: "UPDATE import SQL must contain WHERE clause", LangMn: "UPDATE import SQL нь WHERE нөхцөлтэй байх ёстой", LangJa: "UPDATE インポートSQLには WHERE 句が必要です", LangKo: "UPDATE 가져오기 SQL에는 WHERE 절이 필요합니다", LangMs: "SQL import UPDATE mesti mengandungi klausa WHERE"},
	"exportTemplateImportSQLParamRequired":           {LangZh: "导入SQL必须使用命名参数", LangEn: "Import SQL must use named parameters", LangMn: "Import SQL нь нэрлэсэн параметр ашиглах ёстой", LangJa: "インポートSQLは名前付きパラメータを使用する必要があります", LangKo: "가져오기 SQL은 명명된 매개변수를 사용해야 합니다", LangMs: "SQL import mesti menggunakan parameter bernama"},
	"exportTemplateExcelDataNotEnough":               {LangZh: "Excel 数据不足，至少需要标题行和一行数据", LangEn: "Excel data is insufficient, at least one header row and one data row are required", LangMn: "Excel өгөгдөл дутуу, дор хаяж гарчиг мөр ба нэг өгөгдлийн мөр шаардлагатай", LangJa: "Excelデータが不足しています。少なくともヘッダー行とデータ行が1行必要です", LangKo: "Excel 데이터가 부족합니다. 최소한 헤더 행 1개와 데이터 행 1개가 필요합니다", LangMs: "Data Excel tidak mencukupi, sekurang-kurangnya satu baris tajuk dan satu baris data diperlukan"},
	"requestTooFrequent":                             {LangZh: "请求过于频繁，请稍后再试", LangEn: "Too many requests, please try again later", LangMn: "Хүсэлт хэт олон байна, дараа дахин оролдоно уу", LangJa: "リクエストが多すぎます。しばらくしてからお試しください", LangKo: "요청이 너무 많습니다. 잠시 후 다시 시도해 주세요", LangMs: "Terlalu banyak permintaan, sila cuba lagi kemudian"},
	"visitorHeartbeatTooFrequent":                    {LangZh: "心跳请求过于频繁，请稍后再试", LangEn: "Heartbeat requests are too frequent, please try again later", LangMn: "Зочны heartbeat хүсэлт хэт олон байна, дараа дахин оролдоно уу", LangJa: "ハートビートのリクエストが多すぎます。しばらくしてからお試しください", LangKo: "하트비트 요청이 너무 많습니다. 잠시 후 다시 시도해 주세요", LangMs: "Permintaan heartbeat terlalu kerap, sila cuba lagi kemudian"},
	"invalidID":                                      {LangZh: "ID参数错误", LangEn: "Invalid ID parameter", LangMn: "ID параметр буруу"},
	"invalidIDs":                                     {LangZh: "IDs参数错误", LangEn: "Invalid IDs parameter", LangMn: "IDs параметр буруу"},
	"invalidTaskNo":                                  {LangZh: "taskNo参数错误", LangEn: "Invalid task number parameter", LangMn: "taskNo параметр буруу"},
	"invalidOrder":                                   {LangZh: "订单参数错误", LangEn: "Invalid order parameter", LangMn: "Захиалгын параметр буруу"},
	"userIDRequired":                                 {LangZh: "用户ID不能为空", LangEn: "User ID is required", LangMn: "Хэрэглэгчийн ID хоосон байж болохгүй"},
	"visitorIDRequired":                              {LangZh: "visitorId不能为空", LangEn: "visitorId is required", LangMn: "visitorId хоосон байж болохгүй"},
	"dateRequired":                                   {LangZh: "date不能为空", LangEn: "date is required", LangMn: "date хоосон байж болохгүй"},
	"configGroupRequired":                            {LangZh: "configGroup不能为空", LangEn: "configGroup is required", LangMn: "configGroup хоосон байж болохгүй"},
	"configKeyRequired":                              {LangZh: "configKey不能为空", LangEn: "configKey is required", LangMn: "configKey хоосон байж болохгүй"},
	"configNotPublic":                                {LangZh: "该配置不对外开放", LangEn: "This config is not public", LangMn: "Энэ тохиргоо олон нийтэд нээлттэй биш"},
	"noPermission":                                   {LangZh: "无权限操作", LangEn: "Permission denied", LangMn: "Эрх хүрэлцэхгүй", LangJa: "権限がありません", LangKo: "권한이 없습니다", LangMs: "Tiada kebenaran"},
	"adjustSuccess":                                  {LangZh: "调整成功", LangEn: "Adjusted successfully", LangMn: "Амжилттай тохирууллаа"},
	"adjustAmountMustPositive":                       {LangZh: "调整数量必须大于0", LangEn: "Adjustment amount must be greater than 0", LangMn: "Тохируулгын хэмжээ 0-оос их байх ёстой"},
	"adjustReasonRequired":                           {LangZh: "调整原因不能为空", LangEn: "Adjustment reason is required", LangMn: "Тохируулгын шалтгаан хоосон байж болохгүй"},
	"pointChangeTypeInvalid":                         {LangZh: "增减类型必须是 increase 或 decrease", LangEn: "Change type must be increase or decrease", LangMn: "Өөрчлөлтийн төрөл нь increase эсвэл decrease байх ёстой"},
	"pointRecordRequiredFields":                      {LangZh: "用户ID、增减类型和积分变化不能为空", LangEn: "User ID, change type and point change are required", LangMn: "Хэрэглэгчийн ID, өөрчлөлтийн төрөл, онооны өөрчлөлт хоосон байж болохгүй"},
	"assetTypeMustPointOrTryonPoint":                 {LangZh: "资产类型必须是 point 或 tryon_point", LangEn: "Asset type must be point or tryon_point", LangMn: "Хөрөнгийн төрөл нь point эсвэл tryon_point байх ёстой"},
	"pointInsufficient":                              {LangZh: "积分不足，无法扣除", LangEn: "Insufficient points, cannot deduct", LangMn: "Оноо хүрэлцэхгүй тул суутгах боломжгүй", LangJa: "ポイント不足のため、差し引きできません", LangKo: "포인트가 부족하여 차감할 수 없습니다", LangMs: "Mata tidak mencukupi, tidak boleh ditolak"},
	"tryonPointInsufficient":                         {LangZh: "试衣币不足，无法扣除", LangEn: "Insufficient try-on points, cannot deduct", LangMn: "Өмсгөл туршилтын оноо хүрэлцэхгүй тул суутгах боломжгүй", LangJa: "試着コイン不足のため、差し引きできません", LangKo: "피팅 코인이 부족하여 차감할 수 없습니다", LangMs: "Syiling cuba pakaian tidak mencukupi, tidak boleh ditolak"},
	"submitSuccess":                                  {LangZh: "提交成功", LangEn: "Submitted successfully", LangMn: "Амжилттай илгээлээ", LangJa: "送信に成功しました", LangKo: "제출에 성공했습니다", LangMs: "Berjaya dihantar"},
	"cancelSuccess":                                  {LangZh: "取消成功", LangEn: "Cancelled successfully", LangMn: "Амжилттай цуцаллаа", LangJa: "キャンセルに成功しました", LangKo: "취소에 성공했습니다", LangMs: "Berjaya dibatalkan"},
	"confirmSuccess":                                 {LangZh: "确认成功", LangEn: "Confirmed successfully", LangMn: "Амжилттай баталгаажууллаа", LangJa: "確認に成功しました", LangKo: "확인에 성공했습니다", LangMs: "Berjaya disahkan"},
	"tryonTaskDone":                                  {LangZh: "试衣任务处理完成", LangEn: "Try-on task completed", LangMn: "Өмсгөл туршилтын даалгавар дууслаа", LangJa: "試着タスクの処理が完了しました", LangKo: "가상 착용 작업 처리가 완료되었습니다", LangMs: "Pemprosesan tugas cuba pakaian telah selesai"},
	"tryonTaskProcessing":                            {LangZh: "试衣任务已提交，处理中", LangEn: "Try-on task submitted and processing", LangMn: "Өмсгөл туршилтын даалгавар илгээгдэж, боловсруулж байна", LangJa: "試着タスクを送信しました。処理中です", LangKo: "가상 착용 작업이 제출되었으며 처리 중입니다", LangMs: "Tugas cuba pakaian telah dihantar dan sedang diproses"},
	"tryonBeautifyDone":                              {LangZh: "智能美肤处理完成", LangEn: "Smart skin beautify completed", LangMn: "Ухаалаг арьс сайжруулалт дууслаа", LangJa: "スマート美肌処理が完了しました", LangKo: "스마트 피부 보정 처리가 완료되었습니다", LangMs: "Pemprosesan cantik kulit pintar telah selesai"},
	"tryonBeautifyProcessing":                        {LangZh: "智能美肤处理中", LangEn: "Smart skin beautify is processing", LangMn: "Ухаалаг арьс сайжруулалт боловсруулж байна", LangJa: "スマート美肌処理中です", LangKo: "스마트 피부 보정 처리 중입니다", LangMs: "Cantik kulit pintar sedang diproses"},
	"tryonBeautifyUnsupported":                       {LangZh: "当前模型不支持智能美肤", LangEn: "Current model does not support smart skin beautify", LangMn: "Одоогийн загвар ухаалаг арьс сайжруулалтыг дэмжихгүй", LangJa: "現在のモデルはスマート美肌に対応していません", LangKo: "현재 모델은 스마트 피부 보정을 지원하지 않습니다", LangMs: "Model semasa tidak menyokong cantik kulit pintar"},
	"tryonBeautifyAlreadyUsed":                       {LangZh: "每个试衣任务仅可使用一次智能美肤", LangEn: "Smart skin beautify can be used only once per try-on task", LangMn: "Өмсгөлийн даалгавар бүрт ухаалаг сайжруулалтыг нэг л удаа ашиглана", LangJa: "スマート美肌は試着タスクごとに1回のみ利用できます", LangKo: "스마트 피부 보정은 가상 착용 작업당 1회만 사용할 수 있습니다", LangMs: "Cantik kulit pintar hanya boleh digunakan sekali bagi setiap tugas cuba pakaian"},
	"tryonTaskNotReadyForBeautify":                   {LangZh: "试衣任务未完成，暂不可进行智能美肤", LangEn: "Try-on task is not completed, beautify is unavailable", LangMn: "Өмсгөлийн даалгавар дуусаагүй тул сайжруулалт хийх боломжгүй", LangJa: "試着タスクが未完了のため、スマート美肌は利用できません", LangKo: "가상 착용 작업이 완료되지 않아 스마트 피부 보정을 사용할 수 없습니다", LangMs: "Tugas cuba pakaian belum selesai, cantik kulit belum tersedia"},
	"requestReusedHistory":                           {LangZh: "请求已受理，返回历史任务", LangEn: "Request reused an existing task", LangMn: "Хүсэлтийг хүлээн авч, өмнөх даалгаврыг буцаалаа", LangJa: "既存タスクを再利用して処理しました", LangKo: "기존 작업을 재사용하여 처리했습니다", LangMs: "Permintaan menggunakan semula tugas sedia ada"},
	"tryonCreateTooFrequent":                         {LangZh: "试衣请求过于频繁，请稍后再试", LangEn: "Try-on requests are too frequent, please try again later", LangMn: "Өмсгөл турших хүсэлт хэт олон байна, дараа дахин оролдоно уу", LangJa: "試着リクエストが多すぎます。しばらくしてからお試しください", LangKo: "가상 피팅 요청이 너무 많습니다. 잠시 후 다시 시도해 주세요", LangMs: "Permintaan cuba pakaian terlalu kerap, sila cuba lagi kemudian"},
	"tryonCreateConcurrencyLimited":                  {LangZh: "当前进行中的试衣任务过多，请等待完成后再试", LangEn: "Too many try-on tasks are processing, please wait and retry", LangMn: "Одоогоор боловсруулж буй өмсгөл туршилтын даалгавар хэт олон байна, дууссаны дараа дахин оролдоно уу", LangJa: "処理中の試着タスクが多すぎます。完了後に再試行してください", LangKo: "처리 중인 가상 피팅 작업이 너무 많습니다. 완료 후 다시 시도해 주세요", LangMs: "Terlalu banyak tugasan cuba pakaian sedang diproses, sila tunggu dan cuba semula"},
	"orderNotFound":                                  {LangZh: "订单不存在", LangEn: "Order not found", LangMn: "Захиалга олдсонгүй", LangJa: "注文が存在しません", LangKo: "주문을 찾을 수 없습니다", LangMs: "Pesanan tidak ditemui"},
	"orderIDRequired":                                {LangZh: "订单ID不能为空", LangEn: "Order ID is required", LangMn: "Захиалгын ID хоосон байж болохгүй"},
	"orderStatusInvalid":                             {LangZh: "状态错误", LangEn: "Invalid order status", LangMn: "Захиалгын төлөв буруу"},
	"orderCreateTooFrequent":                         {LangZh: "下单请求过于频繁，请稍后再试", LangEn: "Order requests are too frequent, please try again later", LangMn: "Захиалга үүсгэх хүсэлт хэт олон байна, дараа дахин оролдоно уу"},
	"orderPendingLimitExceeded":                      {LangZh: "待支付订单过多，请先完成或取消后再下单", LangEn: "Too many pending orders, please complete or cancel first", LangMn: "Төлөгдөөгүй захиалга хэт олон байна, эхлээд дуусгах эсвэл цуцлаад дахин оролдоно уу"},
	"orderStateInvalidForUpdate":                     {LangZh: "当前订单状态不支持修改", LangEn: "Current order status does not support updates", LangMn: "Одоогийн захиалгын төлөвт өөрчлөх боломжгүй"},
	"orderExpiredRecreate":                           {LangZh: "订单已超时，请重新下单", LangEn: "Order expired, please create a new order", LangMn: "Захиалга хугацаа хэтэрсэн тул дахин захиална уу"},
	"orderStateInvalidForConfirmPayment":             {LangZh: "当前订单状态不支持确认支付", LangEn: "Current order status does not support payment confirmation", LangMn: "Одоогийн захиалгын төлөв төлбөр баталгаажуулахыг дэмжихгүй"},
	"orderStateInvalidForConfirmReceive":             {LangZh: "当前订单状态不支持确认收货", LangEn: "Current order status does not support confirming receipt", LangMn: "Одоогийн захиалгын төлөв бараа хүлээн авалт баталгаажуулахыг дэмжихгүй"},
	"orderStateInvalidForCancel":                     {LangZh: "当前订单状态不支持取消", LangEn: "Current order status does not support cancellation", LangMn: "Одоогийн захиалгын төлөвт цуцлах боломжгүй"},
	"orderStateInvalidForRefund":                     {LangZh: "当前订单状态不支持退款", LangEn: "Current order status does not support refund", LangMn: "Одоогийн захиалгын төлөв буцаалт хийхийг дэмжихгүй"},
	"orderStateInvalidForSubmitPayment":              {LangZh: "当前订单状态不支持提交付款确认", LangEn: "Current order status does not support payment submission", LangMn: "Одоогийн захиалгын төлөв төлбөр илгээхийг дэмжихгүй"},
	"orderPayMethodNotQrcodeForSubmit":               {LangZh: "仅扫码支付订单可提交付款确认", LangEn: "Only QR-code payment orders can submit payment confirmation", LangMn: "Зөвхөн QR төлбөртэй захиалга төлбөрийн баталгаажуулалт илгээж болно"},
	"orderNoGenFail":                                 {LangZh: "订单号生成失败，请稍后重试", LangEn: "Failed to generate order number, please try again later", LangMn: "Захиалгын дугаар үүсгэж чадсангүй, дараа дахин оролдоно уу"},
	"orderCouponUnavailable":                         {LangZh: "优惠券不可用", LangEn: "Coupon is unavailable", LangMn: "Купон ашиглах боломжгүй"},
	"orderCouponExpired":                             {LangZh: "优惠券已过期", LangEn: "Coupon has expired", LangMn: "Купон хугацаа нь дууссан"},
	"orderCouponNotApplicable":                       {LangZh: "当前订单不可使用此券", LangEn: "This coupon is not applicable to the current order", LangMn: "Одоогийн захиалгад энэ купоныг ашиглах боломжгүй"},
	"orderGoodInfoQueryFailed":                       {LangZh: "查询商品信息失败", LangEn: "Failed to query product information", LangMn: "Барааны мэдээлэл шалгахад алдаа гарлаа"},
	"orderPointsUseTimesLimitReached":                {LangZh: "商品积分抵扣次数已达上限", LangEn: "Points discount usage limit reached for this product", LangMn: "Энэ барааны онооны хөнгөлөлтийн хэрэглээ дээд хязгаарт хүрсэн"},
	"orderContainsGoodsNotSupportPoints":             {LangZh: "订单中存在不支持积分抵扣的商品", LangEn: "Order contains items that do not support points discount", LangMn: "Захиалгад онооны хөнгөлөлт дэмждэггүй бараа байна"},
	"orderInventoryInsufficient":                     {LangZh: "库存不足", LangEn: "Insufficient inventory", LangMn: "Нөөц хүрэлцэхгүй байна"},
	"orderInventoryInsufficientRetry":                {LangZh: "库存不足，请重试", LangEn: "Insufficient inventory, please retry", LangMn: "Нөөц хүрэлцэхгүй байна, дахин оролдоно уу"},
	"orderCartEmpty":                                 {LangZh: "购物车为空", LangEn: "Cart is empty", LangMn: "Сагс хоосон байна"},
	"orderOnlyUserCanSubmitPendingConfirm":           {LangZh: "仅用户可提交付款确认", LangEn: "Only users can submit payment confirmation", LangMn: "Зөвхөн хэрэглэгч төлбөрийн баталгаажуулалт илгээж болно"},
	"orderConfirmPaymentStateInvalid":                {LangZh: "当前订单状态不支持确认收款", LangEn: "Current order status does not support payment confirmation", LangMn: "Одоогийн захиалгын төлөв төлбөр баталгаажуулахыг дэмжихгүй"},
	"orderRefundNotPaid":                             {LangZh: "订单未支付，无法申请退款", LangEn: "Order is unpaid and cannot apply for refund", LangMn: "Захиалга төлөгдөөгүй тул буцаалт хүсэх боломжгүй"},
	"orderRefundPendingConfirmUnsupported":           {LangZh: "订单待后台确认，暂不支持退款", LangEn: "Order is pending backend confirmation and refund is not supported", LangMn: "Захиалга админ баталгаажуулалт хүлээж байгаа тул буцаалт дэмжихгүй"},
	"orderRefundCanceled":                            {LangZh: "订单已取消，无法申请退款", LangEn: "Order is canceled and cannot apply for refund", LangMn: "Захиалга цуцлагдсан тул буцаалт хүсэх боломжгүй"},
	"orderAlreadyRefunded":                           {LangZh: "订单已退款", LangEn: "Order already refunded", LangMn: "Захиалгад буцаалт хийгдсэн"},
	"orderRefundProcessing":                          {LangZh: "退款申请处理中", LangEn: "Refund request is processing", LangMn: "Буцаалтын хүсэлт боловсруулагдаж байна"},
	"orderStateInvalidForApplyRefund":                {LangZh: "当前订单状态不允许退款", LangEn: "Current order status does not allow refund", LangMn: "Одоогийн захиалгын төлөв буцаалт зөвшөөрөхгүй"},
	"orderRefundStateNotApplying":                    {LangZh: "订单未处于退款申请中", LangEn: "Order is not in refund applying state", LangMn: "Захиалга буцаалт хүсэж буй төлөвт байхгүй"},
	"payMethodRequired":                              {LangZh: "支付方式不能为空", LangEn: "Payment method is required", LangMn: "Төлбөрийн арга хоосон байж болохгүй", LangJa: "支払い方法は必須です", LangKo: "결제 수단은 필수입니다", LangMs: "Kaedah pembayaran diperlukan"},
	"loginRequired":                                  {LangZh: "请先登录", LangEn: "Please login first", LangMn: "Эхлээд нэвтэрнэ үү", LangJa: "先にログインしてください", LangKo: "먼저 로그인해 주세요", LangMs: "Sila log masuk dahulu"},
	"tryonRechargeAmountFormatError":                 {LangZh: "充值金额格式不正确", LangEn: "Invalid recharge amount format", LangMn: "Цэнэглэх дүнгийн формат буруу"},
	"tryonRechargeOrderNoGenFail":                    {LangZh: "订单号生成失败，请稍后重试", LangEn: "Failed to generate order number, please try again later", LangMn: "Захиалгын дугаар үүсгэж чадсангүй, дараа дахин оролдоно уу"},
	"tryonRechargePlanNotConfigured":                 {LangZh: "充值套餐未配置", LangEn: "Recharge plans are not configured", LangMn: "Цэнэглэх багцууд тохируулагдаагүй байна"},
	"tryonRechargePlanConfigInvalid":                 {LangZh: "充值套餐配置格式错误", LangEn: "Recharge plan configuration format is invalid", LangMn: "Цэнэглэх багцын тохиргооны формат буруу"},
	"tryonRechargePlanChanged":                       {LangZh: "充值套餐已变更，请刷新页面后重试", LangEn: "Recharge plans changed, please refresh and try again", LangMn: "Цэнэглэх багц өөрчлөгдсөн тул хуудсыг шинэчлээд дахин оролдоно уу"},
	"tryonRechargePointsMustPositive":                {LangZh: "充值点数必须大于0", LangEn: "Recharge points must be greater than 0", LangMn: "Цэнэглэх оноо 0-оос их байх ёстой"},
	"tryonRechargeAmountMustPositive":                {LangZh: "充值金额必须大于0", LangEn: "Recharge amount must be greater than 0", LangMn: "Цэнэглэх дүн 0-оос их байх ёстой"},
	"tryonRechargePendingLimitExceeded":              {LangZh: "待支付充值订单过多，请先完成或取消后再试", LangEn: "Too many pending recharge orders, please complete or cancel first", LangMn: "Төлөгдөөгүй цэнэглэх захиалга хэт олон байна, эхлээд дуусгах эсвэл цуцлаад дахин оролдоно уу"},
	"tryonRechargeOrderExpired":                      {LangZh: "充值订单已超时，请重新下单", LangEn: "Recharge order expired, please create a new order", LangMn: "Цэнэглэх захиалга хугацаа хэтэрсэн тул дахин захиална уу"},
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
	"loginSuccess":           {LangZh: "登录成功", LangEn: "Login successful", LangMn: "Амжилттай нэвтэрлээ", LangJa: "ログインに成功しました", LangKo: "로그인에 성공했습니다", LangMs: "Log masuk berjaya"},
	"loginFail":              {LangZh: "用户名或密码错误", LangEn: "Invalid username or password", LangMn: "Хэрэглэгчийн нэр эсвэл нууц үг буруу", LangJa: "ユーザー名またはパスワードが正しくありません", LangKo: "사용자 이름 또는 비밀번호가 올바르지 않습니다", LangMs: "Nama pengguna atau kata laluan tidak betul"},
	"phoneLoginFail":         {LangZh: "手机号不存在或密码错误", LangEn: "Phone number not found or wrong password", LangMn: "Утасны дугаар эсвэл нууц үг буруу", LangJa: "電話番号が存在しないか、パスワードが正しくありません", LangKo: "전화번호가 없거나 비밀번호가 올바르지 않습니다", LangMs: "Nombor telefon tidak wujud atau kata laluan salah"},
	"captchaError":           {LangZh: "验证码错误", LangEn: "Invalid captcha", LangMn: "Баталгаажуулах код буруу", LangJa: "認証コードが正しくありません", LangKo: "인증 코드가 올바르지 않습니다", LangMs: "Kod pengesahan tidak sah"},
	"codeEmpty":              {LangZh: "code不能为空", LangEn: "Code is required", LangMn: "Код хоосон байж болохгүй", LangJa: "code は必須です", LangKo: "code는 필수입니다", LangMs: "code diperlukan"},
	"openidFail":             {LangZh: "获取openid失败", LangEn: "Failed to get OpenID", LangMn: "OpenID авч чадсангүй", LangJa: "OpenID の取得に失敗しました", LangKo: "OpenID 가져오기에 실패했습니다", LangMs: "Gagal mendapatkan OpenID"},
	"tokenFail":              {LangZh: "获取token失败", LangEn: "Failed to get token", LangMn: "Token авч чадсангүй", LangJa: "token の取得に失敗しました", LangKo: "token 가져오기에 실패했습니다", LangMs: "Gagal mendapatkan token"},
	"loginStatusFail":        {LangZh: "设置登录状态失败", LangEn: "Failed to set login status", LangMn: "Нэвтрэлтийн төлөв тохируулж чадсангүй", LangJa: "ログイン状態の設定に失敗しました", LangKo: "로그인 상태 설정에 실패했습니다", LangMs: "Gagal menetapkan status log masuk"},
	"jwtBlacklistFail":       {LangZh: "jwt作废失败", LangEn: "Failed to revoke JWT", LangMn: "JWT хүчингүй болгож чадсангүй", LangJa: "JWT の失効に失敗しました", LangKo: "JWT 폐기에 실패했습니다", LangMs: "Gagal membatalkan JWT"},
	"passwordMismatch":       {LangZh: "两次输入的密码不一致", LangEn: "Passwords do not match", LangMn: "Нууц үг таарахгүй байна", LangJa: "入力した2つのパスワードが一致しません", LangKo: "두 번 입력한 비밀번호가 일치하지 않습니다", LangMs: "Dua kata laluan tidak sepadan"},
	"usernameOrPhoneExists":  {LangZh: "用户名或手机号码已存在", LangEn: "Username or phone already exists", LangMn: "Хэрэглэгчийн нэр эсвэл утасны дугаар аль хэдийн бүртгэгдсэн байна", LangJa: "ユーザー名または電話番号は既に存在します", LangKo: "사용자 이름 또는 전화번호가 이미 존재합니다", LangMs: "Nama pengguna atau telefon sudah wujud"},
	"phoneAlreadyRegistered": {LangZh: "该手机号已注册", LangEn: "Phone number already registered", LangMn: "Энэ утасны дугаар аль хэдийн бүртгэгдсэн байна", LangJa: "この電話番号は既に登録されています", LangKo: "이 전화번호는 이미 등록되었습니다", LangMs: "Nombor telefon ini sudah didaftarkan"},
	"signInSuccess":          {LangZh: "签到成功", LangEn: "Sign-in successful", LangMn: "Амжилттай тэмдэглэлээ", LangJa: "チェックインに成功しました", LangKo: "출석 체크에 성공했습니다", LangMs: "Daftar masuk berjaya"},
	"signInAlreadyToday":     {LangZh: "今日已签到", LangEn: "Already signed in today", LangMn: "Өнөөдөр аль хэдийн тэмдэглэсэн байна", LangJa: "本日は既にチェックイン済みです", LangKo: "오늘 이미 출석 체크했습니다", LangMs: "Sudah daftar masuk hari ini"},
	"cannotModify":           {LangZh: "无法修改", LangEn: "Cannot modify", LangMn: "Өөрчлөх боломжгүй"},
	"seriesNotFound":         {LangZh: "剧集不存在", LangEn: "Series not found", LangMn: "Цуврал олдсонгүй"},
	"loginLocked":            {LangZh: "登录失败次数过多，请稍后再试", LangEn: "Too many failed attempts, please try again later", LangMn: "Нэвтрэлт хэт олон удаа амжилтгүй, дараа дахин оролдоно уу", LangJa: "ログイン失敗回数が多すぎます。しばらくしてからお試しください", LangKo: "로그인 실패 횟수가 너무 많습니다. 잠시 후 다시 시도해 주세요", LangMs: "Terlalu banyak percubaan gagal, sila cuba lagi kemudian"},
	"accountBanned":          {LangZh: "账号已被封禁", LangEn: "Account has been banned", LangMn: "Бүртгэл хориглогдсон"},
	"registerIPLimit":        {LangZh: "该IP注册次数已达上限", LangEn: "Registration limit reached for this IP", LangMn: "Энэ IP-ээс бүртгэх хязгаарт хүрсэн", LangJa: "このIPの登録回数が上限に達しました", LangKo: "이 IP의 가입 횟수가 한도에 도달했습니다", LangMs: "Had pendaftaran IP telah dicapai"},
	"areaCodeInvalid":        {LangZh: "区号无效或未启用", LangEn: "Invalid or disabled area code", LangMn: "Бүсийн код буруу эсвэл идэвхгүй", LangJa: "国番号が無効または未有効化です", LangKo: "국가번호가 유효하지 않거나 비활성화되어 있습니다", LangMs: "Kod kawasan tidak sah atau dinyahaktif"},
	"phoneFormatError":       {LangZh: "手机号格式不正确", LangEn: "Invalid phone number format", LangMn: "Утасны дугаарын формат буруу", LangJa: "電話番号の形式が正しくありません", LangKo: "전화번호 형식이 올바르지 않습니다", LangMs: "Format nombor telefon tidak sah"},
	"changeSuccess":          {LangZh: "修改成功", LangEn: "Changed successfully", LangMn: "Амжилттай өөрчлөлөө"},
	"phoneRequired":          {LangZh: "手机号不能为空", LangEn: "Phone number is required", LangMn: "Утасны дугаар хоосон байж болохгүй"},
	"passwordError":          {LangZh: "密码错误", LangEn: "Incorrect password", LangMn: "Нууц үг буруу", LangJa: "パスワードが正しくありません", LangKo: "비밀번호가 올바르지 않습니다", LangMs: "Kata laluan salah"},
	"userNotExist":           {LangZh: "用户不存在", LangEn: "User not found", LangMn: "Хэрэглэгч олдсонгүй", LangJa: "ユーザーが存在しません", LangKo: "사용자를 찾을 수 없습니다", LangMs: "Pengguna tidak ditemui"},

	// ========== JWT 中间件 ==========
	"notLogin":     {LangZh: "未登录或非法访问，请登录", LangEn: "Please login to continue", LangMn: "Нэвтэрнэ үү", LangJa: "未ログインまたは不正アクセスです。ログインしてください", LangKo: "로그인되지 않았거나 비정상 접근입니다. 로그인해 주세요", LangMs: "Belum log masuk atau akses tidak sah, sila log masuk"},
	"tokenInvalid": {LangZh: "您的帐户异地登陆或令牌失效", LangEn: "Your session is invalid, please login again", LangMn: "Таны токен хүчингүй болсон", LangJa: "アカウントが他所でログインしたか、トークンが無効です", LangKo: "계정이 다른 곳에서 로그인되었거나 토큰이 유효하지 않습니다", LangMs: "Akaun log masuk di tempat lain atau token tidak sah"},
	"tokenExpired": {LangZh: "登录已过期，请重新登录", LangEn: "Session expired, please login again", LangMn: "Нэвтрэлт дууссан, дахин нэвтэрнэ үү", LangJa: "ログインの有効期限が切れました。再度ログインしてください", LangKo: "로그인이 만료되었습니다. 다시 로그인해 주세요", LangMs: "Sesi log masuk telah tamat, sila log masuk semula"},
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

func languageFallbackChain(lang string) []string {
	switch lang {
	case LangVi, LangAr, LangJa, LangKo, LangMs:
		return []string{lang, LangEn, LangZh}
	default:
		return []string{lang, LangZh}
	}
}

// HasKey 判断 i18n key 是否已定义。
func HasKey(key string) bool {
	key = strings.TrimSpace(key)
	if key == "" {
		return false
	}
	_, ok := messages[key]
	return ok
}

// T 根据消息 key 和当前请求语言返回翻译文本
func T(c *gin.Context, key string) string {
	lang := GetLang(c)
	if m, ok := messages[key]; ok {
		for _, fallbackLang := range languageFallbackChain(lang) {
			if s, exists := m[fallbackLang]; exists && s != "" {
				return s
			}
		}
		for _, text := range m {
			if text != "" {
				return text
			}
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
