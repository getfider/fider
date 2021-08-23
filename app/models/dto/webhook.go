package dto

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/webhook"
)

type WebhookTriggerResult struct {
	Webhook    *entity.Webhook `json:"webhook"`
	Props      webhook.Props   `json:"props"`
	Success    bool            `json:"success"`
	Url        string          `json:"url"`
	Content    string          `json:"content"`
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Error      string          `json:"error"`
}

type WebhookPreviewResult struct {
	Url     PreviewedField `json:"url"`
	Content PreviewedField `json:"content"`
}

type PreviewedField struct {
	Value   string `json:"value,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
