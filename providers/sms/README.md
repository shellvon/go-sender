# go-sender SMS Providers

> **Status: In Development**
>
> This SMS module is under active development. APIs, features, and supported providers may change in the future. Contributions and feedback are welcome!

## Supported SMS Providers

| Provider Name (品牌)      | Website/Docs                                                           | Source Code                          | Bulk SMS | International |
| ------------------------- | ---------------------------------------------------------------------- | ------------------------------------ | -------- | ------------- |
| Aliyun (阿里云, Mainland) | [Aliyun SMS](https://help.aliyun.com/document_detail/419273.html)      | [aliyun.go](./aliyun.go)             | ✅       | ❌            |
| Aliyun Intl (阿里云国际)  | [Aliyun Intl](https://help.aliyun.com/document_detail/108146.html)     | [aliyun_globe.go](./aliyun_globe.go) | ✅       | ✅            |
| Huawei Cloud (华为云)     | [Huawei SMS](https://support.huaweicloud.com/intl/en-us/api-msgsms/)   | [huawei.go](./huawei.go)             | ✅       | ✅            |
| Luosimao (螺丝帽)         | [Luosimao](https://luosimao.com/docs/api/)                             | [luosimao.go](./luosimao.go)         | ✅       | ❌            |
| 253 Yun (创蓝 253)        | [253 云通讯](https://doc.253.com/)                                     | [cl253.go](./cl253.go)               | ✅       | ❌ (WIP)      |
| Juhe (聚合数据)           | [Juhe Data](https://www.juhe.cn/docs/api/id/54)                        | [juhe.go](./juhe.go)                 | ❌       | ❌            |
| SMSBao (短信宝)           | [SMSBao](https://www.smsbao.com/openapi/213.html)                      | [smsbao.go](./smsbao.go)             | ✅       | ❌            |
| UCP (云之讯)              | [UCP](https://doc.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sms:index) | [ucp.go](./ucp.go)                   | ✅       | ❌            |

**Legend:**

- ✅ = Supported
- ❌ = Not Supported
- WIP = Work in Progress

## Planned/Not Yet Implemented Providers

- Tencent Cloud SMS (腾讯云) ([docs](https://cloud.tencent.com/document/product/382/5976))
- Yunpian (云片) ([docs](https://www.yunpian.com/official/document/sms/zh_CN/))
- Submail (赛邮) ([docs](https://www.mysubmail.com/chs/documents/developer/index))
- Volcano Engine (火山引擎) ([docs](https://www.volcengine.com/docs/6348/70138))

## Notes

- Some providers support only single or bulk SMS, or only domestic/international. See the table above for details.
- This module is under active development. APIs and features may change.
- For questions or contributions, please open an issue or pull request.
