package cmd

type StoreBlob struct {
	Key         string
	Content     []byte
	ContentType string
}

type DeleteBlob struct {
	Key string
}
