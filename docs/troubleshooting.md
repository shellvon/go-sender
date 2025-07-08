# Troubleshooting

## Common Issues

- **Provider returns an error**: Verify API keys, template IDs, and network connectivity. If the problem persists, inspect the provider’s dashboard and the go-sender logs.
- **Rate limit exceeded**: Adjust rate limiter settings or check provider limits.
- **Message not delivered**: Check the provider dashboard, go-sender logs and any callback/webhook endpoints. Some platforms silently drop messages that fail template or signature review.
- **Invalid parameters**: Double-check all required fields for the provider (for example template ID, sign name, recipient format).

## Debugging Tips

- Enable go-sender debug logging.
- Inspect the provider dashboard and API documentation.
- Check go-sender metrics to spot anomalies.
- Reproduce with the smallest possible message to isolate the problem.

## Getting Help

- [Open an issue on GitHub](https://github.com/shellvon/go-sender/issues)
- [Read the full documentation](../README.md)
  Consult the provider’s official API documentation for error codes and troubleshooting advice.
