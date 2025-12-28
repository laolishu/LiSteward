<!--
 * @Descripttion: 
 * @version: 
 * @Author: lfzxs@qq.com
 * @Date: 2025-12-26 15:16:16
 * @LastEditors: lfzxs@qq.com
 * @LastEditTime: 2025-12-26 15:39:38
-->
# 提案：添加开机自启动功能 (add-autostart)
## 为什么

用户希望在启动操作系统时，LiSteward 应用能自动运行，无需手动启动。这在需要长期管理 Hosts 配置或监控特定网络设置的场景中很有用。当前应用需要手动启动，降低了使用体验。通过提供跨平台开机自启功能，应用可以更好地支持用户的工作流。

## 变更内容

本变更为桌面应用增加了以下功能：

1. **后端跨平台自启提供者**：新增 `internal/autostart` 包，包含抽象接口与 Windows/macOS/Linux 平台实现。
2. **用户配置持久化**：在 `UserConfig` 中添加 `AutoStart` 字段，允许用户设置被保存。
3. **Wails 接口暴露**：在 `app.go` 添加 `GetAutoStart()` 与 `SetAutoStart()` 方法供前端调用。
4. **前端 UI**：在 Settings 页面添加"开机自启动"开关，允许用户启用/禁用此功能。
5. **国际化支持**：为中文与英文添加相应标签文本。
概述
----
为桌面应用添加“开机自启动（开机自启 / 启动时运行）”功能，允许用户在设置中启用或禁用应用随系统启动自动运行。功能需跨平台支持 Windows、macOS、Linux（常见发行版），并提供安全、可回退的实现。

目标
----
- 为前端提供简单的开关（启用/禁用）接口。
- 后端实现平台特定的启用/禁用逻辑，并返回当前状态。
- 保持最小权限原则：实现无需提升管理员权限（Windows 下通过当前用户的启动项或快捷方式实现，而非修改全局注册表项）。
- 提供合理的失败回退和日志，供前端显示错误信息。

约束
----
- 不自动在安装时开启（默认关闭），由用户在设置中主动启用。
- 在实现过程中尽可能使用标准用户级方法：
  - Windows：创建快捷方式到 `shell:startup`（当前用户的启动文件夹）或使用注册表 HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run（仅当快捷方式不可用时）
  - macOS：使用 Launch Agent (`~/Library/LaunchAgents/`) 或 Apple Script 注册（以用户级启动项为准）
  - Linux：创建 `~/.config/autostart/*.desktop` 文件，或在无图形环境下记录为不支持。

验收标准（高层）
----
- 在 Windows/macOS/Linux 上：调用后端的 `SetAutoStart(true)` 后，下次登录时应用自动启动。
- `IsAutoStartEnabled()` 能准确报告当前用户级启动项状态。
- 前端在设置页提供开关与友好提示（为何需要开机自启、如何撤销）。

交付物
----
- openspec 相关文档（proposal/tasks/design/specs）。
- 后端接口原型建议：`GetAutoStart()` / `SetAutoStart(bool)`（Wails 暴露）。
- 测试计划与验证步骤（见 tasks.md）。

负责人
----
提案编写者: 自动生成（需要分配实现者）。
