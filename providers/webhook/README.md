# Webhook Provider

> Unified HTTP webhook messaging for any endpoint.

[⬅️ Back to project README](../../README.md)

---

## Features

- Any HTTP method: `GET` / `POST` / `PUT` / `PATCH` / `DELETE`
- Path parameters & query parameters for dynamic URLs
- Custom headers & raw body support (JSON / form / plain text / binary)
- Multiple endpoints with built-in **round-robin / weighted / random** strategies
- Flexible response validation by `StatusCode` or **JSON / XML / Text** content
- Seamless integration with framework middlewares: retry, rate-limit, circuit-break, queue, metrics

---

## Quick Start

```go
import (
    "context"
    "net/http"

    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

cfg := &webhook.Config{
    ProviderMeta: core.ProviderMeta{            // generic meta
        Strategy: core.StrategyRoundRobin,      // pick endpoints round-robin
    },
    Items: []*webhook.Endpoint{
        {
            Name:   "primary",
            URL:    "https://api.example.com/webhook",
            Method: http.MethodPost,            // default POST if empty
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
            ResponseConfig: &webhook.ResponseConfig{ // optional
                CheckBody: true,
                BodyType:  core.BodyTypeJSON,
                Path:      "status",
                Expect:    "ok",
                Mode:      core.MatchEq,         // eq / contains / regex / gt etc.
            },
        },
    },
}

provider, _ := webhook.New(cfg)

msg := webhook.Webhook().
    Body([]byte(`{"foo":"bar"}`)).
    Header("Content-Type", "application/json").
    Build()

_ = provider.Send(context.Background(), msg, nil)
```

---

## Builder API

```go
msg := webhook.Webhook().
    Method(http.MethodPatch).        // custom HTTP method
    PathParam("user_id", "123").     // replace {user_id} in URL
    Query("token", "abcdef").        // append ?token=abcdef
    Header("X-Trace-Id", "uuid").
    Body([]byte(`{"name":"Alice"}`)).
    Build()
```

> Final URL = `Endpoint.URL` + PathParam + Query.

---

## Integrate with `Sender`

```go
sender := gosender.NewSender()

provider, _ := webhook.New(cfg)
sender.RegisterProvider(core.ProviderTypeWebhook, provider, nil)

if err := sender.Send(context.Background(), msg); err != nil {
    log.Fatal(err)
}
```

### Specify Endpoint with SendVia

```go
_ = sender.SendVia(context.Background(), "primary", msg)
```

`SendVia` only selects an endpoint **inside the Webhook provider**. It does not let you reuse a message across providers.

---

## ResponseConfig Reference

| Field          | Type                | Description                                                         |
| -------------- | ------------------- | ------------------------------------------------------------------- |
| `AcceptStatus` | `[]int`             | Additional HTTP status codes to accept; empty slice ⇒ any 2xx is OK |
| `CheckBody`    | `bool`              | Whether to parse & validate the response body                       |
| `BodyType`     | `core.BodyType`     | `json` / `xml` / `text` / `raw` / `none` (auto-detect)              |
| `Path`         | `string`            | JSON/XML path or regex; empty = whole body                          |
| `Expect`       | `string`            | Expected value / regex pattern / numeric threshold                  |
| `Mode`         | `core.MatchMode`    | `eq` / `contains` / `regex` / `gt` / `gte` / `lt` / `lte`           |
| `CodePath`     | `string`            | (optional) JSON/XML path for error code                             |
| `MsgPath`      | `string`            | (optional) JSON/XML path for error message                          |
| `CodeMap`      | `map[string]string` | (optional) error-code → friendly message mapping                    |

> **Dot notation**: use dots to access nested fields, e.g. `data.status.code`. If a key itself contains a dot, escape it with a backslash: `data\.status`.

---

## Examples

### Slack

```go
cfg := &webhook.Config{Items: []*webhook.Endpoint{
    {
        Name:   "slack",
        URL:    "https://hooks.slack.com/services/XXX/YYY/ZZZ",
        Method: http.MethodPost,
        ResponseConfig: &webhook.ResponseConfig{
            CheckBody: true,
            BodyType:  core.BodyTypeText,
            Path:      "",              // whole text
            Expect:    "ok",           // Slack returns "ok"
        },
    },
}}
```

### Discord

```go
cfg := &webhook.Config{Items: []*webhook.Endpoint{
    {
        Name:   "discord",
        URL:    "https://discord.com/api/webhooks/XXX/YYY",
        Method: http.MethodPost,
        ResponseConfig: &webhook.ResponseConfig{
            CheckBody: true,
            BodyType:  core.BodyTypeJSON,
            Path:      "id",           // success if field exists
            Mode:      core.MatchContains,
            Expect:    "",             // any non-empty value
        },
    },
}}
```

### Bark (iOS)

```go
cfg := &webhook.Config{Items: []*webhook.Endpoint{
    {
        Name:   "bark",
        URL:    "https://api.day.app/YOUR_KEY/Hello",
        Method: http.MethodGet,
    },
}}
```

> More samples can be found in the `examples/` folder and unit tests.

---

## FAQ

1. **Self-signed HTTPS?**  
   Use `Sender.SetDefaultHTTPClient(&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}})`.
2. **PathParam not replaced?**  
   Ensure the placeholder exists in `Endpoint.URL`, e.g. `https://api/{version}/user/{id}`.
3. **Custom retry / rate-limit / queue?**  
   Provide a `core.SenderMiddleware` when registering the provider or call the setters on `Sender`.

---
