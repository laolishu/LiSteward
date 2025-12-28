# 提案：优化 README.md 并添加中英文切换支持

**变更ID**: optimize-readme-with-language-switch

**状态**: 已实现完成

**优先级**: 中

**创建日期**: 2025-12-28

## 摘要

当前项目有两个独立的 README 文档：`README.md`（英文）和 `README_CN.md`（中文）。用户无法在文档中快速切换语言，需要手动在两个文件之间导航。同时，现有文档内容较为简洁，未能充分展示 Hosts 和 NVM 两个核心模块的功能与使用方法。本提案通过优化 README 内容并在两个文档中添加语言切换入口，提升用户体验和项目可发现性。

## 问题陈述

### 当前问题

1. **语言切换困难**：用户无法直接在文档中切换语言，需要手动点击文件链接或返回仓库根目录
2. **功能介绍不足**：现有 README 对两个核心模块（Hosts 管理、NVM 版本管理）的介绍过于简洁
   - 未详细说明 Hosts 管理的功能（配置管理、备份、导入/导出、快速切换等）
   - 未详细说明 NVM 管理的功能（版本检测、切换、环境管理等）
3. **缺乏使用场景**：文档没有说明用户为什么需要这两个功能，应用场景不清楚
4. **功能发现性低**：新用户可能不清楚应用包含哪些功能，可能遗漏某些有用的特性

## 设计目标

1. 在两个 README 文件中添加**语言切换链接**（顶部或合适位置）
2. 优化并扩充**功能介绍部分**，详细描述：
   - Hosts 管理器的核心功能
   - NVM 版本管理器的核心功能
   - 各功能的使用场景和优势
3. 改进**快速开始**部分，使其更易理解
4. 增加**功能特性列表**，列举主要特点
5. 保持文档结构一致，两个版本（中文/英文）内容对等

## 范围

### 包括

1. **README.md（英文）**
   - 顶部添加语言切换提示和链接（指向 README_CN.md）
   - 扩展"Key features"部分，详细介绍两个模块
   - 添加"Features"或"Capabilities"部分，列举应用功能
   - 优化"Build & Run"和"Build (release)"部分的说明

2. **README_CN.md（中文）**
   - 顶部添加语言切换提示和链接（指向 README.md）
   - 扩展"主要功能"部分，详细介绍两个模块
   - 添加"核心特性"部分，列举应用功能
   - 优化"构建与运行"部分的说明

3. **文档同步**
   - 确保中文和英文版本内容对等
   - 跨引用相关功能规范（如 i18n、user-config 等）

### 不包括

- 修改应用代码或功能实现
- 更改项目的技术栈或架构
- 创建新的文档文件（仅优化现有 README）

## 背景

### 当前情况

**README.md** 的结构：
- 简单的概述（2-3 行）
- 快速开始（2 行）
- 构建与运行（开发/发布两部分）
- 许可和赞赏信息
- 贡献指南

**README_CN.md** 的结构与英文版基本相同，但有部分信息不完整（如缺少 Hosts 管理的介绍）

### 模块功能分析

#### Hosts 管理模块 (`HostsManager.vue`)
- 可视化编辑 hosts 文件
- 多配置方案管理（方案切换）
- 本地备份和恢复
- 导入/导出功能
- 条目快速启用/禁用
- 搜索和过滤功能
- 支持多个 IP 地址和备用 IP

#### NVM 版本管理模块 (`NodeManager.vue`)
- 检测和显示已安装的 Node 版本
- 快速切换 Node 版本
- NVM 环境管理
- 显示当前 Node 版本
- 支持安装新版本（打开 `nvm list available`）
- 跨平台支持（Windows/macOS/Linux）

### 文档问题根因

1. 原有 README 按照标准开源项目模板，内容通用且简洁
2. 两个模块的功能和使用方法没有在 README 中充分体现
3. 没有明确的语言切换机制

## 实现策略

### 结构规划

#### README.md（英文）新结构：
```
1. Language Switch (顶部)
   - [中文版本](README_CN.md)

2. Overview
   - LiSteward 的定义、用途

3. Key Features
   - Hosts Configuration Management
     - Visual editor, profile management, backup/restore, import/export
   - Node Version Management (NVM)
     - Version detection, quick switching, environment management
   
4. Use Cases
   - 什么场景需要使用这个应用

5. Quick Start
   - Installation / Getting Started

6. Build & Run
   - Development 和 Release 构建说明

7. Documentation
   - 指向详细规范

8. License, Support, Contributing
   - 保持现有内容
```

#### README_CN.md（中文）新结构：
结构与英文版对应，内容对等

### 关键信息补充

1. **Hosts 管理的优势**：
   - 快速管理多个开发环境配置
   - 方便的备份和恢复机制
   - 直观的可视化界面，避免手动编辑错误

2. **NVM 版本管理的优势**：
   - 在多个 Node 版本之间快速切换
   - 项目特定版本管理
   - 跨平台支持

3. **典型使用场景**：
   - 前端开发：管理 hosts 映射、快速切换 Node 版本
   - 项目管理：为不同项目维护不同的开发配置

## 验收标准

- [ ] README.md 和 README_CN.md 顶部清晰显示语言切换链接
- [ ] 功能介绍部分准确描述 Hosts 和 NVM 两个模块的功能
- [ ] 中英文版本内容对等，逻辑结构一致
- [ ] 快速开始和构建说明更清晰易懂
- [ ] 所有链接有效（包括到其他规范、代码的链接）
- [ ] 文档符合 Markdown 最佳实践

## 相关资源

- [Hosts 管理规范](../specs/hosts-management/)
- [i18n 核心规范](../specs/i18n-core/)
- [项目上下文](../project.md)

## 下一步

1. 审核本提案，确认范围和目标
2. 完成 design.md 和 tasks.md 的详细规划
3. 生成规范增量（如需）
4. 进入 apply 阶段实现文档优化
