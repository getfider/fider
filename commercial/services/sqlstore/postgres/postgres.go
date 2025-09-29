package postgres

import (
	"github.com/getfider/fider/app/pkg/bus"
)

func init() {
	bus.Register(CommercialService{})
}

type CommercialService struct{}

func (s CommercialService) Name() string {
	return "Commercial PostgreSQL"
}

func (s CommercialService) Category() string {
	return "commercial-sqlstore"
}

func (s CommercialService) Enabled() bool {
	return true
}

func (s CommercialService) Init() {
	// Register commercial moderation handlers
	// These should override the open source stub handlers
	bus.AddHandler(ApprovePost)
	bus.AddHandler(DeclinePost)
	bus.AddHandler(ApproveComment)
	bus.AddHandler(DeclineComment)
	bus.AddHandler(BulkApproveItems)
	bus.AddHandler(BulkDeclineItems)
	bus.AddHandler(GetModerationItems)
	bus.AddHandler(GetModerationCount)
	bus.AddHandler(VerifyUser)
}
