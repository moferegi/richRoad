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
	"ok":                 {LangZh: "操作成功", LangEn: "Success", LangMn: "Амжилттай"},
	"fail":               {LangZh: "操作失败", LangEn: "Operation failed", LangMn: "Амжилтгүй"},
	"success":            {LangZh: "成功", LangEn: "Success", LangMn: "Амжилттай"},
	"createSuccess":      {LangZh: "创建成功", LangEn: "Created successfully", LangMn: "Амжилттай үүсгэлээ"},
	"createFail":         {LangZh: "创建失败", LangEn: "Create failed", LangMn: "Үүсгэж чадсангүй"},
	"deleteSuccess":      {LangZh: "删除成功", LangEn: "Deleted successfully", LangMn: "Амжилттай устгалаа"},
	"deleteFail":         {LangZh: "删除失败", LangEn: "Delete failed", LangMn: "Устгаж чадсангүй"},
	"batchDeleteSuccess": {LangZh: "批量删除成功", LangEn: "Batch deleted", LangMn: "Бөөнөөр устгалаа"},
	"batchDeleteFail":    {LangZh: "批量删除失败", LangEn: "Batch delete failed", LangMn: "Бөөнөөр устгаж чадсангүй"},
	"updateSuccess":      {LangZh: "更新成功", LangEn: "Updated successfully", LangMn: "Амжилттай шинэчлэлээ"},
	"updateFail":         {LangZh: "更新失败", LangEn: "Update failed", LangMn: "Шинэчлэж чадсангүй"},
	"querySuccess":       {LangZh: "查询成功", LangEn: "Query successful", LangMn: "Хайлт амжилттай"},
	"queryFail":          {LangZh: "查询失败", LangEn: "Query failed", LangMn: "Хайлт амжилтгүй"},
	"getSuccess":         {LangZh: "获取成功", LangEn: "Retrieved successfully", LangMn: "Амжилттай авлаа"},
	"getFail":            {LangZh: "获取失败", LangEn: "Retrieve failed", LangMn: "Авч чадсангүй"},
	"setSuccess":         {LangZh: "设置成功", LangEn: "Set successfully", LangMn: "Амжилттай тохируулалаа"},
	"setFail":            {LangZh: "设置失败", LangEn: "Set failed", LangMn: "Тохируулж чадсангүй"},

	// ========== 认证/登录 ==========
	"loginSuccess":     {LangZh: "登录成功", LangEn: "Login successful", LangMn: "Амжилттай нэвтэрлээ"},
	"loginFail":        {LangZh: "用户名不存在或者密码错误", LangEn: "Invalid username or password", LangMn: "Хэрэглэгчийн нэр эсвэл нууц үг буруу"},
	"captchaError":     {LangZh: "验证码错误", LangEn: "Invalid captcha", LangMn: "Баталгаажуулах код буруу"},
	"codeEmpty":        {LangZh: "code不能为空", LangEn: "Code is required", LangMn: "Код хоосон байж болохгүй"},
	"openidFail":       {LangZh: "获取openid失败", LangEn: "Failed to get OpenID", LangMn: "OpenID авч чадсангүй"},
	"tokenFail":        {LangZh: "获取token失败", LangEn: "Failed to get token", LangMn: "Token авч чадсангүй"},
	"loginStatusFail":  {LangZh: "设置登录状态失败", LangEn: "Failed to set login status", LangMn: "Нэвтрэлтийн төлөв тохируулж чадсангүй"},
	"jwtBlacklistFail": {LangZh: "jwt作废失败", LangEn: "Failed to revoke JWT", LangMn: "JWT хүчингүй болгож чадсангүй"},
	"passwordMismatch": {LangZh: "两次输入的密码不一致", LangEn: "Passwords do not match", LangMn: "Нууц үг таарахгүй байна"},
	"cannotModify":     {LangZh: "无法修改", LangEn: "Cannot modify", LangMn: "Өөрчлөх боломжгүй"},
	"seriesNotFound":   {LangZh: "剧集不存在", LangEn: "Series not found", LangMn: "Цуврал олдсонгүй"},
	"loginLocked":      {LangZh: "登录失败次数过多，请稍后再试", LangEn: "Too many failed attempts, please try again later", LangMn: "Нэвтрэлт хэт олон удаа амжилтгүй, дараа дахин оролдоно уу"},
	"accountBanned":    {LangZh: "账号已被封禁", LangEn: "Account has been banned", LangMn: "Бүртгэл хориглогдсон"},
	"registerIPLimit":  {LangZh: "该IP注册次数已达上限", LangEn: "Registration limit reached for this IP", LangMn: "Энэ IP-ээс бүртгэх хязгаарт хүрсэн"},
	"areaCodeInvalid":  {LangZh: "区号无效或未启用", LangEn: "Invalid or disabled area code", LangMn: "Бүсийн код буруу эсвэл идэвхгүй"},
	"phoneFormatError": {LangZh: "手机号格式不正确", LangEn: "Invalid phone number format", LangMn: "Утасны дугаарын формат буруу"},
	"changeSuccess":    {LangZh: "修改成功", LangEn: "Changed successfully", LangMn: "Амжилттай өөрчлөлөө"},
	"phoneRequired":    {LangZh: "手机号不能为空", LangEn: "Phone number is required", LangMn: "Утасны дугаар хоосон байж болохгүй"},
	"passwordError":    {LangZh: "密码错误", LangEn: "Incorrect password", LangMn: "Нууц үг буруу"},
	"userNotExist":     {LangZh: "用户不存在", LangEn: "User not found", LangMn: "Хэрэглэгч олдсонгүй"},

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
