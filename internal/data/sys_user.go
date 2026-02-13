package data

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

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

func (r *UserRepo) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).Preload("Roles").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserList(page, pageSize int) ([]*User, int64, error) {
	var users []*User
	var total int64
	err := r.db.Model(&User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Preload("Roles").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *UserRepo) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) UpdateUser(id uint, updates map[string]interface{}, roleIds []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		if roleIds != nil {
			var user User
			user.ID = id
			var roles []Role
			if len(roleIds) > 0 {
				if err := tx.Find(&roles, roleIds).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepo) DeleteUser(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
