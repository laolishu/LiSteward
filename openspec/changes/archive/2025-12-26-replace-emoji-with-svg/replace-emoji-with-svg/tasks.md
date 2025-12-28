# 任务清单：用 SVG 图标替换 Emoji 图标

## 任务列表

### Phase 1：SVG 图标创建

#### 1. 创建菜单图标
- [x] 1.1 创建 `frontend/public/images/menu-hosts.svg`（Hosts 图标）
- [x] 1.2 创建 `frontend/public/images/menu-node.svg`（Node 图标）
- [x] 1.3 创建 `frontend/public/images/menu-settings.svg`（设置图标）
- [x] 1.4 验证三个菜单图标的颜色和大小一致

#### 2. 创建工具栏操作图标
- [x] 2.1 创建 `frontend/public/images/toolbar-menu.svg`（菜单按钮）
- [x] 2.2 创建 `frontend/public/images/toolbar-search.svg`（搜索图标）
- [x] 2.3 创建 `frontend/public/images/toolbar-refresh.svg`（刷新按钮）
- [x] 2.4 创建 `frontend/public/images/toolbar-apply.svg`（应用/保存按钮，绿色）

#### 3. 创建消息图标
- [x] 3.1 创建 `frontend/public/images/message-success.svg`（成功消息图标，绿色）
- [x] 3.2 创建 `frontend/public/images/message-error.svg`（错误消息图标，红色）

#### 4. 图标验证
- [x] 4.1 验证所有 SVG 文件大小 < 1KB
- [x] 4.2 验证所有 SVG 文件可在浏览器中正确显示
- [x] 4.3 验证图标的清晰度和可辨识性

### Phase 2：Sidebar 组件更新

#### 5. 修改 Sidebar.vue
- [x] 5.1 更新菜单项定义，将 Emoji 图标替换为 SVG 文件路径
  - 原：`{ id: 'hosts', icon: '🌐', label: t('menu.hosts') }`
  - 新：`{ id: 'hosts', icon: '/images/menu-hosts.svg', label: t('menu.hosts') }`
- [x] 5.2 在模板中将 `<span class="nav-icon">{{ item.icon }}</span>` 改为 `<img :src="item.icon" />`
- [x] 5.3 更新 CSS 类 `.nav-icon` 适应图像显示（width、height、display 等）
- [x] 5.4 验证菜单图标在不同状态下的显示（正常、悬停、活跃）

### Phase 3：HostsManager 组件更新

#### 6. 替换工具栏图标
- [x] 6.1 将菜单按钮（☰）替换为 `<img src="/images/toolbar-menu.svg" />`
- [x] 6.2 将搜索图标（🔍）替换为 `<img src="/images/toolbar-search.svg" />`
- [x] 6.3 将刷新按钮（🔄）替换为 `<img src="/images/toolbar-refresh.svg" />`
- [x] 6.4 将应用按钮中的 Emoji 替换为 `<img src="/images/toolbar-apply.svg" />`

#### 7. 替换消息图标
- [x] 7.1 将成功消息图标（✓）替换为 `<img src="/images/message-success.svg" />`
- [x] 7.2 将错误消息图标（✗）替换为 `<img src="/images/message-error.svg" />`

#### 7.5. 替换表格操作列图标
- [x] 7.5.1 创建 `action-swap.svg`、`action-edit.svg`、`action-delete.svg` 图标
- [x] 7.5.2 将表格操作列中的交换（⇄）替换为 SVG 图标
- [x] 7.5.3 将表格操作列中的编辑（✏️）替换为 SVG 图标
- [x] 7.5.4 将表格操作列中的删除（🗑️）替换为 SVG 图标
- [x] 7.5.5 将侧栏方案编辑（✏️）和删除（🗑️）按钮替换为 SVG 图标

#### 8. 样式调整
- [x] 8.1 更新 `.btn-icon` CSS 以适应 SVG 图像
- [x] 8.2 更新 `.message-icon` CSS 以适应 SVG 图像
- [x] 8.3 更新 `.search-icon` CSS 以适应 SVG 图像
- [x] 8.4 验证所有按钮的悬停和活跃状态样式

### Phase 4：集成测试和验证

#### 9. 视觉验证
- [x] 9.1 验证菜单图标在不同屏幕分辨率下的显示
- [x] 9.2 验证工具栏图标的大小和对齐
- [x] 9.3 验证消息图标的颜色对比度
- [x] 9.4 验证图标在暗光和明亮背景下的可见性

#### 10. 功能验证
- [x] 10.1 验证菜单项点击功能正常
- [x] 10.2 验证所有工具栏按钮功能正常
- [x] 10.3 验证消息显示功能正常
- [x] 10.4 前端编译无错误：`npm run build`

#### 11. 跨平台测试
- [ ] 11.1 在 Windows 上验证所有图标显示正确
- [ ] 11.2 在 macOS 上验证所有图标显示正确（如可用）
- [ ] 11.3 在 Linux 上验证所有图标显示正确（如可用）

#### 12. 性能验证
- [ ] 12.1 检查应用加载时间无明显增加
- [ ] 12.2 检查浏览器开发工具中的网络请求（SVG 文件应被缓存）
- [ ] 12.3 验证内存占用无异常增长

### Phase 5：文档和交付

#### 13. 文档更新
- [ ] 13.1 更新 `ARCHITECTURE.md` 关于 UI 资源的部分
- [ ] 13.2 在 `frontend/README.md` 中添加图标管理说明

#### 14. 提交和验证
- [x] 14.1 所有文件已提交
- [x] 14.2 运行 `openspec-cn validate replace-emoji-with-svg --strict` 通过
- [ ] 14.3 准备合并到主分支

## 依赖关系

- Phase 1 必须完成后才能进行 Phase 2 和 Phase 3
- Phase 2 和 Phase 3 可以并行进行
- Phase 4 需要 Phase 2 和 Phase 3 完成后才能进行

## 工作量估计

| Phase | 任务数 | 预期时间 |
|-------|--------|---------|
| 1 | 7 | 2 小时 |
| 2 | 4 | 1 小时 |
| 3 | 8 | 2 小时 |
| 4 | 4 | 1.5 小时 |
| 5 | 2 | 0.5 小时 |
| **总计** | **25** | **7 小时** |
