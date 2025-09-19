package wecombot

import (
	//nolint:gosec // compatibility with legacy system, not for security
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shellvon/go-sender/utils"
)

// ImageBuilder 提供了一个流式 API 来构建企业微信图片消息。
//
// 推荐使用便捷的静态构造函数（返回错误信息）：
//   - [ImageFromFile] - 从文件路径创建，自动处理所有步骤
//   - [ImageFromBytes] - 从字节数据创建，自动编码和计算MD5
//   - [ImageFromBase64] - 从Base64字符串创建，自动计算MD5（支持data URL前缀）
//
// 或者使用传统的链式调用（需要手动处理验证）：
//
//	msg := wecombot.Image().Base64(imgB64).MD5(imgMD5).Build()
type ImageBuilder struct {
	base64 string
	md5    string
}

// Image 创建一个新的 ImageBuilder 实例。
// 基于企业微信机器人 API 的 SendImageParams
// 参考：https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
// 返回值：*ImageBuilder - 新创建的 ImageBuilder 实例，用于构建 ImageMessage。
func Image() *ImageBuilder {
	return &ImageBuilder{}
}

// Base64 设置图片内容的 Base64 编码。
// 参数：data string - 图片的 Base64 编码字符串
// 返回值：*ImageBuilder - 返回 ImageBuilder 实例以支持链式调用.
func (b *ImageBuilder) Base64(data string) *ImageBuilder {
	b.base64 = data
	return b
}

// MD5 设置原始图片的 MD5 哈希值。
// 参数：hash string - 图片的 MD5 哈希值
// 返回值：*ImageBuilder - 返回 ImageBuilder 实例以支持链式调用.
func (b *ImageBuilder) MD5(hash string) *ImageBuilder {
	b.md5 = hash
	return b
}

// ImageFromFile 从文件路径创建图片消息，自动处理编码和MD5计算。
// 支持常见的图片格式：JPG, JPEG, PNG, GIF, BMP, WEBP等。
//
// 参数：filePath string - 图片文件的路径
// 返回值：
//   - *ImageMessage - 创建的图片消息实例
//   - error - 如果出错返回具体的错误信息
func ImageFromFile(filePath string) (*ImageMessage, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("图片文件不存在: %s", filePath)
	}

	// 验证文件扩展名
	if !isValidImageExtension(filePath) {
		return nil, fmt.Errorf("不支持的图片格式: %s，支持的格式: jpg, jpeg, png, gif, bmp, webp", filepath.Ext(filePath))
	}

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取图片文件失败: %w", err)
	}

	return ImageFromBytes(fileData)
}

// ImageFromBytes 从字节数据创建图片消息，自动处理Base64编码和MD5计算。
//
// 参数：data []byte - 图片的原始字节数据
// 返回值：
//   - *ImageMessage - 创建的图片消息实例
//   - error - 如果出错返回具体的错误信息
func ImageFromBytes(data []byte) (*ImageMessage, error) {
	if len(data) == 0 {
		return nil, errors.New("图片数据不能为空")
	}

	if len(data) > maxImageSizeBytes {
		return nil, fmt.Errorf("图片大小超过 %dMB 限制，当前大小: %d 字节", maxImageSizeBytes/bytesPerMB, len(data))
	}

	// 计算MD5和Base64
	md5Hash := utils.HashHex(md5.New, data)
	base64Data := base64.StdEncoding.EncodeToString(data)

	return Image().Base64(base64Data).MD5(md5Hash).Build(), nil
}

// ImageFromBase64 从Base64字符串创建图片消息，自动计算MD5。
// 支持带前缀的Base64字符串（如 "data:image/jpeg;base64,..."），会自动提取纯Base64内容。
//
// 参数：base64Data string - 图片的 Base64 编码字符串，支持带或不带前缀
// 返回值：
//   - *ImageMessage - 创建的图片消息实例
//   - error - 如果出错返回具体的错误信息
func ImageFromBase64(base64Data string) (*ImageMessage, error) {
	// 清理Base64字符串
	cleanBase64 := cleanBase64String(base64Data)
	if cleanBase64 == "" {
		return nil, errors.New("Base64 数据不能为空")
	}

	// 验证并解码Base64以计算MD5
	rawData, err := base64.StdEncoding.DecodeString(cleanBase64)
	if err != nil {
		return nil, fmt.Errorf("无效的 Base64 编码: %w", err)
	}

	// 检查解码后的数据大小
	if len(rawData) > maxImageSizeBytes {
		return nil, fmt.Errorf("图片大小超过 %dMB 限制，当前大小: %d 字节", maxImageSizeBytes/bytesPerMB, len(rawData))
	}

	// 计算MD5
	md5Hash := utils.HashHex(md5.New, rawData)

	return Image().Base64(cleanBase64).MD5(md5Hash).Build(), nil
}

// Build 构建并返回一个 ImageMessage 实例。
// 返回值：*ImageMessage - 基于 ImageBuilder 配置创建的图片消息实例.
func (b *ImageBuilder) Build() *ImageMessage {
	return &ImageMessage{
		BaseMessage: newBaseMessage(TypeImage),
		Image: ImageContent{
			Base64: b.base64,
			MD5:    b.md5,
		},
	}
}

// cleanBase64String 清理Base64字符串，移除可能的data URL前缀
// 支持格式：
//   - "data:image/jpeg;base64,/9j/4AAQ..." -> "/9j/4AAQ..."
//   - "/9j/4AAQ..." -> "/9j/4AAQ..." (不变)
func cleanBase64String(data string) string {
	// 直接查找并移除 "base64," 前缀及其之前的所有内容
	if idx := strings.Index(data, "base64,"); idx != -1 {
		data = data[idx+7:] // "base64," 长度为7
	}
	return strings.TrimSpace(data)
}

// isValidImageExtension 检查文件扩展名是否为支持的图片格式.
func isValidImageExtension(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
