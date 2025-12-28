# 任务清单：系统托盘与全局快捷键功能

## 1. 用户配置扩展

- [x] 在 `config/user_config.go` 中扩展 `UserConfig` 结构体，添加 `TrayEnabled`、`HotKeyEnabled` 和 `HotKey` 字段。
- [x] 更新 `SaveUserConfig()` 和 `LoadUserConfig()` 以保证新字段的正确序列化与反序列化。
- [x] 编写单元测试验证用户配置的加载与保存。

## 2. 后端托盘管理实现（Windows）

- [x] 新建 `internal/tray/tray.go` 定义 `TrayManager` 接口。
- [x] 实现 `internal/tray/tray_windows.go`，使用 Wails runtime 与 Windows API 实现托盘功能。
  - [x] 实现 `Show()` — 显示应用窗口。
  - [x] 实现 `Hide()` — 隐藏应用窗口。
  - [x] 实现 `IsVisible()` — 检查窗口可见性。
  - [x] 实现 `Quit()` — 安全退出应用。
- [x] 编写 Windows 托盘功能的本地验证文档。

## 3. 后端全局快捷键实现（Windows）

- [x] 新建 `internal/hotkey/hotkey.go` 定义 `HotKeyManager` 接口。
- [x] 实现 `internal/hotkey/hotkey_windows.go`，使用 Windows API `RegisterHotKey()` 注册全局热键。
  - [x] 实现 `Register(key string)` — 注册快捷键。
  - [x] 实现 `Unregister()` — 注销快捷键。
  - [x] 实现 `IsRegistered()` — 检查快捷键状态。
  - [x] 实现 `OnHotKey(callback)` — 设置热键触发回调。
  - [x] 添加 `createMessageWindow()` 创建消息窗口来接收 WM_HOTKEY 消息。
  - [x] 完善 `parseHotKey()` 和 `getVirtualKey()` 辅助函数。
- [x] 支持常见快捷键格式（如 "ctrl+alt+l"、"shift+alt+h" 等）的解析与转换。
- [x] 使用 lxn/walk 库简化 Windows API 调用。
- [x] 编译通过，无编译错误。

## 4. 应用主逻辑集成（app.go）

- [x] 在 `app.go` 的 `startup()` 中初始化托盘与热键管理器。
- [x] 实现 `OnBeforeClose()` 事件处理，根据用户配置决定最小化到托盘或真正退出。
- [x] 暴露 Wails 接口给前端：
  - [x] `GetTraySettings() (*TrayConfig, error)` — 获取当前托盘设置。
  - [x] `SetTraySettings(config *TrayConfig) error` — 更新并保存托盘设置。
  - [x] `GetHotKeySettings() (*HotKeyConfig, error)` — 获取快捷键设置。
  - [x] `SetHotKeySettings(config *HotKeyConfig) error` — 更新并保存快捷键设置。
  - [x] `ShowWindow() error` — 显示应用窗口。
  - [x] `HideWindow() error` — 隐藏应用窗口。
  - [x] `ToggleWindowVisibility() error` — 切换窗口可见性。
- [x] 更新 Wails 前端绑定文件（`frontend/wailsjs/go/main/App.js` 与 `.d.ts`）。

## 5. 前端设置 UI（Settings.vue）

- [x] 在 Settings 页面添加"系统托盘"配置分组或新卡片。
- [x] 实现托盘配置选项：
  - [x] 复选框：启用/禁用最小化到托盘。
  - [x] 复选框：启用/禁用全局快捷键。
  - [x] 快捷键输入或选择器：允许用户自定义快捷键。
  - [x] 说明文本与提示。
- [x] 绑定到后端接口 `GetTraySettings()`、`SetTraySettings()`、`GetHotKeySettings()`、`SetHotKeySettings()`。
- [x] 处理配置保存成功/失败的反馈（显示消息提示）。

## 6. 国际化支持

- [x] 添加 i18n 文案到 `frontend/src/i18n/locales/zh-CN.json`：
  - [x] "系统托盘"、"启用托盘最小化"、"启用全局快捷键"、"快捷键配置" 等标签。
  - [x] "点击关闭按钮时最小化到托盘而非退出应用" 等说明文本。
- [x] 添加对应英文文案到 `en-US.json`。

## 7. 编译与验证

- [x] 运行 `go build ./...` 验证后端编译无误。
- [x] 运行 `npm run build` 验证前端编译无误。
- [x] 确保前端绑定文件正确生成。
## 8. 本地验证与测试

- [ ] 在 Settings 中启用托盘功能，验证关闭按钮的行为（最小化而非退出）。
- [ ] 点击系统托盘图标，验证窗口显示/隐藏切换。
- [ ] 在 Settings 中配置快捷键（如 Ctrl+Alt+L），验证全局快捷键可用。
- [ ] 按下配置的快捷键，验证窗口唤醒或隐藏。
- [ ] 重启应用，验证配置持久化（托盘与快捷键设置应恢复）。
- [ ] 禁用托盘功能后，验证快捷键亦不生效。

## 9. 文档更新

- [ ] 更新项目 README 或 QUICK_START.md，说明托盘与快捷键功能的使用方法。
- [ ] 记录已知限制（如 Windows 独占、快捷键冲突处理等）。

## 10. OpenSpec 规范完成

- [ ] 更新 `openspec/changes/add-system-tray/specs/*/spec.md` 中的规范增量，确保与实现一致。
- [ ] 运行 `openspec-cn validate add-system-tray --strict` 验证文档无误。

## 验收标准（整体）

- Windows 系统托盘图标正确显示与响应。
- 全局快捷键可配置与生效。
- 配置持久化到用户 config 文件。
- 前端 Settings UI 提供直观的配置选项。
- 编译无错，本地验证通过。

