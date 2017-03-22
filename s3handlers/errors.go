package s3handlers

import (
	"encoding/xml"
	"net/http"
	"the-interceptor/api"
)

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error"`

	Code      string
	Message   string
	Resource  string
	RequestId string
}

func SendNoSuchBucketError(bucketName string, w http.ResponseWriter, r *http.Request) {
	resp := ErrorResponse{
		Code:      "NoSuchBucket",
		Message:   "The specified bucket does not exist",
		Resource:  bucketName,
		RequestId: "NotImplementedYet",
	}
	api.SendNotFoundXml(w, resp)
}

func SendNoSuchKeyError(key string, w http.ResponseWriter, r *http.Request) {
	resp := ErrorResponse{
		Code:      "NoSuchKey",
		Message:   "The specified bucket key does not exist",
		Resource:  key,
		RequestId: "NotImplementedYet",
	}
	api.SendNotFoundXml(w, resp)
}

func SendInternalError(message string, w http.ResponseWriter, r *http.Request) {
	resp := ErrorResponse{
		Code:      "InternalError",
		Message:   message,
		Resource:  "",
		RequestId: "NotImplementedYet",
	}
	api.SendInternalErrorXml(w, resp)
}
