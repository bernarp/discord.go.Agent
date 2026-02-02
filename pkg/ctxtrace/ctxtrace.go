package ctxtrace

import (
	"context"
)

type ctxKey struct{}

var corridKey = ctxKey{}

func WithCorrelationID(
	ctx context.Context,
	corrid string,
) context.Context {
	return context.WithValue(ctx, corridKey, corrid)
}

func Extract(ctx context.Context) string {
	if id, ok := ctx.Value(corridKey).(string); ok {
		return id
	}
	return ""
}
