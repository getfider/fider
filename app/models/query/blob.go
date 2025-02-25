package query

import "github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"

type ListBlobs struct {
	Prefix string

	Result []string
}

type GetBlobByKey struct {
	Key string

	Result *dto.Blob
}
