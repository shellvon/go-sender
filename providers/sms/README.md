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

## Quick Start

```go
import "github.com/shellvon/go-sender/providers/sms"

// Create a text message using Aliyun
msg := sms.Aliyun().NewTextMessage(
    []string{"***REMOVED***"},
    "Your verification code is: 1234",
)

// Create a voice message
msg := sms.Aliyun().NewVoiceMessage(
    []string{"***REMOVED***"},
    "Your verification code is: 1234",
)

// Create a template message
msg := sms.Aliyun().NewTextMessage(
    []string{"***REMOVED***"},
    "",
    sms.WithTemplateID("SMS_123456"),
    sms.WithTemplateParams(map[string]string{"code": "1234"}),
)
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

## Development Status

This project is currently under active development. Please note:

- APIs may change without notice
- Some features may be incomplete or experimental
- Documentation may be outdated
- Breaking changes are expected in future releases

For production use, please ensure you thoroughly test the APIs and monitor for updates.
