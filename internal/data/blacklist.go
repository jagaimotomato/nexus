package data

import (
	"time"

	"gorm.io/gorm"
)

type Blacklist struct {
	BaseModel
	IP        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"ip"` // IP 地址
	Reason string `gorm:"type:varchar(255);" json:"reason"`
	ExpiresAt *time.Time `gorm:"index" json:"expiresAt"` // 过期时间 null 表示永久
	IsActive bool `gorm:"default:true" json:"isActive"`
}

func (Blacklist) TableName() string {
	return "blacklist"
}

func InitBlacklist(db *gorm.DB) {
	err := db.AutoMigrate(&Blacklist{})
	if err != nil {
		panic("黑名单表初始化失败" + err.Error())
	}
}