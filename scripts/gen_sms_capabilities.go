//go:build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/shellvon/go-sender/providers/sms"
)

var verbose = flag.Bool("verbose", false, "Enable verbose output")

type ProviderMeta struct {
	Name           string
	Website        string
	APIDoc         string
	Implementation string
	NewFunc        func() sms.SMSProviderInterface
	Constructor    string
	Type           sms.ProviderType
	TypeConst      string
}

// SMSProviderConfigs 定义所有SMS providers的配置
var SMSProviderConfigs = []struct {
	Type        string
	TypeConst   string
	Constructor string
	Name        string
	Enabled     bool
}{
	{"aliyun", "ProviderTypeAliyun", "NewAliyunProvider", "阿里云", true},
	{"cl253", "ProviderTypeCl253", "NewCl253Provider", "创蓝253", true},
	{"yuntongxun", "ProviderTypeYuntongxun", "NewYuntongxunProvider", "云通讯", true},
	{"juhe", "ProviderTypeJuhe", "NewJuheProvider", "聚合数据", true},
	{"huawei", "ProviderTypeHuawei", "NewHuaweiProvider", "华为云", true},
	{"volc", "ProviderTypeVolc", "NewVolcProvider", "火山引擎", true},
	{"smsbao", "ProviderTypeSmsbao", "NewSmsbaoProvider", "短信宝", true},
	{"submail", "ProviderTypeSubmail", "NewSubmailProvider", "Submail", true},
	{"ucp", "ProviderTypeUcp", "NewUcpProvider", "UCP", true},
	{"luosimao", "ProviderTypeLuosimao", "NewLuosimaoProvider", "螺丝帽", true},
	{"yunpian", "ProviderTypeYunpian", "NewYunpianProvider", "云片", true},
	{"tencent", "ProviderTypeTencent", "NewTencentProvider", "腾讯云", true},
}

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)
	scriptDir := filepath.Dir(filename)
	smsDir := filepath.Join(scriptDir, "../providers/sms")

	outputPath := filepath.Join(smsDir, "capabilities.md")
	absolutePath, _ := filepath.Abs(outputPath)

	providers := scanProviderFiles()
	generateCapabilityMatrix(providers, absolutePath)
	generateProviderRegistry(providers, smsDir)

	fmt.Printf("✅ SMS capabilities generation completed. Output: %s\n", absolutePath)
}

func scanProviderFiles() []ProviderMeta {
	var result []ProviderMeta

	// 构造器映射
	constructorMap := map[string]interface{}{
		"NewAliyunProvider": sms.NewAliyunProvider, "NewCl253Provider": sms.NewCl253Provider,
		"NewYuntongxunProvider": sms.NewYuntongxunProvider, "NewJuheProvider": sms.NewJuheProvider,
		"NewHuaweiProvider": sms.NewHuaweiProvider, "NewVolcProvider": sms.NewVolcProvider,
		"NewSmsbaoProvider": sms.NewSmsbaoProvider, "NewSubmailProvider": sms.NewSubmailProvider,
		"NewUcpProvider": sms.NewUcpProvider, "NewLuosimaoProvider": sms.NewLuosimaoProvider,
		"NewYunpianProvider": sms.NewYunpianProvider,
		"NewTencentProvider": sms.NewTencentProvider,
	}

	reName := regexp.MustCompile(`^.*@ProviderName:\s*(.+)$`)
	reWebsite := regexp.MustCompile(`^.*@Website:\s*(.+)$`)
	reAPIDoc := regexp.MustCompile(`^.*@APIDoc:\s*(.+)$`)

	for _, config := range SMSProviderConfigs {
		if !config.Enabled {
			continue
		}
		goFile := filepath.Join("providers/sms", config.Type+".go")
		absGoFile, _ := filepath.Abs(goFile)
		name, website := config.Name, ""
		apiDocs := []string{}
		if data, err := ioutil.ReadFile(absGoFile); err == nil {
			lines := strings.Split(string(data), "\n")
			header := lines[:min(80, len(lines))]
			headerStr := strings.Join(header, "\n")
			fmt.Printf("[debug] %s: headerStr preview:\n%s\n", config.Type, headerStr)
			// 逐行匹配@ProviderName
			for _, l := range header {
				if m := reName.FindStringSubmatch(l); len(m) > 1 {
					name = strings.TrimSpace(m[1])
					break
				}
			}
			// 逐行匹配@Website
			websiteFound := false
			for _, l := range header {
				if m := reWebsite.FindStringSubmatch(l); len(m) > 1 {
					website = strings.TrimSpace(m[1])
					fmt.Printf("[debug] %s: @Website matched: %s\n", config.Type, website)
					websiteFound = true
					break
				}
			}
			if !websiteFound {
				fmt.Printf("[debug] %s: @Website not found\n", config.Type)
			}
			// 逐行匹配@APIDoc主行
			apidocFound := false
			for _, l := range header {
				if m := reAPIDoc.FindStringSubmatch(l); len(m) > 1 {
					apiDocs = append(apiDocs, strings.TrimSpace(m[1]))
					fmt.Printf("[debug] %s: @APIDoc matched: %s\n", config.Type, m[1])
					apidocFound = true
					break
				}
			}
			if !apidocFound {
				fmt.Printf("[debug] %s: @APIDoc not found\n", config.Type)
			}
		} else {
			fmt.Printf("[debug] %s: failed to read file: %v\n", config.Type, err)
		}
		meta := &ProviderMeta{
			Name: name, Website: website, APIDoc: strings.Join(apiDocs, "\n"),
			Implementation: config.Type + ".go", Constructor: config.Constructor,
			Type: sms.ProviderType(config.Type), TypeConst: config.TypeConst,
		}
		if constructor, ok := constructorMap[config.Constructor]; ok {
			c := constructor
			t := sms.ProviderType(config.Type)
			n := config.Type
			meta.NewFunc = func() sms.SMSProviderInterface {
				cfg := sms.SMSProvider{Type: t, Name: n}
				args := []reflect.Value{reflect.ValueOf(cfg)}
				results := reflect.ValueOf(c).Call(args)
				return results[0].Interface().(sms.SMSProviderInterface)
			}
		}
		result = append(result, *meta)
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateCapabilityMatrix(providers []ProviderMeta, outputPath string) {
	if *verbose {
		fmt.Printf("[gen] Capability matrix will be written to: %s\n", outputPath)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "# SMS Provider Capability Matrix / 短信服务商能力矩阵")
	fmt.Fprintln(f, "")
	fmt.Fprintf(f, "> Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintln(f, "> This matrix is automatically generated by a script. For details, see the source code at [scripts/gen_sms_capabilities.go](../scripts/gen_sms_capabilities.go).")
	fmt.Fprintln(f, "> 本能力矩阵由脚本自动生成，具体生成逻辑请参考 [scripts/gen_sms_capabilities.go](../scripts/gen_sms_capabilities.go) 源码。")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "| Provider / 提供商 | Website / 官网 | API Doc | Implementation / 实现文件 | SMS(国内/国际) | MMS(国内/国际) | Voice(国内/国际) | Notes / 备注 |")
	fmt.Fprintln(f, "|------------------|----------------|--------|-------------------------|--------------|--------------|----------------|-------------|")

	for _, p := range providers {
		if *verbose {
			fmt.Printf("[gen] Processing provider: %s\n", p.Implementation)
		}

		website := p.Website
		if website != "" {
			website = fmt.Sprintf("[%s](%s)", stripURL(website), p.Website)
		}
		apiDoc := ""
		if p.APIDoc != "" {
			apiDocLines := strings.Split(p.APIDoc, "\n")
			if len(apiDocLines) == 1 {
				apiDoc = fmt.Sprintf("[API](%s)", apiDocLines[0])
			} else {
				for i, line := range apiDocLines {
					if i == 0 {
						apiDoc += fmt.Sprintf("- [API](%s)\n", line)
					} else {
						apiDoc += fmt.Sprintf("- [API](%s)\n", line)
					}
				}
				apiDoc = strings.TrimRight(apiDoc, "\n")
			}
		}

		smsStatus, mmsStatus, voiceStatus, notes := "🚧/🚧", "🚧/🚧", "🚧/🚧", ""
		cap := (*sms.Capabilities)(nil)
		errMsg := ""

		if p.NewFunc != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						errMsg = "不兼容"
						if *verbose {
							fmt.Printf("[gen] Provider %s panic: %v\n", p.Implementation, r)
						}
					}
				}()
				inst := p.NewFunc()
				if inst != nil {
					cap = inst.GetCapabilities()
					if *verbose {
						fmt.Printf("[gen] Provider %s GetCapabilities() returned: %+v\n", p.Implementation, cap)
					}
				}
			}()

			if cap != nil && errMsg == "" {
				smsStatus = getStatus(cap.SMS.Domestic) + "/" + getStatus(cap.SMS.International)
				mmsStatus = getStatus(cap.MMS.Domestic) + "/" + getStatus(cap.MMS.International)
				voiceStatus = getStatus(cap.Voice.Domestic) + "/" + getStatus(cap.Voice.International)
				notes = joinNotes(cap)
				if *verbose {
					fmt.Printf("[gen] Provider %s status: sms=%s, mms=%s, voice=%s\n", p.Implementation, smsStatus, mmsStatus, voiceStatus)
				}
			} else if errMsg != "" {
				smsStatus, mmsStatus, voiceStatus, notes = "❌/❌", "❌/❌", "❌/❌", errMsg
				if *verbose {
					fmt.Printf("[gen] Provider %s marked as incompatible.\n", p.Implementation)
				}
			} else if *verbose {
				fmt.Printf("[gen] Provider %s cap is nil, marking as developing.\n", p.Implementation)
			}
		} else if *verbose {
			fmt.Printf("[gen] Provider %s NewFunc is nil, marking as developing.\n", p.Implementation)
		}

		fmt.Fprintf(f, "| %s | %s | %s | [%s](./%s) | %s | %s | %s | %s |\n",
			ifEmpty(p.Name, p.Implementation), website, apiDoc, p.Implementation, p.Implementation, smsStatus, mmsStatus, voiceStatus, notes)
	}

	fmt.Fprintln(f, "\n## Legend / 图例")
	fmt.Fprintln(f, "- ✅ Supported / 支持")
	fmt.Fprintln(f, "- ❌ Not Supported / 不支持")
	fmt.Fprintln(f, "- 🚧 Under Development / 开发中")
	fmt.Fprintln(f, "- 不兼容: 文件存在但实现不兼容或未实现能力接口")
}

func generateProviderRegistry(providers []ProviderMeta, smsDir string) {
	outputPath := filepath.Join(smsDir, "provider_registry.go")
	absolutePath, _ := filepath.Abs(outputPath)

	if *verbose {
		fmt.Printf("[gen] Provider registry will be written to: %s\n", absolutePath)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	tmpl := `// Code generated by scripts/gen_sms_capabilities.go. DO NOT EDIT.

package sms

// Auto-generated provider registry
func init() {
	{{range .}}
	RegisterProviderConstructor({{.TypeConst}}, {{.Constructor}})
	{{end}}
}
`

	t, err := template.New("registry").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	err = t.Execute(f, providers)
	if err != nil {
		panic(err)
	}
}

func getStatus(r sms.RegionCapability) string {
	if r.Single {
		return "✅"
	}
	if r.Desc == "WIP" || r.Desc == "开发中" || strings.Contains(r.Desc, "开发中") {
		return "🚧"
	}
	return "❌"
}

func joinNotes(cap *sms.Capabilities) string {
	notes := []string{}
	if cap.SMS.Domestic.Desc != "" {
		notes = append(notes, "SMS(国内): "+cap.SMS.Domestic.Desc)
	}
	if cap.SMS.International.Desc != "" {
		notes = append(notes, "SMS(国际): "+cap.SMS.International.Desc)
	}
	if cap.MMS.Domestic.Desc != "" {
		notes = append(notes, "MMS(国内): "+cap.MMS.Domestic.Desc)
	}
	if cap.MMS.International.Desc != "" {
		notes = append(notes, "MMS(国际): "+cap.MMS.International.Desc)
	}
	if cap.Voice.Domestic.Desc != "" {
		notes = append(notes, "Voice(国内): "+cap.Voice.Domestic.Desc)
	}
	if cap.Voice.International.Desc != "" {
		notes = append(notes, "Voice(国际): "+cap.Voice.International.Desc)
	}
	if cap.SMS.Limits.MaxBatchSize > 0 {
		notes = append(notes, fmt.Sprintf("Batch: %d", cap.SMS.Limits.MaxBatchSize))
	}
	if cap.SMS.Limits.MaxContentLen > 0 {
		notes = append(notes, fmt.Sprintf("Length: %d", cap.SMS.Limits.MaxContentLen))
	}
	return join(notes, "; ")
}

func join(arr []string, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	out := arr[0]
	for _, s := range arr[1:] {
		out += sep + s
	}
	return out
}

func stripURL(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return url
}

func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
