# 提案：Node 子模块添加安装新版本功能 (add-install-node-version)

## 摘要

在 Node 管理模块的"Node 仓库"行添加"安装"按钮，点击时自动弹出系统命令行窗口并执行 `nvm list available` 命令，让用户能够查看可用的 Node 版本并进行安装。

## 问题陈述

当前 Node 管理模块提供了版本列表和切换功能，但用户无法直观地看到有哪些可用的 Node 版本可以安装。虽然已有"安装"相关后端功能（`InstallNodeVersion`），前端需要：

1. 一个更友好的方式让用户发现和选择可用版本
2. 用户可以快速打开命令行查看可用版本列表（`nvm list available`）
3. 按钮应与现有 UI 布局和对齐方式保持一致

## 范围

### 包括
- 在 NodeManager.vue 的"Node 仓库"行添加"安装"按钮
- 按钮放置在文本框后面，保持与上一行对齐
- 缩小"Node 仓库"文本框宽度，为按钮腾出空间
- 点击按钮时，后端打开系统命令行窗口并执行 `nvm list available`
- 添加相应的 i18n 文本支持

### 不包括
- 改变 Node 版本列表界面
- 修改现有的版本切换逻辑
- 在应用内集成可用版本列表显示（外部命令行方式）

## 设计概要

### 前端改动
1. **NodeManager.vue 模板修改**：
   - 缩小"Node 仓库"行的文本框宽度（约 70%）
   - 在文本框后添加"安装"按钮（约 30%）
   - 按钮使用一致的样式（`btn-open` 或类似）

2. **事件处理**：
   - 添加 `installNewVersion()` 方法
   - 调用后端的 `OpenNodeAvailableList()` 接口（新建）

### 后端改动
1. **App.go**：
   - 新增 `OpenNodeAvailableList()` 方法
   - 调用 `internal/nvm` 中的对应实现

2. **internal/nvm/nvm.go**：
   - 新增 `OpenAvailableList()` 函数
   - 打开系统命令行窗口并执行 `nvm list available`
   - Windows 和类 Unix 平台分别处理

### i18n 新增
- `node.installButton`: "安装" （zh-CN）, "Install" （en-US）

## 实现方式

### 核心逻辑
1. 点击"安装"按钮 → 调用 `installNewVersion()`
2. 前端调用 `OpenNodeAvailableList()` 后端接口
3. 后端在 Windows 上调用 `cmd.exe` 打开新窗口运行 `nvm list available`
4. 在类 Unix 上调用相应的终端命令

### UI 布局
```
NVM 版本:  [显示版本/安装按钮]

Node 仓库:  [仓库路径---][安装]  ← 按钮在这里
```

## 依赖关系
- 依赖现有的 `internal/nvm` 包功能
- 依赖现有的 Wails 绑定机制
- 不依赖其他提案

## 验收标准
- [ ] 按钮正确显示在"Node 仓库"行
- [ ] 按钮与上一行元素对齐
- [ ] 点击按钮时打开系统命令行窗口
- [ ] Windows 和 macOS/Linux 都能正常工作
- [ ] i18n 文本正确显示
- [ ] 无编译错误和类型错误

## 当前状态
**已完成** - 所有功能已实现并通过验证
