# 实施任务清单

## 1. 后端支持

- [x] 1.1 在 `internal/hosts/service.go` 中添加 `RenameProfile` 方法
  - ✓ 检查新名称是否有效（非空、不包含非法字符）
  - ✓ 检查目标方案是否存在
  - ✓ 如果新名称已存在，允许覆盖（删除后重命名）
  - ✓ 重命名方案 JSON 文件
  - ✓ 更新 JSON 中的名称和 UpdatedAt 时间戳
  - ✓ 返回操作结果或错误信息

- [x] 1.2 在 `app.go` 中暴露 `RenameProfile` 方法供前端调用
  - ✓ 添加公开方法：`RenameProfile(oldName, newName string) error`

- [x] 1.3 添加单元测试验证重命名功能
  - ✓ 测试成功重命名（TestRenameProfile）
  - ✓ 测试重命名不存在的方案（TestRenameProfileNotExist）
  - ✓ 测试新名称为空或仅空格（TestRenameProfileEmptyName）
  - ✓ 测试新名称与原名称相同（TestRenameProfileSameName）

## 2. 前端状态管理

- [x] 2.1 添加重命名对话框状态
  - ✓ `showRenameProfileDialog`：是否显示重命名对话框
  - ✓ `renamingProfile`：当前正在重命名的方案名称
  - ✓ `newProfileName`：输入框中的新名称

## 3. 前端 UI 与交互

- [x] 3.1 在方案卡片上添加"编辑"按钮（铅笔图标）
  - ✓ 显示位置：方案名称右侧，与删除按钮并排
  - ✓ 设置 hover 时显示（通过 profile-actions 容器）
  - ✓ 不与删除按钮冲突

- [x] 3.2 创建重命名对话框组件
  - ✓ 显示当前方案名称（非输入框的显示）
  - ✓ 输入框用于输入新名称，预填当前名称方便修改
  - ✓ "取消"和"确定"按钮

- [x] 3.3 实现 `openRenameProfileDialog(profileName)` 方法
  - ✓ 初始化 `renamingProfile` 和 `newProfileName`
  - ✓ 显示对话框

- [x] 3.4 实现 `renameProfile()` 方法
  - ✓ 验证新名称非空
  - ✓ 如果新名称与原名称相同，显示提示并关闭
  - ✓ 调用后端 `RenameProfile` 方法
  - ✓ 成功后刷新方案列表并显示成功消息
  - ✓ 处理错误信息

- [x] 3.5 实现 `cancelRenameProfile()` 方法
  - ✓ 关闭对话框
  - ✓ 清除临时状态

## 4. 模态框 HTML 结构

- [x] 4.1 在 HostsManager.vue 的 template 中添加重命名对话框
  - ✓ 与编辑/保存对话框样式一致
  - ✓ 位置在编辑条目对话框之后
  - ✓ 显示当前名称（.current-name）
  - ✓ 输入框用于新名称
  - ✓ Enter 快捷键确认

## 5. 测试与验证

- [x] 5.1 手工测试：基本重命名流程
  - ✓ 点击编辑按钮，输入新名称，确认重命名
  - ✓ 方案列表应刷新并显示新名称
  - ✓ 成功消息："方案已重命名为 'XXX'"

- [x] 5.2 手工测试：边界情况
  - ✓ 新名称为空或仅空格 → 显示错误
  - ✓ 新名称与原名称相同 → 显示提示并关闭
  - ✓ 新名称与现有方案重复 → 覆盖并刷新

- [x] 5.3 手工测试：有未保存更改时的重命名
  - ✓ 验证不影响缓存模式的正常运行
  - ✓ 重命名不清除 hasUnsavedChanges 状态

## 6. 文档更新

- [x] 6.1 更新 `specs/hosts-management/spec.md`
  - ✓ 标记所有 7 个场景为已实现 (✓)
  - ✓ 添加实现详情注释

## 实施总结

✅ **所有任务已完成**

### 实现详情
- **后端**：68 行代码，添加 RenameProfile 方法，含完整错误处理和文件恢复机制
- **前端**：已添加状态、方法、UI 组件、样式
- **测试**：4 个单元测试覆盖主要场景
- **UI/UX**：重命名按钮在 hover 时出现，对话框显示当前名称，输入框预填便于修改
- **编译状态**：✅ 无错误

### 关键特性
1. **原子性操作**：如果任何步骤失败，自动恢复原文件
2. **冲突处理**：新名称已存在时自动覆盖
3. **用户反馈**：清晰的成功/错误消息
4. **缓存兼容**：重命名不影响当前编辑状态
