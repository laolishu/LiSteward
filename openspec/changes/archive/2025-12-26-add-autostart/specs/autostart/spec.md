## 新增需求

### 需求：后端提供自启开关接口

后端应提供两个 Wails 可调用的方法支持自启功能：`GetAutoStart() (bool, error)` 与 `SetAutoStart(enabled bool) error`。后端应返回当前用户会话下的自启状态，并在启用时写入用户级启动项，禁用时删除用户级启动项。

#### 场景：查询当前自启状态

- 前置条件：应用已安装并可运行。
- 操作：前端调用 `GetAutoStart()`。
- 期望结果：返回 `true`（若用户级自启已设置）或 `false`（未设置）；出错时返回错误信息供前端显示。

#### 场景：启用开机自启

- 前置条件：`GetAutoStart()` 返回 `false`。
- 操作：前端在设置页中切换开关为"启用"，调用 `SetAutoStart(true)`。
- 期望结果：方法返回 nil，且系统在下次登录/启动时自动运行应用（用户级）。如果写入失败，方法返回错误，前端应显示提示并恢复开关状态。

#### 场景：禁用开机自启

- 前置条件：`GetAutoStart()` 返回 `true`。
- 操作：前端在设置页中切换开关为"禁用"，调用 `SetAutoStart(false)`。
- 期望结果：方法返回 nil，系统不再在下次登录/启动时自动运行应用；若删除失败，返回错误并提示用户。

### 需求：平台特定的自启实现

应用必须提供 Windows、macOS 与 Linux 三个平台的自启实现，每个实现应当使用用户级启动机制，无需管理员权限。

#### 场景：Windows 自启实现

- 实现应在 `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup` 文件夹创建/删除启动脚本或快捷方式。
- 若文件系统受限或快捷方式创建失败，可回退到 HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run 注册表写入。
- 期望结果：启用后，下次登录时应用自动启动；禁用后，应用不再自启。

#### 场景：macOS 自启实现

- 实现应在 `~/Library/LaunchAgents/com.laolishu.listeward.plist` 写入 LaunchAgent 配置。
- 无需立即通过 `launchctl` 加载，允许在下一次登录时生效。
- 期望结果：启用后，下次登录时应用自动启动；禁用后删除 plist 文件。

#### 场景：Linux 自启实现

- 实现应在 `~/.config/autostart/listeward.desktop` 写入 .desktop 文件。
- 确保 `Exec` 字段指向应用的绝对路径。
- 在无 GUI 环境可返回不支持或 noop，并记录日志。
- 期望结果：启用后，下次登录时应用自动启动（若环境支持）；禁用后删除 .desktop 文件。

### 需求：前端设置界面集成

应用的设置页面应包含"开机自启动"开关，允许用户启用/禁用此功能。

#### 场景：显示当前自启状态

- 期望结果：设置页面加载时，开关应初始化为当前用户的自启状态（由 `GetAutoStart()` 提供）。

#### 场景：用户切换自启开关

- 操作：用户在设置页面中点击/切换"开机自启动"开关。
- 期望结果：前端调用 `SetAutoStart(enabled)`，若成功显示确认提示，若失败显示错误信息并恢复开关状态。

### 需求：错误处理与提示

若自启配置写入或读取失败，后端应返回结构化错误，前端应向用户显示可操作建议。

#### 场景：权限不足时的错误提示

- 期望结果：若因权限或文件系统限制导致失败，返回的错误信息应建议用户"权限不足，请以管理员运行或手动添加启动项"。

### 需求：可测试与可验证

应用的自启功能应在每个平台上通过手动测试验证其持久化文件或注册表项的存在与正确性。

#### 场景：验证 Windows 启动项

- 期望结果：启用后，应在 Startup 文件夹找到对应的启动脚本或在注册表 HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run 找到应用条目。

#### 场景：验证 macOS LaunchAgent

- 期望结果：启用后，`~/Library/LaunchAgents/com.laolishu.listeward.plist` 应存在且包含正确的 ProgramArguments 配置。

#### 场景：验证 Linux .desktop 文件

- 期望结果：启用后，`~/.config/autostart/listeward.desktop` 应存在且 `Exec` 字段指向应用可执行路径。
