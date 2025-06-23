package core

import (
	"context"
)

type (
	strategyKey struct{}
	itemNameKey struct{}
	metadataKey struct{}
)

// WithCtxStrategy provides a specific instance of a selection strategy.
func WithCtxStrategy(ctx context.Context, strategy SelectionStrategy) context.Context {
	return context.WithValue(ctx, strategyKey{}, strategy)
}

// GetStrategyFromCtx retrieves the strategy instance from the context.
func GetStrategyFromCtx(ctx context.Context) SelectionStrategy {
	if strategy, ok := ctx.Value(strategyKey{}).(SelectionStrategy); ok {
		return strategy
	}
	return nil
}

// WithCtxItemName specifies the name of a specific item (e.g., a bot or account) to use.
func WithCtxItemName(ctx context.Context, itemName string) context.Context {
	return context.WithValue(ctx, itemNameKey{}, itemName)
}

// GetItemNameFromCtx retrieves the specific item name from the context.
func GetItemNameFromCtx(ctx context.Context) string {
	if name, ok := ctx.Value(itemNameKey{}).(string); ok {
		return name
	}
	return ""
}

func WithCtxSendMetadata(ctx context.Context, metadata map[string]interface{}) context.Context {
	return context.WithValue(ctx, metadataKey{}, metadata)
}

// GetSendMetadataFromCtx 从context读取metadata
func GetSendMetadataFromCtx(ctx context.Context) map[string]interface{} {
	if m, ok := ctx.Value(metadataKey{}).(map[string]interface{}); ok {
		return m
	}
	return nil
}
