package tasks

import (
	"context"
	"fmt"
	"html"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/worker"
)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

// link returns an HTML anchor that displays its own URL. The result is
// rendered unescaped in email templates, so the URL is HTML-escaped here.
func link(baseURL, path string, args ...any) string {
	url := html.EscapeString(baseURL + fmt.Sprintf(path, args...))
	return fmt.Sprintf("<a href='%[1]s'>%[1]s</a>", url)
}

// linkWithText returns an HTML anchor with the given plain text as its label.
// The result is rendered unescaped in email templates, so both the URL and the
// text are HTML-escaped here. Callers must pass plain text, not HTML.
func linkWithText(text, baseURL, path string, args ...any) string {
	url := html.EscapeString(baseURL + fmt.Sprintf(path, args...))
	return fmt.Sprintf("<a href='%s'>%s</a>", url, html.EscapeString(text))
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
