package s3handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"the-interceptor/api"
)

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
	key := v["object"]

	fmt.Println("Write ", bucket, key)
	defer r.Body.Close()

	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		SendInternalError("Cannot Read Request Body", w, r)
	}
	fmt.Println("Body Length ", len(b))

	// TODO: Set ETag
	w.Header().Set("Content-Length", fmt.Sprintf("%d", 0))
	api.SendNoBodySuccess(w)
}
