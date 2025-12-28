# Wails Dev vs Build 的差异和修复说明

## 问题描述
使用 `wails build` 编译出的程序在打开 Node 模块时会连续弹出多个命令行窗口，然后自动关闭。

## 根本原因分析

### Wails Dev 与 Wails Build 的差异

| 特性 | `wails dev` | `wails build` |
|------|-----------|-------------|
| **运行环境** | 开发环境（Hot-reload） | 生产环境 |
| **调试工具** | DevTools 默认启用 | DevTools 禁用（除非指定 -devtools） |
| **进程隔离** | 与 Go 后端进程共享控制台 | 独立的可执行程序 |
| **窗口创建** | 通过 Wails 容器运行 | 直接系统执行 |
| **命令执行** | 可能继承父进程属性 | 完全独立进程 |
| **错误重试** | 开发时易见 | 后台不可见，但仍会反复执行 |

### 为什么会弹出多个窗口？

1. **Windows 特性**：`cmd /c start cmd /k "nvm list available"` 在编译版本中，`start` 命令的执行方式与开发模式不同
2. **继承属性**：编译后的程序创建子进程时，可能会创建**可见的**中间窗口
3. **错误处理**：如果命令失败（如 NVM 未安装），Go 程序可能会重试，导致反复弹窗

## 实现的修复方案

### 1. 隐藏中间窗口 (Windows)
```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    CreationFlags: 0x08000000, // CREATE_NO_WINDOW - 隐藏窗口
}
```

**作用**：
- `CREATE_NO_WINDOW` 标志会隐藏 `cmd /c start` 这个中间过程的窗口
- 只保留最终打开的 nvm 命令窗口对用户可见
- 防止反复弹窗的问题

### 2. 改进的跨平台处理

#### Windows
```go
cmd := ExecCommand("cmd", "/c", "start", "cmd", "/k", "nvm list available")
cmd.SysProcAttr = &syscall.SysProcAttr{
    CreationFlags: 0x08000000, // 隐藏中间窗口
}
cmd.Start() // 不等待，完全独立
```

#### macOS
```go
cmd := ExecCommand("open", "-a", "Terminal", "-n")
cmd.Start()
```

#### Linux
```go
cmd := ExecCommand("gnome-terminal", "--", "bash", "-c", "nvm list available; exec bash")
cmd.Start()
```

### 3. 错误处理
```go
// 简化为直接返回，不阻塞主程序
return cmd.Start()
```

## 验证方式

### 1. 使用 `wails dev` 测试
```bash
wails dev
# 打开 Node 模块，点击"安装"按钮
# 应该只弹出一个 nvm 窗口，无其他窗口
```

### 2. 使用 `wails build` 编译测试
```bash
# Windows
wails build -platform windows/amd64
# 运行编译后的程序：build/bin/LiSteward.exe
# 打开 Node 模块，点击"安装"按钮
# 应该只弹出一个 nvm 窗口，不再有反复弹窗

# macOS
wails build -platform darwin/universal
# 运行编译后的程序

# Linux
wails build -platform linux/amd64
```

## 技术细节

### Windows 进程创建标志
- `0x08000000` = `CREATE_NO_WINDOW` = 隐藏新窗口
- `0x00000200` = `CREATE_NEW_PROCESS_GROUP` = 新进程组（可选）
- `0x00000100` = `CREATE_NEW_CONSOLE` = 新控制台（不使用）

### 为什么只有 Build 版本有问题？

1. **Wails Dev**：前端开发服务器运行在同一进程中，命令执行通过 Wails 的事件系统管理
2. **Wails Build**：编译后的程序是独立的二进制文件，直接与系统交互
3. **进程隔离**：Build 版本的子进程没有 Wails 框架的保护和管理

## 后续优化建议

### 1. 异步处理优化
前端可以添加状态提示：
```vue
async function installNewVersion() {
    try {
        // 显示"正在打开..."的提示
        isOpeningTerminal.value = true
        await OpenNodeAvailableList()
        // 窗口已打开，清除提示
        isOpeningTerminal.value = false
    } catch (e) {
        console.error('Failed to open terminal:', e)
        isOpeningTerminal.value = false
    }
}
```

### 2. 事件通知（可选）
```go
// 添加事件通知，让前端知道窗口已打开
runtime.EventsEmit(ctx, "terminal-opened", map[string]string{
    "command": "nvm list available",
    "platform": runtime.GOOS,
})
```

### 3. 日志记录（可选）
```go
import "log"

func OpenAvailableList() error {
    // 记录打开操作
    log.Printf("[Node] Opening terminal for: nvm list available on %s", runtime.GOOS)
    // ... 执行逻辑
}
```

## 总结

| 改进点 | 修复前 | 修复后 |
|------|------|------|
| 窗口数量 | 多个弹窗 | 单一窗口 |
| 闪退问题 | 频繁闪退 | 窗口保持打开 |
| 用户体验 | 混乱、不可控 | 清晰、可交互 |
| 跨平台兼容 | 不稳定 | 完整覆盖 W/M/L |
| 代码质量 | 无状态管理 | 完全独立进程 |

修复后，`wails build` 编译的程序应该与 `wails dev` 表现一致。
