<!--
 * @Descripttion: 
 * @version: 
 * @Author: lfzxs@qq.com
 * @Date: 2025-12-26 18:15:18
 * @LastEditors: lfzxs@qq.com
 * @LastEditTime: 2025-12-26 18:50:11
-->
# 变更：删除系统托盘和快捷键功能

## 为什么

当前应用中的系统托盘和全局快捷键功能处于未完成状态（add-system-tray 变更仅完成 24/34 任务），存在以下问题：

- **功能不稳定**：快捷键系统存在可靠性问题，用户反馈功能经常不响应
- **维护成本高**：系统托盘和快捷键涉及复杂的 Windows API 调用，需要特殊权限和错误处理
- **使用率低**：这些功能对核心 Hosts 管理业务不是必需的
- **代码负担**：相关代码分散在多个文件中，增加了应用复杂性
- **用户需求**：用户的主要需求是稳定可靠的 Hosts 管理，而非系统集成功能

## 变更内容

### 删除的功能

1. **系统托盘功能**
   - 删除将应用最小化到系统托盘的功能
   - 删除托盘菜单（显示、隐藏、退出）
   - 删除相关 UI 控件和设置选项

2. **全局快捷键功能**
   - 删除全局热键注册和处理逻辑
   - 删除快捷键设置界面
   - 删除相关配置存储

### 影响范围

#### 后端代码删除
- `internal/tray/` - 整个托盘管理包
- `internal/hotkey/` - 整个快捷键管理包
- `app.go` 中的托盘和快捷键相关方法：
  - `GetTraySettings()`
  - `SetTraySettings()`
  - `GetHotKeySettings()`
  - `SetHotKeySettings()`

#### 前端代码删除
- `frontend/src/views/Settings.vue` 中的托盘和快捷键设置面板
- 相关的 i18n 翻译字符串
- Wails 生成的绑定类型定义

#### 配置文件修改
- `config/user_config.go` 中删除 `TrayEnabled` 和 `HotKeyEnabled` 字段
- `build/bin/data/config/user.json` 中删除相关配置字段
- `HOTKEY_USAGE.md` 文档删除

### 向后兼容性

- 现有用户配置文件中的 `tray_enabled` 和`hotkey_*` 字段将被忽略
- 应用启动时自动清理相关配置，不影响现有的 Hosts 管理功能
- 仅影响设置面板 UI，核心 Hosts 管理功能保持完整

## 相关变更

- **撤销**：add-system-tray 变更（24/34 完成）
- **影响**：无其他变更直接依赖这些功能

## 验收标准

- ✅ 所有托盘相关代码已删除
- ✅ 所有快捷键相关代码已删除  
- ✅ 设置界面不再显示相关选项
- ✅ 应用正常编译和运行
- ✅ Hosts 管理核心功能保持完整
- ✅ 无遗留的 import 或引用
- ✅ 相关测试已清理或更新
