# 系统集成规范增量：移除系统托盘和快捷键

## 概述
本规范增量描述从 LiSteward 中移除系统托盘（System Tray）和全局快捷键（Global Hotkey）功能所需的规范修改。

## 移除需求

### 需求：移除系统托盘功能

#### 场景：系统托盘图标显示被移除
应用不再在系统托盘中显示图标。用户无法通过系统托盘最小化应用。

**验收标准**：
- 应用启动时不初始化托盘管理器
- internal/tray/ 包的代码被完全移除
- app.go 中不包含 trayMgr 相关代码
- 应用可以正常启动和关闭

#### 场景：托盘设置 UI 被移除
Settings.vue 中不再包含托盘启用/禁用的设置选项。

**验收标准**：
- Settings.vue 中删除了所有托盘相关的 UI 组件
- 配置文件（UserConfig）中删除 TrayEnabled 字段
- i18n 文件中删除 settings.trayEnabled 和 settings.trayEnabledDesc 等翻译键
- Settings 页面加载正常，不显示残留 UI

### 需求：移除全局快捷键功能

#### 场景：全局快捷键唤醒功能被移除
应用不再支持通过全局快捷键唤醒。用户无法使用快捷键（如 Ctrl+Alt+H）唤醒应用。

**验收标准**：
- 应用启动时不初始化快捷键管理器
- internal/hotkey/ 包的代码被完全移除
- app.go 中不包含 hotKeyMgr 相关代码
- app.go 中不包含快捷键相关的方法（如 GetHotKeySettings、SetHotKeySettings）
- 应用可以正常启动和关闭

#### 场景：快捷键配置 UI 被移除
Settings.vue 中不再包含快捷键配置的设置面板。

**验收标准**：
- Settings.vue 中删除了所有快捷键相关的 UI 组件和配置选项
- 配置文件（UserConfig）中删除 HotKeySettings 字段
- i18n 文件中删除 settings.hotKeySettings 等快捷键相关的翻译键
- Settings 页面加载正常，不显示残留的快捷键配置

## 保留的需求

以下系统集成需求保留：
- ✅ 应用作为独立窗口正常启动
- ✅ 应用窗口可以最小化到任务栏
- ✅ 用户可以通过 Settings 调整应用设置
- ✅ 应用支持国际化

## 规范更新指导

当主规范库（`openspec/specs/system-integration/spec.md`）存在时，应：
1. 删除所有提及"系统托盘"相关的需求和验收标准
2. 删除所有提及"全局快捷键"相关的需求和验收标准
3. 删除"Windows 特定功能"或"系统集成"章节中的这些内容
4. 更新"支持的功能"或类似部分，移除这两项

## 实现影响

### 后端 API 移除
```
移除的 Go 函数（从 app.go）：
- GetHotKeySettings() -> string
- SetHotKeySettings(settings string)
- GetTraySettings() -> bool
- SetTraySettings(enabled bool)
```

### 前端 i18n 移除
```
移除的翻译键（i18n locale files）：
- settings.trayEnabled
- settings.trayEnabledDesc
- settings.hotKeySettings
- settings.hotKeySettingsDesc
- settings.hotKeyPlaceholder
- settings.hotKeyHelp
```

### 配置结构移除
```go
// 从 UserConfig 结构移除：
type UserConfig struct {
    TrayEnabled    bool           // 删除
    HotKeySettings HotKeySettings // 删除
}
```

## 验收标准

✅ **必须**：所有快捷键代码从 `internal/hotkey/` 删除
✅ **必须**：所有托盘代码从 `internal/tray/` 删除（如存在）
✅ **必须**：从 Settings.vue 删除所有托盘/快捷键 UI 元素
✅ **必须**：从 i18n 配置删除所有相关翻译
✅ **必须**：应用编译无错误
✅ **必须**：Settings 页面显示正常，无残留 UI 元素
✅ **必须**：无代码中残留的快捷键/托盘引用

## 相关变更
- **remove-tray-hotkey 变更** - 主变更提案，包含代码删除计划
