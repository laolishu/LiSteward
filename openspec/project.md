# 项目 上下文

## 目的
LiSteward 是一个基于 Wails 框架的桌面应用程序，使用 Golang 后端和 Vue3 前端构建。

## 技术栈

### 后端
- **Golang** (版本 1.23)
- **Wails v2** (v2.11.0) - Go 桌面应用框架
- **依赖管理**：Go Modules

### 前端
- **Vue 3** (Composition API)
- **Vite** - 构建工具
- **Tailwind CSS** - UI 样式框架
- **开发服务器**：http://localhost:34115

### 运行时环境
- **操作系统**：Windows
- **权限要求**：Administrator 权限（强制要求）
- **Webview**：go-webview2 (Windows WebView2)

## 项目约定

### 代码风格

#### 后端 (Golang)
- 遵循 Go 标准代码风格
- 使用 `gofmt` 格式化代码
- 遵循 Go 命名约定：
  - 导出函数/变量：大写开头 (PascalCase)
  - 私有函数/变量：小写开头 (camelCase)
  - 包名：小写单词

#### 前端 (Vue3)
- 使用 **Composition API**（不使用 Options API）
- 组件命名：PascalCase
- 使用 Tailwind CSS 进行样式设计
- TypeScript 可选，但建议使用类型注释

### 架构模式

#### 应用结构
```
LiSteward/
├── app.go              # Wails 应用程序绑定
├── main.go             # 入口点
├── frontend/           # Vue3 前端代码
│   ├── src/
│   └── package.json
├── build/              # 构建输出
└── wails.json          # Wails 项目配置
```

#### 后端架构
- Wails 绑定层：Go 方法通过 Wails 暴露给前端
- 服务层：业务逻辑处理
- 数据层：文件系统、配置管理

#### 前端架构
- 组件化：使用 Vue 组件
- 状态管理：根据需要使用 Vue 响应式 API
- 与后端通信：通过 Wails 运行时调用 Go 方法

### 开发流程

#### 本地开发
```bash
# 运行开发模式（热重载）
wails dev

# 浏览器调试
# 访问 http://localhost:34115 调用 Go 方法
```

#### 构建发布
```bash
# 构建生产版本
wails build
```

### 测试策略
- **后端测试**：使用 Go 标准测试库 (`testing`)
- **前端测试**：根据需要添加 Vitest/Vue Test Utils
- **集成测试**：测试 Wails 绑定和前后端交互

### Git工作流
- **主分支**：`main` - 稳定版本
- **功能分支**：`feature/<功能名>` 或 `add-<功能名>`
- **修复分支**：`fix/<问题描述>`
- **提交信息**：使用中文，清晰描述变更内容

## 领域上下文

### Windows 管理工具
- LiSteward 需要 Windows Administrator 权限运行
- 可能涉及系统级操作和配置
- 需要处理 Windows 特定的 API 和服务

### Wails 特性
- Go 后端方法自动绑定到前端 JavaScript
- 通过 `runtime` 包访问系统功能
- 事件系统：后端和前端之间的事件通信

## 重要约束

### 权限约束
- **必须**以 Windows Administrator 权限运行
- 需要在代码中处理权限检查和提升

### 平台约束
- **仅支持 Windows**（使用 WebView2）
- 需要用户安装 WebView2 运行时（或应用内嵌）

### 技术约束
- Wails v2 框架限制
- Go 1.23+ 要求
- Vue 3 Composition API 约定

## 外部依赖

### 核心框架
- **Wails v2.11.0**：应用框架
- **WebView2**：Windows WebView 渲染引擎

### Go 依赖
- `github.com/wailsapp/wails/v2` - 主框架
- `github.com/wailsapp/go-webview2` - WebView2 绑定
- 其他工具库（详见 go.mod）

### 前端依赖
- **Vue 3** - 前端框架
- **Vite** - 构建工具和开发服务器
- **Tailwind CSS** - 样式框架

### 开发工具
- Wails CLI：项目管理和构建
- Go 工具链：编译和测试
- Node.js/npm：前端依赖管理
