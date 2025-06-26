# SMS Provider Support Matrix / 短信服务商支持矩阵

> **⚠️ Development Notice / 开发说明**: This module is currently under active development. APIs may be unstable and subject to change. Please use with caution in production environments.
>
> 本模块正在积极开发中，API 可能不稳定且会发生变化。请在生产环境中谨慎使用。
>
> **本 Provider 力求支持所有主流平台的短信、彩信、语音短信三种能力，具体每个服务商的能力支持情况请查阅下方[能力矩阵](./capabilities.md)。**

## Supported Providers / 已支持的提供商

The following SMS service providers are currently implemented based on official documentation, but not all have been fully tested:

以下短信服务提供商目前已按官方文档实现，尚未全部经过实际测试：

| Provider / 提供商            | Website / 官网                                             | API Docs / API 文档                                                                                                                   | Implementation / 实现文件        |
| ---------------------------- | ---------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------- |
| Aliyun / 阿里云              | [https://www.aliyun.com](https://www.aliyun.com)           | [API Docs](https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms)                                        | [aliyun.go](./aliyun.go)         |
| Cl253 (Chuanglan) / 创蓝 253 | [https://www.253.com](https://www.253.com)                 | [API Docs](https://www.253.com/api)                                                                                                   | [cl253.go](./cl253.go)           |
| Huawei Cloud / 华为云        | [https://www.huaweicloud.com](https://www.huaweicloud.com) | [API Docs](https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html)                                                    | [huawei.go](./huawei.go)         |
| Juhe / 聚合数据              | [https://www.juhe.cn](https://www.juhe.cn)                 | [API Docs](https://www.juhe.cn/docs)                                                                                                  | [juhe.go](./juhe.go)             |
| Luosimao / 螺丝帽            | [https://luosimao.com](https://luosimao.com)               | [API Docs](https://luosimao.com/docs)                                                                                                 | [luosimao.go](./luosimao.go)     |
| Smsbao / 短信宝              | [https://www.smsbao.com](https://www.smsbao.com)           | [API Docs](https://www.smsbao.com/openapi)                                                                                            | [smsbao.go](./smsbao.go)         |
| Submail / 赛邮               | [https://www.mysubmail.com/](https://www.mysubmail.com/)   | [API Docs](https://www.mysubmail.com/documents)                                                                                       | [submail.go](./submail.go)       |
| Tencent / 腾讯云             | [https://cloud.tencent.com](https://cloud.tencent.com)     | [SMS API](https://cloud.tencent.com/document/product/382/55981)<br>[Voice API](https://cloud.tencent.com/document/product/1128/51559) | [tencent.go](./tencent.go)       |
| UCP / 云之讯                 | [https://www.ucpaas.com](https://www.ucpaas.com)           | [API Docs](http://docs.ucpaas.com)                                                                                                    | [ucp.go](./ucp.go)               |
| Volcengine / 火山引擎        | [https://www.volcengine.com](https://www.volcengine.com)   | [API Docs](https://www.volcengine.com/docs/63933)                                                                                     | [volc.go](./volc.go)             |
| Yunpian / 云片               | [https://www.yunpian.com](https://www.yunpian.com)         | [API Docs](https://www.yunpian.com/official/document/sms/zh_CN/domestic_list)                                                         | [yunpian.go](./yunpian.go)       |
| Yuntongxun / 云讯通          | [https://www.yuntongxun.com](https://www.yuntongxun.com)   | [API Docs](https://www.yuntongxun.com/developer-center)                                                                               | [yuntongxun.go](./yuntongxun.go) |

## Planned Providers / 计划支持的提供商

The following providers are planned for future implementation:

以下提供商计划在未来实现：

- NetEase Cloud Communication / 网易云信
- Baidu Cloud SMS / 百度云短信
- Qiniu Cloud SMS / 七牛云短信
- More international providers... / 更多国际提供商...

## Capability Matrix / 能力矩阵

For detailed capability information including SMS, MMS, and Voice support (domestic/international), please see the [Capability Matrix](./capabilities.md).

有关短信、彩信和语音支持（国内/国际）的详细能力信息，请参阅[能力矩阵](./capabilities.md)。

The capability matrix is automatically generated and includes:
能力矩阵是自动生成的，包括：

- ✅ Supported features / 支持的功能
- ❌ Unsupported features / 不支持的功能
- 🚧 Features under development / 开发中的功能
- Detailed notes and limitations for each provider / 每个提供商的详细说明和限制

## Quick Start / 快速开始

```go
import "github.com/shellvon/go-sender/providers/sms"

// Configure SMS provider / 配置短信提供商
config := sms.Config{
    Providers: []sms.SMSProvider{
        {
            Name:      "aliyun",
            Type:      sms.ProviderTypeAliyun,
            AppID:     "your-app-id",
            AppSecret: "your-app-secret",
            SignName:  "your-sign-name",
        },
    },
}

// Create provider / 创建提供商
provider, err := sms.New(config)
if err != nil {
    log.Fatal(err)
}

// Send SMS / 发送短信
msg := &sms.Message{
    Type:        sms.SMSText, // 普通短信
    Category:    sms.CategoryVerification,
    TemplateID:  "SMS_123456789",
    TemplateParams: map[string]string{"code": "123456"},
    Mobiles:     []string{"***REMOVED***"},
    RegionCode:  86,
}

// Send Voice SMS / 发送语音短信（如腾讯云、阿里云等支持的语音验证码/通知）
voiceMsg := &sms.Message{
    Type:        sms.Voice, // 语音短信
    Category:    sms.CategoryVerification, // 或 sms.CategoryNotification
    TemplateID:  "123456",
    TemplateParams: map[string]string{"code": "654321"},
    Mobiles:     []string{"***REMOVED***"},
    RegionCode:  86,
}

err = provider.Send(context.Background(), msg)
err = provider.Send(context.Background(), voiceMsg)
```

## Contributing / 贡献

If you'd like to add support for a new SMS provider or improve existing implementations, please:

如果您想添加新的短信提供商支持或改进现有实现，请：

1. Check the [capability matrix](./capabilities.md) for current status / 查看[能力矩阵](./capabilities.md)了解当前状态
2. Review existing provider implementations for reference / 参考现有提供商实现
3. Follow the established patterns for capability definition and error handling / 遵循既定的能力定义和错误处理模式
4. Include proper documentation links and website references / 包含适当的文档链接和网站引用
5. Submit a pull request with tests / 提交包含测试的拉取请求

## Notes / 说明

- All providers implement the `SMSProviderInterface` / 所有提供商都实现 `SMSProviderInterface`
- Capabilities are defined using the `Capabilities` struct / 能力使用 `Capabilities` 结构体定义
- Error handling follows consistent patterns across providers / 错误处理遵循跨提供商的一致模式
- International SMS support varies by provider / 国际短信支持因提供商而异
- Voice and MMS support is limited to specific providers / 语音和彩信支持仅限于特定提供商

For detailed implementation information, please refer to individual provider files and their inline documentation.

有关详细实现信息，请参阅各个提供商文件及其内联文档。
