package s3handlers

import (
	"net/http"
	"the-interceptor/api"
)

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	api.SendSuccess(w, []byte("Call2"), "Test")
}
