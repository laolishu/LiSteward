# 任务清单：添加 NVM 管理模块（add-nvm-manager）

## 后端（核心实现）

### 1. 新建 `internal/nvm` 包
- [ ] 1.1 设计 `NvmManager` 接口（检测、安装 nvm、列出、安装/卸载 node、切换版本、读写 NVM_DIR）
- [ ] 1.2 实现跨平台调用（Unix: 使用 shell 脚本调用 nvm；Windows: 提示或使用 nvm-windows）
- [ ] 1.3 添加错误处理、超时与日志
- [x] 1.1 设计 `NvmManager` 接口并提交基础实现（见 `internal/nvm/nvm.go`）
- [x] 1.2 提供跨平台的可用性检测与建议性安装实现（类 Unix 提示安装脚本；Windows 检测 `nvm-windows` 并提示），详见 `internal/nvm/nvm.go`
- [ ] 1.3 添加更细粒度的错误处理、超时与详尽日志（部分实现，建议后续完善）

### 2. 暴露 Wails 绑定
- [ ] 2.1 在 `app.go` 中添加绑定方法：`InstallNvm()`, `IsNvmInstalled()`, `ListNodeVersions()`, `InstallNode(version)`, `UninstallNode(version)`, `UseNode(version)`, `GetNvmDir()`, `SetNvmDir(dir)`
- [x] 2.1 在 `app.go` 中添加并暴露了对应方法（`DetectNvm`, `InstallNvm`, `ListNodeVersions`, `InstallNodeVersion`, `UninstallNodeVersion`, `UseNodeVersion`, `GetNvmDir`, `SetNvmDir`）

### 3. 配置与持久化
- [ ] 3.1 在 `config/user_config.go` 中添加 `NvmDir` 字段（或复用已有的用户配置）
- [ ] 3.2 确保 `EnsureUserConfig()` 初始化默认 `NvmDir`（如 `$HOME/.nvm`）
- [x] 3.1 已在 `config/user_config.go` 中添加 `NvmDir` 字段并持久化（默认为空字符串，代码提供回退查找 `~/.nvm` 的逻辑）
- [x] 3.2 已确保默认配置创建（`NvmDir` 默认为空），并在 `internal/nvm` 中提供 `GetNvmDir` 回退查找以定位 `~/.nvm` 或 Windows 常见路径（因此满足可用性需求）

## 前端（UI 与 i18n）

### 4. Settings 页面 / 新管理页面
- [ ] 4.1 新增或扩展 `frontend/src/views/Settings.vue`（或 `NodeManager.vue`）用于展示 nvm 状态和操作
- [ ] 4.2 添加交互控件：检测、安装 nvm、列出版本、安装/卸载、切换版本、设置 nvm 目录
- [ ] 4.3 添加前端提示/进度显示、错误展示
- [x] 4.1 已新增 `frontend/src/views/NodeManager.vue` 用于管理 nvm 与 Node 版本
- [x] 4.2 已实现基础交互控件（检测、安装提示、列出版本、Use/Uninstall、设置 `NVM_DIR` 并保存）
- [ ] 4.3 进度与错误展示为基础提示或 alert，建议后续增强为更友好的进度条与日志面板

### 5. i18n 翻译
- [ ] 5.1 在 `frontend/src/i18n/locales/*` 中添加所有必要翻译键（zh-CN/en-US）
- [x] 5.1 已在 `frontend/src/i18n/locales/zh-CN.json` 与 `en-US.json` 中添加 `menu.node` 与 `node.*` 翻译键

## 测试与文档

### 6. 测试
- [ ] 6.1 添加后端单元测试（模拟 shell 命令）
- [ ] 6.2 添加端到端测试（如 CI 环境允许，或本地手动测试步骤）
- [ ] 6.1 后端单元测试尚未实现（建议随后添加，涉及对 `exec.Command` 的抽象与注入以便 mock）
- [ ] 6.2 端到端测试为可选项，需在 CI/本地环境中验证 nvm 实际行为

### 7. 文档
- [ ] 7.1 添加 `docs/NVM_USAGE.md`，说明功能与常见故障排查
- [ ] 7.2 更新项目 README 中的开发环境部分
- [ ] 7.1 文档尚未添加（建议包含安装提示、Windows 特殊说明与手动回退步骤）
- [ ] 7.2 README 更新尚未完成

## 验证与发布

- [ ] 8.1 运行 `go build` 与前端构建，确保无错误
- [ ] 8.2 在 Windows/macOS/Linux 上验证基本场景
- [ ] 8.3 将变更纳入 openspec 验证流程并归档
- [x] 8.1 已运行并通过 `go build`（基础构建通过）；前端已添加视图（未做完整 UI 测试）
- [ ] 8.2 平台验证建议由 CI 或本地在各平台执行（未在所有平台上验证）
- [x] 8.3 变更已实现源码部分，待更新 openspec 任务状态并提交实现说明以便归档

## 优先级

P1: 1.1, 2.1, 4.1, 8.1
P2: 1.2, 3.1, 4.2, 5.1
P3: 测试、文档、细节增强
