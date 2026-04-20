package shop

import (
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type PopupService struct{}

// CreatePopup 创建弹窗
func (s *PopupService) CreatePopup(popup *shop.Popup) (err error) {
	err = global.GVA_DB.Create(popup).Error
	return err
}

// DeletePopup 删除弹窗
func (s *PopupService) DeletePopup(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Popup{}, "id = ?", ID).Error
	return err
}

// UpdatePopup 更新弹窗
func (s *PopupService) UpdatePopup(popup shop.Popup) (err error) {
	err = global.GVA_DB.Model(&shop.Popup{}).Where("id = ?", popup.ID).Updates(&popup).Error
	return err
}

// GetPopup 根据ID获取弹窗
func (s *PopupService) GetPopup(ID string) (popup shop.Popup, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&popup).Error
	return
}

// GetPopupList 分页获取弹窗列表(后台管理)
func (s *PopupService) GetPopupList(info shopReq.PopupSearch) (list []shop.Popup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.Popup{})

	if info.Position != "" {
		db = db.Where("position = ?", info.Position)
	}
	if info.ClientType != "" {
		db = db.Where("client_type = ? OR client_type = 'all'", info.ClientType)
	}
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
	}
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 排序支持
	orderClause := "sort ASC"
	if info.OrderBy != "" {
		allowedCols := map[string]bool{"created_at": true}
		if allowedCols[info.OrderBy] {
			dir := "asc"
			if info.OrderDir == "desc" {
				dir = "desc"
			}
			orderClause = info.OrderBy + " " + dir
		}
	}
	err = db.Order(orderClause).Find(&list).Error
	return
}

// GetActivePopups 获取当前有效弹窗(客户端用)
func (s *PopupService) GetActivePopups(position string, clientType string, page string) (list []shop.Popup, err error) {
	now := time.Now()
	db := global.GVA_DB.Where("is_enabled = ?", true)
	if clientType != "" {
		db = db.Where("client_type = ? OR client_type = 'all'", clientType)
	}
	db = db.Where("(start_time IS NULL OR start_time <= ?) AND (end_time IS NULL OR end_time >= ?)", now, now)
	err = db.Order("sort ASC").Find(&list).Error
	if err != nil {
		return
	}

	// 按页面路径过滤：优先使用 pages 字段，兼容旧 position 字段
	if page != "" || position != "" {
		matchPath := page
		if matchPath == "" {
			matchPath = position
		}
		var filtered []shop.Popup
		for _, p := range list {
			if p.Pages != "" {
				// 新逻辑：pages 字段逗号分隔匹配
				pages := splitPages(p.Pages)
				for _, pg := range pages {
					if pg == matchPath || pg == "all" {
						filtered = append(filtered, p)
						break
					}
				}
			} else if p.Position != "" {
				// 兼容旧逻辑：position="home" 匹配首页路径
				if p.Position == matchPath || p.Position == "all" || isPositionMatch(p.Position, matchPath) {
					filtered = append(filtered, p)
				}
			} else {
				// 无页面限制，全部匹配
				filtered = append(filtered, p)
			}
		}
		list = filtered
	}
	return
}

// isPositionMatch 旧版 position 值与新版页面路径的兼容映射
func isPositionMatch(position string, pagePath string) bool {
	positionMap := map[string]string{
		"home":     "/pages/tabBar/index",
		"category": "/pages/tabBar/category",
		"cart":     "/pages/tabBar/cart",
		"mine":     "/pages/tabBar/mine",
	}
	if mapped, ok := positionMap[position]; ok {
		return mapped == pagePath
	}
	return false
}

// splitPages 将逗号分隔的页面路径拆分为切片
func splitPages(pages string) []string {
	var result []string
	for _, p := range strings.Split(pages, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
