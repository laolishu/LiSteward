# 设计文档：删除系统托盘和快捷键功能

## 1. 概述

本设计文档描述如何从 LiSteward 应用中安全移除系统托盘（System Tray）和全局快捷键（Global Hotkey）功能。

## 2. 移除理由

### 2.1 稳定性问题
- 系统托盘和快捷键是高耦合的 Windows 原生功能，容易引入系统级错误
- 当前实现存在稳定性问题，导致应用偶尔无响应

### 2.2 维护成本
- 这两个功能涉及 Windows API 的深度集成，增加维护负担
- 跨平台支持成本高（仅 Windows 可用）
- 代码复杂度高，需要特殊的 Windows 消息循环和系统级钩子

### 2.3 功能优先级
- 核心功能（Hosts 管理、缓存模式、解析增强）更重要
- 系统托盘和快捷键是辅助功能，非必需

### 2.4 用户反馈
- 用户更关注主应用功能可靠性而非便利快捷方式

## 3. 移除范围

### 3.1 后端代码移除

#### 3.1.1 删除整个包
- **`internal/hotkey/`** - 全局快捷键实现包
  - `hotkey_windows.go` - Windows 快捷键注册和消息处理
  - `hotkey.go` - 快捷键接口（如存在）

- **`internal/tray/`** - 系统托盘实现包（如存在）
  - 所有托盘相关实现文件

#### 3.1.2 删除 app.go 中的方法
```
# 快捷键方法
- GetHotKeySettings()    // 获取快捷键设置
- SetHotKeySettings()    // 设置快捷键
- RegisterGlobalHotKey() // 注册全局快捷键（如存在）

# 托盘方法
- GetTraySettings()      // 获取托盘设置
- SetTraySettings()      // 设置托盘
- ShowTrayIcon()         // 显示托盘图标（如存在）
- HideTrayIcon()         // 隐藏托盘图标（如存在）
```

#### 3.1.3 配置结构变更
**原始结构**（`config/user_config.go`）：
```go
type UserConfig struct {
    TrayEnabled      bool        // 删除
    HotKeySettings   HotKeySettings  // 删除或简化
    // ... 其他字段
}

type HotKeySettings struct {
    // ... 快捷键配置字段（全部删除）
}
```

**迁移后结构**：
```go
type UserConfig struct {
    // TrayEnabled 和 HotKeySettings 移除
    // ... 保留其他字段
}
```

#### 3.1.4 import 语句更新
- 从 `app.go` 删除对 `internal/hotkey` 的 import
- 从 `app.go` 删除对 `internal/tray` 的 import（如存在）

### 3.2 前端代码移除

#### 3.2.1 Settings.vue 修改
**删除以下 UI 组件**：
- 托盘启用/禁用的 toggle 开关
- 快捷键配置面板
- 快捷键输入框和绑定按钮
- 相关的 help 文本和提示

**删除以下逻辑**：
- `toggleTrayEnabled()` 方法
- `updateHotKeySettings()` 方法
- `onHotKeyChange()` 事件处理器
- `trayEnabled` 响应式数据变量
- `hotKeySettings` 响应式数据变量

**删除以下 import**：
- 与 tray 相关的 Wails 绑定
- 与 hotkey 相关的 Wails 绑定

#### 3.2.2 i18n 翻译清理
**从 `frontend/src/i18n/locales/zh-CN.json` 删除**：
```json
{
  "settings": {
    "trayEnabled": "...",
    "hotKeySettings": "...",
    // 所有托盘和快捷键相关的翻译键
  }
}
```

**从 `frontend/src/i18n/locales/en-US.json` 删除**：
- 相同的翻译键及其英文版本

### 3.3 配置文件清理

#### 3.3.1 user.json 清理
从 `build/bin/data/config/user.json` 移除：
```json
{
  "tray_enabled": false,
  "hotkey_settings": {...}
}
```

## 4. 架构影响

### 4.1 应用启动流程
**前**：
```
启动 → 初始化托盘 → 注册快捷键 → 主窗口
```

**后**：
```
启动 → 主窗口（简化）
```

### 4.2 窗口生命周期
- 移除系统托盘最小化逻辑
- 移除通过快捷键唤醒应用的逻辑
- 简化窗口管理代码

### 4.3 Go-Vue 通信
- 减少 Wails IPC 调用（删除 4 个快捷键和托盘相关方法）
- 简化前端 GetAppInfo 返回值（无需包含 tray/hotkey 状态）

### 4.4 依赖关系
**删除依赖**：
- Windows API 快捷键注册相关的 syscall 调用
- 任何托盘库依赖

## 5. 迁移策略

### 5.1 用户配置迁移
- 现有用户配置中的 `trayEnabled` 和 `hotkey_settings` 字段将被忽略
- 下次保存配置时，这些字段不会被重新创建
- **无需版本迁移**：旧配置与新版本兼容（只是新字段被忽略）

### 5.2 向后兼容性
- 这是**不向后兼容**的更改（功能移除）
- 建议在发布说明中明确说明
- 用户需要手动重新学习（无托盘，无快捷键）

## 6. 实现步骤

### Phase 1：代码清理（高优先级）
1. 删除 `internal/hotkey/` 包
2. 删除 `internal/tray/` 包
3. 从 `config/user_config.go` 删除相关字段
4. 从 `app.go` 删除相关方法和 import
5. 重新生成 Wails 绑定

### Phase 2：前端清理（中优先级）
6. 更新 `Settings.vue`，删除 UI 和逻辑
7. 从 i18n 配置文件删除翻译键
8. 验证前端构建成功

### Phase 3：验证和文档（低优先级）
9. 完整编译测试（`go build`、`wails build`）
10. 功能测试（应用启动、Settings 显示正常）
11. 代码搜索验证（无残留引用）
12. 删除文档（`HOTKEY_USAGE.md`）

## 7. 编译验证检查表

- [ ] `go build ./...` 通过，无错误
- [ ] `cd frontend && npm run build` 通过
- [ ] `wails build` 成功生成可执行文件
- [ ] 没有 unused import 警告
- [ ] 没有 undefined reference 错误

## 8. 功能验证检查表

- [ ] 应用正常启动
- [ ] Settings 页面显示正常（托盘/快捷键选项已移除）
- [ ] 其他 Settings 功能正常（语言、缓存、etc）
- [ ] Hosts 管理功能完全可用
- [ ] 应用正常关闭
- [ ] 无控制台错误或警告

## 9. 代码搜索验证

使用以下命令确保无残留代码：
```bash
# 搜索快捷键残留
grep -r "hotkey" . --include="*.go" --include="*.vue" --include="*.js"
grep -r "hotKey" . --include="*.go" --include="*.vue" --include="*.js"
grep -r "HotKey" . --include="*.go" --include="*.vue" --include="*.js"

# 搜索托盘残留
grep -r "tray" . --include="*.go" --include="*.vue" --include="*.js"
grep -r "Tray" . --include="*.go" --include="*.vue" --include="*.js"
grep -r "sysTray" . --include="*.go" --include="*.vue" --include="*.js"

# 搜索 import 残留
grep -r "internal/hotkey" . --include="*.go"
grep -r "internal/tray" . --include="*.go"
```

## 10. 风险评估

### 10.1 低风险
- 代码删除（非修改）
- 功能完全独立（不影响其他功能）
- 用户配置兼容（旧字段被忽略）

### 10.2 中风险
- 可能遗漏某些引用（需要仔细搜索）
- 前端类型定义需要重新生成

### 10.3 缓解措施
- 完整的代码搜索验证
- 自动化编译测试
- 手动功能测试覆盖

## 11. 相关规范

本变更涉及以下规范的修改：
- **系统集成规范** - 移除系统托盘和快捷键相关需求

## 12. 验收标准

✅ 所有后端代码已删除
✅ 所有前端代码已删除
✅ 配置结构已更新
✅ 编译无错误
✅ 功能测试通过
✅ 无残留代码
✅ 文档已更新
