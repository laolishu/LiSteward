# 规范增量：ModuleHeader 组件

## 新增需求

### 需求：ModuleHeader 组件必须提供基础的页头结构和标题显示功能

ModuleHeader 是一个可复用的页头组件，必须用于统一应用各模块的页面标题显示风格。

#### 场景：标题属性接收和显示
- 组件必须接收 `title` 属性（必需），类型为 String
- 组件必须在页头区域以大号字体（24px，字重 600）显示标题
- 标题颜色必须为 #24292e（深灰色）

#### 场景：副标题的可选显示
- 组件必须接收 `subtitle` 属性（可选），类型为 String
- 组件必须在标题下方显示副标题，字体大小 14px，颜色 #586069
- 当 `subtitle` 未提供时，禁止显示副标题区域

#### 场景：页头容器的样式设置
- 页头区域必须有白色背景 (#ffffff)
- 下方必须有 1px 实心边框，颜色为 #e1e4e8
- 内边距必须为上下 24px，左右 24px
- 布局必须使用 flexbox 垂直布局

### 需求：ModuleHeader 组件必须支持多种类型的消息提示显示

ModuleHeader 必须能够接收和显示不同类型的提示信息，包括成功、错误、警告和信息类型。

#### 场景：成功提示的显示
- 组件必须接收 `hint` 属性（可选），结构为 `{ type: 'success', message: '成功消息' }`
- 当 `hint.type === 'success'` 时，必须在页头下方显示绿色背景提示框
- 背景色必须为 #f0f9ff，文字色必须为 #0369a1，边框必须为 1px solid #bfdbfe
- 必须显示对应的图标 (message-success.svg)

#### 场景：错误提示的显示
- 当 `hint.type === 'error'` 时，必须在页头下方显示红色背景提示框
- 背景色必须为 #fef2f2，文字色必须为 #991b1b，边框必须为 1px solid #fecaca
- 必须显示对应的图标 (message-error.svg)

#### 场景：警告和信息提示的支持
- 必须支持 `hint.type === 'warning'`，黄色背景 (#fffbeb)，文字色 #92400e
- 必须支持 `hint.type === 'info'`，浅蓝背景 (#f0f9ff)，文字色 #0369a1
- 各类型提示禁止缺少对应的图标

#### 场景：提示的显示/隐藏控制
- 当 `showHint` 属性为 false 时，必须不显示任何提示
- 默认值 `showHint` 必须为 true
- 提示框必须包含关闭按钮 (×)，用户点击可手动关闭

#### 场景：提示的自动消退
- 成功和信息提示必须在显示 5 秒后自动消失
- 错误和警告提示禁止自动消退，需用户手动关闭或长期保留
- 用户手动关闭提示后，计时器必须停止

### 需求：ModuleHeader 组件必须具有清晰的视觉布局和过渡效果

ModuleHeader 的布局必须清晰分层，样式必须符合设计规范，过渡必须平滑。

#### 场景：整体页头布局结构
- 页头必须包含两部分: 标题区和提示区
- 标题区必须居顶部，高度自适应
- 提示区必须在标题下方，禁止在 hint 不存在时显示
- 整个页头区域宽度必须为 100%，必须与容器对齐

#### 场景：提示框的样式设置
- 提示框宽度必须为 100%
- 高度必须自适应，内边距必须为左右 16px，上下 12px
- 图标和文字之间的间距必须为 12px
- 提示文字必须使用 14px 字号

#### 场景：平滑的转场效果
- 提示框的显示/隐藏必须使用 CSS transition
- 必须实现淡入淡出效果，过渡时间必须为 0.3s
- 提示框消退时必须具有平滑动画效果

### 需求：ModuleHeader 组件必须支持国际化集成

ModuleHeader 组件必须正确处理国际化需求，支持多语言显示。

#### 场景：标题属性的国际化支持
- `title` 和 `subtitle` 属性必须支持直接传入中文字符串
- 组件必须支持传入 i18n key，由父组件使用 `$t()` 翻译后传入
- 组件本身禁止执行 i18n 翻译，必须保持父子职责分离

#### 场景：内置文本的国际化处理
- 关闭按钮的 aria-label 必须支持国际化
- 提示框内容必须由父组件控制，必须支持 i18n

### 需求：ModuleHeader 组件必须具有良好的可访问性支持

ModuleHeader 必须遵循 Web 可访问性标准，确保屏幕阅读器用户能够正确使用。

#### 场景：语义化标签的使用
- 标题必须使用 `<h1>` 或 `<h2>` 标签（根据嵌套层级）
- 提示框必须使用 `role="alert"` 或 `role="status"` 属性

#### 场景：ARIA 属性的完整配置
- 关闭按钮必须包含 `aria-label="关闭提示"`
- 提示框必须包含 `aria-live="polite"` 或 `aria-live="assertive"`（根据类型）
- 图片必须使用 `alt` 属性

## 修改需求

### 需求：NodeManager.vue 必须修改以支持 ModuleHeader 的集成

NodeManager 模块必须使用 ModuleHeader 组件来统一其页面标题显示。

#### 场景：替换现有页头
- NodeManager.vue 必须删除 `.settings-header` div 和相关样式
- 必须使用 `<ModuleHeader :title="$t('node.title')" />` 替换
- 必须保持原有的 `#f6f8fa` 背景容器

#### 场景：错误提示的集成显示
- 若 NodeManager 中存在错误提示，必须通过 `hint` 属性传入 ModuleHeader
- 示例: `:hint="{ type: 'error', message: nvmError }"`

### 需求：Settings.vue 必须修改以支持 ModuleHeader 的集成

Settings 模块必须使用 ModuleHeader 组件来统一其页面标题显示。

#### 场景：替换页头结构
- Settings.vue 必须删除 `.settings-header` div 和相关样式
- 必须使用 `<ModuleHeader :title="$t('settings.title')" />` 替换

#### 场景：语言切换提示的显示
- 语言切换成功时，必须通过 `hint` 属性显示成功提示
- 示例: `:hint="languageChanged ? { type: 'success', message: $t('settings.languageChanged') } : null"`

### 需求：HostsManager.vue 必须添加 ModuleHeader 组件用于页面标题显示

HostsManager 模块必须添加 ModuleHeader 组件来显示模块标题和管理消息提示。

#### 场景：添加 ModuleHeader 页头
- HostsManager.vue 必须添加 ModuleHeader，显示 hosts 管理模块标题
- 若 HostsManager 中有全局消息提示，必须集成到 ModuleHeader

## 重命名需求

无

## 移除需求

无（旧页头样式随 .vue 文件重构而逐步移除）

## 依赖关系

### 内部依赖

- Vue 3 Composition API: 组件的基础框架
- vue-i18n: 国际化支持（父组件侧）

### 外部依赖

- 无新增外部依赖

### 与其他规范的关系

- 与 `i18n-core` 规范协同: 确保 title 和 subtitle 的 i18n key 正确
- 与 UI 风格规范协同: 样式保持与现有 LiSteward 设计一致

## 验证条件

- [ ] ModuleHeader 组件正确渲染所有部分（title、subtitle、hint）
- [ ] 提示框样式与设计文档一致
- [ ] 各提示类型 (success/error/warning/info) 样式正确
- [ ] 自动消退逻辑正确运行
- [ ] 用户手动关闭提示时，计时器停止
- [ ] NodeManager.vue 集成后，页头效果与重构前相同
- [ ] Settings.vue 集成后，页头和消息提示正常工作
- [ ] HostsManager.vue 集成后，如有全局消息，正确显示
- [ ] 国际化文本正确显示，无缺失
- [ ] 可访问性属性 (aria-label, role) 正确设置
- [ ] 没有样式回归
