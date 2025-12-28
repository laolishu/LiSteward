# 设计文档：Node 安装按钮功能

## 架构概述

### 核心流程

```
用户点击"安装"按钮
    ↓
前端调用 OpenNodeAvailableList()
    ↓
Wails 绑定 → App.OpenNodeAvailableList()
    ↓
后端路由到 nvm.OpenAvailableList()
    ↓
Windows: cmd.exe /k "nvm list available"
Unix:    终端执行 nvm list available
    ↓
系统弹出新窗口，显示可用版本列表
```

## 前端设计

### 布局调整

当前 NodeManager.vue 的"Node 仓库"行结构：
```vue
<div class="settings-item">
    <div class="item-label">{{ $t('node.nodeRepository') }}</div>
    <div class="item-control">
        <div class="nvm-version-display">{{ nvmRoot }}</div>
    </div>
</div>
```

修改后：
```vue
<div class="settings-item">
    <div class="item-label">{{ $t('node.nodeRepository') }}</div>
    <div class="item-control">
        <div class="nvm-version-display" style="flex: 0 0 70%">{{ nvmRoot }}</div>
        <button class="btn-install" @click="installNewVersion" 
                :title="$t('node.install')" aria-label="install-node-version">
            {{ $t('node.installButton') }}
        </button>
    </div>
</div>
```

### 样式设计

```css
/* 文本框容器 - flex 布局 */
.item-control {
    display: flex;
    gap: 12px;
    align-items: center;
}

/* 文本显示框 - 占 70% */
.nvm-version-display {
    flex: 0 0 70%;
    /* 保留现有样式 */
}

/* 安装按钮 - 占 30% */
.btn-install {
    flex: 0 0 auto;
    padding: 8px 16px;
    background-color: #1f6feb;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
}

.btn-install:hover {
    background-color: #388bfd;
}

.btn-install:active {
    background-color: #1f6feb;
}
```

### 事件处理

```javascript
async function installNewVersion() {
    try {
        await OpenNodeAvailableList()
    } catch (e) {
        console.error('Failed to open node available list:', e)
        // 可选：显示错误提示
    }
}
```

## 后端设计

### Wails 绑定 (app.go)

```go
// OpenNodeAvailableList 打开系统命令行窗口显示可用的 Node 版本列表
func (a *App) OpenNodeAvailableList() error {
    return nvm.OpenAvailableList()
}
```

### NVM 包实现 (internal/nvm/nvm.go)

```go
import "os/exec"

// OpenAvailableList 打开系统命令行窗口执行 nvm list available
func OpenAvailableList() error {
    var cmd *exec.Cmd
    
    if runtime.GOOS == "windows" {
        // Windows: 使用 cmd.exe 打开新窗口
        cmd = exec.Command("cmd", "/k", "nvm list available")
    } else {
        // Unix-like: 使用系统终端
        // 支持 macOS (open) 和 Linux (x-terminal-emulator/gnome-terminal)
        term := getTerminalCommand()
        cmd = exec.Command(term, "-e", "bash", "-c", "nvm list available; read")
    }
    
    // 关键：设置 CreationFlags 或其他参数使窗口显示
    // 让用户看到命令执行结果
    
    return cmd.Run()
}

func getTerminalCommand() string {
    // 检测可用的终端程序
    terminals := []string{
        "x-terminal-emulator",  // Debian/Ubuntu
        "gnome-terminal",       // GNOME
        "xfce4-terminal",       // XFCE
        "konsole",              // KDE
        "xterm",                // 通用
    }
    
    for _, term := range terminals {
        if _, err := exec.LookPath(term); err == nil {
            return term
        }
    }
    
    // 默认使用 bash -c 直接运行（会显示在调用终端中）
    return "bash"
}
```

## 平台特性

### Windows
- 使用 `cmd.exe /k` 打开新窗口并保持打开状态
- 用户可以看到完整的 `nvm list available` 输出
- `/k` 保证窗口在命令执行后不会立即关闭

### macOS/Linux
- 优先使用系统默认终端模拟器
- 通过 `bash -c "nvm list available"` 执行
- 添加 `read` 保证窗口不会立即关闭，用户可以查看结果

## i18n 集成

### zh-CN.json
```json
"node": {
    "title": "Node 版本管理",
    "installButton": "安装",
    ...
}
```

### en-US.json
```json
"node": {
    "title": "Node Version Management",
    "installButton": "Install",
    ...
}
```

## 错误处理

1. 如果 NVM 未安装，命令将失败
   - 用户会在命令行中看到相应错误信息
   - 前端可选择添加错误提示

2. 如果终端程序不可用（Linux），改为直接调用 `nvm list available`
   - 输出会显示在后台或调用进程的控制台中

## 测试策略

1. **单元测试**：测试 `getTerminalCommand()` 逻辑
2. **集成测试**：在各平台测试命令执行
3. **UI 测试**：
   - 验证按钮正确显示和对齐
   - 验证点击事件触发
   - 验证 i18n 文本正确

## 风险和考虑

1. **跨平台兼容性**：
   - 不同 Linux 发行版的终端程序差异
   - macOS 的终端处理方式
   - 缓解：优先级检查多个常见终端程序

2. **权限问题**：
   - 可能需要管理员权限运行 nvm 命令
   - 缓解：用户已需要管理员权限运行应用

3. **Windows 中 nvm-windows 的兼容性**：
   - nvm-windows 是 Node.js nvm 的 Windows 端口
   - 命令语法兼容，无额外处理需要

## 未来扩展

1. 可以在应用内显示可用版本列表（需新 UI）
2. 可以添加直接安装按钮（跳过命令行）
3. 可以集成版本检查更新功能
