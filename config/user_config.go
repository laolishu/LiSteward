package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UserConfig 代表用户偏好设置
type UserConfig struct {
	Language  string `json:"language"`   // 用户选择的语言: zh-CN, en-US
	AutoStart bool   `json:"auto_start"` // 是否开机自启动
	// NvmDir: 可选，用户指定的 NVM 根目录（NVM_DIR），用于持久化与查找 Node 版本
	NvmDir string `json:"nvm_dir"`
}

var cachedUserConfig *UserConfig

// getUserConfigPath 返回用户配置文件的路径（相对于可执行程序）
func getUserConfigPath() string {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "data", "config", "user.json")
}

// LoadUserConfig 从 config/user.json 加载用户配置
func LoadUserConfig() (*UserConfig, error) {
	if cachedUserConfig != nil {
		return cachedUserConfig, nil
	}

	path := getUserConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，返回默认配置
			return &UserConfig{Language: "zh-CN", AutoStart: false}, nil
		}
		return nil, err
	}

	var cfg UserConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 设置默认语言
	if cfg.Language == "" {
		cfg.Language = "zh-CN"
	}

	cachedUserConfig = &cfg
	return cachedUserConfig, nil
}

// SaveUserConfig 保存用户配置到 config/user.json
func SaveUserConfig(uc *UserConfig) error {
	if uc == nil {
		return nil
	}

	path := getUserConfigPath()

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 序列化并写入文件
	newData, err := json.MarshalIndent(uc, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, newData, 0644); err != nil {
		return err
	}

	// 更新缓存
	cachedUserConfig = uc

	return nil
}

// EnsureUserConfig 确保用户配置文件存在，如不存在则创建默认配置
func EnsureUserConfig() error {
	path := getUserConfigPath()

	// 检查文件是否存在
	if _, err := os.Stat(path); err == nil {
		// 文件存在，无需做任何操作
		return nil
	} else if !os.IsNotExist(err) {
		// 其他错误
		return err
	}

	// 文件不存在，创建默认配置
	defaultConfig := &UserConfig{
		Language:  "zh-CN",
		AutoStart: false,
		NvmDir:    "",
	}

	return SaveUserConfig(defaultConfig)
}

// ResetUserConfigCache 重置用户配置缓存（用于测试）
func ResetUserConfigCache() {
	cachedUserConfig = nil
}
