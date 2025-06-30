# Core Concepts

go-sender is designed to be simple, decoupled, and extensible. Here are the key concepts every user should know:

## Sender

The main entry point for all message sending.

- Manages all providers and middleware.
- Exposes `Send` and `SendVia` methods for sending messages.
- Handles middleware chains transparently.

## Provider

A pluggable component that knows how to deliver a message to a specific channel (e.g., Aliyun SMS, WeComBot, EmailJS).

- Each provider implements a unified interface.
- Easy to add your own provider for any channel.

## Message

A data structure representing what you want to send.

- Different providers may have different message types (SMS, Email, IM, etc.).
- All messages support common fields: content, recipients, template, etc.
- Message options allow for flexible configuration.

## Middleware

Reusable components that add cross-cutting features:

- Rate Limiting
- Retry Policy
- Circuit Breaker
- Queue
- Metrics

Middleware is applied as a chain, and you can add or remove them as needed.

## Decorator Pattern

go-sender uses the decorator pattern to wrap providers with middleware, so you can add features without changing your business logic.

## HTTP Transformer

For HTTP-based providers, go-sender uses a flexible transformer architecture:

- Converts your message into an HTTP request.
- Supports custom HTTP clients, headers, authentication, etc.

## Extensibility

- Add new providers by implementing the Provider interface.
- Add new middleware by implementing the Middleware interface.
- Customize HTTP behavior with transformers and custom clients.

---

**Next:**

- [Supported Providers & Usage](./providers.md)
- [How to use Middleware](./middleware.md)
