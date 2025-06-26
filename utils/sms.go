package utils

import (
	"strings"
)

// GetSignatureFromContent 从短信内容中提取签名
// 返回签名内容，如果没有找到签名则返回空字符串
// 签名格式：以【开头，以】结尾，且】在合理位置（前20个字符内）
func GetSignatureFromContent(content string) string {
	if content == "" {
		return ""
	}

	// 检查是否以【开头
	if !strings.HasPrefix(content, "【") {
		return ""
	}

	// 查找】的位置
	endIndex := strings.Index(content, "】")
	if endIndex == -1 || endIndex > 20 || endIndex == 1 {
		return ""
	}

	// 提取签名内容（去掉【和】）
	return content[1:endIndex]
}

// HasSignature 检查内容是否已经包含签名
func HasSignature(content string) bool {
	return GetSignatureFromContent(content) != ""
}

// AddSignature 为短信内容添加签名
// 如果内容已经有签名，则直接返回原内容
// 否则在开头添加【signName】
func AddSignature(content, signName string) string {
	if HasSignature(content) || signName == "" {
		return content
	}

	return "【" + signName + "】" + content
}
