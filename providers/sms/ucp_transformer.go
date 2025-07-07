package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shellvon/go-sender/core"
)

// @ProviderName: UCP / 云之讯
// @Website: https://www.ucpaas.com
// @APIDoc: http://docs.ucpaas.com
//
// 官方文档:
//   - 短信API文档: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
//   - 国内外指定模版单发: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms
//   - 国内外指定模版群发: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms_batch
//
// 能力说明:
//   - 国内短信：支持单发和群发，需模板ID。
//   - 国际短信：支持单发和群发，需模板ID。
//   - 彩信/语音：暂不支持。
//
// 注意：支持国内外手机号码，需模板ID。

const (
	ucpDefaultBaseURI = "http://open2.ucpaas.com/sms-server"
	ucpSingleAPI      = "variablesms"
	ucpBatchAPI       = "templatesms"
)

type ucpTransformer struct {
	*BaseTransformer
}

func newUcpTransformer() *ucpTransformer {
	transformer := &ucpTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(core.ProviderTypeSMS),
		string(SubProviderUcp),
		&core.ResponseHandlerConfig{
			SuccessField:      "code",
			SuccessValue:      "000000",
			ErrorCodeField:    "code",
			ErrorMessageField: "msg",
			ErrorField:        "code",
		},
		WithSMSHandler(transformer.transformSMS),
	)
	return transformer
}

func init() {
	RegisterTransformer(string(SubProviderUcp), newUcpTransformer())
}

// transformSMS transforms SMS message to HTTP request
//   - 短信API: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
func (t *ucpTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, NewProviderError(string(SubProviderUcp), "MISSING_PARAM", "mobiles is required")
	}
	if msg.TemplateID == "" {
		return nil, nil, NewProviderError(string(SubProviderUcp), "MISSING_PARAM", "templateID is required")
	}

	// 根据手机号数量选择API
	var apiPath string
	if len(msg.Mobiles) > 1 {
		apiPath = ucpBatchAPI
	} else {
		apiPath = ucpSingleAPI
	}

	params := map[string]interface{}{
		"clientid":   account.APIKey,
		"password":   account.APISecret,
		"templateid": msg.TemplateID,
		"mobile":     strings.Join(msg.Mobiles, ","),
		"uid":        msg.UID,
	}

	// 模板参数处理
	if len(msg.ParamsOrder) > 0 {
		// 模板中的替换参数，如该模板不存在参数则无需传该参数或者参数为空，如果有多个参数则需要写在同一个字符串中，以分号分隔 （如："a;b;c"），参数中不能含有特殊符号"【】"和","
		params["param"] = strings.Join(msg.ParamsOrder, ";")
	}

	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal ucp request body: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      fmt.Sprintf("%s/%s", ucpDefaultBaseURI, apiPath),
		Body:     bodyData,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}
