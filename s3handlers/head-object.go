package s3handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"net/http"
	"the-interceptor/api"
	"the-interceptor/db"
	"the-interceptor/s3client"
)

type headObjectResponseResult struct {
	Key           string
	ContentLength int64
	ContentType   string
	Error         error
	IsReadBucket  bool
}

/**
HEAD Object
see: http://docs.aws.amazon.com/ja_jp/AmazonS3/latest/API/RESTObjectHEAD.html
*/

func HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	bucket, err := GetInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}
	key := v["object"]
	ch := make(chan headObjectResponseResult)
	readBucket := bucket.GetReadBucket()
	writeBucket := bucket.GetWriteBucket()
	requestBuckets := []*db.S3Bucket{readBucket, writeBucket}
	go headObject(readBucket, key, true, ch)
	go headObject(writeBucket, key, false, ch)

	var readResult headObjectResponseResult
	var writeResult headObjectResponseResult
	for range requestBuckets {
		t := <-ch
		if t.IsReadBucket {
			readResult = t
		} else {
			writeResult = t
		}
	}

	if readResult.Error != nil && writeResult.Error != nil {
		SendNoSuchKeyError(key, w, r)
	} else if readResult.Error == nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", readResult.ContentLength))
		api.SendSuccess(w, []byte(""), readResult.ContentType)
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", readResult.ContentLength))
		api.SendSuccess(w, []byte(""), writeResult.ContentType)
	}
}

func headObject(bucket *db.S3Bucket, key string, isReadBucket bool, ch chan<- headObjectResponseResult) {
	// TODO: Support Range
	// TODO: Support PartNumber
	client := s3client.GetS3Client(bucket)
	resp, err := client.HeadObject(&s3sdk.HeadObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		ch <- headObjectResponseResult{
			Key:          key,
			Error:        err,
			IsReadBucket: isReadBucket,
		}
		return
	}

	ch <- headObjectResponseResult{
		Key:           key,
		ContentType:   *resp.ContentType,
		ContentLength: *resp.ContentLength,
		Error:         nil,
		IsReadBucket:  isReadBucket,
	}
}
