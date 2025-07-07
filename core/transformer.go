package core

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// HTTPTransformer is a generic interface for transforming messages to HTTP requests.
// T must implement Selectable interface, usually *Account or other configuration type.
type HTTPTransformer[T Selectable] interface {
	Transform(ctx context.Context, msg Message, config T) (*HTTPRequestSpec, ResponseHandler, error)
	CanTransform(msg Message) bool
}

// HTTPRequestSpec defines the specification for an HTTP request.
type HTTPRequestSpec struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	QueryParams url.Values        `json:"query_params"`
	Body        []byte            `json:"body"`
	BodyType    BodyType          `json:"body_type"` // "raw", "form", "json",
	Timeout     time.Duration     `json:"timeout"`
}

// ResponseHandler processes an HTTP response. Implementations must ensure resp.Body is CLOSED.
// A common pattern is to call utils.ReadAndClose(resp) to obtain body bytes
// then perform validation/decoding.
type ResponseHandler func(resp *http.Response) error
