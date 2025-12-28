# 规范增量：HostsManager 模块 ModuleHeader 集成

## 新增需求

### 需求：HostsManager 必须添加 ModuleHeader 组件用于统一的页面标题显示

为了保持应用界面的一致性并统一页面标题的呈现方式，HostsManager 模块必须采用 ModuleHeader 组件来显示页面标题。

#### 场景：模块标题显示
- HostsManager 组件必须在其主容器顶部添加 ModuleHeader 组件
- ModuleHeader 组件必须显示"Hosts 管理"的本地化文本（i18n key: `hosts.title`）
- 标题显示样式必须与 NodeManager、Settings 模块保持一致
- 标题颜色必须为 #24292e，字号 24px，字重 600

### 需求：HostsManager 必须将消息提示系统完全集成到 ModuleHeader 的 hint 属性

原有的独立消息提示系统必须改进为通过 ModuleHeader 的 hint 属性统一显示，以简化代码并保证视觉风格的一致性。

#### 场景：成功消息显示与自动消退
- HostsManager 执行操作成功时必须通过 ModuleHeader 的 `hint` 属性显示成功消息
- 成功消息的 type 属性必须设置为 'success'
- 成功消息显示的样式必须使用设计规范的颜色（背景: #f0f9ff，文字: #0369a1，边框: #bfdbfe）
- 成功消息必须在 2 秒后自动消退清除

#### 场景：错误消息显示与持久化
- HostsManager 执行操作失败时必须通过 ModuleHeader 的 `hint` 属性显示错误消息
- 错误消息的 type 属性必须设置为 'error'
- 错误消息显示的样式必须使用设计规范的颜色（背景: #fef2f2，文字: #991b1b，边框: #fecaca）
- 错误消息禁止自动消退，必须保持显示直到用户交互或执行新操作

### 需求：HostsManager 必须实现 showMessage 方法用于统一管理所有消息提示

为了统一消息提示的逻辑和减少重复代码，HostsManager 必须提供一个统一的消息管理方法。

#### 场景：消息方法实现与使用
- HostsManager 组件必须添加 `messageHint` 响应式对象，初始值为 null
- HostsManager 组件必须实现 `showMessage(text, type = 'success', duration = 2000)` 方法
- 该方法必须将消息数据设置到 messageHint 对象（结构: `{ type, message: text }`）
- 对于 success 类型，该方法必须在指定的 duration 毫秒后将 messageHint 清除为 null
- 对于 error 类型，该方法禁止自动清除消息，必须保留直到下次操作

### 需求：HostsManager 布局结构必须调整以容纳 ModuleHeader 组件

为了为 ModuleHeader 组件留出合适的位置，容器布局必须进行调整，确保模块标题与内容区域的正确分层。

#### 场景：容器布局和样式
- 容器 `.hosts-manager` 必须使用 flex 列方向布局
- ModuleHeader 组件必须作为第一个子元素
- `.module-layout` 必须作为第二个子元素并设置 flex: 1 自适应高度
- 两面板布局（侧栏和编辑区）禁止变动，必须完全保持不变

### 需求：HostsManager 必须删除所有与旧消息提示系统相关的样式代码

旧的消息提示系统将被 ModuleHeader 的提示功能所替代，相关的样式定义必须清理以避免样式冲突和减少代码冗余。

#### 场景：样式清理与整理
- 必须删除 `.message-toast` 及所有相关的变体样式定义
- 必须删除 `.message-fade` 转场动画样式
- 必须删除 `.message-icon` 样式（已由 ModuleHeader 完全替代）
- 删除后必须确保 ModuleHeader 的样式能够正确应用，禁止存在冲突

## 修改需求

### 需求：HostsManager.vue 组件必须修改以支持 ModuleHeader 的导入和完整使用

为了使 HostsManager 能够使用 ModuleHeader 组件并与新的消息管理系统集成，必须对其源文件进行相应的修改。

#### 场景：导入 ModuleHeader 组件
- HostsManager.vue 必须在 script 部分添加 ModuleHeader 组件的导入语句
- 导入路径必须为 '../components/ModuleHeader.vue'
- 组件必须在 components 对象中声明

#### 场景：修改模板结构
- 模板必须在 `.hosts-manager` 容器中首先添加 `<ModuleHeader :title="$t('hosts.title')" :hint="messageHint" />`
- 必须删除原有的 `<transition name="message-fade">` 及其内部的 message-toast div 结构
- `.module-layout` 必须保持紧跟在 ModuleHeader 之后的位置

#### 场景：添加响应式数据和方法
- 必须在 setup 函数中添加 `const messageHint = ref(null)` 数据
- 必须实现 `const showMessage = (text, type = 'success', duration = 2000) => { ... }` 方法
- showMessage 方法必须设置 messageHint.value，并在需要时设置自动清除计时器

#### 场景：迁移所有消息调用
- 必须查找并替换所有直接设置 message.value 的代码位置
- 所有消息显示必须改为调用 showMessage() 方法
- 主要包括：applyChanges() 的成功/失败提示、reloadHosts() 的完成提示

## 重命名需求

无需求

## 移除需求

### 需求：HostsManager.vue 必须删除旧的消息提示系统的 HTML 和样式代码

#### 场景：删除 HTML 转场结构
- 必须删除 `<transition name="message-fade">` 标签及其完整内容
- 必须删除其内部的 `.message-toast` div 和相关子元素
- 删除后模板结构必须只保留 ModuleHeader 和 module-layout 两个主要部分

## 依赖关系

### 内部依赖
- ModuleHeader 组件必须先在 modularize-module-header 提案中实现完成
- Vue 3 Composition API
- vue-i18n

### 外部依赖
- 无新增外部依赖
- 依赖于 ModuleHeader 组件的完整实现（包括提示框功能）

### 与其他规范的关系
- 参考 "modularize-module-header" 规范：ModuleHeader 组件的 props、事件、样式定义
- 与 i18n 规范协同：确保 `hosts.title` key 的正确定义

## 验证条件

- [ ] ModuleHeader 在 HostsManager 顶部正确显示
- [ ] 模块标题 "Hosts 管理" 或对应的国际化文本显示正确
- [ ] 保存成功时显示成功提示，样式符合设计
- [ ] 保存失败时显示错误提示，样式符合设计
- [ ] 重新加载成功时显示相应提示
- [ ] 成功提示在 2 秒后自动消退
- [ ] 错误提示可手动关闭或长期保留
- [ ] 两面板布局（侧栏 + 编辑区）完全保持不变，无样式或布局破坏
- [ ] 所有原有的消息提示场景都通过 ModuleHeader 正常工作
- [ ] 国际化文本正确显示，中文和英文都无缺失
- [ ] 删除的 CSS 样式不影响其他部分
- [ ] 没有样式回归，与 NodeManager、Settings 风格一致

## 相关工作项

- 更新 i18n 配置，添加 `hosts.title` key（如需要）
- 可选：为其他操作（删除、启用/禁用等）添加提示
- 可选：在 ModuleHeader 中显示当前选中的配置文件名
