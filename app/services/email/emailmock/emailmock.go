package emailmock

import (
	"context"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
)

var MessageHistory = make([]*HistoryItem, 0)

type HistoryItem struct {
	From         dto.Recipient
	To           []dto.Recipient
	TemplateName string
	Props        dto.Props
	Tenant       *entity.Tenant
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
	bus.AddListener(sendMail)
	bus.AddHandler(fetchRecentSupressions)
}

func fetchRecentSupressions(ctx context.Context, c *query.FetchRecentSupressions) error {
	// not implemented for Email Mock
	return nil
}

func sendMail(ctx context.Context, c *cmd.SendMail) {
	if c.Props == nil {
		c.Props = dto.Props{}
	}
	item := &HistoryItem{
		From:         c.From,
		To:           c.To,
		TemplateName: c.TemplateName,
		Props:        c.Props,
	}

	tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if ok {
		item.Tenant = tenant
	}
	MessageHistory = append(MessageHistory, item)
}
