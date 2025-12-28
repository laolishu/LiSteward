# 快速测试指南 - 节点安装功能

## 快速开始

### 场景 1: 开发模式测试 (wails dev)

```bash
# 1. 启动开发服务器
wails dev

# 2. 打开浏览器到 http://localhost:34115 (默认端口)

# 3. 导航到 "Node Manager" 模块

# 4. 点击 "Node 仓库" 行的 "安装" 按钮

# 预期结果:
# ✓ 打开命令行窗口显示 `nvm list available`
# ✓ 窗口保持打开，用户可交互
# ✓ 无错误或崩溃
```

### 场景 2: 编译版本测试 (wails build)

```bash
# 1. 编译程序
wails build -platform windows/amd64

# 2. 运行编译后的程序
# Windows: build/bin/LiSteward.exe
# macOS: build/bin/LiSteward.app
# Linux: build/bin/LiSteward

# 3. 打开应用后进入 "Node Manager" 模块

# 4. 点击 "安装" 按钮

# 预期结果:
# ✓ 打开命令行窗口显示 `nvm list available`
# ✓ NO 反复弹窗现象
# ✓ NO 窗口闪退
# ✓ 仅一个清晰的终端窗口对用户可见
```

## 故障排除

### 问题 1: 仍然看到多个闪窗

**可能原因**:
- Windows 未完全重新编译，旧版本缓存
- 某些杀毒软件阻止了进程创建

**解决方案**:
```bash
# 清除编译缓存
rm -r build/

# 重新编译
wails build -platform windows/amd64

# 禁用杀毒软件实时扫描（临时），重新测试
```

### 问题 2: 窗口打不开

**可能原因**:
- NVM 未安装或未在 PATH 中
- 权限问题

**检查方式**:
```bash
# 验证 NVM 是否可用
nvm --version

# 验证命令是否可运行
nvm list available
```

### 问题 3: macOS 上 Terminal 未打开

**可能原因**:
- Terminal.app 权限问题
- 系统设置阻止应用打开其他应用

**解决方案**:
```bash
# 手动运行等效命令测试
open -a Terminal -n
```

## 开发者调试

### 启用日志输出

编辑 `internal/nvm/nvm.go`，在 `OpenAvailableList()` 添加日志：

```go
import "log"

func OpenAvailableList() error {
    log.Printf("[NVM] Opening available list on %s\n", runtime.GOOS)
    
    if runtime.GOOS == "windows" {
        // ... 其他代码 ...
        err := cmd.Start()
        if err != nil {
            log.Printf("[NVM] Error: %v\n", err)
        }
        return err
    }
    // ...
}
```

### Windows 进程监控

在 Windows 任务管理器中：
1. 打开"详细信息"选项卡
2. 点击"安装"按钮
3. 观察进程树：
   - `LiSteward.exe` (主程序)
   - `cmd.exe` (中间启动器，应该很快消失)
   - `cmd.exe` (最终窗口，对用户可见)

## 验收标准

| 标准 | 满足条件 |
|-----|---------|
| **wails dev** | 窗口打开，正常工作 ✓ |
| **wails build** | 窗口打开，无闪现 ✓ |
| **多次点击** | 无内存泄漏，进程正常退出 ✓ |
| **跨平台** | Windows/macOS/Linux 都能打开终端 ✓ |
| **用户体验** | 清晰、快速、可交互 ✓ |

## 下一步

如果一切正常，将代码合并到主分支：
```bash
git add internal/nvm/nvm.go app.go frontend/src/views/NodeManager.vue frontend/src/i18n/
git commit -m "feat: add install node version button with silent mode"
git push
```

如果有问题，参考 [WAILS_BUILD_FIX.md](WAILS_BUILD_FIX.md) 和 [IMPLEMENTATION_VERIFICATION.md](IMPLEMENTATION_VERIFICATION.md)。
