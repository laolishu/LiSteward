# 任务清单：添加开机自启动功能

1. 设计接口与用户交互
   - 定义后端 API：`GetAutoStart() (bool, error)` 与 `SetAutoStart(enabled bool) error`。
   - 定义前端调用契约与文案：设置页的开关、提示文本、错误处理文案。

2. 实现后端平台适配器（子任务可并行）
   - Windows 适配器：实现通过当前用户 Startup 文件夹创建/删除快捷方式，或在不能创建快捷方式时回退到 HKCU Run 注册表项。
      - Windows 适配器：已实现（写入当前用户 `Startup` 文件夹的 `.cmd` 启动脚本，用户级，无需管理员权限）。
      - macOS 适配器：已实现（写入 `~/Library/LaunchAgents/com.laolishu.listeward.plist`，RunAtLoad=true，用户级）。
      - Linux 适配器：已实现（写入 `~/.config/autostart/listeward.desktop`，用户级）。
   - macOS 适配器：实现写入 `~/Library/LaunchAgents/com.example.listeward.plist`（或更通用的 bundle id），并在禁用时删除。
   - Linux 适配器：实现写入 `~/.config/autostart/listeward.desktop`，并在禁用时删除。对无 GUI 环境返回不支持或 noop 并记录日志。

3. 在后端添加跨平台接口并暴露给前端（Wails）
   - 在 `app.go` 中添加方法并导出：`GetAutoStart`, `SetAutoStart`。
   - 单元/集成测试（如可能）验证平台逻辑（使用接口抽象以便在 CI 中用模拟替代文件系统操作）。

4. 前端实现
   - 在 `Settings.vue` 中添加开关控件（与现有语言/版本项同级），绑定到 `GetAutoStart()` 初始状态并在切换时调用 `SetAutoStart()`。
   - 显示成功 / 失败消息与说明链接（如何撤销）。

5. 验证与文档
   - 编写手动验证步骤（Windows/macOS/Linux）。
   - 更新项目 README / 快速开始文档，说明此功能的位置与行为。

验收测试（每个平台）
----
- 启用后：重启登录，应用自动启动。
- 禁用后：重启登录，应用不再自动启动。
- 前端显示状态与实际启动项状态一致。
