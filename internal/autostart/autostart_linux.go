//go:build linux
// +build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

type linuxProvider struct {
	name string
}

func NewLinuxProvider() AutoStartProvider {
	return &linuxProvider{name: "listeward"}
}

func newPlatformProvider() AutoStartProvider {
	return NewLinuxProvider()
}

func (l *linuxProvider) desktopPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "autostart")
	return filepath.Join(dir, l.name+".desktop"), nil
}

func (l *linuxProvider) IsEnabled() (bool, error) {
	path, err := l.desktopPath()
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

func (l *linuxProvider) Enable() error {
	path, err := l.desktopPath()
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
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Version=1.0
Name=LiSteward
Exec=%s
X-GNOME-Autostart-enabled=true
NoDisplay=false
`, exe)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func (l *linuxProvider) Disable() error {
	path, err := l.desktopPath()
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
