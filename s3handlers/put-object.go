package s3handlers

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"the-interceptor/api"
	"the-interceptor/db"
	"the-interceptor/s3client"
)

type putObjectResponseResult struct {
	Key  string
	ETag string
}

/**
PUT Object
see: http://docs.aws.amazon.com/ja_jp/AmazonS3/latest/API/RESTObjectPUT.html
*/
func PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	bucket, err := GetInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}

	defer r.Body.Close()
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		SendInternalError("Cannot Read Request Body", w, r)
		return
	}
	key := v["object"]
	writeBucket := bucket.GetWriteBucket()
	bodySeeker := bytes.NewReader(b)
	resp, err := putObject(writeBucket, key, bodySeeker)
	if err != nil {
		SendInternalError("Cannot Upload", w, r)
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("%s", resp.ETag))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", 0))
	api.SendNoBodySuccess(w)
}

func putObject(bucket *db.S3Bucket, key string, body io.ReadSeeker) (*putObjectResponseResult, error) {
	client := s3client.GetS3Client(bucket)
	resp, err := client.PutObject(&s3sdk.PutObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return nil, err
	}
	rs := &putObjectResponseResult{
		Key:  key,
		ETag: *resp.ETag,
	}
	return rs, err
}
