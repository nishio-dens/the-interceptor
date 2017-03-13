package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"the-interceptor/db"
	"the-interceptor/s3handlers"
)

func main() {
	db.InitConnection()

	r := mux.NewRouter().SkipClean(true)
	s3handlers.RegisterRoute(r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
