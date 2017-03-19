package s3handlers

import (
	"encoding/xml"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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
	Result   *ListObjectV1Response
	Error    error
	Priority int
}

/**
GET Bucket (List Objects) Version 1
see: http://docs.aws.amazon.com/ja_jp/AmazonS3/latest/API/RESTBucketGET.html
*/
func ListObjectV1Handler(w http.ResponseWriter, r *http.Request) {
	// Future Work
	// TODO: Merge Read and Write Bucket Objects
	// TODO: Support Marker
	// TODO: Support Paging

	v := mux.Vars(r)
	bucket, err := getInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}

	uquery := r.URL.Query()
	readBucket := bucket.GetReadBucket()
	writeBucket := bucket.GetWriteBucket()
	ri := listObjectInput(readBucket, uquery)
	wi := listObjectInput(writeBucket, uquery)
	requestBuckets := []*db.S3Bucket{readBucket, writeBucket}

	rchan := make(chan listObjectV1ResponseResult)
	go getListObjects(readBucket, 100, ri, rchan)
	go getListObjects(writeBucket, 1, wi, rchan)

	results := make([]listObjectV1ResponseResult, 2)
	for i := range requestBuckets {
		results[i] = <-rchan
		if results[i].Error != nil {
			// TODO: Implement Error handling correctly
			SendInternalError("Something Happend. Maybe your bucket settings is wrong.", w, r)
			return
		}
	}

	fr := mergeListObjectResponse(bucket, results)
	api.SendSuccessXml(w, *fr)
}

func getInterceptorBucket(name string) (*db.InterceptorBucket, error) {
	var bucket db.InterceptorBucket
	if db.Conn.Where("name = ?", name).First(&bucket).RecordNotFound() {
		return nil, errors.New("Record Not Found")
	}
	return &bucket, nil
}

func listObjectInput(bucket *db.S3Bucket, uquery url.Values) *s3sdk.ListObjectsInput {
	// TODO: Need Refactor! param should be collect to struct
	delim := "/"
	if len(uquery.Get("delimiter")) > 0 {
		delim = uquery.Get("delimiter")
	}
	prefix := ""
	if len(uquery.Get("prefix")) > 0 {
		prefix = uquery.Get("prefix")
	}
	maxKeys := int64(1000)
	if len(uquery.Get("max-keys")) > 0 {
		k, err := strconv.Atoi(uquery.Get("max-keys"))
		if err == nil {
			maxKeys = int64(k)
		}
	}

	// TODO: Support EncodingType, Marker, RequestPayer
	return &s3sdk.ListObjectsInput{
		Bucket:    aws.String(bucket.BucketName),
		MaxKeys:   aws.Int64(maxKeys),
		Delimiter: aws.String(delim),
		Prefix:    aws.String(prefix),
	}
}

func getListObjects(
	bucket *db.S3Bucket,
	priority int,
	input *s3sdk.ListObjectsInput,
	ch chan<- listObjectV1ResponseResult,
) {
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

	// TODO: implements Prefix, Marker, MaxKeys, IsTruncated
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
		Result:   &rs,
		Error:    nil,
		Priority: priority,
	}
}

func mergeListObjectResponse(
	bucket *db.InterceptorBucket, responses []listObjectV1ResponseResult,
) *ListObjectV1Response {
	contEncountered := map[string]int{}
	conts := map[string]s3.Content{}
	prefixEncountered := map[string]int{}
	prefixes := map[string]s3.CommonPrefix{}

	for _, rr := range responses {
		for _, c := range rr.Result.Contents {
			if contEncountered[c.Key] < rr.Priority {
				contEncountered[c.Key] = rr.Priority
				conts[c.Key] = c
			}
		}

		for _, p := range rr.Result.CommonPrefixes {
			if prefixEncountered[p.Prefix] < rr.Priority {
				prefixEncountered[p.Prefix] = rr.Priority
				prefixes[p.Prefix] = p
			}
		}
	}

	cresults := []s3.Content{}
	for _, v := range conts {
		cresults = append(cresults, v)
	}
	sort.Sort(s3.ContentsSortByKey(cresults))

	presults := []s3.CommonPrefix{}
	for _, v := range prefixes {
		presults = append(presults, v)
	}
	sort.Sort(s3.CommonPrefixSortByPrefix(presults))

	// TODO: implements Prefix, Marker, MaxKeys, IsTruncated
	return &ListObjectV1Response{
		Name:           bucket.Name,
		Prefix:         "",
		Marker:         "",
		MaxKeys:        1000,
		IsTruncated:    false,
		Contents:       cresults,
		CommonPrefixes: presults,
	}
}
