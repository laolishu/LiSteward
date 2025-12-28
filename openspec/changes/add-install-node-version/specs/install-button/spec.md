# 规范增量：Node 子模块安装按钮功能

## 新增需求

### 需求：NodeManager 前端必须在"Node 仓库"行添加"安装"按钮

为了让用户能够便捷地打开命令行查看可用的 Node 版本，前端必须提供一个专门的安装按钮。

#### 场景：按钮的放置和对齐
- "Node 仓库"行的文本框必须缩小宽度至约 70%
- "安装"按钮必须放置在文本框后方，占约 30% 宽度
- 按钮必须与上一行（"NVM 版本"行）的元素保持垂直对齐
- 两个元素之间必须有 12px 的间距

#### 场景：按钮的样式和交互
- 按钮必须使用蓝色背景（#1f6feb）和白色文字
- 按钮必须有圆角（border-radius: 6px）和 padding (8px 16px)
- 鼠标悬停时背景色必须变为 #388bfd
- 按钮必须支持 aria-label 属性以提高可访问性

#### 场景：按钮的事件处理
- 点击按钮必须调用 `installNewVersion()` 方法
- 该方法必须调用后端的 `OpenNodeAvailableList()` 接口
- 必须有适当的错误处理和日志记录

### 需求：NodeManager.vue 模板必须调整 item-control 容器布局

为了容纳新的安装按钮，item-control 容器的布局必须改进为 flexbox。

#### 场景：容器布局的灵活性
- item-control 必须使用 `display: flex` 和 `flex-direction: row`
- 容器必须有 `gap: 12px` 的间距
- 容器必须有 `align-items: center` 的垂直对齐

#### 场景：版本显示框的尺寸
- nvm-version-display 必须使用 `flex: 0 0 70%` 控制宽度
- 必须保持现有的样式和外观

### 需求：后端必须新增 OpenNodeAvailableList 接口用于打开系统命令行

后端必须提供新的 Wails 绑定接口，允许前端触发系统命令行窗口显示可用版本列表。

#### 场景：Wails 绑定接口
- App 必须新增 `OpenNodeAvailableList()` 方法
- 该方法必须调用 `nvm.OpenAvailableList()` 实现
- 方法必须返回 error 类型

#### 场景：前端调用的 JavaScript 绑定
- 前端必须能够导入并调用 `OpenNodeAvailableList()` 函数
- 绑定必须自动生成到 `frontend/wailsjs/go/main/App.js`

### 需求：internal/nvm 包必须实现跨平台的命令行打开功能

internal/nvm 包必须提供 OpenAvailableList 函数，在不同平台上正确打开系统命令行窗口。

#### 场景：Windows 平台的实现
- 必须使用 `cmd.exe /k "nvm list available"` 打开新窗口
- `/k` 参数必须确保窗口在命令执行后保持打开状态
- 用户必须能够看到完整的可用版本列表输出

#### 场景：Unix/Linux 平台的实现
- 必须检测系统可用的终端程序（x-terminal-emulator、gnome-terminal、konsole 等）
- 必须优先级依次尝试不同的终端程序
- 必须执行 `bash -c "nvm list available"` 来显示版本列表
- 命令必须在终端中保持可见，用户能够查看结果

#### 场景：macOS 平台的实现
- 必须使用系统默认终端或 iTerm2
- 必须能够执行 nvm 命令并显示结果
- 终端窗口必须不会在命令执行后立即关闭

#### 场景：错误处理
- 必须处理 NVM 未安装的情况（命令会失败）
- 必须处理终端程序不可用的情况
- 必须返回 error 以供前端处理

### 需求：i18n 配置必须添加新的安装按钮文本

应用的国际化配置必须包含安装按钮的文本。

#### 场景：中文本地化 (zh-CN)
- `node.installButton` 必须设置为 "安装"
- 文本必须添加到 `frontend/src/i18n/locales/zh-CN.json` 的 node 对象中

#### 场景：英文本地化 (en-US)
- `node.installButton` 必须设置为 "Install"
- 文本必须添加到 `frontend/src/i18n/locales/en-US.json` 的 node 对象中

## 修改需求

### 需求：NodeManager.vue 必须完全改进以支持新的安装按钮功能

NodeManager.vue 组件必须从模板到脚本都进行适当的修改以支持新功能。

#### 场景：导入 OpenNodeAvailableList 函数
- 必须在 script setup 部分添加 `OpenNodeAvailableList` 的导入
- 导入路径必须为 `../../wailsjs/go/main/App`

#### 场景：添加事件处理方法
- 必须在 setup 中定义 `installNewVersion()` 异步方法
- 该方法必须调用 `OpenNodeAvailableList()`
- 必须有 try-catch 错误处理

#### 场景：修改模板中的 Node 仓库行
- 必须将纯粹的 div 改为 flex 容器
- 必须在 div 后添加"安装"按钮
- 按钮必须绑定 `@click="installNewVersion"` 事件
- 按钮必须使用正确的 i18n 文本 `$t('node.installButton')`

## 移除需求

无

## 重命名需求

无

## 依赖关系

### 内部依赖
- 依赖于 Wails v2 绑定机制
- 依赖于 internal/nvm 包的现有功能
- 依赖于 Vue 3 Composition API
- 依赖于 vue-i18n

### 外部依赖
- nvm/nvm-windows 命令行工具（用户必须已安装）

### 与其他规范的关系
- 参考 "node-runtime" 规范：NVM 管理的相关接口
- 不依赖其他活跃提案

## 验证条件

- [ ] "安装"按钮正确显示在"Node 仓库"行的文本框后方
- [ ] 按钮与上一行（"NVM 版本"行）元素对齐
- [ ] 文本框宽度约 70%，按钮约 30%
- [ ] 按钮样式符合设计规范（蓝色背景、白色文字、圆角）
- [ ] 点击按钮时打开系统命令行窗口
- [ ] Windows 平台上成功执行 `nvm list available` 并显示结果
- [ ] macOS/Linux 平台也能正常打开命令行并显示结果
- [ ] 中文环境显示"安装"，英文环境显示"Install"
- [ ] 代码无编译错误和类型错误
- [ ] 前端 JavaScript 绑定文件自动更新
