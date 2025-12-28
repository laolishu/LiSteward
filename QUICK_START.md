# 🚀 快速参考指南

## 文件导航

### 核心文件位置
```
应用入口
└─ frontend/src/App.vue                  (50 行)  🟢 简单

主框架容器
└─ frontend/src/layouts/MainLayout.vue   (300 行) 🔵 核心

全局导航
└─ frontend/src/components/Sidebar.vue   (260 行) 🟢 完成

业务模块 (Hosts管理)
└─ frontend/src/views/HostsManager.vue   (1200 行) 🟡 复杂

模块目录 (待开发)
└─ frontend/src/views/modules/
   ├─ Dashboard.vue
   ├─ EnvironmentVariables.vue
   ├─ Registry.vue
   └─ Settings.vue
```

---

## 🎯 常见任务

### ✏️ 添加新菜单项

1. **编辑 Sidebar.vue**
```vue
// 在 menuItems 数组中添加
const menuItems = [
    { id: 'dashboard', icon: '📊', label: '控制台' },
    { id: 'hosts', icon: '🌐', label: 'Hosts 管理' },
    // ... 添加新项
    { id: 'newmodule', icon: '🆕', label: '新模块' },
]
```

2. **编辑 App.vue**
```vue
<NewModule v-if="currentModule === 'newmodule'" />
```

3. **导入新模块**
```javascript
import NewModule from './views/modules/NewModule.vue'
```

---

### 📦 创建新模块

1. **在 views/modules/ 下创建文件**
```bash
touch frontend/src/views/modules/NewModule.vue
```

2. **基础模板**
```vue
<template>
    <div class="new-module">
        <h2>新模块</h2>
        <!-- 模块内容 -->
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// 你的状态和逻辑
const data = ref([])

// 你的方法
const loadData = async () => {
    // 调用后端 Wails 方法
}
</script>

<style scoped>
.new-module {
    padding: 24px;
    background: #f6f8fa;
    height: 100%;
    overflow: auto;
}
</style>
```

3. **在 App.vue 中注册**
```vue
<script setup>
import NewModule from './views/modules/NewModule.vue'
</script>

<template>
    <MainLayout>
        <NewModule v-if="currentModule === 'newmodule'" />
    </MainLayout>
</template>
```

---

### 🎨 自定义侧栏宽度

**编辑 MainLayout.vue**
```typescript
// 第 20-21 行
const minSidebarWidth = 150  // 改为你的最小值
const maxSidebarWidth = 600  // 改为你的最大值
```

---

### 📏 调整窗口最小尺寸

**编辑 wails.json**
```json
"windows": [{
    "minWidth": 1400,      // 改为你的最小宽度
    "minHeight": 900       // 改为你的最小高度
}]
```

---

### 💾 保存侧栏宽度

主框架已自动处理 localStorage，无需额外配置。

**手动查看：**
```javascript
// 在浏览器 DevTools 中
localStorage.getItem('sidebar-width')
```

---

## 🎨 样式规范

### 颜色变量（全局）
```css
/* 背景 */
#F6F8FA  → 浅灰背景
#FFFFFF  → 白色

/* 文本 */
#24292E  → 深灰（主文本）
#57606A  → 中灰（次文本）
#6A737D  → 浅灰（辅助文本）

/* 功能色 */
#34D399  → 翠绿（强调）
#1B7C34  → 绿色（成功）
#0969DA  → 蓝色（链接）
#CB2431  → 红色（危险）

/* 边框 */
#D0D7DE  → 标准边框
#E1E4E8  → 浅边框
```

### 尺寸规范
```
侧栏默认宽度     300px
侧栏最小宽度     200px
侧栏最大宽度     400px
分隔符宽度       4px
顶部栏高度       50px
内容区内边距     16px
组件间隙        8-16px
```

### 圆角规范
```css
按钮/卡片    4-6px
弹窗/模态框  8px
```

---

## 📊 组件交互

### 侧栏折叠/展开
```
用户点击顶部展开按钮
  ↓
MainLayout.toggleSidebar()
  ├─ sidebarWidth = 60px (折叠)
  └─ sidebarWidth = 300px (展开)
```

### 侧栏宽度拖动
```
用户在分隔符按下鼠标
  ↓
MainLayout.startDragging()
  ├─ document.addEventListener('mousemove')
  ├─ 计算 delta = currentX - startX
  ├─ newWidth = startWidth + delta
  └─ 限制在 [200, 400] 范围

释放鼠标
  ↓
localStorage.setItem('sidebar-width', newWidth)
```

### 模块切换
```
用户点击菜单项
  ↓
Sidebar.emit('navigate', moduleName)
  ↓
App.handleNavigate(moduleName)
  └─ currentModule = moduleName
  ↓
MainLayout 自动显示对应模块组件
```

---

## 🔍 调试技巧

### 查看侧栏宽度
```javascript
// 浏览器控制台
console.log(localStorage.getItem('sidebar-width'))
```

### 强制重置侧栏宽度
```javascript
// 浏览器控制台
localStorage.setItem('sidebar-width', '300')
location.reload()
```

### 查看组件树
```
DevTools → Components 标签
App.vue
  └─ MainLayout.vue
      ├─ Sidebar.vue
      └─ HostsManager.vue (或其他模块)
```

### 性能监控
```javascript
// 衡量侧栏拖动性能
performance.mark('drag-start')
// ... 拖动操作
performance.mark('drag-end')
performance.measure('drag', 'drag-start', 'drag-end')
console.log(performance.getEntriesByName('drag')[0].duration) // ms
```

---

## 🐛 常见问题

### Q: 侧栏不响应拖动？
**A:** 检查 MainLayout.vue 是否加载，确认 `isDragging` ref 正常工作
```javascript
console.log('isDragging:', isDragging.value)
```

### Q: 新模块未显示？
**A:** 
1. 确认在 App.vue 中导入了模块
2. 确认菜单项的 `id` 与条件判断一致
3. 检查 `currentModule` 值

### Q: 窗口最小尺寸未生效？
**A:** 需要重新启动应用（开发模式：`wails dev` 重启）

### Q: 侧栏宽度未保存？
**A:** 检查浏览器 localStorage 是否启用
```javascript
try {
    localStorage.setItem('test', 'test')
    localStorage.removeItem('test')
    console.log('localStorage 可用')
} catch (e) {
    console.error('localStorage 不可用:', e)
}
```

### Q: 样式冲突导致显示错乱？
**A:** 确保所有 `<style>` 标签都添加了 `scoped` 属性

---

## 📚 开发流程

### 1️⃣ 启动开发服务器
```bash
cd e:\projects\laolishu\LiSteward
wails dev
```

### 2️⃣ 打开应用
- Wails 自动打开桌面应用窗口
- 或访问 http://localhost:5173/

### 3️⃣ 编辑代码
- 前端代码自动热重载 (HMR)
- 后端代码需要重新启动

### 4️⃣ 测试新功能
- F12 打开 DevTools
- 测试侧栏拖动、折叠等交互
- 验证模块切换

### 5️⃣ 提交更改
```bash
git add .
git commit -m "feat: 描述你的改动"
git push
```

---

## 🔗 相关资源

### 文档
- [完整架构设计](./ARCHITECTURE.md)
- [变更提案](./openspec/changes/ARCH-001-APPLICATION-ARCHITECTURE-REFACTOR.md)
- [项目规范](./openspec/project.md)

### 外部资源
- [Vue 3 文档](https://vuejs.org/)
- [Wails 文档](https://wails.io/)
- [Tailwind CSS](https://tailwindcss.com/)

### 主要文件
- 主入口: `frontend/src/main.ts`
- Vite 配置: `frontend/vite.config.ts`
- Wails 配置: `wails.json`
- Go 后端: `main.go`

---

## ✅ 完成检查表

新建模块时，确保完成：

- [ ] 创建 Vue 文件在 `views/modules/`
- [ ] 编写模板和脚本
- [ ] 添加 scoped 样式
- [ ] 在 Sidebar.vue 中添加菜单项
- [ ] 在 App.vue 中导入并注册模块
- [ ] 在 App.vue 中添加条件渲染
- [ ] 测试模块加载和导航
- [ ] 测试样式和交互
- [ ] 提交代码并记录变更日志

---

**快速参考版本**: 1.0  
**最后更新**: 2025-12-23  
**维护者**: LiSteward 开发团队
