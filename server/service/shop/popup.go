package shop

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type PopupService struct{}

var popupUniPresetPaths = []string{
	"/pages/tabBar/index",
	"/pages/tabBar/clothes/index",
	"/pages/tabBar/my/index",
	"/pages/goodsDetails/goodsDetails",
	"/pages/cart/index",
	"/pages/order/order",
	"/pages/presale/list",
	"/pages/tryon/index",
	"/pages/tryon/history",
	"/pages/user/login",
}

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
		matchPath := strings.TrimSpace(page)
		if matchPath != "" {
			matchPath = normalizePagePath(matchPath)
		} else {
			matchPath = strings.TrimSpace(position)
		}
		var filtered []shop.Popup
		for _, p := range list {
			if p.Pages != "" {
				// 新逻辑：pages 字段逗号分隔匹配
				pages := splitPages(p.Pages)
				for _, pg := range pages {
					if pg == "all" || pg == matchPath {
						filtered = append(filtered, p)
						break
					}
				}
			} else if p.Position != "" {
				// 兼容旧逻辑：position="home" 匹配首页路径
				positionValue := strings.TrimSpace(p.Position)
				if positionValue == matchPath || positionValue == "all" || isPositionMatch(positionValue, matchPath) {
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

// GetPopupPagePathOptions 获取弹窗可选页面路径
func (s *PopupService) GetPopupPagePathOptions(clientType string) (list []string, err error) {
	pathSet := map[string]struct{}{"all": {}}
	addPathsToSet(pathSet, popupPresetPathsByClientType(clientType))

	if shouldLoadUniPagePathOptions(clientType) {
		addPathsToSet(pathSet, loadUniPagePathsFromPagesJSON())
	}

	db := global.GVA_DB.Model(&shop.Popup{}).Select("pages")
	normalizedClientType := strings.ToLower(strings.TrimSpace(clientType))
	if normalizedClientType == "uni" || normalizedClientType == "web" {
		db = db.Where("client_type = ? OR client_type = 'all'", normalizedClientType)
	}

	var records []struct {
		Pages string `gorm:"column:pages"`
	}
	if err = db.Where("pages <> ''").Find(&records).Error; err != nil {
		return nil, err
	}

	for _, row := range records {
		addPathsToSet(pathSet, splitPages(row.Pages))
	}

	list = buildPagePathList(pathSet)
	return
}

// isPositionMatch 旧版 position 值与新版页面路径的兼容映射
func isPositionMatch(position string, pagePath string) bool {
	position = strings.TrimSpace(position)
	pagePath = normalizePagePath(pagePath)
	if pagePath == "" {
		return false
	}

	positionMap := map[string][]string{
		"home":     {"/pages/tabBar/index"},
		"category": {"/pages/tabBar/category", "/pages/tabBar/clothes/index"},
		"clothes":  {"/pages/tabBar/clothes/index"},
		"cart":     {"/pages/tabBar/cart", "/pages/cart/index"},
		"mine":     {"/pages/tabBar/mine", "/pages/tabBar/my/index"},
		"my":       {"/pages/tabBar/my/index"},
	}
	if mappedList, ok := positionMap[position]; ok {
		for _, mapped := range mappedList {
			if normalizePagePath(mapped) == pagePath {
				return true
			}
		}
		return false
	}
	return false
}

// splitPages 将逗号分隔的页面路径拆分为切片
func splitPages(pages string) []string {
	var result []string
	for _, p := range strings.Split(pages, ",") {
		normalized := normalizePagePath(extractPagePathFromToken(p))
		if normalized != "" {
			result = append(result, normalized)
		}
	}
	return result
}

// extractPagePathFromToken 提取页面路径，兼容 "名称|路径" 格式
func extractPagePathFromToken(token string) string {
	raw := strings.TrimSpace(token)
	if raw == "" {
		return ""
	}

	parts := strings.SplitN(raw, "|", 2)
	if len(parts) == 2 {
		candidate := strings.TrimSpace(parts[1])
		if candidate != "" {
			return candidate
		}
	}

	return raw
}

// normalizePagePath 统一页面路径格式，兼容手工输入值
func normalizePagePath(pagePath string) string {
	p := strings.TrimSpace(pagePath)
	if p == "" {
		return ""
	}
	if strings.EqualFold(p, "all") {
		return "all"
	}
	p = strings.ReplaceAll(p, "\\", "/")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = strings.TrimRight(p, "/")
	if p == "" {
		return "/"
	}
	return p
}

func shouldLoadUniPagePathOptions(clientType string) bool {
	normalized := strings.ToLower(strings.TrimSpace(clientType))
	return normalized == "" || normalized == "all" || normalized == "uni"
}

func popupPresetPathsByClientType(clientType string) []string {
	normalized := strings.ToLower(strings.TrimSpace(clientType))
	if normalized == "web" {
		return nil
	}
	return popupUniPresetPaths
}

func addPathsToSet(pathSet map[string]struct{}, paths []string) {
	for _, path := range paths {
		normalized := normalizePagePath(path)
		if normalized == "" {
			continue
		}
		pathSet[normalized] = struct{}{}
	}
}

func buildPagePathList(pathSet map[string]struct{}) []string {
	list := make([]string, 0, len(pathSet))
	hasAll := false
	for path := range pathSet {
		if path == "all" {
			hasAll = true
			continue
		}
		list = append(list, path)
	}

	sort.Strings(list)
	if hasAll {
		return append([]string{"all"}, list...)
	}
	return list
}

func loadUniPagePathsFromPagesJSON() []string {
	type uniPageItem struct {
		Path string `json:"path"`
	}
	type uniTabBarItem struct {
		PagePath string `json:"pagePath"`
	}
	type uniPagesConfig struct {
		Pages  []uniPageItem `json:"pages"`
		TabBar struct {
			List []uniTabBarItem `json:"list"`
		} `json:"tabBar"`
	}

	candidates := []string{
		filepath.Join("uni", "pages.json"),
		filepath.Join("..", "uni", "pages.json"),
		filepath.Join("..", "..", "uni", "pages.json"),
	}

	for _, candidate := range candidates {
		data, readErr := os.ReadFile(candidate)
		if readErr != nil {
			continue
		}

		var cfg uniPagesConfig
		if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
			continue
		}

		var paths []string
		for _, page := range cfg.Pages {
			paths = append(paths, page.Path)
		}
		for _, page := range cfg.TabBar.List {
			paths = append(paths, page.PagePath)
		}
		return buildPagePathList(sliceToPathSet(paths))
	}

	return nil
}

func sliceToPathSet(paths []string) map[string]struct{} {
	pathSet := make(map[string]struct{}, len(paths))
	addPathsToSet(pathSet, paths)
	delete(pathSet, "all")
	return pathSet
}
