# 设计文档：用 SVG 图标替换 Emoji 图标

## 架构设计

### 图标管理策略

#### 存储位置
- **位置**：`frontend/public/images/`
- **格式**：SVG（可缩放矢量图形）
- **命名规则**：`<component>-<action>.svg`，例如：
  - `menu-hosts.svg` - Hosts 菜单图标
  - `menu-node.svg` - Node 菜单图标
  - `menu-settings.svg` - 设置菜单图标
  - `toolbar-menu.svg` - 工具栏菜单按钮
  - `toolbar-search.svg` - 工具栏搜索图标
  - `toolbar-refresh.svg` - 工具栏刷新按钮
  - `toolbar-apply.svg` - 工具栏应用按钮
  - `message-success.svg` - 成功消息图标
  - `message-error.svg` - 错误消息图标

#### 图标规范
- **尺寸**：16x16px（标准）、18x18px（大图标）
- **线条宽度**：1.4-2px（可读性）
- **颜色**：
  - **默认**：#24292e（深灰色，与 Tailwind gray-800 类似）
  - **绿色操作**（确定/应用）：#10B981（Tailwind emerald-500）
  - **红色操作**（删除/错误）：#EF4444（Tailwind red-500）
  - **中性操作**（菜单、搜索、刷新）：#586069（Tailwind gray-600）

#### 使用方式

在组件中使用 `<img>` 标签引用 SVG：
```vue
<img src="/images/menu-hosts.svg" alt="hosts" width="18" height="18" />
```

或在 CSS 中作为背景：
```css
background-image: url('/images/toolbar-search.svg');
```

### 实现计划

#### Phase 1：创建 SVG 图标
1. 为所有必需的操作创建 SVG 文件
2. 确保颜色和样式一致
3. 验证图标清晰度和可辨识性

#### Phase 2：更新组件
1. 修改 `Sidebar.vue`：
   - 将 `item.icon` 从 Emoji 字符串改为 SVG 文件路径
   - 使用 `<img>` 标签显示图标
   - 更新 CSS 以支持图标的对齐和悬停效果

2. 修改 `HostsManager.vue`：
   - 替换所有 Emoji 操作按钮为 SVG 图标
   - 保持按钮的尺寸和交互效果
   - 更新相关的 CSS 样式

#### Phase 3：测试和验证
1. 视觉验证：确保图标显示正确
2. 交互验证：确保按钮功能正常
3. 跨平台验证：在 Windows/macOS/Linux 上测试

### 向后兼容性

- 变更仅影响前端显示层
- 无 API 或数据结构变更
- 用户配置和数据完全不受影响
- 可随时回滚到 Emoji 版本

## 技术决策

### 为什么选择 SVG？
1. **可缩放**：无损缩放，适应不同分辨率
2. **可定制**：支持 CSS 样式和颜色自定义
3. **轻量**：文件小，加载快
4. **版本控制**：文本格式，易于版本管理

### 为什么选择 `<img>` 标签而不是 `<svg>` 元素？
- 简洁：图标文件管理更清晰
- 可缓存：浏览器缓存 SVG 文件
- 性能：使用文件引用避免内联 SVG 导致的 HTML 膨胀

### 颜色选择
- 遵循现有的 Tailwind CSS 配色方案
- 确保与应用整体设计一致
- 支持深色模式扩展（未来功能）

## 潜在问题和解决方案

| 问题 | 解决方案 |
|------|--------|
| SVG 加载时间 | 使用小的、优化的 SVG 文件；浏览器缓存 |
| 跨平台图标显示差异 | 使用标准 SVG 格式；在测试中覆盖多个平台 |
| 图标清晰度（小尺寸） | 使用 16-18px 标准尺寸；避免过细的线条 |
| 样式定制 | 通过外层 `<img>` 容器的 CSS 实现（过滤器、变换等） |

## 性能考虑

- **文件大小**：每个 SVG 通常 < 1KB
- **加载时间**：可忽略不计
- **缓存**：浏览器自动缓存，用户第一次加载后无性能影响
- **渲染**：SVG 渲染性能优于 Emoji（更确定的渲染结果）
