package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/cmd"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

func TestWebSetup(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		trx, ok := c.Value(app.TransactionCtxKey).(*dbx.Trx)
		Expect(ok).IsTrue()
		Expect(trx).IsNotNil()
		return c.NoContent(http.StatusOK)
	})

	Expect(status).Equals(http.StatusOK)
}

func TestWebSetup_Failure(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		return c.Failure(errors.New("Something went wrong..."))
	})

	Expect(status).Equals(http.StatusInternalServerError)
}

func TestWebSetup_Panic(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		panic("Boom!")
	})

	Expect(status).Equals(http.StatusInternalServerError)
}

func TestWebSetup_Logging_Success(t *testing.T) {
	RegisterT(t)

	var infoLogs []*cmd.LogInfo
	bus.AddListener(func(ctx context.Context, q *cmd.LogInfo) error {
		infoLogs = append(infoLogs, q)
		return nil
	})

	var errorLogs []*cmd.LogError
	bus.AddListener(func(ctx context.Context, q *cmd.LogError) error {
		errorLogs = append(errorLogs, q)
		return nil
	})

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		return c.Ok(nil)
	})

	Expect(status).Equals(http.StatusOK)

	Expect(infoLogs).HasLen(2)
	Expect(infoLogs[0].Message).Equals("@{HttpMethod:magenta} @{URL:magenta} started")
	Expect(infoLogs[1].Message).Equals("@{HttpMethod:magenta} @{URL:magenta} finished with @{StatusCode:magenta} in @{ElapsedMs:magenta}ms (@{State})")
	Expect(infoLogs[1].Props["StatusCode"]).Equals(200)
	Expect(infoLogs[1].Props["State"]).Equals("committed")

	Expect(errorLogs).HasLen(0)
}

func TestWebSetup_Logging_RequestCanceled(t *testing.T) {
	RegisterT(t)

	var infoLogs []*cmd.LogInfo
	bus.AddListener(func(ctx context.Context, q *cmd.LogInfo) error {
		infoLogs = append(infoLogs, q)
		return nil
	})

	var errorLogs []*cmd.LogError
	bus.AddListener(func(ctx context.Context, q *cmd.LogError) error {
		errorLogs = append(errorLogs, q)
		return nil
	})

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		return context.Canceled
	})

	Expect(status).Equals(http.StatusOK)

	Expect(infoLogs).HasLen(2)
	Expect(infoLogs[0].Message).Equals("@{HttpMethod:magenta} @{URL:magenta} started")
	Expect(infoLogs[1].Message).Equals("@{HttpMethod:magenta} @{URL:magenta} was canceled after @{ElapsedMs:magenta}ms")

	Expect(errorLogs).HasLen(0)
}

func TestWebSetup_Logging_Error(t *testing.T) {
	RegisterT(t)

	var infoLogs []*cmd.LogInfo
	bus.AddListener(func(ctx context.Context, q *cmd.LogInfo) error {
		infoLogs = append(infoLogs, q)
		return nil
	})

	var errorLogs []*cmd.LogError
	bus.AddListener(func(ctx context.Context, q *cmd.LogError) error {
		errorLogs = append(errorLogs, q)
		return nil
	})

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		return c.Failure(errors.New("Something went wrong..."))
	})

	Expect(status).Equals(http.StatusInternalServerError)

	Expect(infoLogs).HasLen(2)
	Expect(infoLogs[0].Message).Equals("@{HttpMethod:magenta} @{URL:magenta} started")
	Expect(infoLogs[1].Message).Equals("@{HttpMethod:magenta} @{URL:magenta} finished with @{StatusCode:magenta} in @{ElapsedMs:magenta}ms (@{State})")
	Expect(infoLogs[1].Props["StatusCode"]).Equals(500)
	Expect(infoLogs[1].Props["State"]).Equals("next_error")

	Expect(errorLogs).HasLen(1)
	Expect(errorLogs[0].Err.Error()).ContainsSubstring("Something went wrong...")
}

func TestWebSetup_NotQueueTask_OnFailure(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		c.Enqueue(mock.NewNoopTask())
		return c.Failure(errors.New("Something went wrong..."))
	})

	Expect(status).Equals(http.StatusInternalServerError)
	Expect(server.Engine().Worker().Length()).Equals(int64(0))
}

func TestWebSetup_QueueTask_OnSuccess(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	server.Use(middlewares.WebSetup())
	status, _ := server.Execute(func(c *web.Context) error {
		c.Enqueue(mock.NewNoopTask())
		return c.Ok(web.Map{})
	})

	Expect(status).Equals(http.StatusOK)
	Expect(server.Engine().Worker().Length()).Equals(int64(1))
}

func TestWorkerSetup(t *testing.T) {
	RegisterT(t)

	c := worker.NewContext(context.Background(), "0", worker.Task{Name: "Any Task"})
	mw := middlewares.WorkerSetup()
	err := mw(func(c *worker.Context) error {
		Expect(c.Value(app.TransactionCtxKey)).IsNotNil()
		return nil
	})(c)
	Expect(err).IsNil()
}

func TestWorkerSetup_Failure(t *testing.T) {
	RegisterT(t)

	c := worker.NewContext(context.Background(), "0", worker.Task{Name: "Any Task"})
	mw := middlewares.WorkerSetup()
	err := mw(func(c *worker.Context) error {
		Expect(c.Value(app.TransactionCtxKey)).IsNotNil()
		return errors.New("Not Found")
	})(c)
	Expect(err).IsNotNil()
}
