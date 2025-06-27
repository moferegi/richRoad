package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type CategoryService struct {
}

// CreateCategory 创建商品分类记录
// Author [piexlmax](https://github.com/piexlmax)
func (categoryService *CategoryService) CreateCategory(category *shop.Category) (err error) {
	err = global.GVA_DB.Create(category).Error
	return err
}

// DeleteCategory 删除商品分类记录
// Author [piexlmax](https://github.com/piexlmax)
func (categoryService *CategoryService) DeleteCategory(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Category{}, "id = ?", ID).Error
	return err
}

// DeleteCategoryByIds 批量删除商品分类记录
// Author [piexlmax](https://github.com/piexlmax)
func (categoryService *CategoryService) DeleteCategoryByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Category{}, "id in ?", IDs).Error
	return err
}

// UpdateCategory 更新商品分类记录
// Author [piexlmax](https://github.com/piexlmax)
func (categoryService *CategoryService) UpdateCategory(category shop.Category) (err error) {
	err = global.GVA_DB.Model(&shop.Category{}).Where("id = ?", category.ID).Updates(&category).Error
	return err
}

// GetCategory 根据ID获取商品分类记录
// Author [piexlmax](https://github.com/piexlmax)
func (categoryService *CategoryService) GetCategory(ID string) (category shop.Category, err error) {
	err = global.GVA_DB.Preload("Goods").Where("id = ?", ID).First(&category).Error
	return
}

// GetCategoryInfoList 分页获取商品分类记录
// Author [piexlmax](https://github.com/piexlmax)
func (categoryService *CategoryService) GetCategoryInfoList(info shopReq.CategorySearch) (list []shop.Category, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Category{})
	var categorys []shop.Category
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&categorys).Error
	return makeCategoryTree(categorys), total, err
}

func makeCategoryTree(list []shop.Category) []shop.Category {
	// 1. 创建一个map，用于存放所有的分类
	categoryMap := make(map[uint]*shop.Category)
	// 2. 遍历所有的分类，将每个分类的id作为key，分类本身作为value，存放到map中
	for i := range list {
		categoryMap[list[i].ID] = &list[i]
	}
	// 3. 遍历所有的分类，将每个分类添加到其父分类的children字段中
	for i := range list {
		// 3.1 如果当前分类是一级分类，直接跳过
		if list[i].ParentID == 0 {
			continue
		}
		// 3.2 如果当前分类是二级分类，找到其父分类，将其添加到父分类的children字段中
		parent, ok := categoryMap[list[i].ParentID]
		if ok {
			parent.Children = append(parent.Children, list[i])
		}
	}
	// 4. 创建一个切片，用于存放所有的一级分类
	var categoryTree []shop.Category
	// 5. 遍历所有的分类，将所有的一级分类添加到切片中
	for i := range list {
		if list[i].ParentID == 0 {
			categoryTree = append(categoryTree, list[i])
		}
	}
	return categoryTree
}

func (categoryService *CategoryService) GetCategoryMobile(parentID int) (list []shop.Category, err error) {
	err = global.GVA_DB.Find(&list, "parent_id = ?", parentID).Error
	return
}

func (categoryService *CategoryService) GetChildrenCategoryAndProduct(parentID int) (list []shop.Category, err error) {
	// 查询当前分类及其商品
	var currentCategory shop.Category
	err = global.GVA_DB.Preload("Goods").Where("id = ?", parentID).First(&currentCategory).Error
	if err != nil {
		return nil, err
	}

	// 将当前分类添加到结果列表
	list = append(list, currentCategory)

	// 查询子分类及其商品
	var childCategories []shop.Category
	err = global.GVA_DB.Preload("Goods").Where("parent_id = ?", parentID).Find(&childCategories).Error
	if err != nil {
		return nil, err
	}

	// 将子分类添加到结果列表
	list = append(list, childCategories...)

	return list, nil
}
