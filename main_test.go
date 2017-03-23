package main

import (
	"testing"
	"the-interceptor/db"
)

func TestDbInitConnection(t *testing.T) {
	db.InitConnection()
	db.Conn.Exec("select 1 from interceptor_buckets;")
}
