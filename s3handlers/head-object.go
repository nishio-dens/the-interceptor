package s3handlers

import (
	"net/http"
	"the-interceptor/api"
)

// TODO: In order to use aws s3 cp s3://some/bucket/path - commands,
//       We need to implements HEAD
func HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: not implemented yet
	w.Header().Set("Content-Length", "100") // FIXME: dummy
	api.SendSuccess(w, []byte(""), "")
}
