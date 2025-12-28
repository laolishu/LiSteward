//go:build windows
// +build windows

package autostart

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type windowsProvider struct {
	name string
}

// NewWindowsProvider 返回一个简单的 Windows 启动项提供者。
// 它在当前用户的 Startup 文件夹中写入一个 .cmd 脚本，脚本通过 start 启动可执行文件。
func NewWindowsProvider() AutoStartProvider {
	return &windowsProvider{name: "LiSteward"}
}

func newPlatformProvider() AutoStartProvider {
	return NewWindowsProvider()
}

func (w *windowsProvider) startupDir() (string, error) {
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		return "", errors.New("APPDATA not set")
	}
	return filepath.Join(appdata, "Microsoft", "Windows", "Start Menu", "Programs", "Startup"), nil
}

func (w *windowsProvider) filePath() (string, error) {
	dir, err := w.startupDir()
	if err != nil {
		return "", err
	}
	fname := fmt.Sprintf("%s-autostart.cmd", strings.ToLower(w.name))
	return filepath.Join(dir, fname), nil
}

func (w *windowsProvider) IsEnabled() (bool, error) {
	path, err := w.filePath()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	// 简单检查内容中是否包含可执行路径（如果需要更严格的校验，可读取内容并匹配）
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	// 文件存在且非空视为启用
	return len(data) > 0, nil
}

func (w *windowsProvider) Enable() error {
	path, err := w.filePath()
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// 脚本内容：使用 start 启动可执行文件，这样不会阻塞 shell
	content := fmt.Sprintf("@echo off\nstart \"\" \"%s\"\r\n", exePath)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func (w *windowsProvider) Disable() error {
	path, err := w.filePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(path)
}
