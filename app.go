package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	osruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"LiSteward/config"
	"LiSteward/internal/autostart"
	"LiSteward/internal/hosts"
	"LiSteward/internal/language"
	"LiSteward/internal/nvm"
	"LiSteward/internal/version"
)

// App struct
type App struct {
	ctx          context.Context
	hostsService *hosts.Service
	cfg          *config.AppConfig
	appDir       string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化应用数据目录（以可执行程序为根目录）
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取可执行程序路径失败: %v\n", err)
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	appDataDir := filepath.Join(exeDir, "data")

	// 保存程序目录，供前端“打开目录”功能使用
	a.appDir = exeDir

	// 确保配置目录存在
	configDir := filepath.Join(appDataDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("创建配置目录失败: %v\n", err)
	}

	// 确保用户配置文件存在
	if err := config.EnsureUserConfig(); err != nil {
		fmt.Printf("初始化用户配置失败: %v\n", err)
	}

	// 读取应用配置并缓存到 App 中，供前端查询
	if cfg, err := config.Load(); err == nil {
		a.cfg = cfg

		// 如果配置中没有语言设置，自动检测系统语言
		if a.cfg.Language == "" {
			if syslang, err := language.GetSystemLanguage(); err == nil {
				a.cfg.Language = syslang
			} else {
				a.cfg.Language = "zh-CN"
			}
		}

		// 处理 sponsorQrCode：如果是相对路径，则尝试读取文件并以 Data URL 返回，确保前端可以直接显示图片
		if a.cfg.SponsorQrCode != "" {
			path := a.cfg.SponsorQrCode
			// 相对路径以 ./ 或 data/ 开头时，相对于可执行程序目录
			if strings.HasPrefix(path, ".") || strings.HasPrefix(path, "data") {
				exePath, err := os.Executable()
				if err == nil {
					exeDir := filepath.Dir(exePath)
					absPath := filepath.Join(exeDir, path)
					if data, err := os.ReadFile(absPath); err == nil {
						encoded := base64.StdEncoding.EncodeToString(data)
						// 假设 png（如果需要更严格的检测，可根据文件头判断 MIME）
						a.cfg.SponsorQrCode = "data:image/png;base64," + encoded
					}
				}
			}
		}
	}

	// 初始化 Hosts 服务
	hostsService, err := hosts.NewService(appDataDir)
	if err != nil {
		fmt.Printf("初始化 Hosts 服务失败: %v\n", err)
	}
	a.hostsService = hostsService

}

// GetAppConfig 返回应用的配置信息（供前端调用）
func (a *App) GetAppConfig() (*config.AppConfig, error) {
	if a.cfg != nil {
		return a.cfg, nil
	}
	return config.Load()
}

// GetVersion 返回应用的版本号
func (a *App) GetVersion() (string, error) {
	if a.cfg != nil && a.cfg.Version != "" {
		return a.cfg.Version, nil
	}

	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	return cfg.Version, nil
}

// OpenAppDir 打开当前程序所在目录（跨平台）
func (a *App) OpenAppDir() error {
	dir := a.appDir
	if dir == "" {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("无法获取可执行文件路径: %w", err)
		}
		dir = filepath.Dir(exePath)
	}

	var cmd *exec.Cmd
	switch osruntime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("打开目录失败: %w", err)
	}

	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ==================== Hosts 管理方法 ====================

// GetHostsEntries 获取所有 Hosts 条目
func (a *App) GetHostsEntries() ([]hosts.HostEntry, error) {
	return a.hostsService.ReadHostsFile()
}

// SaveHostsEntries 保存 Hosts 条目
func (a *App) SaveHostsEntries(entries []hosts.HostEntry) error {
	return a.hostsService.WriteHostsFile(entries)
}

// SaveHostsEntriesWithReadOnlyRestore 保存 Hosts 条目，自动处理只读属性
func (a *App) SaveHostsEntriesWithReadOnlyRestore(entries []hosts.HostEntry) error {
	return a.hostsService.WriteHostsFileWithReadOnlyRestore(entries)
}

// ValidateHostEntry 验证 Hosts 条目
func (a *App) ValidateHostEntry(entry hosts.HostEntry) error {
	return a.hostsService.ValidateEntry(entry)
}

// CreateBackup 创建备份
func (a *App) CreateBackup() (*hosts.BackupInfo, error) {
	return a.hostsService.BackupHosts()
}

// ListBackups 列出所有备份
func (a *App) ListBackups() ([]hosts.BackupInfo, error) {
	return a.hostsService.ListBackups()
}

// RestoreBackup 从备份恢复
func (a *App) RestoreBackup(backupID string) error {
	return a.hostsService.RestoreFromBackup(backupID)
}

// DeleteBackup 删除备份
func (a *App) DeleteBackup(backupID string) error {
	return a.hostsService.DeleteBackup(backupID)
}

// SwapIPs 交换条目的主 IP 和备用 IP
func (a *App) SwapIPs(entries []hosts.HostEntry, index int) ([]hosts.HostEntry, error) {
	return a.hostsService.SwapIPs(entries, index)
}

// SaveProfile 保存配置方案
func (a *App) SaveProfile(name string, entries []hosts.HostEntry) error {
	return a.hostsService.SaveProfile(name, entries)
}

// LoadProfile 加载配置方案
func (a *App) LoadProfile(name string) (*hosts.HostsProfile, error) {
	return a.hostsService.LoadProfile(name)
}

// ListProfiles 列出所有配置方案
func (a *App) ListProfiles() ([]hosts.HostsProfile, error) {
	return a.hostsService.ListProfiles()
}

// DeleteProfile 删除配置方案
func (a *App) DeleteProfile(name string) error {
	return a.hostsService.DeleteProfile(name)
}

// RenameProfile 重命名配置方案
func (a *App) RenameProfile(oldName, newName string) error {
	return a.hostsService.RenameProfile(oldName, newName)
}

// ApplyProfile 应用配置方案
func (a *App) ApplyProfile(name string) error {
	return a.hostsService.ApplyProfile(name)
}

// ==================== 国际化方法 ====================

// GetSystemLanguage 获取系统语言
func (a *App) GetSystemLanguage() (string, error) {
	return language.GetSystemLanguage()
}

// ==================== NVM 管理器绑定 ====================

// DetectNvm 检测系统是否已安装 nvm（或存在 NVM_DIR）
func (a *App) DetectNvm() (bool, string, error) {
	ok, dir := nvm.Detect()
	return ok, dir, nil
}

// GetNvmVersion 获取 NVM 的版本号
func (a *App) GetNvmVersion() (string, error) {
	return nvm.GetNvmVersion()
}

// InstallNvm 返回安装提示或启动安装（仅提供建议）
func (a *App) InstallNvm() (string, error) {
	return nvm.InstallNvm()
}

// ListNodeVersions 列出可用/已安装的 Node 版本
func (a *App) ListNodeVersions() ([]nvm.NodeVersion, error) {
	// 优先使用用户配置中的 NvmDir
	uc, err := config.LoadUserConfig()
	if err != nil {
		// 若无法加载用户配置，仍尝试无参数列出
		return nvm.ListVersions("")
	}
	return nvm.ListVersions(uc.NvmDir)
}

// InstallNodeVersion 使用 nvm 安装指定版本
func (a *App) InstallNodeVersion(version string) (string, error) {
	return nvm.InstallVersion(version)
}

// UseNodeVersion 切换到指定 Node 版本
func (a *App) UseNodeVersion(version string) (string, error) {
	return nvm.UseVersion(version)
}

// UninstallNodeVersion 卸载指定 Node 版本
func (a *App) UninstallNodeVersion(version string) (string, error) {
	return nvm.UninstallVersion(version)
}

// GetNvmDir 返回当前配置的 NVM_DIR（若未设置返回空）
func (a *App) GetNvmDir() (string, error) {
	uc, err := config.LoadUserConfig()
	if err != nil {
		return "", err
	}
	return uc.NvmDir, nil
}

// SetNvmDir 保存用户指定的 NVM_DIR 到用户配置
func (a *App) SetNvmDir(dir string) error {
	uc, err := config.LoadUserConfig()
	if err != nil {
		uc = &config.UserConfig{Language: "zh-CN"}
	}
	uc.NvmDir = dir
	return config.SaveUserConfig(uc)
}

// GetNvmRoot 获取 Node 仓库路径
func (a *App) GetNvmRoot() (string, error) {
	return nvm.GetNvmRoot()
}

// OpenNodeAvailableList 打开系统命令行窗口显示可用的 Node 版本列表
func (a *App) OpenNodeAvailableList() error {
	return nvm.OpenAvailableList()
}

// GetUserLanguagePreference 获取用户语言偏好
func (a *App) GetUserLanguagePreference() string {
	if a.cfg != nil {
		return a.cfg.GetLanguage()
	}
	return "zh-CN"
}

// SetUserLanguagePreference 设置用户语言偏好
func (a *App) SetUserLanguagePreference(lang string) error {
	if a.cfg == nil {
		return fmt.Errorf("config not initialized")
	}
	return a.cfg.SaveLanguage(lang)
}

// GetAutoStart 返回当前用户配置中的开机自启动设置
func (a *App) GetAutoStart() (bool, error) {
	// 优先使用平台 provider 检测实际系统启动项状态
	if prov, err := autostart.NewProvider(); err == nil {
		if enabled, err := prov.IsEnabled(); err == nil {
			return enabled, nil
		}
		// 若 provider 检测失败，则回退到用户配置
	}

	uc, err := config.LoadUserConfig()
	if err != nil {
		return false, err
	}
	return uc.AutoStart, nil
}

// SetAutoStart 设置并保存用户配置中的开机自启动选项
func (a *App) SetAutoStart(enabled bool) error {
	uc, err := config.LoadUserConfig()
	if err != nil {
		// 若加载失败，则创建默认配置并继续
		uc = &config.UserConfig{Language: "zh-CN"}
	}
	uc.AutoStart = enabled

	// 先保存用户配置
	if err := config.SaveUserConfig(uc); err != nil {
		return err
	}

	// 再尝试更新系统级启动项（如果平台支持）
	prov, err := autostart.NewProvider()
	if err != nil {
		// 平台不支持系统启动项，已保存用户意图
		return nil
	}

	if enabled {
		return prov.Enable()
	}
	return prov.Disable()
}

// ShowWindow 显示窗口
func (a *App) ShowWindow() error {
	runtime.WindowShow(a.ctx)
	return nil
}

// HideWindow 隐藏窗口
func (a *App) HideWindow() error {
	runtime.WindowHide(a.ctx)
	return nil
}

// ToggleWindowVisibility 切换窗口可见性
func (a *App) ToggleWindowVisibility() error {
	// 检查窗口是否最小化
	isMinimised := runtime.WindowIsMinimised(a.ctx)

	if isMinimised {
		runtime.WindowUnminimise(a.ctx)
		runtime.WindowShow(a.ctx)
	} else {
		runtime.WindowHide(a.ctx)
	}
	return nil
}

// GetAppInfo 获取应用版本和信息
func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"version":     version.Version,
		"author":      version.Author,
		"email":       version.AuthorEmail,
		"repository":  version.Repository,
		"buildTime":   version.BuildTime,
		"buildCommit": version.BuildCommit,
		"description": version.Description,
		"goVersion":   osruntime.Version(),
		"os":          osruntime.GOOS,
		"arch":        osruntime.GOARCH,
	}
}
