package core

import (
	"context"
	"time"
)

// HTTPTransformer 泛型接口，用于转换消息到 HTTP 请求
// T 必须实现 Selectable 接口，通常是 *Account 或其他配置类型.
type HTTPTransformer[T Selectable] interface {
	Transform(ctx context.Context, msg Message, config T) (*HTTPRequestSpec, ResponseHandler, error)
	CanTransform(msg Message) bool
}

// HTTPRequestSpec defines the specification for an HTTP request.
type HTTPRequestSpec struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	Body        []byte            `json:"body"`
	BodyType    string            `json:"body_type"` // "raw", "form", "json"
	Timeout     time.Duration     `json:"timeout"`
}

// ResponseHandler defines the interface for handling HTTP responses.
type ResponseHandler func(statusCode int, body []byte) error
