# Examples

Explore real-world usage patterns with go-sender.

## Multi-channel Fallback

```go
// Try SMS, then fallback to WeComBot if SMS fails
err := sender.SendVia("aliyun", msg)
if err != nil {
    _ = sender.SendVia("wecombot", msg)
}
```

## Batch Sending

```go
for _, mobile := range mobiles {
    msg := sms.Aliyun().NewTextMessage([]string{mobile}, "Hello", ...)
    _ = sender.Send(ctx, msg)
}
```

## Asynchronous Queue

```go
sender.SetQueue(myQueue)
sender.Send(ctx, msg) // Will be enqueued and sent asynchronously
```

## Custom HTTP Client

```go
sender.Send(ctx, msg, core.WithSendHTTPClient(myClient))
```

## Advanced: Dynamic Provider Selection

```go
// Use a strategy to select provider based on region or message type
if isInternational(mobile) {
    sender.SendVia("yunpian", msg)
} else {
    sender.SendVia("aliyun", msg)
}
```
