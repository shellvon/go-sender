# Examples

Explore real-world usage patterns with go-sender.

## Multi-channel Fallback

```go
// Try first SMS account, then fallback to second SMS account if fails
err := sender.SendVia("aliyun-account-1", msg)
if err != nil {
    _ = sender.SendVia("aliyun-account-2", msg)
}
```

## Batch Sending

### Method 1: Provider-Supported Batch Sending

Some providers support sending to multiple recipients in a single API call:

```go
// SMS providers that support multiple recipients
msg := sms.Aliyun().
    To("***REMOVED***,13800138001,13800138002").  // Multiple numbers separated by comma
    Content("Hello everyone!").
    Build()

err := sender.Send(ctx, msg)
```

### Method 2: Manual Loop for Providers Without Batch Support

For providers that don't support batch sending, use a loop:

```go
mobiles := []string{"***REMOVED***", "13800138001", "13800138002"}

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

mobiles := []string{"***REMOVED***", "13800138001", "13800138002"}
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
    core.WithSendCallback(func(sendErr error) {
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

### 1. Timeout Control - 防止请求挂起

```go
import (
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// 设置30秒超时，防止请求无限等待
client := &http.Client{
    Timeout: 30 * time.Second,
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 2. Proxy Support - 企业代理支持

```go
import (
    "net/http"
    "net/url"
    "github.com/shellvon/go-sender/core"
)

// 通过企业代理服务器路由请求
proxyURL, _ := url.Parse("http://proxy.company.com:8080")
client := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 3. TLS Configuration - 自定义证书和安全设置

```go
import (
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "os"
    "github.com/shellvon/go-sender/core"
)

// 加载自定义CA证书
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

### 4. Connection Pooling - 连接池优化

```go
import (
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// 优化高并发场景下的连接复用
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,              // 最大空闲连接数
        MaxIdleConnsPerHost: 10,               // 每个主机最大空闲连接数
        IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
        DisableCompression:  false,            // 启用压缩
    },
}

err = sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### 5. Retry Logic - 内置重试机制

```go
import (
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// 自定义重试策略
retryPolicy := &core.RetryPolicy{
    MaxAttempts:   3,
    InitialDelay:  1 * time.Second,
    MaxDelay:      10 * time.Second,
    BackoffFactor: 2.0,
}

client := &http.Client{
    Timeout: 30 * time.Second,
}

// 使用自定义重试策略
err = sender.Send(ctx, msg,
    core.WithSendHTTPClient(client),
    core.WithSendRetryPolicy(retryPolicy),
)
```

### 6. Load Balancing - 负载均衡

```go
import (
    "net/http"
    "net/url"
    "github.com/shellvon/go-sender/core"
)

// 自定义Transport实现负载均衡
type LoadBalancedTransport struct {
    endpoints []string
    current   int
}

func (t *LoadBalancedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // 轮询选择端点
    endpoint := t.endpoints[t.current%len(t.endpoints)]
    t.current++

    // 修改请求URL
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

### 7. Authentication - 自定义认证

```go
import (
    "net/http"
    "github.com/shellvon/go-sender/core"
)

// 添加自定义认证头
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

### 8. Monitoring - 请求监控

```go
import (
    "log"
    "net/http"
    "time"
    "github.com/shellvon/go-sender/core"
)

// 自定义Transport添加监控
type MonitoredTransport struct {
    base http.RoundTripper
}

func (t *MonitoredTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    start := time.Now()

    // 记录请求开始
    log.Printf("Request started: %s %s", req.Method, req.URL)

    resp, err := t.base.RoundTrip(req)

    // 记录请求完成
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

### 9. Caching - 响应缓存

```go
import (
    "net/http"
    "sync"
    "time"
    "github.com/shellvon/go-sender/core"
)

// 简单的内存缓存实现
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

    // 检查缓存
    t.mu.RLock()
    if entry, exists := t.cache[cacheKey]; exists && time.Now().Before(entry.expires) {
        t.mu.RUnlock()
        return entry.response, nil
    }
    t.mu.RUnlock()

    // 执行实际请求
    resp, err := t.base.RoundTrip(req)
    if err == nil {
        // 缓存响应（5分钟）
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

### 综合示例 - 生产环境配置

```go
import (
    "crypto/tls"
    "net/http"
    "net/url"
    "time"
    "github.com/shellvon/go-sender/core"
)

// 生产环境综合配置
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        // 代理配置
        Proxy: http.ProxyURL(&url.URL{
            Scheme: "http",
            Host:   "proxy.company.com:8080",
        }),

        // TLS配置
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: false,
            MinVersion:         tls.VersionTLS12,
        },

        // 连接池配置
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,

        // 其他优化
        DisableCompression: false,
        ForceAttemptHTTP2:  true,
    },
}

// 使用自定义重试策略
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
    sender.SendVia("yunpian-account", msg)
} else {
    sender.SendVia("aliyun-account", msg)
}
```
