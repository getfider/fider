package cmd

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/webhook"
)

type TestWebhook struct {
	ID int

	Result *dto.WebhookTriggerResult
}

type TriggerWebhooks struct {
	Type  enum.WebhookType
	Props webhook.Props
}

type PreviewWebhook struct {
	Type    enum.WebhookType
	Url     string
	Content string

	Result *dto.WebhookPreviewResult
}

type GetWebhookProps struct {
	Type enum.WebhookType

	Result webhook.Props
}
