package service

import (
	"nexus/internal/data"
	"nexus/internal/logger"
	"nexus/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitData() {
	logger.Log.Info("正在检查并初始化基础数据...")

	// 1. 初始化角色
	adminRole := data.Role{
		Name:   "超级管理员",
		Key:    "admin",
		Sort:   0,
		Status: 1,
	}
	
	// FirstOrCreate: 如果按 Key 找到了就不创建，没找到就创建
	if err := data.DB.Where(data.Role{Key: "admin"}).FirstOrCreate(&adminRole).Error; err != nil {
		logger.Log.Error("初始化角色失败", zap.Error(err))
		return
	}

	// 2. 初始化管理员用户
	var user data.User
	err := data.DB.Where(data.User{Username: "admin"}).Preload("Roles").First(&user).Error
	
	// 如果用户不存在，则创建
	if err == gorm.ErrRecordNotFound {
		// A. 密码加密
		hashPwd, _ := utils.HashPassword("123456") // 默认密码

		// B. 构建用户对象
		newUser := data.User{
			Username: "admin",
			Name:     "系统管理员",
			Password: hashPwd,
			Status:   1, // 启用
			Roles:    []*data.Role{&adminRole}, // 关联角色
		}

		// C. 插入数据库
		if err := data.DB.Create(&newUser).Error; err != nil {
			logger.Log.Error("初始化管理员账号失败", zap.Error(err))
		} else {
			logger.Log.Info("初始化管理员账号成功: admin / 123456")
		}
	} else if err != nil {
		logger.Log.Error("查询管理员账号失败", zap.Error(err))
	} else {
		logger.Log.Info("管理员账号已存在，跳过初始化")
	}
}