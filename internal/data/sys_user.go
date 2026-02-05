package data

import (
	"errors"

	"gorm.io/gorm"
)
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:100;not null;comment:用户名" json:"username"`
	Name     string `gorm:"name;comment:姓名" json:"name"`
	Password string `gorm:"size:100;not null;comment:密码" json:"password"`
	Avatar   string `gorm:"size:255;comment:头像" json:"avatar"`
	Phone    string `gorm:"size:20;comment:手机号" json:"phone"`
	Status   int    `gorm:"default:1;comment:状态 1:启用 0:停用" json:"status"`
	Email    string `gorm:"size:100;comment:邮箱" json:"email"`
	// 多对多关联：一个用户拥有多个角色
	// gorm 会自动创建中间表 sys_user_roles
	Roles []*Role `gorm:"many2many:sys_user_role;" json:"roles"`
}

func (User) TableName() string {
	return "sys_user"
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).Preload("Roles").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}