# 任务清单：删除系统托盘和快捷键功�?

## 后端清理

### 1. 删除快捷键包
- [x] 1.1 删除 `internal/hotkey/` 目录及所有文�?
- [x] 1.2 删除 `app.go` 中的快捷键相关方法：
  - [ ] 1.2.1 `GetHotKeySettings()`
  - [ ] 1.2.2 `SetHotKeySettings()`

### 2. 删除托盘�?
- [x] 2.1 删除 `internal/tray/` 目录及所有文�?
- [x] 2.2 删除 `app.go` 中的托盘相关方法�?
  - [ ] 2.2.1 `GetTraySettings()`
  - [ ] 2.2.2 `SetTraySettings()`

### 3. 修改配置结构
- [x] 3.1 修改 `config/user_config.go`
  - [ ] 3.1.1 删除 `TrayEnabled` 字段
  - [ ] 3.1.2 删除 `HotKeySettings` 字段或整�?HotKey 相关字段
  - [ ] 3.1.3 更新配置初始化逻辑

### 4. 删除相关 import
- [x] 4.1 检�?`app.go` 中删除对 `internal/tray` �?`internal/hotkey` �?import
- [x] 4.2 检查其他文件中是否有相�?import 并删�?

### 5. 清理配置文件
- [x] 5.1 清理 `build/bin/data/config/user.json`，删�?`tray_enabled` 和相关热键配�?

## 前端清理

### 6. 修改 Settings 页面
- [x] 6.1 �?`frontend/src/views/Settings.vue` 中删除：
  - [ ] 6.1.1 托盘设置区块（UI 元素和绑定）
  - [ ] 6.1.2 快捷键设置区块（UI 元素和绑定）
  - [ ] 6.1.3 相关�?import 语句（GetTraySettings, SetTraySettings, GetHotKeySettings, SetHotKeySettings�?
  - [ ] 6.1.4 相关的数据变量（trayEnabled, hotKeySettings 等）
  - [ ] 6.1.5 相关的方法（toggleTrayEnabled, updateHotKeySettings 等）

### 7. 删除 i18n 翻译
- [x] 7.1 �?`frontend/src/i18n/locales/zh-CN.json` 删除�?
  - [ ] 7.1.1 `settings.trayEnabled`
  - [ ] 7.1.2 `settings.hotKeySettings`
  - [ ] 7.1.3 其他相关的托盘和快捷键翻译键
- [x] 7.2 �?`frontend/src/i18n/locales/en-US.json` 删除相同的翻译键

### 8. 重新生成 Wails 绑定
- [x] 8.1 运行 `wails generate bindings` 更新前端类型定义

## 文档清理

### 9. 删除相关文档
- [x] 9.1 删除 `HOTKEY_USAGE.md`
- [x] 9.2 删除或清�?`openspec/changes/add-system-tray/` 目录下的文件（作为参考留存，或完全删除）

## 测试和验�?

### 10. 编译验证
- [x] 10.1 运行 `go build` 确保没有编译错误
- [x] 10.2 运行 `cd frontend && npm run build` 或相关构建命令确保前端无错误
- [x] 10.3 运行 `wails build` 完整构建应用

### 11. 功能验证
- [x] 11.1 启动应用，验�?Settings 页面不显示托�?快捷键选项
- [x] 11.2 验证 Hosts 管理核心功能完全可用
- [x] 11.3 验证应用正常启动和关�?
- [x] 11.4 验证没有控制台错误或警告

### 12. 代码搜索验证
- [x] 12.1 搜索 `internal/tray` 确保无残留引�?
- [x] 12.2 搜索 `internal/hotkey` 确保无残留引�?
- [x] 12.3 搜索 `tray` �?`hotkey` 确保无遗留代�?

## 优先�?

| 优先�?| 任务 | 原因 |
|------|------|------|
| P1 | 1-4, 6, 10 | 核心代码清理和编译验�?|
| P2 | 5, 7, 8, 11 | 配置和前端清�?|
| P3 | 9, 12 | 文档清理和最终验�?|

## 依赖关系

```
1 (删除快捷�? 并行 2 (删除托盘)
  �?
3 (修改配置)
  �?
4 (删除 import)
  �?
5 (清理配置文件)
  �?
6 (修改 Settings)
  �?
7 (删除 i18n)
  �?
8 (重新生成绑定)
  �?
9 (删除文档)
  �?
10 (编译验证)
  �?
11 (功能验证)
  �?
12 (代码搜索验证)
```

## 工作量估�?

- 后端清理�?-2 小时
- 前端清理�? 小时
- 测试验证�? 小时
- **总计**�?-4 小时
