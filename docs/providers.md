# Provider Reference

go-sender supports a comprehensive ecosystem of notification providers across multiple categories. This reference provides an overview of capabilities, status, and integration details.

> **Quick Start**: Looking for usage examples? See [Getting Started](./getting-started.md) for hands-on tutorials.

## Provider Ecosystem Overview

| **Category** | **Providers** | **Key Features** | **Use Cases** |
|-------------|---------------|------------------|---------------|
| **SMS** | Aliyun, Tencent, Huawei, CL253, Volc | Template support, signature management, multi-region | OTP, alerts, marketing |
| **Email** | SMTP, EmailJS, Resend, Mailgun | Direct SMTP or API-based delivery | Transactional, newsletters |
| **IM/Bot** | WeCom, DingTalk, Lark, Telegram | Rich media, markdown, interactive cards | Team notifications, bots |
| **Webhook** | Generic HTTP, Custom APIs | Universal integration for any HTTP API | Custom integrations, push services |

**Status Legend:**
- ✅ **Production Ready**: Fully tested and production-ready
- 🚧 **In Development**: Functional but may have limitations  
- 📋 **Planned**: On roadmap for future implementation

## SMS & Voice

| Provider                  | Website                                        | API Docs                                                                                                                             | Provider Doc                                    |
| ------------------------- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------- |
| Aliyun (阿里云)           | [aliyun.com](https://www.aliyun.com)           | [API](https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms)                                            | [Aliyun README](../providers/sms/README.md)     |
| Tencent Cloud (腾讯云)    | [cloud.tencent.com](https://cloud.tencent.com) | [SMS API](https://cloud.tencent.com/document/product/382/55981) / [Voice API](https://cloud.tencent.com/document/product/1128/51559) | [Tencent README](../providers/sms/README.md)    |
| Huawei Cloud (华为云)     | [huaweicloud.com](https://www.huaweicloud.com) | [API](https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html)                                                        | [Huawei README](../providers/sms/README.md)     |
| Volcano Engine (火山引擎) | [volcengine.com](https://www.volcengine.com)   | [API](https://www.volcengine.com/docs/63933)                                                                                         | [Volc README](../providers/sms/README.md)       |
| Yunpian (云片)            | [yunpian.com](https://www.yunpian.com)         | [API](https://www.yunpian.com/official/document/sms/zh_CN/domestic_list)                                                             | [Yunpian README](../providers/sms/README.md)    |
| CL253 (创蓝 253)          | [253.com](https://www.253.com)                 | [API](https://www.253.com/api)                                                                                                       | [CL253 README](../providers/sms/README.md)      |
| Submail (赛邮)            | [mysubmail.com](https://www.mysubmail.com/)    | [API](https://www.mysubmail.com/documents)                                                                                           | [Submail README](../providers/sms/README.md)    |
| UCP (云之讯)              | [ucpaas.com](https://www.ucpaas.com)           | [API](http://docs.ucpaas.com)                                                                                                        | [UCP README](../providers/sms/README.md)        |
| Juhe (聚合数据)           | [juhe.cn](https://www.juhe.cn)                 | [API](https://www.juhe.cn/docs)                                                                                                      | [Juhe README](../providers/sms/README.md)       |
| SMSBao (短信宝)           | [smsbao.com](https://www.smsbao.com)           | [API](https://www.smsbao.com/openapi)                                                                                                | [SMSBao README](../providers/sms/README.md)     |
| Yuntongxun (云讯通)       | [yuntongxun.com](https://www.yuntongxun.com)   | [API](https://www.yuntongxun.com/developer-center)                                                                                   | [Yuntongxun README](../providers/sms/README.md) |

## Email Providers

| **Provider** | **Protocol** | **Key Features** | **Status** | **Documentation** |
|-------------|--------------|------------------|------------|-------------------|
| **SMTP** | SMTP | Direct server connection, full control | ✅ Production | [Email README](../providers/email/README.md) |
| **EmailJS** | HTTP API | Browser-friendly, template support | ✅ Production | [EmailAPI README](../providers/emailapi/README.md) |
| **Resend** | HTTP API | Developer-focused, webhooks | ✅ Production | [EmailAPI README](../providers/emailapi/README.md) |
| **Mailgun** | HTTP API | Enterprise features, analytics | ✅ Production | [EmailAPI README](../providers/emailapi/README.md) |

## IM/Bot/Enterprise Notification

- WeCom App (企业微信应用) ([README](../providers/wecomapp/README.md))
- WeCom Bot (企业微信机器人) ([README](../providers/wecombot/README.md))
- DingTalk Bot (钉钉机器人) ([README](../providers/dingtalk/README.md))
- Lark/Feishu (飞书/国际版) ([README](../providers/lark/README.md))
- Telegram ([README](../providers/telegram/README.md))
- ServerChan ([README](../providers/serverchan/README.md))

## Webhook/Universal Push

- Webhook (generic HTTP integration) ([README](../providers/webhook/README.md))

  | Service    | Website                                          | Docs / API Refrence                                                 |
  | ---------- | ------------------------------------------------ | ------------------------------------------------------------------- |
  | ntfy       | [ntfy.sh](https://ntfy.sh/)                      | [Docs](https://docs.ntfy.sh/publish/)                               |
  | Bark       | [Bark](https://github.com/Finb/Bark)             | [API](https://github.com/Finb/Bark#http-api)                        |
  | PushDeer   | [PushDeer](https://github.com/easychen/pushdeer) | [API](https://github.com/easychen/pushdeer#api)                     |
  | PushPlus   | [PushPlus](https://pushplus.hxtrip.com/)         | [Docs](https://pushplus.hxtrip.com/message.html)                    |
  | IFTTT      | [IFTTT](https://ifttt.com/)                      | [Webhooks](https://ifttt.com/maker_webhooks)                        |
  | PushAll    | [PushAll](https://pushall.ru/)                   | [API](https://pushall.ru/api/)                                      |
  | PushBack   | [PushBack](https://pushback.io/)                 | [Docs](https://docs.pushback.io/)                                   |
  | Pushy      | [Pushy](https://pushy.me/)                       | [API](https://pushy.me/docs/api/send-notifications)                 |
  | Pushbullet | [Pushbullet](https://www.pushbullet.com/)        | [API](https://docs.pushbullet.com/#create-push)                     |
  | Gotify     | [Gotify](https://gotify.net/)                    | [API](https://gotify.net/docs/api/push/)                            |
  | OneBot     | [OneBot](https://github.com/botuniverse/onebot)  | [Docs](https://github.com/botuniverse/onebot/blob/master/README.md) |
  | Push       | [Push](https://push.techulus.com/)               | [API](https://docs.push.techulus.com/)                              |
  | Pushjet    | [Pushjet](https://pushjet.io/)                   | [API](https://pushjet.io/docs/api/)                                 |
  | Pushsafer  | [Pushsafer](https://www.pushsafer.com/)          | [API](https://www.pushsafer.com/en/pushapi)                         |
  | Pushover   | [Pushover](https://pushover.net/)                | [API](https://pushover.net/api)                                     |
  | Simplepush | [Simplepush](https://simplepush.io/)             | [API](https://simplepush.io/api)                                    |
  | Zulip      | [Zulip](https://zulip.com/)                      | [API](https://zulip.com/api/send-message)                           |
  | Mattermost | [Mattermost](https://mattermost.com/)            | [API](https://api.mattermost.com/)                                  |
  | Discord    | [Discord](https://discord.com/)                  | [Webhooks](https://discord.com/developers/docs/resources/webhook)   |

> **Looking for a provider that is not listed?**
> - Use the generic **Webhook provider** for any HTTP API
> - Create a **custom provider** - see [Advanced Usage: Custom Providers](./advanced.md#custom-providers) for the complete guide

---

_See each provider’s README for configuration examples and advanced usage._
