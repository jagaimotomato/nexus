package service

import (
	"errors"
	"nexus/internal/data"
)

type MenuService struct {
	repo *data.MenuRepo
}

func NewMenuService(d *data.Data) *MenuService {
	return &MenuService{repo: data.NewMenuRepo(d.DB)}
}

type MenuInput struct {
	Pid       uint   `json:"pid"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	Type      int    `json:"type"`
	Hidden    bool   `json:"hidden"`
	KeepAlive bool   `json:"keepAlive"`
	Perms     string `json:"perms"`
	Redirect  string `json:"redirect"`
}

func (s *MenuService) CreateMenu(req *MenuInput) error {
	menu := &data.Menu{
		Pid:       req.Pid,
		Name:      req.Name,
		Title:     req.Title,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Type:      req.Type,
		Hidden:    req.Hidden,
		KeepAlive: req.KeepAlive,
		Perms:     req.Perms,
		Redirect:  req.Redirect,
	}
	return s.repo.CreateMenu(menu)
}

func (s *MenuService) UpdateMenu(id uint, req *MenuInput) error {
	updates := map[string]interface{}{
		"pid":        req.Pid,
		"name":       req.Name,
		"title":      req.Title,
		"path":       req.Path,
		"component":  req.Component,
		"icon":       req.Icon,
		"sort":       req.Sort,
		"type":       req.Type,
		"hidden":     req.Hidden,
		"keep_alive": req.KeepAlive,
		"perms":      req.Perms,
		"redirect":   req.Redirect,
	}
	return s.repo.UpdateMenu(id, updates)
}

func (s *MenuService) DeleteMenu(id uint) error {
	menu, err := s.repo.GetMenuByID(id)
	if err != nil {
		return err
	}
	if menu == nil {
		return errors.New("菜单不存在")
	}

	return s.repo.DeleteMenu(id)
}

func (s *MenuService) GetMenuTree() ([]*data.Menu, error) {
	allMenus, err := s.repo.GetAllMenus()
	if err != nil {
		return nil, err
	}
	return buildTree(allMenus, 0), nil
}

func buildTree(menus []*data.Menu, pid uint) []*data.Menu {
	var tree []*data.Menu
	for _, node := range menus {
		if node.Pid == pid {
			children := buildTree(menus, node.ID)
			node.Children = children
			tree = append(tree, node)
		}
	}
	return tree
}
