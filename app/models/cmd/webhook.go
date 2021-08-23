package cmd

import (
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/webhook"
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
