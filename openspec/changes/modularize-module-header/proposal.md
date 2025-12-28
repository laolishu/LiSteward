<!--
 * @Descripttion: 
 * @version: 
 * @Author: lfzxs@qq.com
 * @Date: 2025-12-27 11:40:42
 * @LastEditors: lfzxs@qq.com
 * @LastEditTime: 2025-12-27 13:01:44
-->
# 提案：模块页头组件化

**变更ID**: modularize-module-header

**状态**: 已实施完成

**优先级**: 中

**创建日期**: 2025-12-27

## 摘要

当前各模块视图（HostsManager、NodeManager、Settings 等）都有重复的页头结构，包含模块标题（Title）、子标题（SubTitle）和提示信息（Hint）。本提案通过创建统一的 `ModuleHeader` 组件来消除这种重复，提高代码可维护性和样式一致性。

## 背景

### 当前状况

- **NodeManager.vue**: 包含 `settings-header` 和 `h1` 标签
- **HostsManager.vue**: 包含消息提示 toast
- **Settings.vue**: 包含 `settings-header` 和 `h1` 标签
- 样式分散在各个组件的 `<style>` 块中，存在重复

### 痛点

1. **代码重复**: 三个模块都有相似的页头结构
2. **样式分散**: 页头相关样式在多个文件中重复定义
3. **维护困难**: 修改页头样式或结构需要同时更新多个文件
4. **不一致性**: 各模块的页头实现细节不统一

## 提案范围

### 新增功能

1. **ModuleHeader 组件**: 统一的模块页头组件
   - 支持 Title（模块名）
   - 支持 SubTitle（模块子标题）
   - 支持 Hint（各类提示，如错误/成功/警告）
   - 集中管理样式

### 受影响的组件

- `frontend/src/components/ModuleHeader.vue` (新建)
- `frontend/src/views/NodeManager.vue` (重构)
- `frontend/src/views/HostsManager.vue` (重构，如需要)
- `frontend/src/views/Settings.vue` (重构)

## 实施方案要点

1. 创建可复用的 `ModuleHeader` 组件
2. 支持通过 props 配置 title、subtitle、hint 等内容
3. 集成国际化支持（i18n）
4. 整合消息提示功能
5. 逐步在各模块中替换旧的页头实现

## 验收标准

- [ ] ModuleHeader 组件正确渲染所有部分
- [ ] 各模块使用新组件后样式保持一致
- [ ] 国际化文本正确显示
- [ ] 消息提示功能正常工作
- [ ] 没有样式回归

## 相关规范

- 前端组件规范：使用 Composition API + Vue 3 最佳实践
- 样式规范：使用 Tailwind CSS 和内联样式
- 国际化规范：使用 vue-i18n

## 后续步骤

1. 创建详细设计文档（design.md）
2. 编写规范增量（spec.md）
3. 编写任务列表（tasks.md）
4. 获得反馈和批准
5. 实施阶段开始
