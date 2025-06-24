# go-sender SMS Providers

> **Status: In Development**
>
> This SMS module is under active development. APIs, features, and supported providers may change in the future. Contributions and feedback are welcome!

## Supported SMS Providers

| Provider Name (品牌)           | Website/Docs                                                                                                                  | Source Code                  | Bulk SMS | International |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------------- | ---------------------------- | -------- | ------------- |
| Aliyun (阿里云, Mainland/Intl) | [Aliyun SMS](https://help.aliyun.com/document_detail/419273.html) [Intl](https://help.aliyun.com/document_detail/108084.html) | [aliyun.go](./aliyun.go)     | ✅       | ✅            |
| Huawei Cloud (华为云)          | [Huawei SMS](https://support.huaweicloud.com/intl/en-us/api-msgsms/)                                                          | [huawei.go](./huawei.go)     | ✅       | ✅            |
| Luosimao (螺丝帽)              | [Luosimao](https://luosimao.com/docs/api/)                                                                                    | [luosimao.go](./luosimao.go) | ✅       | ❌            |
| 253 Yun (创蓝 253)             | [253 云通讯](https://doc.253.com/)                                                                                            | [cl253.go](./cl253.go)       | ✅       | ❌            |
| Juhe (聚合数据)                | [Juhe Data](https://www.juhe.cn/docs/api/id/54)                                                                               | [juhe.go](./juhe.go)         | ❌       | ❌            |
| SMSBao (短信宝)                | [SMSBao](https://www.smsbao.com/openapi/213.html)                                                                             | [smsbao.go](./smsbao.go)     | ✅       | ❌            |
| UCP (云之讯)                   | [UCP](https://doc.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sms:index)                                                        | [ucp.go](./ucp.go)           | ✅       | ❌            |
| Yunpian (云片)                 | [Yunpian](https://www.yunpian.com/official/document/sms/zh_CN/domestic_list)                                                  | [yunpian.go](./yunpian.go)   | ✅       | ✅            |
| Submail (赛邮)                 | [Submail](https://www.mysubmail.com/documents/FppOR3)                                                                         | [submail.go](./submail.go)   | ✅       | ✅            |
| Volcengine (火山引擎)          | [Volcengine](https://www.volcengine.com/docs/6361/67380)                                                                      | [volc.go](./volc.go)         | ✅       | ❌            |
| Tencent Cloud SMS (腾讯云)     | [Tencent Cloud SMS](https://cloud.tencent.com/document/product/382/5976)                                                      | (WIP)                        | ✅       | ✅            |

**Legend:**

- ✅: Supported
- ❌: Not supported
- WIP: Work in progress (only for 未实现/部分实现)

- Some providers support only single or bulk SMS, or only domestic/international. See the table above for details.
- This module is under active development. APIs and features may change.
- For questions or contributions, please open an issue or pull request.

---

## Planned / Not Yet Implemented Providers

- Tencent Cloud SMS (腾讯云)
- Ucloud
- Qiniu Cloud (七牛云)
- SendCloud
- Yuntongxun (容联云通讯)
- ihuyi (互亿无线)
- Baidu Cloud (百度云)
- Huaxin (华信短信平台)
- Chuanglan Cloud Intelligence (创蓝云智)
- RongCloud (融云)
- Tianyi Wireless (天毅无线)
- Netease Yunxin (网易云信)
- KXTON (凯信通)
- UE35.net
- Tiniyo
- Moduyun (摩杜云)
- Zhutong (融合云/助通)
- Zhizhuyun (蜘蛛云)
- Ronghe Yunxin (融合云信)
- Tianruiyun (天瑞云)
- Era Interconnect (时代互联)
- CTWing (电信天翼云)
- Twilio (国际)
