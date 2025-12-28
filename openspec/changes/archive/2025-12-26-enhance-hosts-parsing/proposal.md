# 变更：增强 Hosts 文件解析方案

## 为什么

当前的 Hosts 管理系统只支持单一的 IP 和注释。在实际使用中，用户经常需要：

- **多 IP 支持**：为同一个域名维护主 IP 和备选 IP（如主生产 IP 和备用 IP）
- **灵活切换**：快速在主 IP 和备选 IP 之间切换，而无需编辑
- **扩展注释**：在保持 Hosts 文件标准格式的基础上，支持结构化的元数据

例如开发者可能需要：
- 在 `127.0.0.1`（开发本地） 和 `192.168.1.100`（开发机器） 之间切换
- 保留域名的历史 IP 地址以供快速恢复
- 为不同的 IP 配置保存备注说明

当前系统无法满足这些需求，导致用户必须手动编辑 Hosts 文件。

## 变更内容

### 核心功能

#### 1. 增强 Hosts 文件格式支持

引入扩展的 Hosts 格式，在注释部分支持结构化数据：

```
IP    DOMAIN    # ALTERNATE_IP#DESCRIPTION
```

示例：
- 标准格式：`127.0.0.1 www.example.com # 本地测试`
- 扩展格式：`127.0.0.1 www.example.com # 192.168.1.100#生产环境备选 IP`
- 纯 IP：`127.0.0.1 www.example.com # 192.168.1.100#`（备注为空）
- 仅备注：`127.0.0.1 www.example.com # #生产备注`（备选 IP 为空）

#### 2. 数据模型扩展

在 `HostEntry` 结构体中添加：
- `AlternateIP`: 备选 IP 地址（可选，可为 null/空）
- `Description`: 详细描述（可选，可为 null/空）

#### 3. 前端交互

在 Hosts 管理界面添加新功能：
- **显示备选 IP**：在表格中显示备选 IP 和描述
- **交换 IP 操作**：一键交换主 IP 和备选 IP
- **编辑弹窗**：支持编辑备选 IP 和描述

### 向后兼容性

- 现有的纯注释格式继续有效：`127.0.0.1 domain # 注释` 仍被解析为注释
- 如果未提供备选 IP 或描述，显示为空白而不是 "null"
- 旧格式的 Hosts 文件不需要迁移

### 用户体验

- 最小化 UI 变更，仅在表格中添加一列或在编辑弹窗中添加输入框
- 交换 IP 操作显示为单个按钮或快捷操作
- 验证备选 IP 格式，确保其为有效的 IPv4/IPv6 地址或空

## 影响

### 受影响规范

- **修改**：`specs/hosts-management/` - 增强现有的 Hosts 文件读写规范

### 受影响代码

#### 后端

- `internal/hosts/models.go` - 扩展 `HostEntry` 结构体
- `internal/hosts/service.go` - 修改 `parseLine()` 方法以支持新格式
- `internal/hosts/service.go` - 修改 `generateHostsContent()` 方法以支持新格式
- `app.go` - 可能需要添加交换 IP 的 Wails 绑定方法

#### 前端

- `frontend/src/views/HostsManager.vue` - 修改表格以显示备选 IP 和描述
- `frontend/src/views/HostsManager.vue` - 修改编辑弹窗逻辑以支持新字段
- `frontend/src/views/HostsManager.vue` - 添加交换 IP 按钮和逻辑

### 测试覆盖

- 单元测试：解析新格式的 Hosts 行
- 单元测试：生成新格式的 Hosts 行
- 集成测试：从文件读取并序列化新格式
- 前端测试：编辑、交换 IP 等交互逻辑

### 风险评估

- **低风险**：
  - 解析逻辑添加到现有代码，不影响现有功能
  - 向后兼容现有格式
  - Hosts 文件标准格式保持不变，仅在注释部分扩展
  
- **需要注意**：
  - 确保备选 IP 验证正确（不能是无效地址）
  - 生成 Hosts 文件时正确转义特殊字符
  - 迁移旧版本时，旧数据正常加载且不丢失

## 依赖关系

- 无新的外部依赖
- 基于现有的 Hosts 管理模块功能
