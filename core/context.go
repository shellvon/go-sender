package core

import (
	"context"
)

type (
	metadataKey struct{}
	routeKey    struct{}
)

// RouteInfo carries per-send routing preferences such as explicit account
// selection or a custom strategy instance.
type RouteInfo struct {
	AccountName  string
	StrategyType StrategyType
}

// WithRoute stores the provided RouteInfo into context and returns the derived
// context. Pass nil to leave the original context untouched.
func WithRoute(ctx context.Context, info *RouteInfo) context.Context {
	if info == nil {
		return ctx
	}
	return context.WithValue(ctx, routeKey{}, info)
}

// GetRoute extracts RouteInfo from context. Returns nil when absent.
func GetRoute(ctx context.Context) *RouteInfo {
	if ri, ok := ctx.Value(routeKey{}).(*RouteInfo); ok {
		return ri
	}
	return nil
}

// GetSendMetadataFromCtx retrieves the metadata map stored internally in the context.
func GetSendMetadataFromCtx(ctx context.Context) map[string]interface{} {
	if m, ok := ctx.Value(metadataKey{}).(map[string]interface{}); ok {
		return m
	}
	return nil
}
