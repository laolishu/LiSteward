//go:build darwin
// +build darwin

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

type darwinProvider struct {
	label string
}

func NewDarwinProvider() AutoStartProvider {
	return &darwinProvider{label: "com.laolishu.listeward"}
}

func newPlatformProvider() AutoStartProvider {
	return NewDarwinProvider()
}

func (d *darwinProvider) plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Library", "LaunchAgents")
	return filepath.Join(dir, d.label+".plist"), nil
}

func (d *darwinProvider) IsEnabled() (bool, error) {
	path, err := d.plistPath()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *darwinProvider) Enable() error {
	path, err := d.plistPath()
	if err != nil {
		return err
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 简单 plist，指定 ProgramArguments 以确保可执行被启动
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>%s</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
</dict>
</plist>
`, d.label, exe)

	return os.WriteFile(path, []byte(content), 0644)
}

func (d *darwinProvider) Disable() error {
	path, err := d.plistPath()
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
