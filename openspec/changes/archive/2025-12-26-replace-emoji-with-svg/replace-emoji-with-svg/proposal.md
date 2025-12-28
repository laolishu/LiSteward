# 提案：用 SVG 图标替换 Emoji 图标

## 为什么

目前应用中的菜单和操作按钮使用 Emoji 表情符号（如 🌐、🧩、⚙️、🔍、🔄 等），存在以下问题：
- Emoji 渲染不一致，跨平台显示可能不同
- Emoji 无法统一自定义样式（颜色、大小、悬停效果）
- 无法满足设计规范的一致性要求
- 不够专业，不适合桌面应用的企业级外观

使用 SVG 图标可以：
- 提供一致的、可缩放的图形显示
- 支持主题化和样式自定义
- 更专业的应用外观
- 便于维护和版本控制

## 变更内容

### 范围

1. **菜单图标**：将左侧导航栏的 Emoji 图标替换为 SVG 图标
   - Hosts 管理：🌐 → hosts.svg
   - Node 管理：🧩 → node.svg
   - 设置：⚙️ → settings.svg

2. **Hosts Manager 操作按钮**：将 HostsManager.vue 中的 Emoji 图标替换为 SVG 图标
   - 菜单按钮：☰ → menu.svg
   - 搜索图标：🔍 → search.svg
   - 刷新按钮：🔄 → refresh.svg
   - 应用按钮：💾 → apply.svg
   - 消息图标：✓/✗ → success.svg / error.svg

3. **SVG 文件存储**：所有 SVG 图标存储在 `frontend/public/images/` 目录

### 验收标准

- **必须**：创建所有必需的 SVG 图标文件并保存在 `frontend/public/images/` 中
- **必须**：更新 Sidebar.vue 使用 SVG 图标替换 Emoji
- **必须**：更新 HostsManager.vue 使用 SVG 图标替换 Emoji
- **必须**：应用前后的视觉效果保持一致或改进
- **必须**：所有 SVG 图标使用一致的配色方案
- **应该**：为图标配置悬停状态的样式改变

### 影响范围

- 前端：修改 `frontend/src/components/Sidebar.vue`、`frontend/src/views/HostsManager.vue`
- 静态资源：新增 SVG 图标文件到 `frontend/public/images/`
- 无后端变更
- 无 i18n 变更

### 风险与注意事项

- 需要确保 SVG 图标的加载不影响应用启动性能
- 跨平台测试（Windows/macOS/Linux）确保图标正确显示
- 需要保证现有的鼠标悬停和交互效果仍然有效

## 下一步

起草 `tasks.md`、`design.md` 与规范增量 `specs/ui-icons/spec.md`，并运行 `openspec-cn validate replace-emoji-with-svg --strict`。
