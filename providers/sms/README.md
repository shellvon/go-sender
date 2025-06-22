# SMS Provider

This provider supports sending SMS messages via multiple SMS service providers. The implementation references the [smsBomb](https://github.com/opendream/SmsBomb) project.

## Features

- **Multiple Provider Support**: Configure multiple SMS service providers with different strategies (round-robin, random, weighted)
- **Provider Types**: Support for 11 different SMS service providers
- **Template Support**: Support for both template SMS and content SMS
- **Fallback Strategy**: Automatic fallback when one provider fails
- **Signature Support**: Support for SMS signature configuration
- **Based on smsBomb**: Implementation references the [smsBomb](https://github.com/opendream/SmsBomb) project

## Supported Providers

| Provider              | Type       | Official Website                                    | API Documentation                                                                     | Status     |
| --------------------- | ---------- | --------------------------------------------------- | ------------------------------------------------------------------------------------- | ---------- |
| **Tencent Cloud SMS** | `tencent`  | [腾讯云短信](https://cloud.tencent.com/product/sms) | [API 文档](https://cloud.tencent.com/document/product/382)                            | ✅ Active  |
| **Aliyun SMS**        | `aliyun`   | [阿里云短信](https://www.aliyun.com/product/sms)    | [API 文档](https://help.aliyun.com/document_detail/101414.html)                       | ✅ Active  |
| **Yunpian**           | `yunpian`  | [云片网](https://www.yunpian.com/)                  | [API 文档](https://www.yunpian.com/doc/zh_CN/api/single_send.html)                    | ✅ Active  |
| **UCP**               | `ucp`      | [云之讯](https://www.ucpaas.com/)                   | [API 文档](https://www.ucpaas.com/doc/)                                               | ✅ Active  |
| **CL253**             | `cl253`    | [蓝创 253](https://www.cl253.com/)                  | [API 文档](https://www.cl253.com/doc/)                                                | ✅ Active  |
| **SMSBao**            | `smsbao`   | [短信宝](https://www.smsbao.com/)                   | [API 文档](https://www.smsbao.com/openapi/)                                           | ✅ Active  |
| **Juhe**              | `juhe`     | [聚合数据](https://www.juhe.cn/)                    | [API 文档](https://www.juhe.cn/docs/api/id/54)                                        | ✅ Active  |
| **Luosimao**          | `luosimao` | [螺丝帽](https://luosimao.com/)                     | [API 文档](https://luosimao.com/docs/api/)                                            | ✅ Active  |
| **Netease**           | `netease`  | [网易云信](https://yunxin.163.com/)                 | [API 文档](https://dev.yunxin.163.com/docs/product/IM即时通讯/服务端API文档/短信服务) | ✅ Active  |
| **Normal**            | `normal`   | -                                                   | -                                                                                     | ✅ Generic |

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/sms"
)

// Create SMS configuration
config := sms.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyWeighted,
    },
    Providers: []sms.SMSProvider{
        {
            Name:      "tencent",
            Type:      sms.ProviderTypeTencent,
            AppID:     "your_tencent_app_id",
            AppSecret: "your_tencent_app_secret",
            SignName:  "your_sign_name",
            Weight:    100,
            Disabled:  false,
        },
        {
            Name:      "aliyun",
            Type:      sms.ProviderTypeAliyun,
            AppID:     "your_aliyun_access_key",
            AppSecret: "your_aliyun_access_secret",
            SignName:  "your_sign_name",
            Weight:    80,
            Disabled:  false,
        },
    },
}

// Create provider
provider, err := sms.New(config)
if err != nil {
    log.Fatalf("Failed to create SMS provider: %v", err)
}
```

## Message Types

### 1. Content SMS (Non-template)

```go
// Simple content SMS
msg := sms.NewMessage("13800138000",
    sms.WithContent("您的验证码是123456，5分钟内有效"),
    sms.WithSignName("您的签名"),
)
```

### 2. Template SMS

```go
// Template SMS with parameters
msg := sms.NewMessage("13800138000",
    sms.WithTemplateCode("SMS_123456789"),
    sms.WithTemplateParams(map[string]string{
        "code": "123456",
        "time": "5分钟",
    }),
    sms.WithSignName("您的签名"),
)
```

### 3. Advanced Message Configuration

```go
// Message with all options
msg := sms.NewMessage("13800138000",
    sms.WithContent("这是一条测试短信"),
    sms.WithTemplateCode("SMS_TEMPLATE_CODE"),
    sms.WithTemplateParams(map[string]string{
        "param1": "value1",
        "param2": "value2",
    }),
    sms.WithSignName("公司名称"),
)
```

## Provider-Specific Configuration

### Tencent Cloud SMS (腾讯云短信)

**Official Website**: [https://cloud.tencent.com/product/sms](https://cloud.tencent.com/product/sms)  
**API Documentation**: [https://cloud.tencent.com/document/product/382](https://cloud.tencent.com/document/product/382)  
**API Endpoint**: `https://yun.tim.qq.com/v5/tlssmssvr/sendsms`

```go
{
    Name:      "tencent",
    Type:      sms.ProviderTypeTencent,
    AppID:     "your_sdkappid",        // SDK AppID from Tencent Cloud Console
    AppSecret: "your_app_key",         // App Key from Tencent Cloud Console
    SignName:  "your_sign_name",       // SMS signature (approved by Tencent)
    Weight:    100,
}
```

**Notes**:

- Requires SDK AppID and App Key from Tencent Cloud Console
- Supports both template and content SMS
- Template parameters are sent as JSON array
- SignName must be pre-approved by Tencent Cloud

### Aliyun SMS (阿里云短信)

**Official Website**: [https://www.aliyun.com/product/sms](https://www.aliyun.com/product/sms)  
**API Documentation**: [https://help.aliyun.com/document_detail/101414.html](https://help.aliyun.com/document_detail/101414.html)  
**API Endpoint**: `https://dysmsapi.aliyuncs.com`

```go
{
    Name:      "aliyun",
    Type:      sms.ProviderTypeAliyun,
    AppID:     "your_access_key_id",     // Access Key ID from Aliyun Console
    AppSecret: "your_access_key_secret", // Access Key Secret from Aliyun Console
    SignName:  "your_sign_name",         // SMS signature (approved by Aliyun)
    Weight:    100,
}
```

**Notes**:

- Requires Access Key ID and Access Key Secret from Aliyun Console
- Template parameters are sent as JSON string
- SignName is required for all messages and must be pre-approved
- Supports both template and content SMS

### Yunpian (云片网)

**Official Website**: [https://www.yunpian.com/](https://www.yunpian.com/)  
**API Documentation**: [https://www.yunpian.com/doc/zh_CN/api/single_send.html](https://www.yunpian.com/doc/zh_CN/api/single_send.html)  
**API Endpoint**: `https://sms.yunpian.com/v2/sms/single_send.json`

```go
{
    Name:      "yunpian",
    Type:      sms.ProviderTypeYunpian,
    AppID:     "your_api_key",          // API Key from Yunpian Console
    AppSecret: "",                      // Not used for Yunpian
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

**Notes**:

- Requires API Key from Yunpian Console
- Supports both template and content SMS
- Template parameters are sent as key-value pairs
- SignName is included in the message content

### UCP (云之讯)

**Official Website**: [https://www.ucpaas.com/](https://www.ucpaas.com/)  
**API Documentation**: [https://www.ucpaas.com/doc/](https://www.ucpaas.com/doc/)  
**API Endpoint**: `https://open.ucpaas.com/ol/sms/sendsms`

```go
{
    Name:      "ucp",
    Type:      sms.ProviderTypeUcp,
    AppID:     "your_sid",              // SID from UCP Console
    AppSecret: "your_token",            // Token from UCP Console
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

### CL253 (蓝创 253)

**Official Website**: [https://www.cl253.com/](https://www.cl253.com/)  
**API Documentation**: [https://www.cl253.com/doc/](https://www.cl253.com/doc/)  
**API Endpoint**: `https://smssh1.253.com/msg/v1/send/json`

```go
{
    Name:      "cl253",
    Type:      sms.ProviderTypeCl253,
    AppID:     "your_account",          // Account from CL253 Console
    AppSecret: "your_password",         // Password from CL253 Console
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

### SMSBao (短信宝)

**Official Website**: [https://www.smsbao.com/](https://www.smsbao.com/)  
**API Documentation**: [https://www.smsbao.com/openapi/](https://www.smsbao.com/openapi/)  
**API Endpoint**: `https://api.smsbao.com/sms`

```go
{
    Name:      "smsbao",
    Type:      sms.ProviderTypeSmsbao,
    AppID:     "your_username",         // Username from SMSBao Console
    AppSecret: "your_password",         // Password from SMSBao Console
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

### Juhe (聚合数据)

**Official Website**: [https://www.juhe.cn/](https://www.juhe.cn/)  
**API Documentation**: [https://www.juhe.cn/docs/api/id/54](https://www.juhe.cn/docs/api/id/54)  
**API Endpoint**: `http://v.juhe.cn/sms/send`

```go
{
    Name:      "juhe",
    Type:      sms.ProviderTypeJuhe,
    AppID:     "your_key",              // API Key from Juhe Console
    AppSecret: "",                      // Not used for Juhe
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

### Luosimao (螺丝帽)

**Official Website**: [https://luosimao.com/](https://luosimao.com/)  
**API Documentation**: [https://luosimao.com/docs/api/](https://luosimao.com/docs/api/)  
**API Endpoint**: `https://sms-api.luosimao.com/v1/send.json`

```go
{
    Name:      "luosimao",
    Type:      sms.ProviderTypeLuosimao,
    AppID:     "your_api_key",          // API Key from Luosimao Console
    AppSecret: "",                      // Not used for Luosimao
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

### Netease (网易云信)

**Official Website**: [https://yunxin.163.com/](https://yunxin.163.com/)  
**API Documentation**: [https://dev.yunxin.163.com/docs/product/IM 即时通讯/服务端 API 文档/短信服务](https://dev.yunxin.163.com/docs/product/IM即时通讯/服务端API文档/短信服务)  
**API Endpoint**: `https://api.netease.im/sms/sendcode.action`

```go
{
    Name:      "netease",
    Type:      sms.ProviderTypeNetease,
    AppID:     "your_app_key",          // App Key from Netease Console
    AppSecret: "your_app_secret",       // App Secret from Netease Console
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
}
```

### Normal (通用)

For other SMS providers that follow standard HTTP POST/GET patterns:

```go
{
    Name:      "custom",
    Type:      sms.ProviderTypeNormal,
    AppID:     "your_username",         // Username or API Key
    AppSecret: "your_password",         // Password or API Secret
    SignName:  "your_sign_name",        // SMS signature
    Weight:    100,
    ExtraConfig: map[string]string{
        "url": "https://your-sms-provider.com/api/send",
        "method": "POST",
        "param_mobile": "mobile",
        "param_content": "content",
        "param_sign": "sign",
    },
}
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/sms"
)

// Create sender
s := gosender.NewSender(nil)

// Register SMS provider
smsProvider, err := sms.New(config)
if err != nil {
    log.Fatalf("Failed to create SMS provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeSMS, smsProvider, nil)

// Send SMS
ctx := context.Background()
msg := sms.NewMessage("13800138000",
    sms.WithContent("您的验证码是123456，5分钟内有效"),
)
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send SMS: %v", err)
}
```

## Message Options

### Basic Options

- `WithMobile(mobile string)`: Set the mobile phone number
- `WithContent(content string)`: Set SMS content (for non-template SMS)
- `WithTemplateCode(templateCode string)`: Set template code (for template SMS)
- `WithTemplateParams(params map[string]string)`: Set template parameters
- `WithSignName(signName string)`: Set SMS signature name

### Message Validation

- Mobile number: Required, must be valid phone number
- Content or TemplateCode: At least one must be provided
- TemplateParams: Required when using TemplateCode
- SignName: Optional, can be set at provider level

## API Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Providers`: Array of SMS provider configurations

### SMSProvider

- `Name`: Provider name for identification
- `Type`: SMS service provider type
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this provider is disabled
- `AppID`: App ID/Account/API Key
- `AppSecret`: App Secret/Password/Token
- `SignName`: SMS signature name
- `TemplateCode`: Default template code
- `TemplateParams`: Default template parameters
- `ExtraConfig`: Extra configuration for specific providers

### Message

- `Mobile`: Mobile phone number
- `Content`: SMS content (for non-template SMS)
- `TemplateCode`: Template code (for template SMS)
- `TemplateParams`: Template parameters
- `SignName`: SMS signature name

## Error Handling

The provider automatically handles:

- Provider selection based on strategy
- Fallback to alternative providers on failure
- Rate limiting and retry logic
- Provider-specific error handling

## Rate Limits

Each SMS provider has its own rate limits:

- Respect provider-specific rate limits
- Use multiple providers for high-volume sending
- Implement appropriate retry strategies

## Security Notes

- Store AppID and AppSecret securely
- Use environment variables for sensitive data
- Rotate credentials regularly
- Monitor SMS usage and costs

## Credits and Acknowledgments

This SMS provider implementation references the [smsBomb](https://github.com/opendream/SmsBomb) project, which is another project of the same author.

**Original Project**: [smsBomb](https://github.com/opendream/SmsBomb)  
**License**: [MIT License](https://github.com/opendream/SmsBomb/blob/master/LICENSE)

## API Documentation Links

For detailed API documentation, visit the official websites:

- [Tencent Cloud SMS](https://cloud.tencent.com/document/product/382)
- [Aliyun SMS](https://help.aliyun.com/document_detail/101414.html)
- [Yunpian SMS](https://www.yunpian.com/doc/zh_CN/api/single_send.html)
- [UCP SMS](https://www.ucpaas.com/doc/)
- [CL253 SMS](https://www.cl253.com/doc/)
- [SMSBao](https://www.smsbao.com/openapi/)
- [Juhe SMS](https://www.juhe.cn/docs/api/id/54)
- [Luosimao SMS](https://luosimao.com/docs/api/)
- [Netease SMS](https://dev.yunxin.163.com/docs/product/IM即时通讯/服务端API文档/短信服务)
