package query

import "github.com/getfider/fider/app/models/dto"

type GetBlobByKey struct {
	Key string

	Result *dto.Blob
}
