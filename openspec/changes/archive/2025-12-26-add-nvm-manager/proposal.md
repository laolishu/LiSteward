<!--
 变更: add-nvm-manager
 作者: 自动提案
 日期: 2025-12-26
-->
# 提案：添加 NVM 管理模块（Node 版本管理）

## 背景与动机

当前项目需要在用户环境中管理 Node.js 版本（用于前端构建、插件或开发工具）。引入一个原生支持的 NVM（Node Version Manager）管理模块，可以让用户通过应用界面或命令行：检测是否安装 nvm、安装 nvm、列出已安装的 Node 版本、安装/卸载指定版本、切换当前使用的版本，以及配置 nvm 存储目录（`NVM_DIR`）。

引入该模块可以提升构建可靠性、简化内部测试流程并为高级用户提供便利的环境管理能力。

## 目标

- 在后端添加一个 `internal/nvm` 包，封装对 nvm 的检测、安装、版本管理和目录配置。
- 在后端暴露 Wails 绑定接口供前端调用（例如：`InstallNvm()`、`ListNodeVersions()`、`InstallNode(version)`、`UninstallNode(version)`、`UseNode(version)`、`GetNvmDir()`、`SetNvmDir(dir)`）。
- 在前端 Settings 或单独页面中提供 UI，用于查看与管理 Node 版本和 nvm 设置。
- 提供跨平台工作流（Windows/macOS/Linux）；在 Windows 平台上可选使用 `nvm-windows` 或提示用户使用 WSL/其他方案。
- 可配置 nvm 存储目录，支持用户自定义位置并将配置持久化到用户配置文件。

## 范围（包含/排除）

- 包含：检测/安装 nvm、列举/安装/卸载 Node 版本、切换版本、设置/读取 `NVM_DIR`、前端 UI 与 i18n 条目、必要的单元与集成测试、文档更新。
- 排除：替换系统范围的 Node（不自动修改系统 PATH 除非用户确认）、远程替代方案（如容器化 Node 管理）。

## 验收标准

- 必须：能在受支持平台上检测 nvm 是否已安装并返回准确状态。
- 必须：当未安装 nvm 时，模块能自动安装最新稳定版本（或按用户确认安装）。
- 必须：能列出本地可用/已安装的 Node 版本，并能安装/卸载指定版本。
- 必须：能将当前 Node 版本切换到指定版本（对当前 shell 或通过修改 nvm 管理的 symlink/配置）。
- 必须：提供设置 `NVM_DIR` 的接口并持久化到用户配置（`user.json`）。
- 必须：前端 UI 能调用并展示当前状态、执行安装/切换/卸载操作并显示进度与错误。
- 必须：相关 i18n 条目已添加。

## 影响范围

- 后端：新增 `internal/nvm` 包，修改 `app.go` 暴露绑定方法。
- 前端：添加或扩展 Settings 页面，添加管理界面组件与翻译文本。
- 文档：更新 README 与用户指南，添加 `NVM_USAGE.md`（可选）。

## 风险与兼容性

- Windows 平台存在 `nvm-windows` 的差异，需特别处理；若复杂则提供安装提示或建议使用 WSL。
- 在受限权限环境（受管理的公司机器）安装 nvm 可能失败，需提供回退提示与手动步骤。

---

下一步：起草 `tasks.md`、`design.md` 与规范增量 `specs/node-runtime/spec.md`，并运行 `openspec-cn validate add-nvm-manager --strict`。
