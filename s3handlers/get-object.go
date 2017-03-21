package s3handlers

import (
	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"net/http"
	"the-interceptor/api"
	"the-interceptor/db"
	"the-interceptor/s3client"
	"io/ioutil"
)

/**
GET Object
see: http://docs.aws.amazon.com/ja_jp/AmazonS3/latest/API/RESTObjectGET.html
*/
func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Future Work
	// TODO: Support RangeGet
	// TODO: Support 403 Forbidden (Authorization)
	// TODO: Support Not Found
	v := mux.Vars(r)
	bucket, err := GetInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}
	key := v["object"]
	readBucket := bucket.GetReadBucket()
	writeBucket := bucket.GetWriteBucket()

	ro, _ := getObject(readBucket, key)
	wo, _ := getObject(writeBucket, key)

	// FIXME: Need Refactor
	if ro != nil {
		defer ro.Body.Close()
		b, e := ioutil.ReadAll(ro.Body)
		if e != nil {
			SendInternalError("Something Happend", w, r)
		} else {
			api.SendSuccess(w, b, *ro.ContentType)
		}
	} else if wo != nil {
		defer wo.Body.Close() // TODO: fix memory leak if ro/wo is present
		b, e := ioutil.ReadAll(wo.Body)
		if e != nil {
			SendInternalError("Something Happend", w, r)
		} else {
			api.SendSuccess(w, b, *wo.ContentType)
		}
	} else {
		SendInternalError("404 NotFound. Not Implemented Yet", w, r)
	}
}

func getObject(bucket *db.S3Bucket, key string) (*s3sdk.GetObjectOutput, error) {
	// TODO: Support range Get
	client := s3client.GetS3Client(bucket)
	resp, err := client.GetObject(&s3sdk.GetObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return resp, err
}
