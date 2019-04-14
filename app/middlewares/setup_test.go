package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
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
