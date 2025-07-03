[⬅️ Back to Main README](../../README.md)

# SMS Providers | 短信服务商组件

**⚠️ Warning: This project is under active development. APIs are not stable and may change without notice.**

**⚠️ 警告：本项目处于活跃开发中，API 可能随时变更，请勿用于生产环境。**

This package provides SMS (Short Message Service) functionality with support for multiple SMS service providers.

本包支持多家主流短信服务商，提供统一的短信发送接口。

## Supported SMS Providers | 支持的短信服务商

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

## Message Types | 消息类型

All SMS providers support the following message types:

所有短信服务商均支持以下消息类型：

- **Text SMS 文本短信**: Standard text messages | 标准文本短信
- **Voice SMS 语音短信**: Voice messages (supported by most providers) | 语音短信（大部分服务商支持）
- **MMS 彩信**: Multimedia messages (supported by some providers) | 彩信（部分服务商支持）

## Quick Start (Chainable Builder Style) | 快速开始（链式构建）

```go
import "github.com/shellvon/go-sender/providers/sms"

// English: Aliyun SMS Example
// 中文：阿里云短信示例
msg := sms.Aliyun().
    To("13800138000").
    Content("Hello from go-sender! 你好，世界！").
    TemplateID("SMS_xxx").
    Build()

// English: Tencent SMS Example
// 中文：腾讯云短信示例
msg := sms.Tencent().
    To("13800138000").
    Content("Your code is: 5678").
    TemplateID("123456").
    Sign("YourSign").
    Build()

// English: CL253 Example (with platform-specific params)
// 中文：创蓝253短信示例（带平台参数）
msg := sms.CL253().
    To("13800138000").
    Content("Test message").
    TDFlag(1). // int, see CL253 API doc | int，详见 CL253 API 文档
    Build()
```

## Configuration | 配置说明

Each SMS provider requires specific configuration including API credentials, endpoints, and other provider-specific settings. Please refer to the individual provider documentation for detailed configuration instructions.

每个短信服务商都需要特定的配置（API 密钥、接口地址等），请参考各自的官方文档。

## Features | 功能特性

- **Multi-provider support 多服务商支持**: Switch between different SMS providers easily | 可灵活切换不同服务商
- **Template support 模板支持**: Use pre-approved message templates | 支持模板短信
- **International SMS 国际短信**: Support for international phone numbers | 支持国际短信
- **Callback support 回调支持**: Receive delivery status notifications | 支持回执/状态回调
- **Batch sending 批量发送**: Send messages to multiple recipients | 支持批量发送
- **Error handling 错误处理**: Comprehensive error handling and retry mechanisms | 完善的错误处理与重试机制

## SendVia Usage | SendVia 用法

- `SendVia` is used to specify the account for sending (e.g., different sub-accounts or API keys for the same provider).
- It cannot be used to send the same message object across different message types or providers.
- Example:

- `SendVia` 用于指定发送账号（如同一服务商下不同子账号或 API Key）。
- 不能用于跨类型或跨服务商复用同一消息对象。
- 示例：

```go
// English: Specify account for the same message type
// 中文：同一类型消息指定账号发送
err := sender.SendVia("aliyun-account-1", msg)
if err != nil {
    _ = sender.SendVia("aliyun-account-2", msg)
}
```

## Development Status | 开发状态

This project is currently under active development. Please note:

本项目处于活跃开发阶段，请注意：

- APIs may change without notice | API 可能随时变更
- Some features may be incomplete or experimental | 部分功能尚未完善或为实验性
- Documentation may be outdated | 文档可能未及时更新
- Breaking changes are expected in future releases | 未来版本可能有不兼容变更

For production use, please ensure you thoroughly test the APIs and monitor for updates.

如需生产环境使用，请务必充分测试并关注更新。
