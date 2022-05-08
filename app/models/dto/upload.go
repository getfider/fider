package dto

//ImageUpload is the input model used to upload/remove an image
type ImageUpload struct {
	BlobKey string           `json:"bkey"`
	Upload  *ImageUploadData `json:"upload"`
	Remove  bool             `json:"remove"`
}

//ImageUploadData is the input model used to upload a new logo
type ImageUploadData struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Content     []byte `json:"content"`
}
