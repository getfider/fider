package tasks

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/worker"
)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

func link(baseURL, path string, args ...any) string {
	return fmt.Sprintf("<a href='%[1]s%[2]s'>%[1]s%[2]s</a>", baseURL, fmt.Sprintf(path, args...))
}

func linkWithText(text, baseURL, path string, args ...any) string {
	return fmt.Sprintf("<a href='%s%s'>%s</a>", baseURL, fmt.Sprintf(path, args...), text)
}

func getActiveSubscribers(ctx context.Context, post *entity.Post, channel enum.NotificationChannel, event enum.NotificationEvent) ([]*entity.User, error) {
	q := &query.GetActiveSubscribers{
		Number:  post.Number,
		Channel: channel,
		Event:   event,
	}
	err := bus.Dispatch(ctx, q)
	return q.Result, err
}
