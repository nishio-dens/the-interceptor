package s3

type CommonPrefix struct {
	Prefix string
}

type CommonPrefixSortByPrefix []CommonPrefix

func (c CommonPrefixSortByPrefix) Len() int {
	return len(c)
}

func (c CommonPrefixSortByPrefix) Less(i, j int) bool {
	return c[i].Prefix < c[j].Prefix
}

func (c CommonPrefixSortByPrefix) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
