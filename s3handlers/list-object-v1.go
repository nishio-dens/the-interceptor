package s3handlers

import (
	"encoding/xml"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"net/http"
	"the-interceptor/api"
	"the-interceptor/db"
	"the-interceptor/s3"
	"the-interceptor/s3client"
	"time"
)

type ListObjectV1Response struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult"`

	Name        string
	Prefix      string
	Marker      string
	MaxKeys     int
	IsTruncated bool

	Contents       []s3.Content
	CommonPrefixes []s3.CommonPrefix
}

type listObjectV1ResponseResult struct {
	Result *ListObjectV1Response
	Error  error
}

/**
GET Bucket (List Objects) Version 1
see: http://docs.aws.amazon.com/ja_jp/AmazonS3/latest/API/RESTBucketGET.html
*/
func ListObjectV1Handler(w http.ResponseWriter, r *http.Request) {
	// Future Work
	// TODO: Merge Read and Write Bucket Objects
	// TODO: Support maxKeys
	// TODO: Support Marker
	// TODO: Support Paging

	v := mux.Vars(r)
	uquery := r.URL.Query()
	bucket, err := getInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}

	readBucket := bucket.GetReadBucket()

	delim := "/"
	if len(uquery.Get("delimiter")) > 0 {
		delim = uquery.Get("delimiter")
	}
	prefix := ""
	if len(uquery.Get("prefix")) > 0 {
		prefix = uquery.Get("prefix")
	}
	maxKeys := int64(1000) // TODO: FIXME
	ri := &s3sdk.ListObjectsInput{
		Bucket:    aws.String(readBucket.BucketName),
		MaxKeys:   aws.Int64(maxKeys),
		Delimiter: aws.String(delim),
		Prefix:    aws.String(prefix),
	}
	rchan := make(chan listObjectV1ResponseResult)
	go getListObjects(readBucket, ri, rchan)

	requestBuckets := []*db.S3Bucket{readBucket}
	results := make([]listObjectV1ResponseResult, 2)
	for i := range requestBuckets {
		results[i] = <-rchan
	}

	api.SendSuccessXml(w, results[0].Result)
}

func getInterceptorBucket(name string) (*db.InterceptorBucket, error) {
	var bucket db.InterceptorBucket
	if db.Conn.Where("name = ?", name).First(&bucket).RecordNotFound() {
		return nil, errors.New("Record Not Found")
	}
	return &bucket, nil
}

func getListObjects(bucket *db.S3Bucket, input *s3sdk.ListObjectsInput, ch chan<- listObjectV1ResponseResult) {
	client := s3client.GetS3Client(bucket)
	resp, err := client.ListObjects(input)
	if err != nil {
		// Internal Error
		ch <- listObjectV1ResponseResult{
			Result: nil,
			Error:  err,
		}
		return
	}

	contents := make([]s3.Content, len(resp.Contents))
	for i, c := range resp.Contents {
		o := s3.Owner{
			Id:          string(*c.Owner.ID),
			DisplayName: string(*c.Owner.DisplayName),
		}
		contents[i] = s3.Content{
			Key:          string(*c.Key),
			LastModified: (*c.LastModified).Format(time.RFC3339),
			ETag:         string(*c.ETag),
			Size:         int64(*c.Size),
			StorageClass: string(*c.StorageClass),
			Owner:        o,
		}
	}

	prefixes := make([]s3.CommonPrefix, len(resp.CommonPrefixes))
	for i, p := range resp.CommonPrefixes {
		prefixes[i] = s3.CommonPrefix{Prefix: *p.Prefix}
	}

	rs := ListObjectV1Response{
		Name:           bucket.BucketName,
		Prefix:         "",
		Marker:         "",
		MaxKeys:        1000,
		IsTruncated:    false,
		Contents:       contents,
		CommonPrefixes: prefixes,
	}
	ch <- listObjectV1ResponseResult{
		Result: &rs,
		Error:  nil,
	}
}
