package service

import (
	"nexus/internal/data"
	"nexus/internal/logger"
	"nexus/internal/utils"

	"go.uber.org/zap"
)

func InitData(userRepo *data.UserRepo, roleRepo *data.RoleRepo) {
	logger.Log.Info("正在检查并初始化基础数据...")

	// 1. 初始化角色
	adminRole := data.Role{
		Name:   "超级管理员",
		Key:    "admin",
		Sort:   0,
		Status: 1,
	}

	if err := roleRepo.FirstOrCreateByKey("admin", &adminRole); err != nil {
		logger.Log.Error("初始化角色失败", zap.Error(err))
		return
	}

	user, err := userRepo.GetUserByUsername("admin")
	if err != nil {
		logger.Log.Error("查询管理员账号失败", zap.Error(err))
		return
	}
	if user == nil {
		hashPwd, _ := utils.HashPassword("123456")

		newUser := data.User{
			Username: "admin",
			Name:     "系统管理员",
			Password: hashPwd,
			Status:   1,
			Roles:    []*data.Role{&adminRole},
		}

		if err := userRepo.CreateUser(&newUser); err != nil {
			logger.Log.Error("初始化管理员账号失败", zap.Error(err))
		} else {
			logger.Log.Info("初始化管理员账号成功: admin / 123456")
		}
	} else {
		logger.Log.Info("管理员账号已存在，跳过初始化")
	}
}
