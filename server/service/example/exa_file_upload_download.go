package example

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	defaultUploadMaxSizeMB int64 = 20

	defaultFileListPageSize = 10
	maxFileListPageSize     = 100
	maxFileKeywordLength    = 128

	maxFileNameLength       = 255
	maxFileKeywordLen       = 512
	maxImportBatchSize      = 200
	maxImportURLLength      = 2048
	defaultImportedFileName = "imported-file"

	uploadMaxSizeConfigKey = "security_upload_max_size_mb"
	uploadMaxSizeEnvKey    = "CS_UPLOAD_MAX_SIZE_MB"

	uploadStrictValidationConfigKey = "security_upload_strict_validation_enabled"
	uploadStrictValidationEnvKey    = "CS_UPLOAD_STRICT_VALIDATION"
)

var (
	blockedUploadExts = map[string]struct{}{
		"php":   {},
		"php3":  {},
		"php4":  {},
		"php5":  {},
		"phtml": {},
		"jsp":   {},
		"jspx":  {},
		"asp":   {},
		"aspx":  {},
		"cgi":   {},
		"exe":   {},
		"dll":   {},
		"com":   {},
		"scr":   {},
		"msi":   {},
		"bat":   {},
		"cmd":   {},
		"ps1":   {},
		"sh":    {},
		"bash":  {},
		"zsh":   {},
		"jar":   {},
		"war":   {},
		"vbs":   {},
		"js":    {},
		"mjs":   {},
		"cjs":   {},
		"html":  {},
		"htm":   {},
		"svg":   {},
	}
	blockedUploadMimeTypes = map[string]struct{}{
		"application/x-msdownload": {},
		"application/x-dosexec":    {},
		"application/x-executable": {},
		"application/x-sh":         {},
		"application/x-php":        {},
		"text/x-php":               {},
		"text/x-shellscript":       {},
	}
	blockedUploadMimePrefixes = []string{
		"text/html",
		"application/javascript",
		"text/javascript",
	}
)

func loadUploadMaxSizeBytes() int64 {
	mb := utils.GetInt64Setting(uploadMaxSizeConfigKey, uploadMaxSizeEnvKey, defaultUploadMaxSizeMB)
	if mb <= 0 {
		mb = defaultUploadMaxSizeMB
	}
	return mb * 1024 * 1024
}

func isUploadStrictValidationEnabled() bool {
	return utils.GetBoolSetting(uploadStrictValidationConfigKey, uploadStrictValidationEnvKey, true)
}

func validateUploadHeader(header *multipart.FileHeader) error {
	if header == nil {
		return errors.New("文件不能为空")
	}
	if header.Size <= 0 {
		return errors.New("空文件不允许上传")
	}
	maxSizeBytes := loadUploadMaxSizeBytes()
	if maxSizeBytes > 0 && header.Size > maxSizeBytes {
		return fmt.Errorf("文件大小超过限制(%dMB)", maxSizeBytes/1024/1024)
	}

	fileName := strings.TrimSpace(filepath.Base(header.Filename))
	if fileName == "" || fileName == "." {
		return errors.New("文件名非法")
	}
	if strings.Contains(fileName, "..") {
		return errors.New("文件名非法")
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
	if ext == "" {
		return errors.New("文件扩展名缺失")
	}

	if isUploadStrictValidationEnabled() {
		if _, blocked := blockedUploadExts[ext]; blocked {
			return errors.New("该文件类型禁止上传")
		}

		contentType, detectErr := detectUploadContentType(header)
		if detectErr != nil {
			return errors.New("读取上传文件内容失败")
		}
		if isBlockedUploadMimeType(contentType) {
			return errors.New("该文件内容类型禁止上传")
		}
	}
	return nil
}

func detectUploadContentType(header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, readErr := file.Read(buf)
	if readErr != nil && readErr != io.EOF {
		return "", readErr
	}
	if n <= 0 {
		return "", nil
	}

	return strings.ToLower(strings.TrimSpace(http.DetectContentType(buf[:n]))), nil
}

func isBlockedUploadMimeType(contentType string) bool {
	contentType = strings.TrimSpace(strings.ToLower(contentType))
	if contentType == "" {
		return false
	}
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = strings.TrimSpace(contentType[:idx])
	}
	if _, blocked := blockedUploadMimeTypes[contentType]; blocked {
		return true
	}
	for _, prefix := range blockedUploadMimePrefixes {
		if strings.HasPrefix(contentType, prefix) {
			return true
		}
	}
	return false
}

func sanitizeObjectKeyForDelete(raw string) (string, error) {
	key := strings.TrimSpace(strings.ReplaceAll(raw, "\\", "/"))
	if key == "" {
		return "", errors.New("文件key为空")
	}
	if len(key) > 1024 {
		return "", errors.New("文件key过长")
	}
	if strings.Contains(key, "..") || strings.Contains(key, "\x00") || strings.Contains(key, "://") || strings.ContainsAny(key, "?#") {
		return "", errors.New("文件key非法")
	}

	cleaned := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(key, "/")), "/")
	if cleaned == "" || cleaned == "." {
		return "", errors.New("文件key非法")
	}

	return cleaned, nil
}

func trimRunes(raw string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(raw)
	if len(runes) <= max {
		return raw
	}
	return string(runes[:max])
}

func sanitizeFileName(raw string) string {
	name := strings.TrimSpace(raw)
	name = strings.ReplaceAll(name, "\r", " ")
	name = strings.ReplaceAll(name, "\n", " ")
	name = strings.ReplaceAll(name, "\t", " ")
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	return trimRunes(name, maxFileNameLength)
}

func normalizeImportFileURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("导入URL不能为空")
	}
	if len(value) > maxImportURLLength {
		return "", errors.New("导入URL过长")
	}
	if strings.Contains(value, "\x00") || strings.ContainsAny(value, "\r\n\t") {
		return "", errors.New("导入URL非法")
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		parsed, err := url.Parse(value)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return "", errors.New("导入URL非法")
		}
		parsed.Fragment = ""
		return parsed.String(), nil
	}

	if strings.HasPrefix(value, "/") {
		cleaned := path.Clean(value)
		if cleaned == "/" || cleaned == "." || strings.Contains(cleaned, "..") {
			return "", errors.New("导入URL非法")
		}
		return cleaned, nil
	}

	return "", errors.New("导入URL仅支持http/https或绝对路径")
}

func inferImportFileName(fileURL string) string {
	if strings.TrimSpace(fileURL) == "" {
		return defaultImportedFileName
	}
	if parsed, err := url.Parse(fileURL); err == nil {
		if base := strings.TrimSpace(path.Base(parsed.Path)); base != "" && base != "." && base != "/" {
			return base
		}
	}
	if base := strings.TrimSpace(path.Base(fileURL)); base != "" && base != "." && base != "/" {
		return base
	}
	return defaultImportedFileName
}

func normalizeFileCategory(raw string, fallback string) string {
	category := strings.ToLower(strings.TrimSpace(raw))
	switch category {
	case "shoe", "person", "cloth":
		return category
	default:
		return inferMediaCategory("", "", fallback)
	}
}

func hasGlobalFileManageScope(authorityID uint) bool {
	return authorityID == 888
}

func applyFileOwnershipScope(db *gorm.DB, userID uint, authorityID uint) *gorm.DB {
	if hasGlobalFileManageScope(authorityID) {
		return db
	}
	if userID == 0 {
		return db.Where("1 = 0")
	}
	return db.Where("created_by = ?", userID)
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) Upload(file example.ExaFileUploadAndDownload) error {
	return global.GVA_DB.Create(&file).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: model.ExaFileUploadAndDownload, error

func (e *FileUploadAndDownloadService) FindFile(id uint) (example.ExaFileUploadAndDownload, error) {
	var file example.ExaFileUploadAndDownload
	err := global.GVA_DB.Where("id = ?", id).First(&file).Error
	return file, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFile
//@description: 删除文件记录
//@param: file model.ExaFileUploadAndDownload
//@return: err error

func (e *FileUploadAndDownloadService) DeleteFile(file example.ExaFileUploadAndDownload, operatorUserID uint, operatorAuthorityID uint) (err error) {
	if file.ID == 0 {
		return errors.New("文件ID非法")
	}

	scopedDB := applyFileOwnershipScope(global.GVA_DB, operatorUserID, operatorAuthorityID)
	var fileFromDb example.ExaFileUploadAndDownload
	err = scopedDB.Where("id = ?", file.ID).First(&fileFromDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件不存在或无权限")
		}
		return err
	}

	// key 为空时仅删除数据库记录，避免误删对象存储路径前缀。
	if strings.TrimSpace(fileFromDb.Key) != "" {
		safeKey, keyErr := sanitizeObjectKeyForDelete(fileFromDb.Key)
		if keyErr != nil {
			global.GVA_LOG.Warn("检测到非法文件key，跳过对象存储删除", zap.Uint("fileID", fileFromDb.ID), zap.String("key", fileFromDb.Key), zap.Error(keyErr))
		} else {
			oss := upload.NewOss()
			if err = oss.DeleteFile(safeKey); err != nil {
				return errors.New("文件删除失败")
			}
		}
	}

	err = scopedDB.Where("id = ?", file.ID).Unscoped().Delete(&example.ExaFileUploadAndDownload{}).Error
	return err
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownloadService) EditFileName(file example.ExaFileUploadAndDownload, operatorUserID uint, operatorAuthorityID uint) (err error) {
	if file.ID == 0 {
		return errors.New("文件ID非法")
	}

	name := sanitizeFileName(file.Name)
	if name == "" {
		return errors.New("文件名不能为空")
	}

	result := applyFileOwnershipScope(global.GVA_DB.Model(&example.ExaFileUploadAndDownload{}), operatorUserID, operatorAuthorityID).
		Where("id = ?", file.ID).
		Update("name", name)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("文件不存在或无权限")
	}
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.ExaAttachmentCategorySearch
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info request.ExaAttachmentCategorySearch, operatorUserID uint, operatorAuthorityID uint) (list []example.ExaFileUploadAndDownload, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = defaultFileListPageSize
	}
	if info.PageSize > maxFileListPageSize {
		info.PageSize = maxFileListPageSize
	}

	keyword := trimRunes(strings.TrimSpace(info.Keyword), maxFileKeywordLength)

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := applyFileOwnershipScope(global.GVA_DB.Model(&example.ExaFileUploadAndDownload{}), operatorUserID, operatorAuthorityID)

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	if info.ClassId > 0 {
		db = db.Where("class_id = ?", info.ClassId)
	}

	if normalizedPosition := normalizeUploadPosition(info.Position); normalizedPosition != "" {
		db = db.Where("position = ?", normalizedPosition)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&list).Error
	if err != nil {
		return
	}

	for i := range list {
		item := &list[i]
		changed := false
		if strings.TrimSpace(item.ThumbnailURL) == "" {
			item.ThumbnailURL = item.Url
			changed = true
		}
		if strings.TrimSpace(item.Categories) == "" {
			item.Categories = inferMediaCategory("", "", item.Url)
			changed = true
		}
		if strings.TrimSpace(item.Position) == "" {
			item.Position = inferUploadPosition("", "", item.Url)
			changed = true
		}
		if item.Size < 0 {
			item.Size = 0
			changed = true
		}
		item.Keywords = strings.TrimSpace(item.Keywords)
		if changed {
			_ = global.GVA_DB.Model(&example.ExaFileUploadAndDownload{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
				"thumbnail_url": item.ThumbnailURL,
				"categories":    item.Categories,
				"position":      item.Position,
				"keywords":      item.Keywords,
				"size":          item.Size,
			}).Error
		}
	}
	return list, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file model.ExaFileUploadAndDownload, err error

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string, classId int, folder string, uploadType string, uploadPosition string, operatorUserID uint) (file example.ExaFileUploadAndDownload, err error) {
	if err = validateUploadHeader(header); err != nil {
		return file, err
	}

	oss := upload.NewOss()
	filePath, key, uploadErr := upload.UploadFileToFolder(oss, header, folder)
	if uploadErr != nil {
		return file, uploadErr
	}
	mediaCategory := inferMediaCategory(uploadType, folder, filePath)
	normalizedPosition := normalizeUploadPosition(uploadPosition)
	if normalizedPosition == "" {
		normalizedPosition = inferUploadPosition(uploadType, folder, filePath)
	}
	tag := strings.TrimPrefix(strings.ToLower(filepath.Ext(header.Filename)), ".")
	fileName := sanitizeFileName(filepath.Base(header.Filename))
	if fileName == "" || fileName == "." {
		fileName = sanitizeFileName(filepath.Base(key))
	}
	if fileName == "" {
		fileName = defaultImportedFileName
	}
	f := example.ExaFileUploadAndDownload{
		CreatedBy:    operatorUserID,
		Url:          filePath,
		ThumbnailURL: filePath,
		Categories:   mediaCategory,
		Position:     normalizedPosition,
		Keywords:     "",
		Size:         header.Size,
		Name:         fileName,
		ClassId:      classId,
		Tag:          tag,
		Key:          key,
	}
	if noSave == "0" {
		// 检查是否已存在相同key的记录
		var existingFile example.ExaFileUploadAndDownload
		db := global.GVA_DB.Where("key = ?", key)
		if operatorUserID > 0 {
			db = db.Where("created_by = ?", operatorUserID)
		}
		err = db.First(&existingFile).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return f, e.Upload(f)
		}
		return f, err
	}
	return f, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ImportURL
//@description: 导入URL
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) ImportURL(file *[]example.ExaFileUploadAndDownload, operatorUserID uint) error {
	if file == nil {
		return nil
	}
	if len(*file) == 0 {
		return nil
	}
	if len(*file) > maxImportBatchSize {
		return fmt.Errorf("单次最多导入%d条", maxImportBatchSize)
	}

	for i := range *file {
		item := &(*file)[i]
		item.CreatedBy = operatorUserID

		normalizedURL, urlErr := normalizeImportFileURL(item.Url)
		if urlErr != nil {
			return fmt.Errorf("第%d条URL非法", i+1)
		}
		item.Url = normalizedURL

		if strings.TrimSpace(item.ThumbnailURL) == "" {
			item.ThumbnailURL = item.Url
		} else {
			normalizedThumb, thumbErr := normalizeImportFileURL(item.ThumbnailURL)
			if thumbErr != nil {
				item.ThumbnailURL = item.Url
			} else {
				item.ThumbnailURL = normalizedThumb
			}
		}

		item.Name = sanitizeFileName(item.Name)
		if item.Name == "" {
			item.Name = sanitizeFileName(inferImportFileName(item.Url))
		}
		if item.Name == "" {
			item.Name = defaultImportedFileName
		}

		item.Categories = normalizeFileCategory(item.Categories, item.Url)

		if normalizedPosition := normalizeUploadPosition(item.Position); normalizedPosition != "" {
			item.Position = normalizedPosition
		} else {
			item.Position = inferUploadPosition("", "", item.Url)
		}

		item.Keywords = trimRunes(strings.TrimSpace(item.Keywords), maxFileKeywordLen)
		if item.Keywords == "" {
			item.Keywords = ""
		}
		if item.ClassId < 0 {
			item.ClassId = 0
		}
		if item.Size < 0 {
			item.Size = 0
		}
		if strings.TrimSpace(item.Tag) == "" {
			item.Tag = strings.TrimPrefix(strings.ToLower(path.Ext(item.Url)), ".")
		}
	}

	return global.GVA_DB.Create(file).Error
}

func inferMediaCategory(uploadType string, folder string, filePath string) string {
	parts := []string{uploadType, folder, filePath}
	joined := strings.ToLower(strings.Join(parts, " "))

	if strings.Contains(joined, "shoe") {
		return "shoe"
	}
	if strings.Contains(joined, "person") || strings.Contains(joined, "model") {
		return "person"
	}
	if strings.Contains(joined, "cloth") || strings.Contains(joined, "upper") || strings.Contains(joined, "lower") {
		return "cloth"
	}

	return "cloth"
}

func normalizeUploadPosition(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "", "all":
		return ""
	case "tryon", "try-on", "fitting", "shoe", "shoeroom", "tryon_position":
		return "tryon"
	case "kefu", "cs", "chat", "workbench", "customer_service", "customer-service", "kefu_position":
		return "kefu"
	case "other":
		return "other"
	default:
		return ""
	}
}

func inferUploadPosition(uploadType string, folder string, filePath string) string {
	parts := []string{uploadType, folder, filePath}
	joined := strings.ToLower(strings.Join(parts, " "))

	if strings.Contains(joined, "kefu") || strings.Contains(joined, "workbench") || strings.Contains(joined, "chat") || strings.Contains(joined, "/cs") || strings.Contains(joined, "cs/") {
		return "kefu"
	}

	if strings.Contains(joined, "tryon") || strings.Contains(joined, "try-on") || strings.Contains(joined, "shoe") || strings.Contains(joined, "shoeroom") || strings.Contains(joined, "cloth-on") {
		return "tryon"
	}

	return "other"
}
