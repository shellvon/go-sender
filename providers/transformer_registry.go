package providers

import (
	"sync"

	"github.com/shellvon/go-sender/core"
)

// TransformerRegistry provides a threadsafe map from sub-provider name to
// HTTPTransformer instance. It is generic so individual provider packages
// (SMS, EmailAPI, etc.) can embed their own typed registry without code
// duplication.
type TransformerRegistry[T core.Selectable] struct {
	mu           sync.RWMutex
	transformers map[string]core.HTTPTransformer[T]
}

// NewTransformerRegistry returns an initialised empty registry.
func NewTransformerRegistry[T core.Selectable]() *TransformerRegistry[T] {
	return &TransformerRegistry[T]{
		transformers: make(map[string]core.HTTPTransformer[T]),
	}
}

// Register associates a sub-provider identifier with its transformer.
func (r *TransformerRegistry[T]) Register(subProvider string, transformer core.HTTPTransformer[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transformers[subProvider] = transformer
}

// Get retrieves a transformer by sub-provider name.
func (r *TransformerRegistry[T]) Get(subProvider string) (core.HTTPTransformer[T], bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tr, ok := r.transformers[subProvider]
	return tr, ok
}
