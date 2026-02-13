package data

import (
	"errors"

	"gorm.io/gorm"
)

type Menu struct {
	BaseModel
	Pid       uint    `gorm:"default:0;index;comment:父菜单ID" json:"pid"`
	Name      string  `gorm:"size:100;comment:路由名称" json:"name"`
	Title     string  `gorm:"size:100;comment:显示标题" json:"title"`
	Path      string  `gorm:"size:100;comment:路由路径" json:"path"`
	Component string  `gorm:"size:255;comment:组件路径" json:"component"`
	Icon      string  `gorm:"size:100;comment:图标" json:"icon"`
	Sort      int     `gorm:"default:0;comment:排序" json:"sort"`
	Type      int     `gorm:"default:0;comment:1:目录 2:菜单 3:按钮" json:"type"`
	Hidden    bool    `gorm:"default:false;comment:是否隐藏" json:"hidden"`
	KeepAlive bool    `gorm:"default:false;comment:是否缓存" json:"keepAlive"`
	Perms string `gorm:"size:100;comment:权限标识(如 sys:user:add)" json:"perms"`
	Redirect string  `gorm:"size:255;comment:重定向路径（目录专用）" json:"redirect"`
	Children  []*Menu `gorm:"-" json:"children"`
}

func (Menu) TableName() string {
	return "sys_menu"
}

type MenuRepo struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{db: db}
}

func (r *MenuRepo) GetAllMenus() ([]*Menu, error) {
	var menus []*Menu
	err := r.db.Order("sort asc").Find(&menus).Error
	return menus, err
}

func (r *MenuRepo) GetMenuByID(id uint) (*Menu, error) {
	var menu Menu
	err := r.db.First(&menu, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &menu, err
}

func (r *MenuRepo) CreateMenu(menu *Menu) error {
	return r.db.Create(menu).Error
}

func (r *MenuRepo) UpdateMenu(id uint, updates map[string]interface{}) error {
	return r.db.Model(&Menu{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MenuRepo) DeleteMenu(id uint) error {
	return r.db.Delete(&Menu{}, id).Error
}