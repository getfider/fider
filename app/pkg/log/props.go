package log

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
)

const (
	// PropertyKeySessionID is the session id of current logger
	PropertyKeySessionID = "SessionID"
	// PropertyKeyContextID is the context id of current logger
	PropertyKeyContextID = "ContextID"
	// PropertyKeyUserID is the user id of current logger
	PropertyKeyUserID = "UserID"
	// PropertyKeyTenantID is the tenant id of current logger
	PropertyKeyTenantID = "TenantID"
	// PropertyKeyTag is the tag of current logger
	PropertyKeyTag = "Tag"
)

func GetProperties(ctx context.Context) dto.Props {
	props, ok := ctx.Value(app.LogPropsCtxKey).(*dto.Props)
	if ok {
		return *props
	}
	return dto.Props{}
}

func GetProperty(ctx context.Context, key string) any {
	props := GetProperties(ctx)
	return props[key]
}

func WithProperty(ctx context.Context, key string, value any) context.Context {
	return WithProperties(ctx, dto.Props{
		key: value,
	})
}

func WithProperties(ctx context.Context, kv dto.Props) context.Context {
	props, ok := ctx.Value(app.LogPropsCtxKey).(*dto.Props)
	if !ok {
		props = &dto.Props{}
		ctx = context.WithValue(ctx, app.LogPropsCtxKey, props)
	}

	for key, value := range kv {
		(*props)[key] = value
	}

	return ctx
}
