package s3handlers

import (
	"github.com/gorilla/mux"
)

func RegisterRoute(r *mux.Router) {
	root := r.PathPrefix("/").Subrouter()
	bucket := root.PathPrefix("/{bucket}").Subrouter()

	bucket.Methods("HEAD").Path("/{object:.+}").HandlerFunc(HeadObjectHandler)
	bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(GetObjectHandler)
	bucket.Methods("GET").HandlerFunc(ListObjectV1Handler)
	bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(PutObjectHandler)
}
