package db

import "time"

type S3Bucket struct {
	Id                 int64      `gorm:"column:id;primary_key"`
	BucketName         string     `gorm:"column:bucket_name"`
	BucketAccessKey    string     `gorm:"column:bucket_access_key"`
	BucketAccessSecret string     `gorm:"column:bucket_access_secret"`
	BucketRegion       string     `gorm:"column:bucket_region"`
	BucketDisableSsl   bool       `gorm:"column:bucket_disable_ssl"`
	BucketEndpoint     string     `gorm:"column:bucket_endpoint"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at"`
}
