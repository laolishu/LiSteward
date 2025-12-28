package config

import (
	"LiSteward/internal/language"
	"LiSteward/internal/version"
	"encoding/json"
	"os"
)

type WindowConfig struct {
	Width     int `json:"width"`
	Height    int `json:"height"`
	MinWidth  int `json:"minWidth"`
	MinHeight int `json:"minHeight"`
}

type AppConfig struct {
	ProductName   string            `json:"productName"`
	ProductNameEn string            `json:"productNameEn"`
	Version       string            `json:"version"`
	Description   string            `json:"description"`
	Author        map[string]string `json:"author"`
	Website       string            `json:"website"`
	SponsorQrCode string            `json:"sponsorQrCode"`
	DefaultWindow WindowConfig      `json:"defaultWindow"`
	Language      string            `json:"language"` // 用户选择的语言: zh-CN, en-US
}

var cached *AppConfig

// Load reads config/app.json and unmarshals it
// SponsorQrCode 从 version.go 中读取，而不是从 app.json 中读取
func Load() (*AppConfig, error) {
	if cached != nil {
		return cached, nil
	}

	data, err := os.ReadFile("config/app.json")
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 从 version.go 中读取赞赏二维码，覆盖 app.json 中的值
	cfg.SponsorQrCode = version.SponsorQrCode

	cached = &cfg
	return cached, nil
}

// SaveLanguage 保存用户语言偏好到用户配置文件
func (c *AppConfig) SaveLanguage(lang string) error {
	// 保存到用户配置
	userCfg := &UserConfig{Language: lang}
	if err := SaveUserConfig(userCfg); err != nil {
		return err
	}

	// 同时更新 AppConfig（为兼容旧代码）
	c.Language = lang

	// 也保存到 app.json（向后兼容）
	// 重新读取整个配置文件
	data, err := os.ReadFile("config/app.json")
	if err != nil {
		return err
	}

	// 解析为 map 以保留其他字段
	var cfgMap map[string]interface{}
	if err := json.Unmarshal(data, &cfgMap); err != nil {
		return err
	}

	// 更新语言字段
	cfgMap["language"] = lang

	// 序列化并写入文件
	newData, err := json.MarshalIndent(cfgMap, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("config/app.json", newData, 0644)
}

// GetLanguage 获取用户语言偏好
// 优先级：用户配置 > 系统语言 > 默认(zh-CN)
func (c *AppConfig) GetLanguage() string {
	// 首先尝试加载用户配置
	if userCfg, err := LoadUserConfig(); err == nil && userCfg.Language != "" {
		return userCfg.Language
	}

	// 如果用户设置过（保留兼容旧配置），使用用户设置
	if c.Language != "" {
		return c.Language
	}

	// 检测系统语言
	if syslang, err := language.GetSystemLanguage(); err == nil && syslang != "" {
		return syslang
	}

	// 默认中文
	return "zh-CN"
}
