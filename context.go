package goezyrouting

import (
	"context"
	"net/http"
)

type ctxKey int

const (
	RequestIDKey ctxKey = iota
	RequestURLPath
)

// WithValue is a handy helper function that will update the request with the newest context
func WithValue(r *http.Request, ctx context.Context, key interface{}, value interface{}) context.Context {
	ctx = context.WithValue(ctx, key, value)
	*r = *r.WithContext(ctx)
	return ctx
}
func getReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func getURLPath(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if path, ok := ctx.Value(RequestURLPath).(string); ok {
		return path
	}
	return ""
}
