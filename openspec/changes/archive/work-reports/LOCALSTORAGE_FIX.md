# 🔧 localStorage 缓存清除指南

## 问题说明

应用启动时，从 localStorage 中加载之前保存的侧栏宽度（300px），导致新的默认宽度（250px）被忽略。

## 解决方案

### 方案 1: 浏览器开发者工具清除（最简单）

**步骤**:
1. 按 `F12` 打开开发者工具
2. 选择 "Application" 或 "Storage" 标签
3. 展开 "Local Storage" 
4. 找到当前应用的域名（如 `localhost:5173`）
5. 找到 `sidebar-width` 项
6. **删除** 这一项
7. 刷新页面 `F5`

**结果**: 应用会使用新的默认宽度 250px

---

### 方案 2: 浏览器控制台清除（快速）

**步骤**:
1. 按 `F12` 打开开发者工具
2. 选择 "Console" 标签
3. 复制并执行以下代码：

```javascript
// 删除旧的 sidebar-width
localStorage.removeItem('sidebar-width')

// 刷新页面
window.location.reload()
```

**说明**: 这会立即删除缓存的宽度并刷新页面

---

### 方案 3: 一键清除所有 localStorage（彻底）

如果想完全清除应用的所有本地存储：

```javascript
// 清除所有 localStorage
localStorage.clear()

// 刷新页面
window.location.reload()
```

**warning**: 这会删除所有保存的配置，不仅仅是侧栏宽度

---

## 代码修复

### 修改的部分

**文件**: `frontend/src/layouts/MainLayout.vue`

**改动**: onMounted 钩子现在验证保存的宽度是否在允许范围内 (180-320px)

```typescript
// ❌ 之前 (问题版本)
onMounted(() => {
    const savedWidth = localStorage.getItem('sidebar-width')
    if (savedWidth) {
        sidebarWidth.value = parseInt(savedWidth, 10)  // 无条件使用旧值
    }
})

// ✅ 之后 (修复版本)
onMounted(() => {
    const savedWidth = localStorage.getItem('sidebar-width')
    if (savedWidth) {
        const width = parseInt(savedWidth, 10)
        // 确保宽度在允许的范围内
        if (width >= minSidebarWidth && width <= maxSidebarWidth) {
            sidebarWidth.value = width
        } else {
            // 如果超出范围，重置为默认值并保存
            localStorage.setItem('sidebar-width', '250')
        }
    }
})
```

**优势**:
- ✅ 验证保存的值是否有效
- ✅ 自动修复无效数据
- ✅ 防止将来出现类似问题

---

## 完整步骤

### 第一次更新后（现在）

1. **清除 localStorage**:
   - 打开开发者工具 (F12)
   - 找到 Application → Local Storage
   - 删除 `sidebar-width` 项

2. **刷新页面**:
   - 按 F5 或点击刷新按钮

3. **验证**:
   - 侧栏应该现在显示宽度 ~250px（之前是 300px）
   - 拖动分隔符调整宽度，关闭再打开应该能保存新的宽度

---

## 自动化方案（后续开发）

如果想在应用中添加自动清除旧值的逻辑，可以：

```typescript
// 在 main.ts 或 App.vue 中
// 检查是否需要迁移 localStorage 数据
const validateStoredWidth = () => {
    const saved = localStorage.getItem('sidebar-width')
    if (saved) {
        const width = parseInt(saved, 10)
        if (width === 300) {  // 如果是旧的默认值
            // 迁移到新的默认值
            localStorage.setItem('sidebar-width', '250')
        }
    }
}

// 应用启动时调用
validateStoredWidth()
```

---

## 长期解决方案

为了避免将来出现类似问题，可以：

1. **使用版本控制**: 在 localStorage 中存储数据版本
   ```typescript
   const saveWidth = (width: number) => {
       localStorage.setItem('sidebar-config-v2', JSON.stringify({
           version: 2,
           width
       }))
   }
   ```

2. **使用 Pinia/Vuex**: 不依赖 localStorage 的原始值，统一管理状态

3. **添加默认值验证**: 始终验证加载的值是否在允许范围内

---

## 快速参考

| 操作 | 步骤 |
|------|------|
| **删除单个值** | DevTools → Application → Local Storage → 右键删除 `sidebar-width` |
| **清除所有值** | 控制台执行 `localStorage.clear()` |
| **检查当前值** | 控制台执行 `localStorage.getItem('sidebar-width')` |
| **手动设置新默认值** | 控制台执行 `localStorage.setItem('sidebar-width', '250')` |

---

**修复时间**: 2025-12-23  
**涉及文件**: MainLayout.vue  
**修复类型**: localStorage 验证逻辑
