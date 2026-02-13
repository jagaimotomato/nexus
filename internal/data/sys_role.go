package data

import "gorm.io/gorm"

type Role struct {
	BaseModel
	Name   string `gorm:"size:50;not null;comment:角色名称(汉字)" json:"name"`         // e.g., 超级管理员
	Key    string `gorm:"uniqueIndex;size:50;not null;comment:角色代码" json:"key"` // e.g., admin, editor
	Sort   int    `gorm:"default:0;comment:排序" json:"sort"`
	Status int    `gorm:"default:1;comment:状态" json:"status"` //

	// 多对多关联：一个角色拥有多个菜单权限
	// gorm 会自动创建中间表 sys_role_menus
	Menus []*Menu `gorm:"many2many:sys_role_menus;" json:"menus"`
}

func (Role) TableName() string {
	return "sys_role"
}

func GetRoleList(page, pageSize int) ([]*Role, int64, error) {
	var roles []*Role
	var total int64
	err := DB.Model(&Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error
	return roles, total, err
}

func GetAllRoles() ([]*Role, error) {
	var roles []*Role
	err := DB.Find(&roles).Error
	return roles, err
}

func CreateRole(role *Role) error {
	return DB.Create(role).Error
}

func UpdateRole(id uint, updates map[string]interface{}, menuIds []uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 更新基本信息
		if err := tx.Model(&Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		// 如果提供了菜单ID，更新关联
		if menuIds != nil {
			var role Role
			role.ID = id
			var menus []Menu
			if len(menuIds) > 0 {
				if err := tx.Find(&menus, menuIds).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(&role).Association("Menus").Replace(menus); err != nil {
				return err
			}
		}
		return nil
	})
}

func DeleteRole(id uint) error {
	return DB.Delete(&Role{}).Where("id = ?", id).Error
}
