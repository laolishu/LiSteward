LiSteward — 开发环境配置切换工具
================================

概述
LiSteward 是一款轻量级的开发环境配置切换工具，提供可视化的 `hosts` 管理与 NVM（Node 版本）配置管理。未来计划接入 AI 功能以提供智能建议与自动化操作。
主要功能
![赞赏](#file:wechat-qr.png)
- NVM / 配置文件可视化与切换
- 本地备份、导入/导出配置
- 跨平台桌面应用（Windows / macOS / Linux）

快速开始
1. 仓库根目录包含构建和打包脚本，查看 `frontend/` 与 `build/` 目录了解详细步骤。
2. 按目标平台运行对应的构建脚本生成安装包或可执行文件。

构建与运行（开发）
- 启动前端开发服务器并运行 Go 后端以便快速迭代：

```bash
# 终端 A（前端开发服务器）
cd frontend
npm install
npm run dev

# 终端 B（后端）
cd ..
go run .
```

构建（发布）
- 生成前端生产资源并使用 `wails` 打包原生应用：

```bash
cd frontend
npm install
npm run build

# 在仓库根目录执行
wails build
# 或仅构建后端可执行文件
go build -o LiSteward
```

许可
----
本项目采用 GNU 通用公共许可证 v3.0（GPL-3.0）。详见 `LICENSE` 文件。

赞赏
----
如果本项目对你有帮助，欢迎用小额赞助支持我们：
![赞赏](/res/wechat-qr.png)

参与贡献
----
欢迎通过 GitHub 提交 Issue 或 Pull Request。请保持改动聚焦，并附带必要的说明或验证步骤。
