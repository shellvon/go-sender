# Troubleshooting

## Common Issues

- **Provider returns error**: Check API keys, template codes, and network connectivity. Make sure your credentials and templates are correct. If 问题持续，查看 provider 官方后台和 go-sender 日志。
- **Rate limit exceeded**: Adjust rate limiter settings or check provider limits.
- **Message not delivered**: 检查 provider 官方后台、go-sender 日志、callback URLs。部分平台未通过模板/签名会静默丢弃消息。
- **Invalid parameters**: Double-check required fields for each provider (e.g., template code, sign name, recipient format)。

## Debugging Tips

- 启用 go-sender 日志。
- 查看 provider 官方后台和 API 文档。
- 检查 go-sender 的 Metrics。
- 发送最小化消息定位问题。

## Getting Help

- [Open an issue on GitHub](https://github.com/shellvon/go-sender/issues)
- [Read the full documentation](../README.md)
- 查阅 provider 官方 API 文档获取错误码和排查建议。
