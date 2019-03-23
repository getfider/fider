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

func GetProps(ctx context.Context) dto.Props {
	props, ok := ctx.Value(app.LogPropsCtxKey).(*dto.Props)
	if ok {
		return *props
	}
	return dto.Props{}
}

func GetProperty(ctx context.Context, key string) interface{} {
	props := GetProps(ctx)
	return props[key]
}

func SetProperty(ctx context.Context, key string, value interface{}) context.Context {
	if value == nil {
		return ctx
	}

	props, ok := ctx.Value(app.LogPropsCtxKey).(*dto.Props)
	if ok {
		(*props)[key] = value
		return ctx
	} else {
		p := dto.Props{
			key: value,
		}
		return context.WithValue(ctx, app.LogPropsCtxKey, &p)
	}
}
