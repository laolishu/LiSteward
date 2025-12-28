# LiSteward 应用架构文档

## 📐 总体架构设计

### 三层结构

```
┌─────────────────────────────────────────────────────────────┐
│                    App.vue (应用入口)                         │
├─────────────────────────────────────────────────────────────┤
│                 MainLayout.vue (主框架)                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ Header (50px)                                        │   │
│  │ ├─ 侧栏折叠按钮 ├─ 空白区 ├─ 同步状态指示器            │   │
│  ├──────────────────────────────────────────────────────┤   │
│  │  Sidebar           │  Divider  │   Module Content    │   │
│  │ (可拖动宽度)       │ (拖动分隔) │   (各模块显示)      │   │
│  │  300px(default)    │   4px     │   flex: 1           │   │
│  │                    │           │                     │   │
│  │  - 导航菜单        │           │  HostsManager       │   │
│  │  - 用户信息        │           │  Dashboard          │   │
│  │  - 同步状态        │           │  EnvVariables       │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 📁 文件结构

```
frontend/src/
├── App.vue                          # 全局入口，仅处理路由
├── layouts/
│   └── MainLayout.vue              # 主框架容器 (300行)
│       ├─ 处理窗口布局
│       ├─ 侧栏宽度管理（可拖动分隔符）
│       ├─ 侧栏折叠/展开
│       ├─ 同步状态显示
│       └─ 响应式设计
├── views/
│   ├── HostsManager.vue            # Hosts管理模块 (1200+行)
│   │   ├─ 左侧：配置方案 & 备份历史
│   │   ├─ 右侧：代码编辑器界面
│   │   ├─ 可拖动分隔符（模块内部）
│   │   └─ 所有业务逻辑
│   ├── modules/                     # 其他模块目录（待实现）
│   │   ├── Dashboard.vue
│   │   ├── EnvironmentVariables.vue
│   │   ├── Registry.vue
│   │   └── Settings.vue
├── components/
│   ├── Sidebar.vue                  # 全局导航栏 (260行)
│   │   ├─ 5个菜单项
│   │   ├─ 用户资料卡
│   │   └─ 同步状态指示器
│   └── ...

wails.json                           # Wails配置文件
└─ windows[0]:
   └─ minWidth: 1200, minHeight: 800 # 窗口最小尺寸
```

## 🎯 各层职责

### 1️⃣ **App.vue** - 应用入口层
**职责：**
- 全局路由管理
- 模块切换逻辑
- 引用 MainLayout 容器

**关键代码：**
```vue
<MainLayout ref="mainLayoutRef">
  <HostsManager v-if="currentModule === 'hosts'" />
  <div v-else class="module-placeholder">...</div>
</MainLayout>
```

**尺寸：** ~50 行

---

### 2️⃣ **MainLayout.vue** - 主框架层
**职责：**
- 窗口级别的布局管理
- Sidebar 和 Content 容器的分割与响应式调整
- Sidebar 宽度控制（拖动分隔符）
- Sidebar 折叠/展开切换
- 同步状态显示
- 响应式设计与媒体查询
- localStorage 侧栏宽度缓存

**关键特性：**
| 特性 | 实现 |
|------|-----|
| 侧栏宽度范围 | 200px ~ 400px |
| 默认宽度 | 300px |
| 分隔符宽度 | 4px，可拖动 |
| 分隔符颜色 | #E1E4E8，悬停时 #34D399 |
| 侧栏折叠显示 | 60px (仅图标) |
| 最小窗口大小 | 1200px × 800px |

**核心方法：**
```javascript
toggleSidebar()        // 折叠/展开
startDragging(e)       // 开始拖动分隔符
setSyncStatus(bool)    // 设置同步状态
```

**尺寸：** ~300 行

---

### 3️⃣ **Sidebar.vue** - 全局导航组件
**职责：**
- 全局导航菜单（5个菜单项）
- 用户资料显示
- 同步状态指示
- 菜单项点击事件发出

**菜单项：**
- 📊 控制台 (dashboard)
- 🌐 Hosts 管理 (hosts)
- ⚙️ 环境变量 (env)
- 📝 注册表 (registry)
- ⚙️ 设置 (settings)

**尺寸：** ~260 行

---

### 4️⃣ **HostsManager.vue** - 模块层（示例）
**职责：**
- 完整的 Hosts 文件管理功能
- 左右分屏布局（模块内部）
- 配置方案管理
- 备份历史管理
- Hosts 条目编辑

**内部结构：**
```
HostsManager
├─ 左侧面板 (280px, 可拖动)
│  ├─ 配置方案列表
│  ├─ 方案卡片
│  ├─ 备份历史
│  └─ 恢复按钮
├─ 分隔符 (4px)
└─ 右侧面板 (flex: 1)
   ├─ 搜索框 & 过滤器
   └─ Hosts 条目列表
      ├─ 启用/禁用切换
      ├─ IP & 域名显示
      ├─ 编辑/删除按钮
      └─ 添加新条目栏
```

**关键方法：**
- loadHosts() - 加载 Hosts 文件
- saveHosts() - 保存到系统
- addEntry() - 添加条目
- editEntry() / deleteEntry() - 修改/删除条目
- selectProfile() - 切换方案
- createBackup() / restoreBackup() - 备份管理

**尺寸：** ~1200+ 行

---

## 🎨 设计系统

### 颜色方案（浅色主题）
| 用途 | 颜色代码 | RGB |
|------|----------|-----|
| 背景 | #F6F8FA | 亮灰 |
| 主要内容区 | #FFFFFF | 白色 |
| 文本 | #24292E | 深灰 |
| 次要文本 | #57606A | 中灰 |
| 边框 | #E1E4E8 | 浅灰 |
| 强调色 | #34D399 | 翠绿 |
| 成功 | #1B7C34 | 绿色 |
| 危险 | #CB2431 | 红色 |

### 响应式断点
```scss
@media (max-width: 1400px) { /* 大屏 */ }
@media (max-width: 1000px) { /* 中屏 */ }
@media (max-width: 800px)  { /* 小屏 */ }
```

### 字体
```
系统字体: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue'
代码字体: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas'
```

---

## 🔄 数据流

### 路由与状态管理
```
App.vue
├─ currentModule (ref) ──────────┐
│                                 ↓
├─ handleNavigate(moduleName) ←─ Sidebar 
│  └─ 更新 currentModule
│
└─ MainLayout ──────────┐
   ├─ sidebarWidth     │ 模块独立管理
   ├─ sidebarCollapsed │
   ├─ syncStatus       └─→ 同步状态通知给主框架
   └─ isSyncing
```

### 模块通信
```
模块内部
├─ 独立的 ref() 状态
├─ computed() 计算属性
└─ 异步方法
   └─ Wails 后端调用 (GetHostsEntries, SaveHostsEntries, ...)
```

---

## 🚀 快速开始

### 开发环境
```bash
cd frontend/src
# 组件已自动编译，访问 http://localhost:5173
```

### 添加新模块
1. 在 `views/modules/` 中创建 `NewModule.vue`
2. 在 `Sidebar.vue` 中添加菜单项
3. 在 `App.vue` 中添加路由条件：
```vue
<NewModule v-if="currentModule === 'newmodule'" />
```

### 自定义侧栏宽度
```typescript
// MainLayout.vue
const minSidebarWidth = 150  // 最小宽度
const maxSidebarWidth = 500  // 最大宽度
```

### 窗口最小尺寸配置
```json
// wails.json
"windows": [{
  "minWidth": 1200,
  "minHeight": 800
}]
```

---

## ✨ 关键特性

### ✅ 已实现
- [x] 三层架构（App → MainLayout → Modules）
- [x] 侧栏拖动分隔符（全局）
- [x] 侧栏折叠/展开
- [x] 浅色主题设计
- [x] 响应式布局
- [x] 最小窗口尺寸约束
- [x] 模块化结构
- [x] Hosts 管理完整功能
- [x] localStorage 状态缓存

### ⏳ 待实现
- [ ] Dashboard 模块
- [ ] Environment Variables 模块
- [ ] Registry 模块
- [ ] Settings 模块
- [ ] 深色主题 (可选)
- [ ] 国际化 i18n
- [ ] 单元测试
- [ ] E2E 测试
- [ ] 打包与分发

---

## 📊 性能考虑

### 优化策略
1. **侧栏宽度** - localStorage 缓存，无需重复拖动设置
2. **模块懒加载** - 各模块按需加载（可扩展）
3. **滚动条优化** - 自定义样式，减少重排
4. **CSS 作用域** - 各组件 `scoped` 样式，避免污染

### 监控指标
- 首屏加载时间：< 2s
- 侧栏拖动帧率：> 60fps
- 内存占用：< 150MB

---

## 🔐 安全性

### 权限管理
- Hosts 文件修改需要管理员权限
- Wails 后端在 Go 层验证权限
- 前端不处理敏感操作

### 数据验证
- 条目添加时验证 IP & 域名格式
- 文件保存前创建自动备份
- 备份历史保留 5 个最新版本

---

## 📞 支持与贡献

### 文件位置
- **布局组件**: `frontend/src/layouts/`
- **业务模块**: `frontend/src/views/`
- **通用组件**: `frontend/src/components/`
- **样式系统**: 各组件 `<style scoped>`

### 开发规范
- Vue 3 Composition API
- TypeScript（可选）
- Tailwind CSS 工具类
- BEM 命名规范

---

**最后更新**: 2025年12月23日  
**架构版本**: 1.0 - 三层分离式架构
