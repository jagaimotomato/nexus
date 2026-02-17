package service

import (
	"errors"
	"nexus/internal/data"
	"nexus/internal/utils"
)

type UserService struct {
	repo *data.UserRepo
}

func NewUserService(d *data.Data) *UserService {
	return &UserService{repo: data.NewUserRepo(d.DB)}
}

type UserInput struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	RoleIds  []uint `json:"roleIds"`
}

func (s *UserService) CreateUser(req *UserInput) error {
	exist, _ := s.repo.GetUserByUsername(req.Username)
	if exist != nil {
		return errors.New("用户已存在")
	}
	hashPwd, _ := utils.HashPassword(req.Password)
	user := &data.User{
		Username: req.Username,
		Name:     req.Name,
		Password: hashPwd,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}
	// GORM 的关联创建通常需要先查出 Role 对象，或者使用 Association 模式
	// 这里简化处理，先创建用户，再更新角色
	// 实际项目中建议在 data 层处理事务
	return s.repo.UpdateUser(user.ID, map[string]interface{}{}, req.RoleIds)
}

func (s *UserService) UpdateUser(id uint, req *UserInput) error {
	updates := map[string]interface{}{
		"username": req.Username,
		"name":     req.Name,
		"phone":    req.Phone,
		"email":    req.Email,
		"status":   req.Status,
	}
	if req.Password != "" {
		hashPwd, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		updates["password"] = hashPwd
	}
	return s.repo.UpdateUser(id, updates, req.RoleIds)
}

func (s *UserService) GetUserList(page, pageSize int) ([]*data.User, int64, error) {
	return s.repo.GetUserList(page, pageSize)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}
