package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type InterceptorBucket struct {
	gorm.Model

	Id            int64      `gorm:"column:id;primary_key"`
	Name          string     `gorm:"column:name"`
	ReadBucketId  int64      `gorm:"column:read_bucket_id"`
	WriteBucketId int64      `gorm:"column:write_bucket_id"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`

	// Relations
	ReadBucket  S3Bucket `gorm:"ForeignKey:read_bucket_id"`
	WriteBucket S3Bucket `gorm:"ForeignKey:write_bucket_id"`
}

func (b *InterceptorBucket) GetReadBucket() *S3Bucket {
	var bucket S3Bucket
	Conn.Model(b).Related(&bucket, "ReadBucket")

	return &bucket
}

func (b *InterceptorBucket) GetWriteBucket() *S3Bucket {
	var bucket S3Bucket
	Conn.Model(b).Related(&bucket, "WriteBucket")

	return &bucket
}
