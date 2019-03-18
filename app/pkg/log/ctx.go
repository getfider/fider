package log

import (
	"context"

	"github.com/getfider/fider/app"
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

func GetProps(ctx context.Context) Props {
	props, ok := ctx.Value(app.LogPropsCtxKey).(*Props)
	if ok {
		return *props
	}
	return Props{}
}

func SetProperty(ctx context.Context, key string, value interface{}) context.Context {
	props, ok := ctx.Value(app.LogPropsCtxKey).(*Props)
	if ok {
		(*props)[key] = value
		return ctx
	} else {
		p := Props{
			key: value,
		}
		return context.WithValue(ctx, app.LogPropsCtxKey, &p)
	}
}
