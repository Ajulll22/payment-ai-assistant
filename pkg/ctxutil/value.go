package ctxutil

import (
	"context"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
)

type ContextValue struct {
	Key   constant.ContextKey
	Value any
}

func GetValueFromContext(ctx context.Context, key constant.ContextKey) any {
	return ctx.Value(key)
}

func SetValueToContext(ctx context.Context, values ...ContextValue) context.Context {
	for _, value := range values {
		ctx = context.WithValue(ctx, value.Key, value.Value)
	}

	return ctx
}
