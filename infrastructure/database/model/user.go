package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"column:username"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
