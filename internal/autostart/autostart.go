package autostart

import (
	"fmt"
	"runtime"
)

// AutoStartProvider 抽象接口
type AutoStartProvider interface {
	IsEnabled() (bool, error)
	Enable() error
	Disable() error
}

// NewProvider 根据当前平台返回实现
func NewProvider() (AutoStartProvider, error) {
	// newPlatformProvider 在每个平台的实现文件中定义（windows/darwin/linux）
	if prov := newPlatformProvider(); prov != nil {
		return prov, nil
	}
	return nil, fmt.Errorf("autostart not implemented for %s", runtime.GOOS)
}
