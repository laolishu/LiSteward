# 实施任务清单

## 1. 后端支持

- [x] 1.1 �?`internal/hosts/service.go` 中添�?`SaveWithReadOnlyRestore` 方法
  - 保存前检查文件只读属�?
  - 临时取消只读
  - 执行保存
  - 恢复只读属�?
  - 返回操作结果和错误信�?
  
- [x] 1.2 �?`app.go` 中暴露新方法 `SaveHostsEntriesWithReadOnlyRestore` 供前端调�?

- [x] 1.3 添加单元测试验证只读属性的正确处理

## 2. 前端状态管�?

- [x] 2.1 �?`HostsManager.vue` 中添加缓存状�?
  - 添加 `hasUnsavedChanges` 响应式标�?
  - 记录操作前的初始条目副本

- [x] 2.2 修改 `addEntry` 方法
  - 条目仅添加到内存，不调用 `SaveHostsEntries`
  - 设置 `hasUnsavedChanges = true`

- [x] 2.3 修改 `deleteEntry` 方法
  - 条目仅从内存中删除，不调�?`SaveHostsEntries`
  - 设置 `hasUnsavedChanges = true`

- [x] 2.4 修改 `saveEdit` 方法
  - 编辑后仅更新内存，不调用 `SaveHostsEntries`
  - 设置 `hasUnsavedChanges = true`

- [x] 2.5 修改 `toggleEntry` 方法
  - 保留即时响应式更新（�?UI�?
  - 设置 `hasUnsavedChanges = true`

## 3. 前端 UI 与交�?

- [x] 3.1 添加"应用"按钮到工具栏
  - 显示条件：`hasUnsavedChanges === true`
  - 点击时调�?`applyChanges()` 方法

- [x] 3.2 添加"刷新"按钮到工具栏
  - 始终显示
  - 点击时调�?`reloadHosts()` 方法

- [x] 3.3 实现 `applyChanges()` 方法
  - 调用后端�?`SaveHostsEntriesWithReadOnlyRestore`
  - 成功后设�?`hasUnsavedChanges = false`
  - 显示成功提示，包含备份信�?

- [x] 3.4 实现 `reloadHosts()` 方法
  - �?`hasUnsavedChanges === true`，显示确认对话框
  - 确认后调�?`loadHosts()`
  - 设置 `hasUnsavedChanges = false`

- [x] 3.5 添加视觉指示
  - "应用"按钮可见时，使用强调颜色（如红色或橙色）
  - 工具栏或标题栏显�?存在未保存更�?提示

## 4. 测试与验�?

- [x] 4.1 手工测试：缓存模式下的添�?删除/编辑
  - 添加条目后不立即出现�?Hosts 文件
  - 点击"刷新"后，内存中未保存的更改被丢弃
  - 点击"应用"后，更改写入 Hosts 文件

- [x] 4.2 手工测试：只�?Hosts 文件处理
  - �?Hosts 文件设置为只�?
  - 点击"应用"，系统自动处理权限，写入成功
  - Hosts 文件在保存后恢复为只�?

- [x] 4.3 手工测试：UI 反馈
  - 修改条目后，"应用"按钮变为可见
  - 点击"刷新"后，按钮隐藏

- [x] 4.4 集成测试（可选）
  - 验证前后端交互正确�?
  - 验证错误处理流程

## 5. 文档更新

- [x] 5.1 更新 `specs/hosts-management/spec.md`
  - 添加"缓存模式"相关需�?
  - 添加"应用/刷新"操作场景
  - 添加"只读属性处�?场景
