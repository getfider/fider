package paddle

import "encoding/json"

type PaddleResponse struct {
	IsSuccess bool            `json:"success"`
	Response  json.RawMessage `json:"response"`
	Error     struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}
