package data

import (
	// "time"

	"gorm.io/gorm"
)

type BaseModel struct {
	// ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	// CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt time.Time `json:"updatedAt"`
	// DeletedAt gorm.DeletedAt `json:"deletedAt"`
	gorm.Model
}