# 设计说明：开机自启动（AutoStart）

目标要点
----
- 提供用户可控的开机自启动功能，跨 Windows/macOS/Linux 支持。
- 优先使用用户级方案以避免需要管理员权限或影响其他用户。
- 保持实现简单、可测试、可回退。

设计原则
----
1. 最小权限与用户级：始终使用当前用户的启动机制（Startup folder / LaunchAgents / ~/.config/autostart）。
2. 可检测与可回退：实现读取方法 `GetAutoStart()` 用于检测当前是否启用；若写入失败返回清晰错误信息并写日志。
3. 抽象化：后端实现采用适配器模式（interface AutoStartProvider），每个平台实现具体逻辑，便于单元测试与模拟。

后端架构（建议）
----
定义接口：

```go
type AutoStartProvider interface {
    IsEnabled() (bool, error)
    Enable() error
    Disable() error
}
```

在运行时，根据 `runtime.GOOS` 实例化具体实现：
- windowsProvider：通过 `os` 与 COM 或直接创建 `.lnk` 快捷方式到 `shell:startup`（使用 `github.com/akavel/rsrc` 或调用 PowerShell），作为回退可写入 HKCU Run。
- darwinProvider：生成 `~/Library/LaunchAgents/com.laolishu.listeward.plist`，内容指定可执行路径与参数；写入后调用 `launchctl bootstrap`（可选）或让用户在下次登录生效。
- linuxProvider：生成 `.desktop` 文件至 `~/.config/autostart/`，确保 `Exec` 指向当前可执行文件并设置 `X-GNOME-Autostart-enabled=true`。

实现注意事项
----
- 路径与可执行文件：后端需获取当前可执行路径 `os.Executable()` 并使用绝对路径写入启动项。
- 权限：写入用户目录通常无需提升权限，但在受限环境可能失败，需捕获并向前端显示友好错误。
- 单元测试：通过抽象文件系统或模拟 provider 来测试逻辑，而不是直接修改用户文件系统。

前端集成
----
- 在 `Settings.vue` 添加一行“开机自启动”开关，初始值由 `GetAutoStart()` 提供。
- 用户切换开关时调用 `SetAutoStart(enabled)` 并显示操作结果。

安全与隐私
----
- 在 UI 中说明为何需要开机自启以及如何撤销，避免误操作。
