package model

import (
	"time"
)

type TelemetryLog struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	StartTime     time.Time `gorm:"column:start_time;not null"`
	EndTime       time.Time `gorm:"column:end_time;not null"`
	LocationCount int       `gorm:"column:location_count;not null;default:0"`
	IsPublic      bool      `gorm:"column:is_public;not null;default:false"`
	CreatedAt     time.Time `gorm:"column:created_at;"`
	UpdatedAt     time.Time `gorm:"column:updated_at;"`
}
