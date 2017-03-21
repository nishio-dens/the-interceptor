package s3handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"the-interceptor/api"
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
	_, err := GetInterceptorBucket(v["bucket"])
	if err != nil {
		SendNoSuchBucketError(v["bucket"], w, r)
		return
	}

	api.SendSuccess(w, []byte("Call2"), "Test")
}
