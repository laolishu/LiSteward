package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureUserConfig(t *testing.T) {
	// 重置缓存
	ResetUserConfigCache()

	// 测试创建新的用户配置
	err := EnsureUserConfig()
	if err != nil {
		t.Fatalf("EnsureUserConfig 失败: %v", err)
	}

	// 验证配置文件存在
	path := getUserConfigPath()
	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("用户配置文件不存在: %v", err)
	}
}

func TestLoadUserConfig(t *testing.T) {
	// 重置缓存
	ResetUserConfigCache()

	// 创建用户配置
	err := EnsureUserConfig()
	if err != nil {
		t.Fatalf("EnsureUserConfig 失败: %v", err)
	}

	// 加载用户配置
	cfg, err := LoadUserConfig()
	if err != nil {
		t.Fatalf("LoadUserConfig 失败: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadUserConfig 返回 nil")
	}

	if cfg.Language == "" {
		t.Fatal("用户配置语言为空")
	}
}

func TestSaveUserConfig(t *testing.T) {
	// 重置缓存
	ResetUserConfigCache()

	// 创建测试配置
	testCfg := &UserConfig{
		Language: "en-US",
	}

	// 保存用户配置
	err := SaveUserConfig(testCfg)
	if err != nil {
		t.Fatalf("SaveUserConfig 失败: %v", err)
	}

	// 重置缓存，重新加载
	ResetUserConfigCache()

	// 加载并验证
	cfg, err := LoadUserConfig()
	if err != nil {
		t.Fatalf("LoadUserConfig 失败: %v", err)
	}

	if cfg.Language != "en-US" {
		t.Fatalf("保存的语言不正确: 期望 en-US, 获得 %s", cfg.Language)
	}

	// 清理测试数据
	path := getUserConfigPath()
	os.Remove(path)
}

func TestAppConfigGetLanguage(t *testing.T) {
	// 重置缓存
	ResetUserConfigCache()

	// 创建用户配置设置语言
	testCfg := &UserConfig{
		Language: "en-US",
	}
	SaveUserConfig(testCfg)

	// 创建应用配置
	appCfg := &AppConfig{
		Language: "zh-CN",
	}

	// 应该返回用户配置中的语言（优先级更高）
	lang := appCfg.GetLanguage()
	if lang != "en-US" {
		t.Fatalf("GetLanguage 应该返回用户配置中的语言: 期望 en-US, 获得 %s", lang)
	}

	// 清理测试数据
	ResetUserConfigCache()
	path := getUserConfigPath()
	os.Remove(path)
}

func TestCreateConfigDirectory(t *testing.T) {
	// 测试配置目录结构（相对于可执行程序）
	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("获取可执行程序路径失败: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	configDir := filepath.Join(exeDir, "data", "config")

	// 这些目录应该存在（通过 EnsureUserConfig 或应用启动创建）
	err = EnsureUserConfig()
	if err != nil {
		t.Fatalf("EnsureUserConfig 失败: %v", err)
	}

	_, err = os.Stat(configDir)
	if err != nil {
		t.Fatalf("config 目录不存在: %v", err)
	}
}
