package apiv1_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestAllCommentsHandler(t *testing.T) {
	RegisterT(t)

	// Create 10 test comment refs
	testComments := make([]*entity.CommentRef, 10)
	for i := 0; i < 10; i++ {
		testComments[i] = &entity.CommentRef{
			ID:        i + 1,
			CreatedAt: time.Now(),
			UserID:    1,
			PostID:    1,
		}
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentRefs) error {
		q.Result = testComments
		return nil
	})

	code, query := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecuteAsJSON(apiv1.AllComments())

	Expect(code).Equals(http.StatusOK)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(10)
}
