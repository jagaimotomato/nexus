package data

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:100;not null;comment:用户名"`
	Name     string `gorm:"name;comment:姓名"`
	Password string `gorm:"size:100;not null;comment:密码"`
	Avatar   string `gorm:"size:255;comment:头像"`
	Phone    string `gorm:"size:20;comment:手机号"`
	Status   int    `gorm:"default:1;comment:状态 1:启用 0:停用"`
	Email    string `gorm:"size:100;comment:邮箱"`
	// 多对多关联：一个用户拥有多个角色
	// gorm 会自动创建中间表 sys_user_roles
	Roles []*Role `gorm:"many2many:sys_user_role;"`
}

func (User) TableName() string {
	return "sys_user"
}
