package db

import "time"

type User struct {
	Id           int64      `gorm:"column:id;primary_key"`
	Username     string     `gorm:"column:username"`
	AccessKey    string     `gorm:"column:access_key"`
	AccessSecret string     `gorm:"column:access_secret"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}
