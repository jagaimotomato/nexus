package service

import (
	"nexus/internal/data"
)

type RoleService struct {
	repo *data.RoleRepo
}

func NewRoleService(d *data.Data) *RoleService {
	return &RoleService{repo: data.NewRoleRepo(d.DB)}
}

type RoleInput struct {
	Name    string `json:"name" binding:"required"`
	Key     string `json:"key" binding:"required"`
	Sort    int    `json:"sort"`
	Status  int    `json:"status"`
	MenuIds []uint `json:"menuIds"` // 关联的菜单ID列表
}

func (s *RoleService) CreateRole(req *RoleInput) error {
	role := &data.Role{
		Name:   req.Name,
		Key:    req.Key,
		Sort:   req.Sort,
		Status: req.Status,
	}
	// 创建角色基础信息
	if err := s.repo.CreateRole(role); err != nil {
		return err
	}
	// 关联菜单权限
	return s.repo.UpdateRole(role.ID, map[string]interface{}{}, req.MenuIds)
}

func (s *RoleService) UpdateRole(id uint, req *RoleInput) error {
	updates := map[string]interface{}{
		"name":   req.Name,
		"key":    req.Key,
		"sort":   req.Sort,
		"status": req.Status,
	}
	return s.repo.UpdateRole(id, updates, req.MenuIds)
}

func (s *RoleService) GetRoleList(page, pageSize int) ([]*data.Role, int64, error) {
	return s.repo.GetRoleList(page, pageSize)
}

func (s *RoleService) GetAllRoles() ([]*data.Role, error) {
	return s.repo.GetAllRoles()
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.repo.DeleteRole(id)
}
