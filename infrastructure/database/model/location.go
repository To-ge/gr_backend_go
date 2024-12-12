package model

import "time"

type Location struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Latitude  float64   `gorm:"column:latitude"`
	Longitude float64   `gorm:"column:longitude"`
	Altitude  float64   `gorm:"column:altitude"`
	Speed     float64   `gorm:"column:speed"`
	Timestamp int64     `gorm:"column:timestamp"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
