package example

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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

func (e *FileUploadAndDownloadService) DeleteFile(file example.ExaFileUploadAndDownload) (err error) {
	var fileFromDb example.ExaFileUploadAndDownload
	fileFromDb, err = e.FindFile(file.ID)
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = global.GVA_DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}

// EditFileName 编辑文件名或者备注
func (e *FileUploadAndDownloadService) EditFileName(file example.ExaFileUploadAndDownload) (err error) {
	var fileFromDb example.ExaFileUploadAndDownload
	return global.GVA_DB.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.ExaAttachmentCategorySearch
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info request.ExaAttachmentCategorySearch) (list []example.ExaFileUploadAndDownload, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&example.ExaFileUploadAndDownload{})

	if len(info.Keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+info.Keyword+"%")
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

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string, classId int, folder string, uploadType string, uploadPosition string) (file example.ExaFileUploadAndDownload, err error) {
	var filePath, key string
	var uploadErr error

	// 优先使用默认上传云（ExternalLinkDomain.defaultUpload 或 isDefault）
	cloudSvc := clientService.CloudStorageService{}
	defaultDomain, _ := cloudSvc.GetDefaultUploadDomain()
	if defaultDomain != nil {
		// 走默认云端上传，返回相对路径（不含域名）
		filePath, key, uploadErr = cloudSvc.UploadFileToCloud(*defaultDomain, header, folder)
	} else {
		// 回退到 config.yaml 的 oss-type
		oss := upload.NewOss()
		filePath, key, uploadErr = upload.UploadFileToFolder(oss, header, folder)
	}

	if uploadErr != nil {
		return file, uploadErr
	}

	// 同步上传到其他标记了 sync_upload 的云（异步，不阻塞主流程）
	go syncUploadToClouds(cloudSvc, header, folder)

	mediaCategory := inferMediaCategory(uploadType, folder, filePath)
	normalizedPosition := normalizeUploadPosition(uploadPosition)
	if normalizedPosition == "" {
		normalizedPosition = inferUploadPosition(uploadType, folder, filePath)
	}
	s := strings.Split(header.Filename, ".")
	f := example.ExaFileUploadAndDownload{
		Url:          filePath,
		ThumbnailURL: filePath,
		Categories:   mediaCategory,
		Position:     normalizedPosition,
		Keywords:     "",
		Size:         header.Size,
		Name:         header.Filename,
		ClassId:      classId,
		Tag:          s[len(s)-1],
		Key:          key,
	}
	if noSave == "0" {
		// 检查是否已存在相同key的记录
		var existingFile example.ExaFileUploadAndDownload
		err = global.GVA_DB.Where(&example.ExaFileUploadAndDownload{Key: key}).First(&existingFile).Error
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

func (e *FileUploadAndDownloadService) ImportURL(file *[]example.ExaFileUploadAndDownload) error {
	if file == nil {
		return nil
	}

	for i := range *file {
		item := &(*file)[i]
		if strings.TrimSpace(item.ThumbnailURL) == "" {
			item.ThumbnailURL = item.Url
		}
		if strings.TrimSpace(item.Categories) == "" {
			item.Categories = inferMediaCategory("", "", item.Url)
		}
		if strings.TrimSpace(item.Position) == "" {
			item.Position = inferUploadPosition("", "", item.Url)
		}
		if item.Keywords == "" {
			item.Keywords = ""
		}
		if item.Size < 0 {
			item.Size = 0
		}
	}

	return global.GVA_DB.Create(&file).Error
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

// syncUploadToClouds 异步同步上传到所有标记了 sync_upload 的云
func syncUploadToClouds(cloudSvc clientService.CloudStorageService, header *multipart.FileHeader, folder string) {
	syncDomains, err := cloudSvc.GetSyncUploadDomains()
	if err != nil {
		global.GVA_LOG.Error("获取同步上传云列表失败", zap.Error(err))
		return
	}
	if len(syncDomains) == 0 {
		return
	}

	for _, domain := range syncDomains {
		d := domain
		go func() {
			_, _, uploadErr := cloudSvc.UploadFileToCloud(d, header, folder)
			if uploadErr != nil {
				global.GVA_LOG.Error("同步上传失败",
					zap.String("cloud_name", d.Name),
					zap.String("cloud_type", d.CloudType),
					zap.Error(uploadErr))
			} else {
				global.GVA_LOG.Info("同步上传成功",
					zap.String("cloud_name", d.Name),
					zap.String("cloud_type", d.CloudType))
			}
		}()
	}
}
