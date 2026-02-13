package data

import "gorm.io/gorm"

type Role struct {
	BaseModel
	Name   string `gorm:"size:50;not null;comment:角色名称(汉字)" json:"name"`        // e.g., 超级管理员
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

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetRoleList(page, pageSize int) ([]*Role, int64, error) {
	var roles []*Role
	var total int64
	err := r.db.Model(&Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error
	return roles, total, err
}

func (r *RoleRepo) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) CreateRole(role *Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepo) FirstOrCreateByKey(key string, role *Role) error {
	return r.db.Where(Role{Key: key}).FirstOrCreate(role).Error
}

func (r *RoleRepo) UpdateRole(id uint, updates map[string]interface{}, menuIds []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
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

func (r *RoleRepo) DeleteRole(id uint) error {
	return r.db.Delete(&Role{}, id).Error
}
