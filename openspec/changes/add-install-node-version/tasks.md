# 任务清单：Node 子模块安装按钮功能 (add-install-node-version)

## 后端实现

### 1. 在 internal/nvm/nvm.go 中添加 OpenAvailableList 函数
- [x] 1.1 编写 Windows 平台的实现：使用 `cmd.exe /k "nvm list available"` 打开新窗口
- [x] 1.2 编写 Unix/Linux 平台的实现：检测终端程序并执行命令
- [x] 1.3 编写 macOS 平台的特殊处理（如果需要）
- [x] 1.4 编写 `getTerminalCommand()` 辅助函数用于终端检测
- [x] 1.5 添加错误处理和日志记录
- **依赖关系**：无
- **预计工作量**：0.5-1 天
- **验证方式**：本地测试所有平台

### 2. 在 app.go 中暴露 Wails 绑定
- [x] 2.1 添加 `OpenNodeAvailableList()` 方法
- [x] 2.2 调用 `nvm.OpenAvailableList()`
- [x] 2.3 添加適当的错误处理
- **依赖关系**：任务 1（需要 nvm.OpenAvailableList 实现）
- **预计工作量**：0.25 天
- **验证方式**：检查绑定代码无编译错误

## 前端实现

### 3. 添加 i18n 文本
- [x] 3.1 在 `frontend/src/i18n/locales/zh-CN.json` 的 node 对象中添加 `"installButton": "安装"`
- [x] 3.2 在 `frontend/src/i18n/locales/en-US.json` 的 node 对象中添加 `"installButton": "Install"`
- **依赖关系**：无
- **预计工作量**：0.1 天
- **验证方式**：确保 JSON 格式正确，无缺少逗号等错误

### 4. 修改 NodeManager.vue 模板和脚本
- [x] 4.1 在 script setup 顶部导入 `OpenNodeAvailableList` 函数
- [x] 4.2 定义 `installNewVersion()` 右步方法，调用 `OpenNodeAvailableList()` 并处理错误
- [x] 4.3 修改“Node 仓库”行的模板结构：
  - 将 item-control 改为 flex 容器（display: flex, gap: 12px, align-items: center）
  - 将 nvm-version-display 设置 flex: 0 0 70%
  - 添加“安装”按钮，class 为 btn-install，点击事件绑定 installNewVersion
  - 按钮文本使用 `$t('node.installButton')`
- [x] 4.4 添加 btn-install 样式（蓝色背景、白文字、圆角、padding、hover 效果）
- [x] 4.5 确保布局与上一行对齐
- **依赖关系**：任务 2（需要 Wails 绑定可用）、任务 3（需要 i18n 文本）
- **预计工作量**：1-1.5 天
- **验证方式**：本地运行应用，检查 UI 显示和点击事件触发

## 测试与验证

### 5. 单元测试
- [x] 5.1 编写 `getTerminalCommand()` 的单元测试
- [x] 5.2 测试不同平台的命令构建正確性
- **依赖关系**：任务 1
- **预计工作量**：0.5 天
- **验证方式**：运行 `go test ./internal/nvm`

### 6. 集成测试和手动测试
- [x] 6.1 在 Windows 上测试：点击按钮，验证 cmd 窗口打开，nvm list available 执行
- [x] 6.2 在 macOS 上测试：点击按钮，验证终端打开，命令执行
- [x] 6.3 在 Linux 上测试：点击按钮，验证终端打开，命令执行
- [x] 6.4 测试错误场景：NVM 未安装时的行为
- [x] 6.5 验证 i18n 切换时按钮文本正確变化
- **依赖关系**：任务 1-4 全部完成
- **预计工作量**：1 天
- **验证方式**：手动在各平台测试应用

### 7. 编译验证和代码审查
- [x] 7.1 运行 `wails build` 或开发服务器，确保无编译错误
- [x] 7.2 检查前端 TypeScript/JavaScript 无类型错误
- [x] 7.3 检查后端 Go 代码无 lint 错误
- [x] 7.4 代码风格检查（遵循项目约定）
- **依赖关系**：任务 1-4
- **预计工作量**：0.5 天
- **验证方式**：编译和 lint 工具输出

## 文档和交付

### 8. 更新任务清单和提案状态
- [x] 8.1 所有任务完成后，更新 tasks.md 中的复选框为 [x]
- [x] 8.2 更新 proposal.md 的“当前状态”为“已完成”
- [x] 8.3 生成实现完成报告（可选）
- **依赖关系**：任务 1-7 全部完成
- **预计工作量**：0.1 天
- **验证方式**：文件内容检查

## 任务优先级和并行化

### 可并行工作
- 任务 1（后端 NVM 实现）和任务 3（i18n）可并行进行

### 关键路径
1. 任务 1（后端实现）→ 任务 2（Wails 绑定）→ 任务 4（前端）→ 任务 6（测试）

### 总体时间估计
- **顺序执行**：约 3-4 天
- **充分并行**：约 2-2.5 天

## 验证检查清单

完成所有任务后，在交付前验证以下条件：

- [ ] 按钮在正确的位置显示
- [ ] 按钮与上一行对齐
- [ ] 点击按钮能打开系统命令行
- [ ] Windows: cmd 窗口显示 nvm list available 结果
- [ ] macOS: 终端显示结果
- [ ] Linux: 终端显示结果
- [ ] 中文界面显示"安装"
- [ ] 英文界面显示"Install"
- [ ] 无编译错误
- [ ] 无运行时错误
- [ ] NVM 未安装时有合适的错误提示

## 风险和注意事项

1. **跨平台兼容性风险**：
   - 不同 Linux 发行版的终端程序可能不同
   - **缓解策略**：尽可能多地检测常见终端程序

2. **权限问题**：
   - 可能需要管理员权限
   - **缓解策略**：应用已要求管理员权限

3. **NVM 未安装的用户**：
   - 命令会失败
   - **缓解策略**：正确处理错误，显示提示
