package emailmock

import (
	"context"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/email"
)

var MessageHistory = make([]*HistoryItem, 0)

type HistoryItem struct {
	From         string
	To           []email.Recipient
	TemplateName string
	Params       email.Params
	Tenant       *models.Tenant
}

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Mock"
}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	MessageHistory = make([]*HistoryItem, 0)
	bus.AddEventListener(sendEmail)
}

func sendEmail(ctx context.Context, cmd *email.SendMessageCommand) {
	if cmd.Params == nil {
		cmd.Params = email.Params{}
	}
	item := &HistoryItem{
		From:         cmd.From,
		To:           cmd.To,
		TemplateName: cmd.TemplateName,
		Params:       cmd.Params,
	}

	tenant, ok := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	if ok {
		item.Tenant = tenant
	}
	MessageHistory = append(MessageHistory, item)
}
