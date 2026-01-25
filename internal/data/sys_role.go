package data

type Role struct {
	BaseModel
	Name   string `gorm:"size:50;not null;comment:角色名称(汉字)"`         // e.g., 超级管理员
	Key    string `gorm:"uniqueIndex;size:50;not null;comment:角色代码"` // e.g., admin, editor
	Sort   int    `gorm:"default:0;comment:排序"`
	Status int    `gorm:"default:1;comment:状态"` //

	// 多对多关联：一个角色拥有多个菜单权限
	// gorm 会自动创建中间表 sys_role_menus
	// Menus []*Menu `gorm:"many2many:sys_role_menus;"`
}

func (Role) TableName() string {
	return "sys_role"
}