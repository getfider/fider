package log_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/log"
)

func TestLog_WithProperty(t *testing.T) {
	RegisterT(t)
	ctx := log.WithProperty(context.Background(), "Name", "Jon Snow")
	Expect(log.GetProperty(ctx, "Name")).Equals("Jon Snow")
}

func TestLog_WithProperties(t *testing.T) {
	RegisterT(t)
	ctx := log.WithProperties(context.Background(), dto.Props{
		"Name": "Jon Snow",
		"Age":  12,
	})
	Expect(log.GetProperty(ctx, "Name")).Equals("Jon Snow")
	Expect(log.GetProperty(ctx, "Age")).Equals(12)

	ctx = log.WithProperty(ctx, "Age", 15)
	Expect(log.GetProperty(ctx, "Name")).Equals("Jon Snow")
	Expect(log.GetProperty(ctx, "Age")).Equals(15)

	Expect(log.GetProperties(ctx)).Equals(dto.Props{
		"Name": "Jon Snow",
		"Age":  15,
	})
}
