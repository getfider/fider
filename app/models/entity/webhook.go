package entity

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/errors"
)

// Webhook represents a webhook
type Webhook struct {
	ID          int                `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Type        enum.WebhookType   `json:"type" db:"type"`
	Status      enum.WebhookStatus `json:"status" db:"status"`
	Url         string             `json:"url" db:"url"`
	Content     string             `json:"content" db:"content"`
	HttpMethod  string             `json:"http_method" db:"http_method"`
	HttpHeaders HttpHeaders        `json:"http_headers" db:"http_headers"`
}

type HttpHeaders map[string]string

func (h HttpHeaders) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *HttpHeaders) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	headers, ok := src.([]byte)
	if !ok {
		return errors.New("Invalid data stored in database")
	}
	return json.Unmarshal(headers, &h)
}
