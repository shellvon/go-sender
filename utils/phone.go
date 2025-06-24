package utils

import "strings"

// IsInternationalMobile checks if the mobile number is international (not China)
func IsInternationalMobile(mobile string) bool {
	cleanMobile := mobile
	if strings.HasPrefix(cleanMobile, "+") {
		cleanMobile = cleanMobile[1:]
	} else if strings.HasPrefix(cleanMobile, "00") {
		cleanMobile = cleanMobile[2:]
	}
	// +86/86/1xxxxxxxxxx 都视为国内
	if strings.HasPrefix(cleanMobile, "86") && len(cleanMobile) == 13 && cleanMobile[2] == '1' {
		return false
	}
	if len(cleanMobile) == 11 && cleanMobile[0] == '1' {
		return false
	}
	return strings.HasPrefix(mobile, "+") || strings.HasPrefix(mobile, "00")
}
