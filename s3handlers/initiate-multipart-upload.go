package s3handlers

import "net/http"

type InitiateMultipartUploadResponseResult struct {
}

/**
Initiate Multipart Upload
see: http://docs.aws.amazon.com/AmazonS3/latest/API/mpUploadInitiate.html
 */
func InitiateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	SendInternalError("Something Happend", w, r)
}