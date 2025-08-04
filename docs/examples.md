# Examples

Explore real-world usage patterns with go-sender.

## Multi-channel Fallback

```go
// Try first SMS account, then fallback to second SMS account if fails
err := sender.Send(context.Background(), msg, core.WithSendAccount("aliyun-account-1"))
if err != nil {
    _ = sender.Send(context.Background(), msg, core.WithSendAccount("aliyun-account-2"))
}
```

## Batch Sending

### Method 1: Provider-Supported Batch Sending

Some providers support sending to multiple recipients in a single API call:

```go
// SMS providers that support multiple recipients
msg := sms.Aliyun().
    To("13800138000,13800138001,13800138002").  // Multiple numbers separated by comma
    Content("Hello everyone!").
    Build()

err := sender.Send(ctx, msg)
```

### Method 2: Manual Loop for Providers Without Batch Support

For providers that don't support batch sending, use a loop:

```go
mobiles := []string{"13800138000", "13800138001", "13800138002"}

for _, mobile := range mobiles {
    msg := sms.Aliyun().
        To(mobile).
        Content("Hello").
        Build()

    if err := sender.Send(ctx, msg); err != nil {
        log.Printf("Failed to send to %s: %v", mobile, err)
        // Continue with next recipient or handle error as needed
    }
}
```

### Method 3: Concurrent Batch Sending

For better performance with manual loops, use goroutines:

```go
import (
    "sync"
    "context"
)

mobiles := []string{"13800138000", "13800138001", "13800138002"}
var wg sync.WaitGroup
errChan := make(chan error, len(mobiles))

for _, mobile := range mobiles {
    wg.Add(1)
    go func(mobile string) {
        defer wg.Done()

        msg := sms.Aliyun().
            To(mobile).
            Content("Hello").
            Build()

        if err := sender.Send(ctx, msg); err != nil {
            errChan <- fmt.Errorf("failed to send to %s: %w", mobile, err)
        }
    }(mobile)
}

wg.Wait()
close(errChan)

// Check for errors
for err := range errChan {
    log.Printf("Error: %v", err)
}
```

## Asynchronous Sending

### Method 1: Using WithSendAsync (Recommended)

The `WithSendAsync()` option is the key to asynchronous sending. If no queue is set, it uses goroutines:

```go
// Send asynchronously using goroutines (no queue needed)
err := sender.Send(ctx, msg, core.WithSendAsync())
if err != nil {
    // This error is from enqueueing, not from actual sending
    log.Printf("Failed to enqueue message: %v", err)
}
```

### Method 2: With Custom Queue

Set a queue for persistent storage and retry capabilities:

```go
// Set memory queue for local persistence
sender.SetQueue(memoryQueue)

// Send asynchronously with queue
err := sender.Send(ctx, msg, core.WithSendAsync())
if err != nil {
    log.Printf("Failed to enqueue message: %v", err)
}
```

### Method 3: With Callback for Local Queue/Goroutine

For local memory queue or goroutine scenarios, you can specify a callback to get sending results:

```go
// Send asynchronously with callback to get actual sending result
err := sender.Send(ctx, msg,
    core.WithSendAsync(),
    core.WithSendCallback(func(res *core.SendResult, sendErr error) {
        if sendErr != nil {
            log.Printf("Message sending failed: %v", sendErr)
        } else {
            log.Printf("Message sent successfully")
        }
    }),
)
```

**Note:** Callback is only effective for local/in-memory queue or async goroutine scenarios. It will not be called in distributed queues (such as Redis).

## Custom HTTP Client

### 1. Timeout Control

```go
import (
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Set a 30-second timeout to prevent hanging requests
client := &http.Client{
    Timeout: 30 * time.Second,
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 2. Proxy Support

```go
import (
    "net/http"
    "net/url"
    "github.com/shellvon/go-sender/core"
)

// Route requests through a corporate proxy
proxyURL, _ := url.Parse("http://proxy.company.com:8080")
client := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 3. TLS Configuration

```go
import (
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "os"
    "github.com/shellvon/go-sender/core"
)

// Load a custom CA certificate
caCert, _ := os.ReadFile("custom-ca.pem")
caCertPool := x509.NewCertPool()
caCertPool.AppendCertsFromPEM(caCert)

client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            RootCAs:            caCertPool,
            InsecureSkipVerify: false,
        },
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 4. Connection Pooling

```go
import (
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Optimize connection reuse for high concurrency
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,              // max idle connections
        MaxIdleConnsPerHost: 10,               // max idle connections per host
        IdleConnTimeout:     90 * time.Second, // idle connection timeout
        DisableCompression:  false,            // enable compression
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 5. Retry Logic

```go
import (
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Custom retry policy
retryPolicy := &core.RetryPolicy{
    MaxAttempts:   3,
    InitialDelay:  1 * time.Second,
    MaxDelay:      10 * time.Second,
    BackoffFactor: 2.0,
}

client := &http.Client{
    Timeout: 30 * time.Second,
}

// Use custom retry policy
err = sender.Send(ctx, msg,
    core.WithSendHTTPClient(client),
    core.WithSendRetryPolicy(retryPolicy),
)
```

### 6. Load Balancing

```go
import (
    "net/http"
    "net/url"
    "github.com/shellvon/go-sender/core"
)

// Custom transport implementing load balancing
type LoadBalancedTransport struct {
    endpoints []string
    current   int
}

func (t *LoadBalancedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Select endpoint in round-robin
    endpoint := t.endpoints[t.current%len(t.endpoints)]
    t.current++

    // Modify request URL
    req.URL.Host = endpoint
    return http.DefaultTransport.RoundTrip(req)
}

client := &http.Client{
    Transport: &LoadBalancedTransport{
        endpoints: []string{"api1.provider.com", "api2.provider.com"},
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 7. Authentication

```go
import (
    "net/http"
    "github.com/shellvon/go-sender/core"
)

// Add custom authentication headers
client := &http.Client{
    Transport: &http.Transport{
        ProxyConnectHeader: http.Header{
            "Authorization": []string{"Bearer your-token"},
            "X-API-Key":     []string{"your-api-key"},
        },
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 8. Monitoring

```go
import (
    "log"
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

if err != nil {
    log.Printf("SMS failed: %v, trying email...", err)

    emailMsg := email.New().
        To("user@example.com").
        Subject("Alert").
        Content("Alert message").
        Build()

    err = sender.Send(ctx, emailMsg)
}

func (t *MonitoredTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    start := time.Now()

    // Log request start
    log.Printf("Request started: %s %s", req.Method, req.URL)

    resp, err := t.base.RoundTrip(req)

    // Log request completed
    duration := time.Since(start)
    log.Printf("Request completed: %s %s in %v", req.Method, req.URL, duration)

    return resp, err
}

client := &http.Client{
    Transport: &MonitoredTransport{
        base: http.DefaultTransport,
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 9. Caching

```go
import (
    "net/http"
    "sync"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Simple in-memory cache implementation
type CachedTransport struct {
    base   http.RoundTripper
    cache  map[string]cacheEntry
    mu     sync.RWMutex
}

type cacheEntry struct {
    response *http.Response
    expires  time.Time
}

func (t *CachedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    cacheKey := req.URL.String()

    // Check cache
    t.mu.RLock()
    if entry, exists := t.cache[cacheKey]; exists && time.Now().Before(entry.expires) {
        t.mu.RUnlock()
        return entry.response, nil
    }
    t.mu.RUnlock()

    // Execute actual request
    resp, err := t.base.RoundTrip(req)
    if err == nil {
        // Cache response for 5 minutes
        t.mu.Lock()
        t.cache[cacheKey] = cacheEntry{
            response: resp,
            expires:  time.Now().Add(5 * time.Minute),
        }
        t.mu.Unlock()
    }

    return resp, err
}

client := &http.Client{
    Transport: &CachedTransport{
        base:  http.DefaultTransport,
        cache: make(map[string]cacheEntry),
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### Comprehensive Example - Production Configuration

```go
import (
    "crypto/tls"
    "net/http"
    "net/url"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Production configuration
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        // Proxy settings
        Proxy: http.ProxyURL(&url.URL{
            Scheme: "http",
            Host:   "proxy.company.com:8080",
        }),

        // TLS settings
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: false,
            MinVersion:         tls.VersionTLS12,
        },

        // Connection pool settings
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,

        // Other optimizations
        DisableCompression: false,
        ForceAttemptHTTP2:  true,
    },
}

// Use custom retry policy
retryPolicy := &core.RetryPolicy{
    MaxAttempts:   3,
    InitialDelay:  1 * time.Second,
    MaxDelay:      10 * time.Second,
    BackoffFactor: 2.0,
}

err = sender.Send(ctx, msg,
    core.WithSendHTTPClient(client),
    core.WithSendRetryPolicy(retryPolicy),
)
```

## Advanced: Dynamic Provider Selection

```go
// Use a strategy to select provider based on region or message type
if isInternational(mobile) {
    sender.Send(context.Background(), msg, core.WithSendAccount("yunpian-account"))
} else {
    sender.Send(context.Background(), msg, core.WithSendAccount("aliyun-account"))
}
```
