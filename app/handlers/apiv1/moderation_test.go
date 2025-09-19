package apiv1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestApprovePost_Success(t *testing.T) {
	RegisterT(t)

	var approvePost *cmd.ApprovePost
	bus.AddHandler(func(ctx context.Context, c *cmd.ApprovePost) error {
		approvePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.ApprovePost(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(approvePost.PostID).Equals(5)
}

func TestApprovePost_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "invalid").
		ExecuteAsJSON(apiv1.ApprovePost())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid post ID")
}

func TestApprovePost_CommandFailure(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.ApprovePost) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.ApprovePost(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestDeclinePost_Success(t *testing.T) {
	RegisterT(t)

	var declinePost *cmd.DeclinePost
	bus.AddHandler(func(ctx context.Context, c *cmd.DeclinePost) error {
		declinePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "10").
		ExecutePost(apiv1.DeclinePost(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(declinePost.PostID).Equals(10)
}

func TestDeclinePost_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "abc").
		ExecuteAsJSON(apiv1.DeclinePost())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid post ID")
}

func TestDeclinePost_CommandFailure(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.DeclinePost) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "10").
		ExecutePost(apiv1.DeclinePost(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApproveComment_Success(t *testing.T) {
	RegisterT(t)

	var approveComment *cmd.ApproveComment
	bus.AddHandler(func(ctx context.Context, c *cmd.ApproveComment) error {
		approveComment = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "15").
		ExecutePost(apiv1.ApproveComment(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(approveComment.CommentID).Equals(15)
}

func TestApproveComment_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "not-a-number").
		ExecuteAsJSON(apiv1.ApproveComment())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid comment ID")
}

func TestApproveComment_CommandFailure(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.ApproveComment) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "15").
		ExecutePost(apiv1.ApproveComment(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestDeclineComment_Success(t *testing.T) {
	RegisterT(t)

	var declineComment *cmd.DeclineComment
	bus.AddHandler(func(ctx context.Context, c *cmd.DeclineComment) error {
		declineComment = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "20").
		ExecutePost(apiv1.DeclineComment(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(declineComment.CommentID).Equals(20)
}

func TestDeclineComment_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "").
		ExecuteAsJSON(apiv1.DeclineComment())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid comment ID")
}

func TestDeclineComment_CommandFailure(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, c *cmd.DeclineComment) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "20").
		ExecutePost(apiv1.DeclineComment(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestDeclinePostAndBlock_Success(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:     1,
		Number: 5,
		Title:  "Test Post",
		User:   mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		if q.PostID == 5 {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	var blockUser *cmd.BlockUser
	bus.AddHandler(func(ctx context.Context, c *cmd.BlockUser) error {
		blockUser = c
		return nil
	})

	var declinePost *cmd.DeclinePost
	bus.AddHandler(func(ctx context.Context, c *cmd.DeclinePost) error {
		declinePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.DeclinePostAndBlock(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(blockUser.UserID).Equals(mock.AryaStark.ID)
	Expect(declinePost.PostID).Equals(5)
}

func TestDeclinePostAndBlock_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "invalid").
		ExecuteAsJSON(apiv1.DeclinePostAndBlock())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid post ID")
}

func TestDeclinePostAndBlock_PostNotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "999").
		ExecutePost(apiv1.DeclinePostAndBlock(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestDeclinePostAndBlock_BlockUserFailure(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:     1,
		Number: 5,
		Title:  "Test Post",
		User:   mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		if q.PostID == 5 {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.BlockUser) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.DeclinePostAndBlock(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestDeclineCommentAndBlock_Success(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{
		ID:      25,
		Content: "Test Comment",
		User:    mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == 25 {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	var blockUser *cmd.BlockUser
	bus.AddHandler(func(ctx context.Context, c *cmd.BlockUser) error {
		blockUser = c
		return nil
	})

	var declineComment *cmd.DeclineComment
	bus.AddHandler(func(ctx context.Context, c *cmd.DeclineComment) error {
		declineComment = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "25").
		ExecutePost(apiv1.DeclineCommentAndBlock(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(blockUser.UserID).Equals(mock.AryaStark.ID)
	Expect(declineComment.CommentID).Equals(25)
}

func TestDeclineCommentAndBlock_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "xyz").
		ExecuteAsJSON(apiv1.DeclineCommentAndBlock())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid comment ID")
}

func TestDeclineCommentAndBlock_CommentNotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "999").
		ExecutePost(apiv1.DeclineCommentAndBlock(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestDeclineCommentAndBlock_BlockUserFailure(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{
		ID:      25,
		Content: "Test Comment",
		User:    mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == 25 {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.BlockUser) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "25").
		ExecutePost(apiv1.DeclineCommentAndBlock(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApprovePostAndVerify_Success(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:     1,
		Number: 5,
		Title:  "Test Post",
		User:   mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		if q.PostID == 5 {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	var approvePost *cmd.ApprovePost
	bus.AddHandler(func(ctx context.Context, c *cmd.ApprovePost) error {
		approvePost = c
		return nil
	})

	var verifyUser *cmd.VerifyUser
	bus.AddHandler(func(ctx context.Context, c *cmd.VerifyUser) error {
		verifyUser = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.ApprovePostAndVerify(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(approvePost.PostID).Equals(5)
	Expect(verifyUser.UserID).Equals(mock.AryaStark.ID)
}

func TestApprovePostAndVerify_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "invalid").
		ExecuteAsJSON(apiv1.ApprovePostAndVerify())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid post ID")
}

func TestApprovePostAndVerify_PostNotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "999").
		ExecutePost(apiv1.ApprovePostAndVerify(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApprovePostAndVerify_ApproveFailure(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:     1,
		Number: 5,
		Title:  "Test Post",
		User:   mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		if q.PostID == 5 {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ApprovePost) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.ApprovePostAndVerify(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApprovePostAndVerify_VerifyUserFailure(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:     1,
		Number: 5,
		Title:  "Test Post",
		User:   mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		if q.PostID == 5 {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ApprovePost) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.VerifyUser) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "5").
		ExecutePost(apiv1.ApprovePostAndVerify(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApproveCommentAndVerify_Success(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{
		ID:      25,
		Content: "Test Comment",
		User:    mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == 25 {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	var approveComment *cmd.ApproveComment
	bus.AddHandler(func(ctx context.Context, c *cmd.ApproveComment) error {
		approveComment = c
		return nil
	})

	var verifyUser *cmd.VerifyUser
	bus.AddHandler(func(ctx context.Context, c *cmd.VerifyUser) error {
		verifyUser = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "25").
		ExecutePost(apiv1.ApproveCommentAndVerify(), `{}`)

	Expect(code).Equals(http.StatusOK)
	Expect(approveComment.CommentID).Equals(25)
	Expect(verifyUser.UserID).Equals(mock.AryaStark.ID)
}

func TestApproveCommentAndVerify_InvalidID(t *testing.T) {
	RegisterT(t)

	code, response := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "invalid").
		ExecuteAsJSON(apiv1.ApproveCommentAndVerify())

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("error")).Equals("Invalid comment ID")
}

func TestApproveCommentAndVerify_CommentNotFound(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "999").
		ExecutePost(apiv1.ApproveCommentAndVerify(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApproveCommentAndVerify_ApproveFailure(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{
		ID:      25,
		Content: "Test Comment",
		User:    mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == 25 {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ApproveComment) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "25").
		ExecutePost(apiv1.ApproveCommentAndVerify(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestApproveCommentAndVerify_VerifyUserFailure(t *testing.T) {
	RegisterT(t)

	comment := &entity.Comment{
		ID:      25,
		Content: "Test Comment",
		User:    mock.AryaStark,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID == 25 {
			q.Result = comment
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ApproveComment) error {
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.VerifyUser) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("id", "25").
		ExecutePost(apiv1.ApproveCommentAndVerify(), `{}`)

	Expect(code).Equals(http.StatusNotFound)
}
