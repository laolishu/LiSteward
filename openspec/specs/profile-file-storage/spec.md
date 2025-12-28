# profile-file-storage Specification

## Purpose
TBD - created by archiving change add-user-config-system. Update Purpose after archive.
## 需求
### 需求：单文件方案存储

Hosts 配置方案必须持久化存储到文件系统，每个方案对应一个独立的 JSON 文件。

#### 场景：保存 Hosts 方案

当用户创建或修改 Hosts 配置方案时，应用必须将该方案保存到 `%APPDATA%/LiSteward/data/profiles/hosts/` 目录下，文件名格式为 `host_<profileName>.json`。

文件内容遵循 HostsProfile 结构：

```json
{
  "name": "work",
  "entries": [
    {
      "ip": "127.0.0.1",
      "domain": "example.com",
      "comment": "Production server",
      "alternate_ip": "",
      "description": "Development environment",
      "enabled": true
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 场景：方案名称与文件名映射

方案名称转换为文件名时应遵循以下规则：

- 将方案名称中的非法字符（除字母、数字、下划线、连字符外）替换为下划线
- 文件名格式：`host_<normalized_name>.json`
- 示例映射：
  - 方案名 "work" → 文件名 "host_work.json"
  - 方案名 "dev server" → 文件名 "host_dev_server.json"
  - 方案名 "test-v1.0" → 文件名 "host_test-v1.0.json"

### 需求：方案加载

应用启动时必须自动加载 `%APPDATA%/LiSteward/data/profiles/hosts/` 目录中所有 `host_*.json` 文件作为可用的 Hosts 方案。

#### 场景：应用启动加载方案

应用启动时必须执行以下步骤：

1. 检查 `%APPDATA%/LiSteward/data/profiles/hosts/` 目录是否存在
2. 如果存在，扫描目录中所有 `host_*.json` 文件
3. 逐个解析文件内容，恢复 HostsProfile 对象
4. 将所有方案加载到内存中供前端使用

#### 场景：无效方案文件处理

如果 `host_*.json` 文件格式无效或解析失败：

- 记录错误日志
- 跳过该文件，继续加载其他方案
- 不中断应用启动流程

### 需求：方案删除

用户删除 Hosts 方案时，必须同时删除对应的文件。

#### 场景：删除方案

当用户删除一个 Hosts 方案时，应用必须删除对应的 `host_<profileName>.json` 文件。

### 需求：方案更新

用户修改 Hosts 方案时，必须更新对应的文件内容。

#### 场景：修改方案

当用户修改一个 Hosts 方案时，应用必须更新 `host_<profileName>.json` 文件中的内容（包括修改 `updated_at` 时间戳）。

