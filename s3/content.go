package s3

type Content struct {
	Key          string
	LastModified string
	ETag         string
	Size         int64
	StorageClass string
	Owner        Owner `xml:"Owner"`
}
