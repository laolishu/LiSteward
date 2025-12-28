# 设计文档：国际化支持

## 架构概述

本变更在现有应用基础上添加完整的国际化体系，涵盖：

1. **后端语言管理**：系统语言检测、用户偏好存储
2. **前端 I18n 框架**：翻译资源、语言切换逻辑
3. **设置界面**：用户语言选择

### 系统交互流程

```
应用启动
  ↓
读取保存的语言偏好
  ├─ 如果存在 → 使用该语言
  └─ 如果不存在 → 检测系统语言 → 使用系统语言或默认中文
  ↓
前端初始化 i18n，加载对应语言资源
  ↓
渲染界面（所有文本都通过 i18n 获取）
  ↓
用户在设置中切换语言
  ↓
实时更新界面 + 保存偏好
```

## 前端 I18n 实现

### 库选择：vue-i18n

**原因：**
- Vue 生态标准方案
- 支持 Vue 3 Composition API
- 功能完整（动态加载、插值、复数等）
- 相对轻量级

### 翻译资源结构

```
frontend/src/i18n/locales/
├── zh-CN.json
└── en-US.json
```

#### 翻译文件格式

**zh-CN.json (中文)**
```json
{
  "app": {
    "name": "李管家"
  },
  "menu": {
    "hosts": "Hosts 管理",
    "settings": "设置"
  },
  "hosts": {
    "title": "Hosts 管理",
    "search": "搜索 IP 或域名...",
    "addEntry": "+ 添加",
    "domain": "域名",
    "ip": "IP 地址",
    "alternateIp": "备用 IP（可选）",
    "description": "描述（可选）",
    "status": "状态",
    "operations": "操作",
    "swap": "交换",
    "edit": "编辑",
    "delete": "删除",
    "enabled": "已启用",
    "disabled": "已禁用",
    "empty": "暂无 Hosts 条目",
    "saving": "保存中...",
    "saved": "✓ 所有更改已保存",
    "unsaved": "● 未保存的更改"
  },
  "settings": {
    "title": "设置",
    "language": "语言",
    "languageChange": "界面语言已切换"
  },
  "dialog": {
    "cancel": "取消",
    "save": "保存",
    "add": "添加",
    "edit": "编辑",
    "delete": "删除"
  },
  "message": {
    "success": "操作成功",
    "error": "操作失败",
    "confirm": "确定要删除吗？"
  }
}
```

**en-US.json (英文)**
```json
{
  "app": {
    "name": "LiSteward"
  },
  "menu": {
    "hosts": "Hosts Manager",
    "settings": "Settings"
  },
  "hosts": {
    "title": "Hosts Manager",
    "search": "Search IP or domain...",
    "addEntry": "+ Add",
    "domain": "Domain",
    "ip": "IP Address",
    "alternateIp": "Alternate IP (Optional)",
    "description": "Description (Optional)",
    "status": "Status",
    "operations": "Operations",
    "swap": "Swap",
    "edit": "Edit",
    "delete": "Delete",
    "enabled": "Enabled",
    "disabled": "Disabled",
    "empty": "No Hosts entries",
    "saving": "Saving...",
    "saved": "✓ All changes saved",
    "unsaved": "● Unsaved changes"
  },
  "settings": {
    "title": "Settings",
    "language": "Language",
    "languageChange": "Interface language changed"
  },
  "dialog": {
    "cancel": "Cancel",
    "save": "Save",
    "add": "Add",
    "edit": "Edit",
    "delete": "Delete"
  },
  "message": {
    "success": "Operation successful",
    "error": "Operation failed",
    "confirm": "Are you sure you want to delete?"
  }
}
```

### I18n 初始化 (frontend/src/i18n/index.js)

```javascript
import { createI18n } from 'vue-i18n'
import zh from './locales/zh-CN.json'
import en from './locales/en-US.json'

export function initI18n() {
  // 从后端获取用户语言偏好，默认 zh-CN
  const locale = localStorage.getItem('app-language') || 'zh-CN'
  
  const i18n = createI18n({
    legacy: false,          // 使用 Composition API 模式
    locale: locale,
    fallbackLocale: 'zh-CN',
    messages: {
      'zh-CN': zh,
      'en-US': en
    }
  })
  
  return i18n
}

export function setLanguage(locale) {
  localStorage.setItem('app-language', locale)
  // 需要同时调用后端保存
}

export function getLanguage() {
  return localStorage.getItem('app-language') || 'zh-CN'
}
```

### 组件使用方式

在 Composition API 中使用：

```vue
<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const handleSwap = () => {
  showMessage(t('hosts.swap') + '成功')
}
</script>

<template>
  <button @click="handleSwap">{{ t('hosts.swap') }}</button>
  <input :placeholder="t('hosts.search')" />
</template>
```

## 后端实现

### 语言检测

在 Windows 系统上通过 Windows API 获取系统语言：

```go
// language/detector.go

package language

import (
  "syscall"
)

// GetSystemLanguage 获取 Windows 系统语言
// 返回格式：zh-CN, en-US 等
func GetSystemLanguage() (string, error) {
  // 使用 Windows API 获取用户默认语言
  // LOCALE_SNAME 返回 zh-CN, en-US 格式
  // 如果失败，返回默认值 zh-CN
}
```

### 配置存储

扩展 config 包：

```go
// config/config.go

type AppConfig struct {
    ProductName    string // "李管家"
    ProductNameEn  string // "LiSteward"
    Language       string // 用户选择的语言: zh-CN, en-US
}

// SaveLanguage 保存用户语言偏好
func (c *AppConfig) SaveLanguage(lang string) error {
    c.Language = lang
    return c.Save()
}

// GetLanguage 获取用户语言偏好
// 优先级：用户设置 > 系统语言 > 默认(zh-CN)
func (c *AppConfig) GetLanguage() string {
    if c.Language != "" {
        return c.Language
    }
    
    // 检测系统语言
    if syslang, err := language.GetSystemLanguage(); err == nil {
        return syslang
    }
    
    return "zh-CN"  // 默认中文
}
```

### 暴露的 Wails 绑定

在 app.go 中添加：

```go
// GetSystemLanguage 获取系统语言
func (a *App) GetSystemLanguage() (string, error) {
    return language.GetSystemLanguage()
}

// GetUserLanguagePreference 获取用户语言偏好
func (a *App) GetUserLanguagePreference() string {
    if a.cfg != nil {
        return a.cfg.GetLanguage()
    }
    return "zh-CN"
}

// SetUserLanguagePreference 设置用户语言偏好
func (a *App) SetUserLanguagePreference(lang string) error {
    if a.cfg == nil {
        return fmt.Errorf("config not initialized")
    }
    return a.cfg.SaveLanguage(lang)
}
```

## 设置模块设计

### Settings 视图

创建 `frontend/src/views/Settings.vue`：

```vue
<template>
  <div class="settings-container">
    <div class="settings-section">
      <h2>{{ $t('settings.title') }}</h2>
      
      <div class="settings-item">
        <label>{{ $t('settings.language') }}</label>
        <select v-model="selectedLanguage" @change="changeLanguage">
          <option value="zh-CN">中文 (Chinese)</option>
          <option value="en-US">English</option>
        </select>
      </div>
      
      <div v-if="languageChanged" class="message-success">
        {{ $t('settings.languageChange') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { GetUserLanguagePreference, SetUserLanguagePreference } from '../../wailsjs/go/main/App'

const { locale } = useI18n()
const selectedLanguage = ref(locale.value)
const languageChanged = ref(false)

const changeLanguage = async () => {
  // 切换 i18n 语言
  locale.value = selectedLanguage.value
  
  // 保存到后端
  await SetUserLanguagePreference(selectedLanguage.value)
  
  // 保存到 localStorage
  localStorage.setItem('app-language', selectedLanguage.value)
  
  // 显示反馈
  languageChanged.value = true
  setTimeout(() => {
    languageChanged.value = false
  }, 2000)
}
</script>
```

## 迁移计划

### 阶段 1：i18n-core 规范
- 集成 vue-i18n
- 创建翻译资源文件
- 后端支持语言检测
- 应用启动时初始化语言

### 阶段 2：i18n-settings 规范
- 创建设置模块
- 实现语言选择界面
- 实现实时切换逻辑
- 完整测试和文档

## 风险和考虑

1. **翻译完整性**：需要确保所有界面文本都有对应翻译
2. **性能**：i18n 库需要适配大型应用，但当前应用规模可以忽略
3. **维护性**：后续新功能需要同时添加中英文翻译
4. **工具链**：考虑后续是否需要翻译管理工具

## 文档更新

- 更新 ARCHITECTURE.md 的国际化部分
- 为开发者提供 i18n 集成指南
- 记录翻译资源维护流程
