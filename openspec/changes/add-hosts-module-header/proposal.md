# 提案：为 Hosts 管理模块添加 ModuleHeader 组件

**变更ID**: add-hosts-module-header

**状态**: 已实施完成

**优先级**: 中

**创建日期**: 2025-12-27

## 摘要

当前 HostsManager 模块拥有独特的布局结构（两面板设计：侧栏 + 编辑区），缺乏统一的模块页头。本提案通过为 HostsManager 添加 ModuleHeader 组件来统一整个应用的页头风格，同时适配其特有的布局和功能需求。

## 为什么

当前存在以下问题：

1. **不一致性**: HostsManager 与其他模块（NodeManager、Settings）的页头样式不统一，用户体验不一致
2. **功能重复**: 独立的 toast 消息提示系统与 ModuleHeader 的提示框功能重复，维护成本高
3. **代码冗余**: 消息提示系统分散在 HostsManager 中，难以维护和扩展
4. **体验不佳**: 模块标题不够突出，用户可能不清楚当前所在模块

## 变更内容

本提案的变更包括：

1. **新增 ModuleHeader 到 HostsManager**: 在模块顶部添加统一的 ModuleHeader 组件
2. **消息提示系统迁移**: 将 toast 消息提示系统替换为 ModuleHeader 的 hint 系统
3. **样式整合**: 确保消息提示样式与其他模块一致
4. **布局保持**: 两面板布局完全保持不变，不影响现有功能

## 背景

### 当前状况

- **HostsManager.vue**: 使用两面板布局，左侧侧栏，右侧编辑区
- **页头缺失**: 没有明确的模块标题区域（不同于 NodeManager 和 Settings）
- **消息提示**: 拥有自己的全局 toast 消息提示系统，位于模块顶部
- **工具栏**: 包含搜索、过滤、刷新等操作工具

### 痛点

1. **不一致性**: HostsManager 与其他模块（NodeManager、Settings）的页头样式不统一
2. **功能重复**: 消息提示系统可与 ModuleHeader 的提示框功能整合
3. **用户体验**: 模块标题不够突出，用户可能不清楚当前所在模块
4. **可维护性**: 独立的消息提示系统不易维护，与通用组件重复

## 提案范围

### 新增功能

1. **添加 ModuleHeader 到 HostsManager**
   - 在模块顶部显示标题（"Hosts 管理"）
   - 集成现有的全局消息提示系统（success/error）
   - 保持与 ModuleHeader 的样式一致性

### 修改功能

1. **HostsManager 消息提示集成**
   - 将现有的 `.message-toast` 消息提示迁移到 ModuleHeader 的 `hint` 系统
   - 保持原有的提示逻辑和用户体验

### 受影响的组件

- `frontend/src/views/HostsManager.vue` (重构消息提示系统)

## 实施方案要点

1. 在 HostsManager 的 `.module-layout` 之前添加 ModuleHeader 组件
2. 将 `message` 状态映射到 ModuleHeader 的 `hint` prop
3. 删除旧的 `.message-toast` 结构和相关样式
4. 保持现有的两面板布局不变
5. 确保消息提示功能完全兼容

## 验收标准

- [ ] ModuleHeader 在 HostsManager 正确显示
- [ ] 模块标题 "Hosts 管理" 清晰可见
- [ ] 消息提示（成功/错误）通过 ModuleHeader 正常显示
- [ ] 消息提示的自动消退、手动关闭功能正常工作
- [ ] 两面板布局不受影响
- [ ] 没有样式回归，与 NodeManager/Settings 风格一致

## 相关规范

- 前端组件规范：使用 Composition API + Vue 3 最佳实践
- 样式规范：使用 Tailwind CSS 和内联样式
- 国际化规范：使用 vue-i18n

## 依赖关系

- **依赖**: modularize-module-header 提案必须先完成
- **受影响**: 无其他模块直接依赖

## 后续步骤

1. 创建详细设计文档（design.md）
2. 编写规范增量（spec.md）
3. 编写任务列表（tasks.md）
4. 获得反馈和批准
5. 实施阶段开始
