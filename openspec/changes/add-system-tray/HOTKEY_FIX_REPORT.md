# Windows 快捷键实现修复记录

## 问题
运行时出现错误：`RegisterHotKey 失败: Invalid window handle. (快捷键可能已被其他应用占用)`

## 根本原因
之前的实现使用了无效的窗口句柄（占位符值 1），导致 Windows API 的 `RegisterHotKey` 调用失败。

## 解决方案

### 改进的实现内容

1. **真实的窗口类注册**
   - 使用 `WNDCLASSEX` 结构体定义窗口类
   - 调用 `RegisterClassExW` API 注册窗口类
   - 每个管理器实例有唯一的类名

2. **创建消息专用窗口**
   - 使用 `CreateWindowExW` 创建 `HWND_MESSAGE` 窗口
   - 消息专用窗口不可见，但可以接收消息
   - 窗口有有效的 HWND，可用于 RegisterHotKey

3. **完整的消息循环**
   - 实现 `messageLoop()` 真正处理 Windows 消息
   - 使用 `GetMessageW` 从消息队列获取消息
   - 使用 `DispatchMessageW` 分发消息给窗口过程

4. **全局窗口过程**
   - 创建全局的 `globalWndProc` 回调函数
   - 使用 `syscall.NewCallback` 将 Go 函数转换为 Windows 回调
   - 使用映射表将 HWND 映射到对应的管理器实例
   - 在 WM_HOTKEY 消息时调用注册的回调函数

### 关键代码改进

**之前（无效）**：
```go
m.hwnd = 1  // 占位符，无效
// RegisterHotKey 因无效 HWND 而失败
```

**现在（有效）**：
```go
// 注册窗口类
procRegisterClassExW.Call(...)

// 创建窗口
procCreateWindowExW.Call(
    0,               // dwExStyle
    className,       // lpClassName
    0,               // lpWindowName
    0,               // dwStyle
    0, 0, 0, 0,      // 位置和大小
    HWND_MESSAGE,    // hWndParent
    0,               // hMenu
    hModule,         // hInstance
    0,               // lpParam
)
// 返回有效的 HWND
```

## 编译验证

✅ `go build ./internal/hotkey` - 通过
✅ `go build -o build\bin\LiSteward.exe .` - 通过

## 预期效果

现在用户在 Settings 中配置快捷键并按下时，应该能正常响应（窗口显示/隐藏）。

## 技术细节

### Windows 消息流
1. 用户按下全局快捷键（如 Ctrl+Alt+L）
2. Windows 操作系统生成 WM_HOTKEY 消息
3. 消息被发送到注册的 HWND（我们的消息窗口）
4. `GetMessageW` 从消息队列中获取消息
5. `DispatchMessageW` 将消息分发给窗口过程
6. `globalWndProc` 处理 WM_HOTKEY 消息
7. 调用注册的回调函数，触发窗口操作

### 线程安全
- 使用 `sync.Mutex` 保护快捷键状态
- 使用全局 `wndProcMutex` 保护 HWND 映射表
- 回调函数在独立的 goroutine 中执行，避免阻塞消息循环

---

**修复日期**：2025-12-26
**版本**：v2 - 真实窗口和消息循环实现
