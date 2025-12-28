package language

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

const (
	// LOCALE_SNAME 获取语言区域的名称（如 zh-CN, en-US）
	LOCALE_SNAME = 0x59
)

// GetSystemLanguage 获取 Windows 系统语言
// 返回格式为 zh-CN, en-US 等
// 如果获取失败，返回默认值 zh-CN
func GetSystemLanguage() (string, error) {
	// 使用 Windows API 获取用户默认语言
	// LOCALE_USER_DEFAULT (0x0400) 表示用户默认区域设置
	langID := uint32(0x0400) // LOCALE_USER_DEFAULT

	// 调用 GetLocaleInfoEx 获取语言代码
	// 使用 NULL 表示获取用户默认区域设置
	lang, err := getLocaleInfoEx(langID)
	if err != nil {
		// 如果获取失败，返回默认的中文
		return "zh-CN", fmt.Errorf("获取系统语言失败: %v", err)
	}

	// 规范化语言代码
	lang = strings.ToLower(lang)
	lang = strings.ReplaceAll(lang, "_", "-")

	// 验证语言代码格式
	if lang == "zh-cn" {
		return "zh-CN", nil
	} else if lang == "en-us" {
		return "en-US", nil
	} else if strings.HasPrefix(lang, "zh") {
		// 其他中文变体，统一用简体中文
		return "zh-CN", nil
	} else if strings.HasPrefix(lang, "en") {
		// 其他英文变体，统一用美国英文
		return "en-US", nil
	}

	// 未知语言，默认中文
	return "zh-CN", nil
}

// getLocaleInfoEx 获取区域设置信息
func getLocaleInfoEx(localeID uint32) (string, error) {
	// 使用 GetLocaleInfoA 获取语言代码
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetLocaleInfo := kernel32.NewProc("GetLocaleInfoA")

	// 获取缓冲区大小
	bufSize := 32
	buf := make([]byte, bufSize)

	// 调用 API
	ret, _, err := procGetLocaleInfo.Call(
		uintptr(localeID),
		uintptr(LOCALE_SNAME),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufSize),
	)

	if ret == 0 {
		return "", fmt.Errorf("GetLocaleInfoA failed: %v", err)
	}

	// 转换为字符串
	lang := strings.TrimRight(string(buf[:ret-1]), "\x00")
	return lang, nil
}
