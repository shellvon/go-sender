//revive:disable:var-naming
package utils

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

// HashSum 返回数据的哈希值原始字节数组。
func HashSum(hasher func() hash.Hash, data []byte) []byte {
	h := hasher()
	h.Write(data)
	return h.Sum(nil)
}

// HashHex 返回数据的哈希值十六进制编码。
func HashHex(hasher func() hash.Hash, data []byte) string {
	return hex.EncodeToString(HashSum(hasher, data))
}

// HashBase64 返回数据的哈希值 base64 编码。
func HashBase64(hasher func() hash.Hash, data []byte) string {
	return Base64EncodeBytes(HashSum(hasher, data))
}

// HMACSum 返回 HMAC 原始字节数组。
func HMACSum(hasher func() hash.Hash, key, data []byte) []byte {
	mac := hmac.New(hasher, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// HMACHex 返回 HMAC 十六进制编码。
func HMACHex(hasher func() hash.Hash, key, data []byte) string {
	return hex.EncodeToString(HMACSum(hasher, key, data))
}

// HMACBase64 返回 HMAC base64 编码。
func HMACBase64(hasher func() hash.Hash, key, data []byte) string {
	return Base64EncodeBytes(HMACSum(hasher, key, data))
}

// Base64EncodeBytes 返回字节数组的 base64 编码。
func Base64EncodeBytes(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
