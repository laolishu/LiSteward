# 设计文档：用户配置系统

## 架构概述

用户配置系统分为两个主要部分：

### 1. 用户配置管理（后端 Go）

在 `config` 包中扩展用户配置管理：

```go
// UserConfig 代表用户偏好设置
type UserConfig struct {
    Language string `json:"language"` // 用户选择的语言: zh-CN, en-US
    // 为未来功能预留空间
}

// LoadUserConfig 从 config/user.json 加载用户配置
func LoadUserConfig() (*UserConfig, error)

// SaveUserConfig 保存用户配置到 config/user.json
func SaveUserConfig(uc *UserConfig) error

// EnsureUserConfig 确保用户配置文件存在，如不存在则创建默认配置
func EnsureUserConfig() error
```

### 2. Hosts 方案文件管理（后端 Go）

在 `internal/hosts` 的 `Service` 中扩展方案文件管理：

```go
// SaveProfile 将配置方案保存到文件
// 文件名格式: host_<profileName>.json
func (s *Service) SaveProfile(profile *HostsProfile) error

// LoadProfile 从文件加载单个方案
func (s *Service) LoadProfile(name string) (*HostsProfile, error)

// LoadAllProfiles 加载目录中所有方案
func (s *Service) LoadAllProfiles() ([]*HostsProfile, error)

// DeleteProfile 删除方案文件
func (s *Service) DeleteProfile(name string) error
```

## 目录结构设计

### Windows AppData 布局

```
%APPDATA%/LiSteward/
├── config/
│   ├── app.json          # 应用配置（只读，打包内置）
│   └── user.json         # 用户配置（读写）
└── data/
    ├── backups/
    │   └── hosts/        # 已有备份文件存储
    └── profiles/
        └── hosts/        # Hosts 方案文件存储（新增）
```

### 文件命名约定

- **用户配置**：`user.json`
- **Hosts 方案**：`host_<profileName>.json`
  - 有效字符：字母、数字、下划线、连字符
  - 示例：`host_work.json`、`host_dev_server.json`

## 初始化流程

应用启动时：

1. 检查 `%APPDATA%/LiSteward/config/user.json` 是否存在
   - 如果存在，加载用户配置
   - 如果不存在，创建默认用户配置
2. 检查 `%APPDATA%/LiSteward/config/app.json` 中是否有旧的 `language` 字段
   - 如果有，迁移到 `user.json`
3. 加载所有 Hosts 方案文件

## 数据持久化

### 保存时机

- **用户配置**：用户在设置中修改语言或其他偏好时立即保存
- **Hosts 方案**：
  - 创建或修改方案时立即保存
  - 删除方案时删除对应的文件

### 文件格式

所有配置文件均使用 JSON 格式，便于编辑和跨平台兼容。

## 扩展性考虑

- 预留用户配置中的其他字段空间（如 `theme`、`autoSync` 等）
- 方案文件采用单文件单方案设计，便于未来支持更大规模方案库
- 目录结构分离 `config` 和 `data`，清晰区分配置和数据

## 后向兼容性

- 应用启动时检查旧配置位置
- 自动迁移存在的旧配置
- 保留旧字段读取能力以支持过渡期
