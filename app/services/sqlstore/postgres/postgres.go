package postgres

import (
	"context"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/stripe/stripe-go/client"
)

var stripeClient *client.API

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "PostgreSQL"
}

func (s Service) Category() string {
	return "sqlstore"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(storeEvent)

	bus.AddHandler(markAllNotificationsAsRead)
	bus.AddHandler(markNotificationAsRead)
	bus.AddHandler(countUnreadNotifications)
	bus.AddHandler(getNotificationByID)
	bus.AddHandler(getActiveNotifications)
	bus.AddHandler(addNewNotification)
	bus.AddHandler(addSubscriber)
	bus.AddHandler(removeSubscriber)
	bus.AddHandler(getActiveSubscribers)

	bus.AddHandler(getTagBySlug)
	bus.AddHandler(getAssignedTags)
	bus.AddHandler(getAllTags)
	bus.AddHandler(addNewTag)
	bus.AddHandler(updateTag)
	bus.AddHandler(deleteTag)
	bus.AddHandler(assignTag)
	bus.AddHandler(unassignTag)

	bus.AddHandler(addVote)
	bus.AddHandler(removeVote)
	bus.AddHandler(listPostVotes)

	bus.AddHandler(addNewPost)
	bus.AddHandler(updatePost)
	bus.AddHandler(getPostByID)
	bus.AddHandler(getPostBySlug)
	bus.AddHandler(getPostByNumber)
	bus.AddHandler(searchPosts)
	bus.AddHandler(getAllPosts)
	bus.AddHandler(countPostPerStatus)
	bus.AddHandler(markPostAsDuplicate)
	bus.AddHandler(setPostResponse)
	bus.AddHandler(postIsReferenced)

	bus.AddHandler(setAttachments)
	bus.AddHandler(getAttachments)

	bus.AddHandler(addNewComment)
	bus.AddHandler(updateComment)
	bus.AddHandler(deleteComment)
	bus.AddHandler(getCommentByID)
	bus.AddHandler(getCommentsByPost)

	bus.AddHandler(countUsers)
}

type SqlHandler func(trx *dbx.Trx, tenant *models.Tenant, user *models.User) error

func using(ctx context.Context, handler SqlHandler) error {
	trx := ctx.Value(app.TransactionCtxKey).(*dbx.Trx)
	tenant, _ := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	user, _ := ctx.Value(app.UserCtxKey).(*models.User)
	return handler(trx, tenant, user)
}
