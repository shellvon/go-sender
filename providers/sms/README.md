[⬅️ Back to Main README](../../README.md)

# SMS Providers

**⚠️ Warning: This project is under active development. APIs are not stable and may change without notice.**

This package provides SMS (Short Message Service) functionality with support for multiple SMS service providers.

## Supported SMS Providers

| Provider       | Description               | Website                                                    |
| -------------- | ------------------------- | ---------------------------------------------------------- |
| **Aliyun**     | Alibaba Cloud SMS Service | [https://www.aliyun.com](https://www.aliyun.com)           |
| **Tencent**    | Tencent Cloud SMS Service | [https://cloud.tencent.com](https://cloud.tencent.com)     |
| **CL253**      | CL253 SMS Service         | [http://www.cl253.com](http://www.cl253.com)               |
| **Huawei**     | Huawei Cloud SMS Service  | [https://www.huaweicloud.com](https://www.huaweicloud.com) |
| **Juhe**       | Juhe SMS Service          | [https://www.juhe.cn](https://www.juhe.cn)                 |
| **Luosimao**   | Luosimao SMS Service      | [https://luosimao.com](https://luosimao.com)               |
| **Smsbao**     | Smsbao SMS Service        | [https://www.smsbao.com](https://www.smsbao.com)           |
| **Submail**    | Submail SMS Service       | [https://www.submail.cn](https://www.submail.cn)           |
| **UCP**        | UCP SMS Service           | [https://www.ucpaas.com](https://www.ucpaas.com)           |
| **Volc**       | Volcengine SMS Service    | [https://www.volcengine.com](https://www.volcengine.com)   |
| **Yuntongxun** | Yuntongxun SMS Service    | [https://www.yuntongxun.com](https://www.yuntongxun.com)   |
| **Yunpian**    | Yunpian SMS Service       | [https://www.yunpian.com](https://www.yunpian.com)         |

## Message Types

All SMS providers support the following message types:

- **Text SMS**: Standard text messages
- **Voice SMS**: Voice messages (supported by most providers)
- **MMS**: Multimedia messages (supported by some providers)

## Quick Start (Chainable Builder Style)

```go
import "github.com/shellvon/go-sender/providers/sms"

// Aliyun SMS Example
msg := sms.Aliyun().
    To("13800138000").
    Content("Hello from go-sender!").
    TemplateID("SMS_xxx").
    Build()

// Tencent SMS Example
msg := sms.Tencent().
    To("13800138000").
    Content("Your code is: 5678").
    TemplateID("123456").
    Sign("YourSign").
    Build()

// CL253 Example (with platform-specific params)
msg := sms.CL253().
    To("13800138000").
    Content("Test message").
    TDFlag(1). // int, see CL253 API doc
    Build()
```

## Configuration

Each SMS provider requires specific configuration including API credentials, endpoints, and other provider-specific settings. Please refer to the individual provider documentation for detailed configuration instructions.

## Features

- **Multi-provider support**: Switch between different SMS providers easily
- **Template support**: Use pre-approved message templates
- **International SMS**: Support for international phone numbers
- **Callback support**: Receive delivery status notifications
- **Batch sending**: Send messages to multiple recipients
- **Error handling**: Comprehensive error handling and retry mechanisms

## SendVia Usage

- `SendVia` is used to specify the account for sending (e.g., different sub-accounts or API keys for the same provider).
- It cannot be used to send the same message object across different message types or providers.
- Example:

```go
// Correct usage: specify account for the same message type
err := sender.SendVia("aliyun-account-1", msg)
if err != nil {
    _ = sender.SendVia("aliyun-account-2", msg)
}
```

## Development Status

This project is currently under active development. Please note:

- APIs may change without notice
- Some features may be incomplete or experimental
- Documentation may be outdated
- Breaking changes are expected in future releases

For production use, please ensure you thoroughly test the APIs and monitor for updates.
