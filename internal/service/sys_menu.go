package service

import (
	"errors"
	"nexus/internal/data"
)

type MenuService struct{}

func CreateMenu(req *data.Menu) error {
	return data.CreateMenu(req)
}

func UpdateMenu(id uint, req *data.Menu) error {
	updates := map[string]interface{} {
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
	return data.UpdateMenu(id, updates)
}

func DeleteMenu(id uint) error {
	menu, err := data.GetMenuByID(id)
	if err != nil {
		return  err
	}
	if menu == nil {
		return errors.New("菜单不存在")
	}

	return data.DeleteMenu(id)
}

func GetMenuTree() ([]*data.Menu, error) {
	allMenus, err := data.GetAllMenus()
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