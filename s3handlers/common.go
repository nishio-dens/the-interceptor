package s3handlers

import (
	"the-interceptor/db"
	"errors"
)

func GetInterceptorBucket(name string) (*db.InterceptorBucket, error) {
	var bucket db.InterceptorBucket
	if db.Conn.Where("name = ?", name).First(&bucket).RecordNotFound() {
		return nil, errors.New("Record Not Found")
	}
	return &bucket, nil
}