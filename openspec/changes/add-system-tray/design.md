# 设计说明：系统托盘与全局快捷键（System Tray & Global Hotkey）

## 目标要点

- 在 Windows 系统托盘中显示应用图标，允许用户最小化应用到托盘。
- 实现全局快捷键注册与事件监听，支持用户在系统任何地方快速唤醒应用。
- 通过清晰的接口抽象为未来跨平台扩展（macOS、Linux）保留可扩展性。
- 保持实现简单、可测试、可配置。

## 设计原则

1. **用户级配置**：托盘与快捷键行为由用户配置控制，默认关闭，用户可在 Settings 中启用/禁用。
2. **最小权限**：不需要提升权限，在应用运行时启用快捷键。
3. **清晰的事件流**：
   - 窗口关闭事件 → 根据配置决定最小化到托盘或真正退出。
   - 托盘点击 → 切换窗口显示/隐藏。
   - 快捷键按下 → 触发唤醒或隐藏窗口。
4. **国际化支持**：托盘菜单项、提示文本需提供中英文。

## 后端架构（建议）

### 1. 扩展用户配置结构

在 `config/user_config.go` 添加托盘与快捷键配置字段：

```go
type UserConfig struct {
    Language       string `json:"language"`
    AutoStart      bool   `json:"auto_start"`
    TrayEnabled    bool   `json:"tray_enabled"`         // 新增：是否启用托盘最小化
    HotKeyEnabled  bool   `json:"hotkey_enabled"`       // 新增：是否启用快捷键
    HotKey         string `json:"hotkey"`               // 新增：快捷键配置，格式如 "ctrl+alt+l"
}
```

### 2. 系统托盘管理

在 `internal/tray/` 包中实现托盘相关逻辑：

```go
// internal/tray/tray.go
type TrayManager interface {
    Show() error              // 显示应用窗口
    Hide() error              // 隐藏应用窗口
    IsVisible() (bool, error) // 检查窗口是否可见
    Quit() error              // 退出应用
}

// Windows 特定实现
type windowsTrayManager struct {
    // 使用 Wails runtime 与 Windows API
}

func NewTrayManager() TrayManager {
    // 根据 runtime.GOOS 返回具体实现
}
```

### 3. 全局热键管理

在 `internal/hotkey/` 包中实现热键相关逻辑：

```go
// internal/hotkey/hotkey.go
type HotKeyManager interface {
    Register(key string) error      // 注册热键
    Unregister() error              // 注销热键
    IsRegistered() (bool, error)    // 检查热键是否已注册
    OnHotKey(callback func())       // 设置热键触发回调
}

// Windows 特定实现
type windowsHotKeyManager struct {
    // 使用 Windows API RegisterHotKey
}

func NewHotKeyManager() HotKeyManager {
    // 根据 runtime.GOOS 返回具体实现
}
```

### 4. 应用主逻辑修改

在 `app.go` 中集成托盘与热键管理：

- 在 `startup()` 中初始化托盘与热键管理器。
- 在 `OnBeforeClose` 事件中根据配置决定最小化或退出。
- 暴露 `GetTraySettings()`、`SetTraySettings()` 与 `GetHotKeySettings()`、`SetHotKeySettings()` 等 Wails 接口供前端调用。
- 暴露 `ShowWindow()`、`HideWindow()`、`ToggleWindowVisibility()` 接口供托盘菜单调用。

## 前端集成

### Settings 页面扩展

在 `frontend/src/views/Settings.vue` 中添加新的配置选项卡或分组：

- "托盘与快捷键"：
  - 复选框：启用/禁用托盘最小化。
  - 复选框：启用/禁用快捷键。
  - 输入框（或快捷键选择器）：配置快捷键（如 Ctrl+Alt+L）。
  - 说明文本：快捷键生效后的行为（唤醒或隐藏）。

### 国际化

添加相应的 i18n 文案到 `frontend/src/i18n/locales/zh-CN.json` 与 `en-US.json`。

## 平台特定实现考虑

### Windows

- **托盘**：使用 Wails 的 `runtime.WindowHide()`、`runtime.WindowShow()`、`runtime.WindowSetAlwaysOnTop()` 等。
- **快捷键**：使用 Windows API `RegisterHotKey()` / `UnregisterHotKey()`，监听 `WM_HOTKEY` 消息。
- **库支持**：可使用 [github.com/lxn/win](https://github.com/lxn/win) 或 [github.com/robotn/gohook](https://github.com/robotn/gohook) 等库简化快捷键实现。

### macOS（预留设计）

- **托盘**：使用 `github.com/getlantern/systray` 或 Wails 的 macOS 原生支持。
- **快捷键**：使用 Carbon Framework 或 `github.com/robotn/gohook`。

### Linux（预留设计）

- **托盘**：使用 `github.com/getlantern/systray` 或 `github.com/tidwall/statichunt`。
- **快捷键**：使用 X11 或 Wayland API（复杂度较高）。

## 错误处理与日志

- 若快捷键注册失败，应返回错误并记录日志，允许用户在设置中重试或更换快捷键。
- 若托盘操作失败（如显示窗口失败），应捕获并向用户显示友好提示。

## 可测试要点

- 在 Settings 中配置托盘与快捷键后，应用应正确注册并响应。
- 点击托盘图标应切换窗口显示/隐藏。
- 按下全局快捷键应唤醒或隐藏应用。
- 配置应持久化到用户 config 文件，重启应用后恢复。

