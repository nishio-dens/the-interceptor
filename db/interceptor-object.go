package db

import "time"

type InterceptorObject struct {
	Id              int64      `gorm:"column:id;primary_key"`
	VirtualBucketId int64      `gorm:"column:virtual_bucket_id"`
	RealBucketId    int64      `gorm:"column:real_bucket_id"`
	Key             string     `gorm:"column:key"`
	Size            int        `gorm:"column:size"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
}
