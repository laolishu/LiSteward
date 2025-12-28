# 设计说明：NVM 管理模块

## 背景

该模块提供对 `nvm`（Node Version Manager）的检测与管理能力，便于用户在本地通过 GUI 或命令行管理 Node 版本，减少手动配置带来的错误。

## 架构概览

- `internal/nvm`：后端实现，暴露 `NvmManager` 接口。
- `app.go`：将 `NvmManager` 的方法通过 Wails 绑定暴露给前端。
- 前端页面：`NodeManager.vue`（或扩展 `Settings.vue`）负责与用户交互并调用后端绑定。

交互流程：
1. 前端请求 `IsNvmInstalled()`。
2. 若未安装，前端可触发 `InstallNvm()`（弹出确认）。
3. 安装完成后，前端调用 `ListNodeVersions()` 展示已安装版本并提供操作按钮（Install/Uninstall/Use）。

## 实现要点

- 跨平台实现：
  - Unix (macOS/Linux)：优先使用官方安装脚本（curl | bash）或 git clone 到 `NVM_DIR`，并在调用时 source 对应脚本以使用 nvm 命令。
  - Windows：检测 `nvm-windows` 并指示用户；如需自动安装，使用 `nvm-setup.zip` 或提示使用 WSL。

- 命令执行：后端使用受限 shell 执行并捕获 stdout/stderr，提供进度回调与超时机制。

- 权限与安全：安装步骤需提示用户，并在受限环境下提供手动安装说明。

- 配置：`NVM_DIR` 可在用户配置中保存，前端可以通过 UI 修改并保存。

## API 设计（建议）

- `IsNvmInstalled() (bool, string)` - 判断是否安装并返回版本或错误信息。
- `InstallNvm(autoConfirm bool) (string, error)` - 安装 nvm，返回安装日志或路径。
- `GetNvmDir() (string, error)` - 获取当前 NVM_DIR
- `SetNvmDir(dir string) error` - 设置并保存 NVM_DIR
- `ListNodeVersions() ([]string, error)` - 列出已安装 Node 版本
- `ListRemoteNodeVersions() ([]string, error)` - 可选：查询远程可用版本
- `InstallNode(version string) error` - 安装 Node 指定版本
- `UninstallNode(version string) error` - 卸载指定版本
- `UseNode(version string) error` - 切换到指定版本

## 风险缓解

- 在 Windows 上优先展示兼容性提示，避免破坏系统 Node 安装。
- 对于自动安装步骤，提供明确的回滚/错误信息。
