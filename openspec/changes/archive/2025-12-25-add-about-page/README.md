# 关于页面变更提案

此目录包含 "add-about-page" 变更的完整提案文档。

## 文件结构

```
add-about-page/
├── proposal.md           # 变更提案概述
├── design.md            # 设计文档和实现细节
├── tasks.md             # 任务清单和验证标准
└── specs/               # 规范增量
    ├── version-info/
    │   └── spec.md      # 版本和应用信息规范
    └── sponsor-info/
        └── spec.md      # 赞赏功能规范
```

## 变更概述

**目标**：在设置模块页面添加应用信息展示，包括：
- 应用版本号
- 官方网址链接
- 作者信息
- 微信赞赏二维码

**影响范围**：
- 前端：Settings.vue 组件
- 后端：config/app.json 配置文件（扩展）
- 资源：赞赏二维码图片

**工作量**：约 9 个任务，预计 2-3 个开发周期

## 关键特性

✅ 最小化实现 - 直接在设置页面添加，无需新增页面
✅ 配置驱动 - 所有信息从配置文件读取
✅ 国际化支持 - 完全支持中英文
✅ 响应式设计 - 适配不同屏幕大小
✅ 向后兼容 - 新字段都是可选的

## 规范快速查看

### 版本信息规范 (version-info/spec.md)
- 显示应用名称和版本号
- 显示官方网址（可点击）
- 显示作者信息和邮箱
- 配置文件包含 version、website、author 字段

### 赞赏功能规范 (sponsor-info/spec.md)
- 显示微信赞赏二维码
- 二维码路径来自配置文件
- 赞赏提示文本国际化
- 支持本地文件、数据 URL、远程 URL

## 接下来的步骤

1. **提案审查**：审核 proposal.md 中的需求和设计理念
2. **规范讨论**：讨论两个规范增量的完整性和准确性
3. **任务评估**：评估任务清单的工作量和依赖关系
4. **批准**：正式批准此变更
5. **实施**：按照 tasks.md 进行实施（apply 阶段）
6. **验证**：按照验收标准进行验证

## 配置示例

扩展后的 config/app.json：

```json
{
  "version": "0.1.0",
  "website": "https://github.com/example/LiSteward",
  "sponsorQrCode": "./data/assets/wechat-qr.png",
  "productName": "李管家",
  "productNameEn": "LiSteward",
  "author": {
    "name": "lfz",
    "email": "lfzxs@qq.com"
  },
  ...
}
```

## 翻译示例

zh-CN.json：
```json
{
  "settings": {
    "about": "关于",
    "officialWebsite": "官方网站",
    "author": "作者",
    "sponsor": "微信赞赏",
    "sponsorTip": "扫描下方二维码，支持开发者继续维护此项目"
  }
}
```

en-US.json：
```json
{
  "settings": {
    "about": "About",
    "officialWebsite": "Official Website",
    "author": "Author",
    "sponsor": "WeChat Sponsor",
    "sponsorTip": "Scan the QR code below to support the developer"
  }
}
```

## 相关变更

此提案与现有变更的兼容性：
- ✅ add-user-config-system - 无冲突，配置系统兼容
- ✅ add-i18n-support - 完全兼容，使用相同 i18n 系统
- ✅ 现有功能 - 不破坏任何现有功能
