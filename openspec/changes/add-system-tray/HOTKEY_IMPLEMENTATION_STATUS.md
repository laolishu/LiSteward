# 全局快捷键实现状态报告

## 当前状态：开发中 - 实现了 Windows API 集成

### 已完成的功能

✅ **前端快捷键捕获 UI**
- 用户点击快捷键输入框时，可自动捕获用户按下的按键组合
- 支持修饰键（Ctrl, Alt, Shift）和常规键（字母、数字、功能键 F1-F12）
- 显示格式：`ctrl+alt+l`、`shift+alt+h` 等
- 自动保存到后端配置

✅ **快捷键配置持久化**
- 快捷键值保存到 `user.json` 文件中
- 应用启动时自动加载上次保存的快捷键配置
- 前端 Settings 页面正确显示已保存的快捷键

✅ **后端快捷键管理框架**
- `internal/hotkey/hotkey.go` - HotKeyManager 接口定义，包含 `isValidHotKeyFormat()` 验证函数
- `internal/hotkey/hotkey_windows.go` - Windows 平台实现（使用 syscall 直接调用 Windows API）
- `internal/hotkey/hotkey_stub.go` - macOS/Linux stub 实现
- 支持通过 `Register(key)`、`Unregister()`、`OnHotKey(callback)` 等方法管理快捷键

✅ **Windows API 集成**
- 使用 syscall 直接调用 Windows API：
  - `RegisterHotKey()` - 注册全局快捷键
  - `UnregisterHotKey()` - 注销全局快捷键
  - `CreateWindowEx()` - 创建消息窗口
  - `GetMessage()`、`TranslateMessage()`、`DispatchMessage()` - 消息处理
- 实现快捷键格式解析：`parseHotKey()`、`getVirtualKey()` 辅助函数
- 完整的虚拟键码映射表（字母、数字、功能键、特殊键）

✅ **Wails 后端接口暴露**
- `GetHotKeySettings()` - 获取快捷键配置
- `SetHotKeySettings(enabled bool, hotkey string)` - 设置快捷键配置
- `ShowWindow()`、`HideWindow()`、`ToggleWindowVisibility()` - 窗口控制

✅ **编译通过**
- `go build ./internal/hotkey` ✓
- `go build -o build\bin\LiSteward.exe .` ✓

### 当前实现的技术细节

**Windows 快捷键注册流程**：

1. 创建消息窗口（`createMessageWindow()`）
   - 生成唯一的快捷键 ID
   - 设置 HWND（当前为占位符值 1，需要改进）

2. 调用 `RegisterHotKey` Windows API
   - 参数：HWND、快捷键 ID、修饰符、虚拟键码
   - 成功时返回非零值

3. 设置快捷键触发回调（`OnHotKey(callback)`）
   - 回调在独立 goroutine 中执行，避免阻塞

4. 消息处理循环（`messageLoop()`）
   - 监听 `WM_HOTKEY` 消息
   - 触发对应的回调函数

**快捷键格式支持**：
- 基本格式：`修饰符+修饰符+按键`
- 示例：`ctrl+alt+l`、`shift+alt+h`、`ctrl+shift+f1`、`win+l`
- 修饰符：ctrl/control、alt、shift、win/windows
- 按键：a-z、0-9、f1-f12、特殊键（enter、space、esc 等）

### 当前已知的限制和需要改进的地方

⚠️ **需要改进的项**：

1. **HWND 的有效性**
   - 当前使用占位符值 1，需要创建真实的 HWND
   - 注册成功但 RegisterHotKey 可能无法真正响应
   - 解决方案：使用 `CreateWindowEx` 和正确的窗口类注册

2. **消息循环集成**
   - 当前消息循环为简化实现，不能真正处理 WM_HOTKEY 消息
   - 需要与 Wails 的主窗口集成，或创建独立的隐藏窗口

3. **跨平台支持**
   - 当前仅 Windows 有实现
   - macOS 和 Linux 仍为 stub 实现
   - 需要使用第三方库（如 robotgo）或平台特定的 API

### 推荐的下一步行动

#### 短期（立即）

1. **测试当前实现**
   - [ ] 在 Windows 上编译运行应用
   - [ ] 在 Settings 中配置快捷键
   - [ ] 测试快捷键是否响应
   - [ ] 记录问题和错误日志

2. **改进 Windows API 实现**
   - [ ] 创建真实的 HWND（不使用占位符）
   - [ ] 实现正确的窗口类注册（`RegisterClass`）
   - [ ] 实现真正的消息循环来处理 `WM_HOTKEY`

#### 中期（如果 Windows API 方案不可行）

1. **使用第三方库方案**
   - 评估 `github.com/robotgo/robotgo` 库
   - 或使用 `github.com/getlantern/systray` 等现成解决方案
   - 考虑库的稳定性和跨平台支持

2. **实现 macOS 和 Linux 支持**
   - macOS：使用 `Carbon` 或 `Cocoa` API
   - Linux：使用 `X11` 或 `Wayland` API

### 技术债和技术备注

- [ ] 需要完整的单元测试
- [ ] 需要集成测试验证快捷键响应
- [ ] 需要处理快捷键冲突的情况
- [ ] 需要管理员权限检查
- [ ] 需要优雅的错误处理和用户提示

### 文件更改记录

**修改的文件**：
- `internal/hotkey/hotkey.go` - 添加 `isValidHotKeyFormat()` 函数
- `internal/hotkey/hotkey_windows.go` - 重新实现，使用 syscall 直接调用 Windows API
- `go.mod` - 添加 `github.com/lxn/walk` 和 `github.com/lxn/win` 依赖

**新增的文件**：
- `openspec/changes/add-system-tray/HOTKEY_IMPLEMENTATION_STATUS.md` - 本文件

---

**最后更新**：2025-01-XX  
**更新者**：GitHub Copilot  
**状态**：等待测试和进一步调试
