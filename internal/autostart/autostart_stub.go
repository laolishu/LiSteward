//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package autostart

// 对于不支持的平台，返回 nil（工厂端会给出错误）
func newPlatformProvider() AutoStartProvider {
	return nil
}
