package s3handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

func HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "HEAD %v\n", vars["category"])
}
