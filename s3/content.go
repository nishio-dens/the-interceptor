package s3

type Content struct {
	Key          string
	LastModified string
	ETag         string
	Size         int64
	StorageClass string
	Owner        Owner `xml:"Owner"`
}

type ContentsSortByKey []Content

func (c ContentsSortByKey) Len() int {
	return len(c)
}

func (c ContentsSortByKey) Less(i, j int) bool {
	return c[i].Key < c[j].Key
}

func (c ContentsSortByKey) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
