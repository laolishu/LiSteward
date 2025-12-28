package nvm

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// ExecCommand 可在测试中替换以 mock 命令执行行为
var ExecCommand = exec.Command

// runCommand 通过 ExecCommand 运行命令并返回合并输出
func runCommand(name string, args ...string) (string, error) {
	cmd := ExecCommand(name, args...)
	//静默运行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// NodeVersion 表示一个已安装的 Node 版本
type NodeVersion struct {
	Version   string `json:"version"`
	Path      string `json:"path"`
	IsCurrent bool   `json:"isCurrent"`
}

// Detect 检测系统中是否存在 nvm 或 nvm-windows，并返回发现标志与可能的 NVM_DIR
func Detect() (bool, string) {
	// 优先检查环境变量 NVM_DIR
	if dir := os.Getenv("NVM_DIR"); dir != "" {
		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			return true, dir
		}
	}

	// 在 PATH 中查找可执行 nvm（在 windows 上通常是 nvm.exe）
	names := []string{"nvm"}
	if runtime.GOOS == "windows" {
		names = []string{"nvm.exe", "nvm"}
		// 常见的 nvm-windows 安装目录（AppData）
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			cand := filepath.Join(appdata, "nvm")
			if stat, err := os.Stat(cand); err == nil && stat.IsDir() {
				return true, cand
			}
		}
	}

	for _, n := range names {
		if _, err := exec.LookPath(n); err == nil {
			// 找到可执行文件，但不一定有 NVM_DIR，返回空字符串
			return true, ""
		}
	}

	// 常见 UNIX nvm 路径 ~/.nvm
	if home, err := os.UserHomeDir(); err == nil {
		cand := filepath.Join(home, ".nvm")
		if stat, err := os.Stat(cand); err == nil && stat.IsDir() {
			return true, cand
		}
	}

	return false, ""
}

// GetNvmVersion 获取 NVM 的版本号
func GetNvmVersion() (string, error) {
	out, err := runCommand("nvm", "version")
	if err != nil {
		return "", err
	}
	// nvm version 返回类似 "v0.39.5" 的输出
	version := strings.TrimSpace(out)
	if version == "" {
		return "", errors.New("无法获取 nvm 版本")
	}
	return version, nil
}

// GetNvmRoot 获取 Node 仓库路径（使用 nvm root 命令）
func GetNvmRoot() (string, error) {
	out, err := runCommand("nvm", "root")
	if err != nil {
		return "", err
	}
	// nvm root 返回 Node 仓库的路径，可能带有 "Current Root:" 前缀
	root := strings.TrimSpace(out)
	if root == "" {
		return "", errors.New("无法获取 nvm root")
	}
	// 移除 "Current Root:" 前缀
	if strings.HasPrefix(root, "Current Root:") {
		root = strings.TrimSpace(strings.TrimPrefix(root, "Current Root:"))
	}
	return root, nil
}

// ListVersions 尝试列出已安装的 Node 版本（尽可能安全地）
func ListVersions(nvmDir string) ([]NodeVersion, error) {
	// 若提供了 nvmDir，则列出其子目录作为版本
	if nvmDir != "" {
		entries, err := os.ReadDir(nvmDir)
		if err == nil {
			var res []NodeVersion
			for _, e := range entries {
				if e.IsDir() {
					name := e.Name()
					// 过滤常见文件夹
					if strings.HasPrefix(name, "v") || strings.HasPrefix(name, "node-v") || strings.Contains(name, "node") {
						res = append(res, NodeVersion{Version: name, Path: filepath.Join(nvmDir, name)})
					}
				}
			}
			return res, nil
		}
	}

	// 回退：如果 nvm 可执行并在 PATH 中，尝试运行 `nvm ls`（在某些系统上可用）
	if path, err := exec.LookPath("nvm"); err == nil && path != "" {
		out, err := runCommand(path, "ls")
		if err != nil {
			return nil, err
		}
		return parseNvmLsOutput(string(out)), nil
	}

	return nil, errors.New("无法列出 Node 版本: 未找到 NVM_DIR 或 nvm")
}

// parseNvmLsOutput 简单解析 nvm ls 的输出，返回版本列表，并标识当前版本
func parseNvmLsOutput(out string) []NodeVersion {
	scanner := bufio.NewScanner(strings.NewReader(out))
	var res []NodeVersion
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 过滤空行和系统版本行
		if line == "" || strings.HasPrefix(line, "system") {
			continue
		}

		// 检测是否为当前版本（以 * 开头）
		isCurrent := strings.HasPrefix(line, "*")

		// 移除行首的 * 或 > 或 -> （表示当前使用或别名）
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, "->")
		line = strings.TrimSpace(line)

		// 提取版本号：从行首到第一个空格或括号
		// 如: "20.19.4" 或 "20.19.4 (Currently using 64-bit executable)"
		versionPart := line
		if idx := strings.IndexAny(versionPart, " ("); idx > 0 {
			versionPart = versionPart[:idx]
		}
		versionPart = strings.TrimSpace(versionPart)

		// 验证是否为有效版本号（以 v 开头或数字开头）
		if len(versionPart) > 0 && (strings.HasPrefix(versionPart, "v") || (versionPart[0] >= '0' && versionPart[0] <= '9')) {
			if !containsVersion(res, versionPart) {
				res = append(res, NodeVersion{Version: versionPart, IsCurrent: isCurrent})
			}
		}
	}
	return res
}

func containsVersion(list []NodeVersion, ver string) bool {
	for _, v := range list {
		if v.Version == ver {
			return true
		}
	}
	return false
}

// InstallNvm 给出指导性实现：在可自动化平台上尝试安装；当前实现返回提示
func InstallNvm() (string, error) {
	if runtime.GOOS == "windows" {
		return "请从 https://github.com/coreybutler/nvm-windows/releases 下载并安装 nvm-windows。", nil
	}
	// 对于类 Unix 系统，建议使用官方安装脚本
	return "请在终端运行: curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash", nil
}

// InstallVersion 尝试使用 nvm 安装指定版本（如果可用）
func InstallVersion(version string) (string, error) {
	// 直接调用命令（通过 ExecCommand，可在测试中 mock）
	out, err := runCommand("nvm", "install", version)
	if err != nil {
		return out, err
	}
	return out, nil
}

// UseVersion 尝试切换到指定版本
func UseVersion(version string) (string, error) {
	out, err := runCommand("nvm", "use", version)
	if err != nil {
		return out, err
	}
	return out, nil
}

// UninstallVersion 尝试卸载指定版本
func UninstallVersion(version string) (string, error) {
	out, err := runCommand("nvm", "uninstall", version)
	if err != nil {
		return out, err
	}
	return out, nil
}

// GetNvmDir 返回传入或环境中的 NVM_DIR
func GetNvmDir(fallback string) string {
	if fallback != "" {
		return fallback
	}
	if dir := os.Getenv("NVM_DIR"); dir != "" {
		return dir
	}
	if home, err := os.UserHomeDir(); err == nil {
		cand := filepath.Join(home, ".nvm")
		if stat, err := os.Stat(cand); err == nil && stat.IsDir() {
			return cand
		}
	}
	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			cand := filepath.Join(appdata, "nvm")
			if stat, err := os.Stat(cand); err == nil && stat.IsDir() {
				return cand
			}
		}
	}
	return ""
}

// getTerminalCommand 检测系统可用的终端程序
func getTerminalCommand() string {
	// 检测可用的终端程序（按优先级）
	terminals := []string{
		"x-terminal-emulator", // Debian/Ubuntu
		"gnome-terminal",      // GNOME
		"xfce4-terminal",      // XFCE
		"konsole",             // KDE
		"xterm",               // 通用
	}

	for _, term := range terminals {
		if _, err := exec.LookPath(term); err == nil {
			return term
		}
	}

	// 默认使用 bash（会在调用进程的控制台中运行）
	return "bash"
}

// OpenAvailableList 打开系统命令行窗口执行 nvm list available
// 窗口完全独立，不受主程序管理
func OpenAvailableList() error {
	if runtime.GOOS == "windows" {
		// Windows: 直接启动新的 cmd 窗口，中间过程隐藏
		cmd := ExecCommand("cmd", "/c", "start", "cmd", "/k", "nvm list available")
		// 隐藏 start 命令本身的窗口
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW - 隐藏窗口
		}
		// 启动进程但不等待 - 窗口独立存在
		return cmd.Start()
	} else if runtime.GOOS == "darwin" {
		// macOS: 使用 open 命令打开 Terminal.app
		// Terminal 会独立打开新窗口，主程序无需管理
		cmd := ExecCommand("open", "-a", "Terminal", "-n")
		// 为 open 命令本身设置隐藏属性（可选）
		cmd.SysProcAttr = &syscall.SysProcAttr{
			// macOS 上的等效隐藏设置
		}
		// 启动进程
		return cmd.Start()
	} else {
		// Unix/Linux: 在后台启动终端程序
		term := getTerminalCommand()
		var cmd *exec.Cmd

		if term == "gnome-terminal" {
			cmd = ExecCommand("gnome-terminal", "--", "bash", "-c", "nvm list available; exec bash")
		} else if term == "xfce4-terminal" {
			cmd = ExecCommand("xfce4-terminal", "-e", "bash -c 'nvm list available; exec bash'")
		} else if term == "konsole" {
			cmd = ExecCommand("konsole", "-e", "bash", "-c", "nvm list available; exec bash")
		} else if term == "xterm" {
			cmd = ExecCommand("xterm", "-e", "bash", "-c", "nvm list available; exec bash")
		} else {
			cmd = ExecCommand("x-terminal-emulator", "-e", "bash", "-c", "nvm list available; exec bash")
		}

		// 启动进程但不等待
		return cmd.Start()
	}
}
