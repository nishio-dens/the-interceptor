package s3handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"the-interceptor/api"
	"the-interceptor/db"
	"the-interceptor/s3client"
	"time"
)

type getObjectResponseResult struct {
	Key           string
	ContentLength int64
	ContentType   string
	LastModified  string
	ETag          string
	Body          []byte
	Error         error
	IsReadBucket  bool
}

/**
GET Object
see: http://docs.aws.amazon.com/ja_jp/AmazonS3/latest/API/RESTObjectGET.html
*/
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Future Work
	// TODO: Support RangeGet
	// TODO: Support 403 Forbidden (Authorization)
	v := mux.Vars(r)
	bucket, err := GetInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}
	key := v["object"]
	ch := make(chan getObjectResponseResult)
	readBucket := bucket.GetReadBucket()
	writeBucket := bucket.GetWriteBucket()
	requestBuckets := []*db.S3Bucket{readBucket, writeBucket}
	go getObject(readBucket, key, true, ch)
	go getObject(writeBucket, key, false, ch)

	var readResult getObjectResponseResult
	var writeResult getObjectResponseResult
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
		writeGetObjectResultHeader(w, &readResult)
		api.SendSuccess(w, readResult.Body, readResult.ContentType)
	} else {
		writeGetObjectResultHeader(w, &writeResult)
		api.SendSuccess(w, writeResult.Body, writeResult.ContentType)
	}
}

func getObject(bucket *db.S3Bucket, key string, isReadBucket bool, ch chan<- getObjectResponseResult) {
	// TODO: Support range Get
	client := s3client.GetS3Client(bucket)
	resp, err := client.GetObject(&s3sdk.GetObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		ch <- getObjectResponseResult{
			Key:          key,
			Error:        err,
			IsReadBucket: isReadBucket,
		}
		return
	}

	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		ch <- getObjectResponseResult{
			Key:          key,
			Error:        err,
			IsReadBucket: isReadBucket,
		}
		return
	}
	ch <- getObjectResponseResult{
		Key:           key,
		ContentLength: *resp.ContentLength,
		ContentType:   *resp.ContentType,
		LastModified:  (*resp.LastModified).Format(time.RFC1123),
		ETag:          *resp.ETag,
		Body:          b,
		Error:         nil,
		IsReadBucket:  isReadBucket,
	}
}

func writeGetObjectResultHeader(w http.ResponseWriter, result *getObjectResponseResult) {
	w.Header().Set("Content-Length", fmt.Sprintf("%d", result.ContentLength))
	w.Header().Set("Last-Modified", result.LastModified)
	w.Header().Set("ETag", result.ETag)
}
