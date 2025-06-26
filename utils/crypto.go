package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// MD5Hex 返回字符串的 MD5 十六进制编码
func MD5Hex(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA1Hex 返回 HMAC-SHA1 十六进制编码
func HMACSHA1Hex(key, data string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACSHA256Hex 返回 HMAC-SHA256 十六进制编码
func HMACSHA256Hex(key, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// Base64Encode 返回 base64 编码字符串
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64EncodeBytes 返回 base64 编码字符串
func Base64EncodeBytes(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// SHA256Hex returns the SHA256 hex encoding of the input data
func SHA256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256Sum returns the SHA256 sum of the input data
func SHA256Sum(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// SHA1Hex returns the SHA1 hex encoding of the input data
func SHA1Hex(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA256 returns the HMAC-SHA256 of data with key (raw bytes)
func HMACSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// HMACSHA1 returns the HMAC-SHA1 of data with key (raw bytes)
func HMACSHA1(key, data []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// HMACSHA1Base64 returns the HMAC-SHA1 of data with key, base64 encoded
func HMACSHA1Base64(key, data string) string {
	mac := HMACSHA1([]byte(key), []byte(data))
	return Base64EncodeBytes(mac)
}

// HMACSHA256Base64 returns the HMAC-SHA256 of data with key, base64 encoded
func HMACSHA256Base64(key, data string) string {
	mac := HMACSHA256([]byte(key), []byte(data))
	return Base64EncodeBytes(mac)
}

// MD5Base64 returns the MD5 of the input string, base64 encoded
func MD5Base64(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
