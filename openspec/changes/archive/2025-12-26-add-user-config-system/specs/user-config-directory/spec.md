# 规范：用户配置目录

## 新增需求

### 需求：用户配置目录标准化

应用必须在 `%APPDATA%/LiSteward/` 下建立标准化的用户配置和数据目录结构。

#### 场景：应用首次启动

当应用首次启动时，必须自动创建以下目录（如不存在）：

- `%APPDATA%/LiSteward/config/` - 用户配置目录
- `%APPDATA%/LiSteward/data/backups/hosts/` - Hosts 文件备份目录
- `%APPDATA%/LiSteward/data/profiles/hosts/` - Hosts 方案目录

#### 场景：配置目录存在

当上述目录已存在时，应直接使用现有目录，无需重新创建。

### 需求：用户配置文件管理

应用必须在 `%APPDATA%/LiSteward/config/` 下管理 `user.json` 用户配置文件。

#### 场景：用户配置不存在

如果 `config/user.json` 不存在，应用必须创建默认用户配置文件，初始内容为：

```json
{
  "language": "zh-CN"
}
```

#### 场景：用户配置已存在

如果 `config/user.json` 已存在，应用必须读取该文件并应用其中的配置。

#### 场景：读取语言偏好

应用必须能正确读取 `config/user.json` 中的 `language` 字段，并按以下优先级使用：

1. 用户在 `config/user.json` 中的设置
2. 系统语言检测（Windows 系统语言）
3. 默认语言（zh-CN）

### 需求：配置文件写入

用户在设置中修改语言偏好时，应用必须将新的语言设置写入到 `config/user.json` 文件中。

#### 场景：保存语言偏好

用户在设置界面选择语言后，应用必须将选择的语言值写入 `config/user.json` 的 `language` 字段中，并立即应用此配置。

### 需求：旧配置迁移

应用必须支持从旧的 `config/app.json` 中的 `language` 字段迁移到新的 `user.json`。

#### 场景：第一次使用新配置系统

当应用第一次运行新配置系统时，必须检查 `config/app.json` 中是否存在 `language` 字段：

- 如果存在，将其值迁移到 `config/user.json`
- 如果 `config/app.json` 中没有 `language` 字段，使用默认值（zh-CN）

## 修改需求

无

## 移除需求

无

## 重命名需求

无
