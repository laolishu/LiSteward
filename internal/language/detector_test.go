package language

import (
	"testing"
)

// TestGetSystemLanguage 测试系统语言检测
func TestGetSystemLanguage(t *testing.T) {
	lang, err := GetSystemLanguage()

	// 检测结果应该不为空
	if lang == "" {
		t.Error("GetSystemLanguage() returned empty string")
	}

	// 检测结果应该是有效的语言代码格式
	validLangs := map[string]bool{
		"zh-CN": true,
		"en-US": true,
	}

	if !validLangs[lang] {
		t.Logf("Warning: GetSystemLanguage() returned: %s, error: %v", lang, err)
		// 不失败，因为可能是其他系统语言
	}

	t.Logf("System language detected: %s", lang)
}

// TestGetSystemLanguageDefaultFallback 测试语言检测失败时的默认值
func TestGetSystemLanguageDefaultFallback(t *testing.T) {
	lang, _ := GetSystemLanguage()

	// 应该总是返回一个有效的语言代码，最坏的情况是默认的 zh-CN
	if lang != "zh-CN" && lang != "en-US" {
		t.Errorf("GetSystemLanguage() should return zh-CN or en-US as fallback, got: %s", lang)
	}
}
