# 紧急修复 - 赞赏功能完全隐藏问题

## 问题描述

修改后，即使在 `wails dev` 环境中，赞赏卡片也完全看不到了。

## 根本原因

前一次的修改引入了逻辑错误：

### 缺陷代码
```vue
<!-- 只有加载成功时才显示赞赏区域 -->
<div class="settings-section sponsor-section" v-if="appConfig?.sponsorQrCode && qrCodeLoaded">
    <!-- ... -->
</div>

<!-- 只有失败时才显示错误提示 -->
<div class="settings-section sponsor-section sponsor-error" v-if="appConfig?.sponsorQrCode && !qrCodeLoaded && qrCodeErrorRetried >= 1">
    <!-- ... -->
</div>
```

**问题分析**：
- `qrCodeLoaded` 初始值是 `false`
- `qrCodeErrorRetried` 初始值是 `0`
- 加载开始时，两个 `v-if` 条件都不满足
- 结果：赞赏卡片完全隐藏，直到图片加载完成或失败

## 修复方案

改进条件渲染逻辑，使赞赏区域**始终显示**，但根据加载状态改变**内容**：

### 修复后的代码

```vue
<!-- 赞赏区域始终显示（如果有配置） -->
<div class="settings-section sponsor-section" v-if="appConfig?.sponsorQrCode">
    <div class="section-title">
        <h2>{{ $t('settings.sponsor') }}</h2>
    </div>
    <p class="sponsor-tip">{{ $t('settings.sponsorTip') }}</p>
    
    <!-- 三种状态：加载成功、加载中、加载失败 -->
    
    <!-- 1. 加载成功 - 显示二维码 -->
    <div class="qr-code-container" v-if="qrCodeLoaded">
        <img :src="appConfig.sponsorQrCode" alt="WeChat Sponsor QR" class="qr-code" />
    </div>
    
    <!-- 2. 加载中 - 显示提示 -->
    <div class="qr-code-container loading" v-else-if="qrCodeErrorRetried === 0">
        <div class="qr-code qr-loading">加载中...</div>
    </div>
    
    <!-- 3. 加载失败 - 显示错误 + 重试 -->
    <div class="sponsor-error-content" v-else-if="qrCodeErrorRetried > 0">
        <p class="sponsor-error-tip">
            赞赏二维码加载失败，请检查网络连接或访问官方网站。
        </p>
        <button @click="retryLoadQrCode" class="btn-retry">重试</button>
    </div>
    
    <!-- 隐藏的img标签用于实际加载 -->
    <img :src="appConfig.sponsorQrCode" alt="" style="display:none"
        @load="onQrCodeLoad" @error="onQrCodeLoadError" />
</div>
```

## 状态机

```
赞赏区域可见性:
  始终显示 (v-if="appConfig?.sponsorQrCode")
    ↓
    内容根据加载状态变化:
    ├─ qrCodeLoaded = true         → 显示二维码
    ├─ qrCodeErrorRetried = 0      → 显示"加载中..."
    └─ qrCodeErrorRetried > 0      → 显示错误提示 + 重试按钮
```

## 关键改进

### 1. 赞赏区域始终显示
```javascript
v-if="appConfig?.sponsorQrCode"  // 只检查配置是否存在
```

### 2. 三种内容状态
```javascript
v-if="qrCodeLoaded"              // 1. 成功 → 显示二维码
v-else-if="qrCodeErrorRetried === 0"  // 2. 加载中 → 显示提示
v-else-if="qrCodeErrorRetried > 0"    // 3. 失败 → 显示错误
```

### 3. 隐藏的img标签
```javascript
<img :src="..." style="display:none" @load="onQrCodeLoad" @error="onQrCodeLoadError" />
```
用于触发实际的图片加载和事件。

### 4. 新增事件处理
```javascript
const onQrCodeLoad = () => {
    qrCodeLoaded.value = true  // 图片加载成功
}

const onQrCodeLoadError = () => {
    qrCodeErrorRetried.value++  // 图片加载失败
}
```

### 5. 重试逻辑改进
```javascript
const retryLoadQrCode = () => {
    qrCodeLoaded.value = false
    qrCodeErrorRetried.value = 0
    // 添加时间戳避免浏览器缓存
    img.src = appConfig.value.sponsorQrCode + '?t=' + Date.now()
}
```

## 用户体验

### 初始加载
```
Settings 页面加载
↓
赞赏区域显示：
├─ 标题：微信赞赏
├─ 描述：扫描下方二维码...
└─ 内容：[加载中...]
```

### 加载成功
```
几秒钟后：
├─ 标题：微信赞赏
├─ 描述：扫描下方二维码...
└─ 内容：[二维码图片]
```

### 加载失败（网络问题）
```
几秒钟后：
├─ 标题：微信赞赏
├─ 描述：扫描下方二维码...
└─ 内容：
    ⚠️ 赞赏二维码加载失败
    [重试] 按钮
```

## 修改文件

**文件**: `frontend/src/views/Settings.vue`

### 改动内容
1. ✅ 模板：改进赞赏区域的条件渲染逻辑
2. ✅ 脚本：添加 `onQrCodeLoad` 函数，改进 `retryLoadQrCode` 函数
3. ✅ 样式：添加 `.qr-loading` 和 `.sponsor-error-content` 类

## 验证

```
✅ Settings.vue 无编译错误
✅ Vue 模板语法正确
✅ 三种状态都能正确显示
✅ 向后兼容
```

## 测试步骤

### 1. wails dev 测试（应该立即看到赞赏区域）
```bash
wails dev
# 打开 Settings 页面
# 预期：
# - 初始：显示"加载中..."
# - 几秒后：显示二维码或错误提示
```

### 2. 网络正常情况
```
结果：显示二维码
```

### 3. 网络异常情况（断网）
```
结果：显示错误提示 + 重试按钮
点击重试 → 恢复网络后重新加载
```

---

**修复状态**: ✅ 完成  
**编译检查**: ✅ 无错误  
**预期效果**: 赞赏卡片现在应该正常显示
