# 设计文档：关于页面

## 架构概述

关于页面作为设置页面的一个新部分，展示应用信息和赞赏选项。主要设计考虑：

1. **最小化实现**：直接在现有 Settings.vue 中添加内容，无需新增页面
2. **配置驱动**：版本、网址等信息从 config/app.json 读取，易于维护
3. **响应式设计**：适配不同屏幕大小
4. **国际化支持**：所有文本通过 i18n 翻译

## 前端实现

### 1. Settings.vue 组件扩展

在现有设置页面中新增"关于"部分：

```vue
<template>
  <div class="settings-section">
    <h2>{{ $t('settings.about') }}</h2>
    
    <!-- 版本信息 -->
    <div class="about-item">
      <div class="app-info">
        <h3>{{ appConfig.productName }} ({{ appConfig.productNameEn }})</h3>
        <p class="version">v{{ appConfig.version }}</p>
      </div>
    </div>
    
    <!-- 官方网址 -->
    <div class="about-item" v-if="appConfig.website">
      <label>{{ $t('settings.officialWebsite') }}</label>
      <a :href="appConfig.website" target="_blank" class="external-link">
        {{ appConfig.website }}
      </a>
    </div>
    
    <!-- 作者信息 -->
    <div class="about-item" v-if="appConfig.author">
      <label>{{ $t('settings.author') }}</label>
      <p>{{ appConfig.author.name }}</p>
      <p class="email">{{ appConfig.author.email }}</p>
    </div>
    
    <!-- 赞赏二维码 -->
    <div class="sponsor-section" v-if="appConfig.sponsorQrCode">
      <h3>{{ $t('settings.sponsor') }}</h3>
      <p>{{ $t('settings.sponsorTip') }}</p>
      <img :src="appConfig.sponsorQrCode" alt="WeChat Sponsor QR" class="qr-code" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetAppConfig } from '../../wailsjs/go/main/App'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const appConfig = ref(null)

onMounted(async () => {
  try {
    appConfig.value = await GetAppConfig()
  } catch (e) {
    console.error('Failed to get app config:', e)
  }
})
</script>

<style scoped>
.about-item {
  margin-bottom: 16px;
  padding: 12px 0;
}

.app-info {
  margin-bottom: 12px;
}

.app-info h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: #24292e;
}

.version {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.email {
  font-size: 12px;
  color: #999;
  margin: 4px 0 0 0;
}

.external-link {
  color: #0366d6;
  text-decoration: none;
  word-break: break-all;
}

.external-link:hover {
  text-decoration: underline;
}

.sponsor-section {
  margin-top: 24px;
  padding: 16px;
  background-color: #f9f9f9;
  border-radius: 6px;
  text-align: center;
}

.sponsor-section h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
}

.sponsor-section p {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #666;
}

.qr-code {
  width: 200px;
  height: 200px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
```

### 2. i18n 翻译

添加到 `zh-CN.json` 和 `en-US.json`：

**zh-CN.json**：
```json
{
  "settings": {
    "about": "关于",
    "officialWebsite": "官方网站",
    "author": "作者",
    "sponsor": "微信赞赏",
    "sponsorTip": "扫描下方二维码，支持开发者继续维护此项目"
  }
}
```

**en-US.json**：
```json
{
  "settings": {
    "about": "About",
    "officialWebsite": "Official Website",
    "author": "Author",
    "sponsor": "WeChat Sponsor",
    "sponsorTip": "Scan the QR code below to support the developer"
  }
}
```

## 后端实现

### 1. 配置扩展 (config/app.json)

添加新字段：

```json
{
  "version": "0.1.0",
  "website": "https://github.com/example/LiSteward",
  "sponsorQrCode": "./data/assets/wechat-qr.png",
  "author": {
    "name": "lfz",
    "email": "lfzxs@qq.com"
  }
}
```

### 2. 无需后端代码修改

因为 `GetAppConfig()` 已经暴露整个应用配置给前端，前端可以直接读取新字段。

## 资源管理

### 赞赏二维码

二维码图片放置在：
- 位置：`./data/assets/wechat-qr.png`
- 格式：PNG 格式
- 尺寸：200×200 像素（或更大，支持缩放）
- 来源：由用户提供的微信赞赏二维码

## 样式设计

- 遵循现有设置页面样式
- 使用相同的颜色和字体
- 响应式设计：在小屏幕上自适应
- 版本号采用灰色，区别于主文本
- 赞赏部分突出显示（浅灰背景）

## 国际化支持

所有用户可见文本都通过 i18n 系统翻译，支持中文和英文。

## 向后兼容性

- 新字段都是可选的，缺失时可正常降级显示
- 现有代码无任何破坏性变更
- 配置文件可逐步升级
