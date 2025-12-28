# 技术设计文档

## 概述

本文档详细说明了增强 Hosts 文件解析方案的技术实现细节。

## 数据模型设计

### 当前模型

```go
type HostEntry struct {
    IP      string `json:"ip"`      // IP 地址
    Domain  string `json:"domain"`  // 域名
    Comment string `json:"comment"` // 注释
    Enabled bool   `json:"enabled"` // 是否启用
}
```

### 增强后的模型

```go
type HostEntry struct {
    IP           string `json:"ip"`            // 主 IP 地址
    Domain       string `json:"domain"`        // 域名
    Comment      string `json:"comment"`       // 注释（向后兼容）
    AlternateIP  string `json:"alternate_ip"`  // 备选 IP 地址（可选）
    Description  string `json:"description"`  // 详细描述（可选）
    Enabled      bool   `json:"enabled"`      // 是否启用
}
```

### 设计理由

1. **向后兼容**：保持现有字段不变，添加新字段
2. **JSON 序列化**：使用 json 标签便于方案持久化
3. **灵活性**：新字段为可选，允许为空字符串

## 文件格式规范

### 格式定义

```
IP    DOMAIN    # [ALTERNATE_IP]#[DESCRIPTION]
```

### 解析规则

1. **基本结构**：`IP DOMAIN # ALTERNATE_IP#DESCRIPTION`
2. **备选 IP**：第一个 `#` 后至下一个 `#` 的内容
3. **描述**：第二个 `#` 后的所有内容
4. **空值处理**：如果为 `null` 或缺失，在显示时作为空字符串处理
5. **向后兼容**：如果 `#` 后只有注释文本（无第二个 `#`），作为旧格式处理，将其放入 `Comment` 字段

### 示例

| 格式 | 解析结果 | 说明 |
|------|--------|------|
| `127.0.0.1 www.test # 192.168.1.1#生产` | IP=127.0.0.1, AlternateIP=192.168.1.1, Description=生产 | 新格式完整 |
| `127.0.0.1 www.test # 192.168.1.1#` | IP=127.0.0.1, AlternateIP=192.168.1.1, Description=空 | 有备选 IP，无描述 |
| `127.0.0.1 www.test # #生产备注` | IP=127.0.0.1, AlternateIP=空, Description=生产备注 | 无备选 IP，有描述 |
| `127.0.0.1 www.test # 生产注释` | IP=127.0.0.1, Comment=生产注释, AlternateIP=空, Description=空 | 旧格式，保持向后兼容 |
| `127.0.0.1 www.test` | IP=127.0.0.1, 所有其他字段为空 | 无注释 |

## 解析算法

### parseLine() 增强逻辑

```
1. 验证 IP 和 DOMAIN（既有逻辑）
2. 检查是否存在 "#"
   如果不存在：解析完成，其他字段为空
   如果存在：进入第 3 步

3. 分离 "#" 之后的部分作为 rawComment
4. 在 rawComment 中查找第二个 "#"
   如果没有第二个 "#"：
     - 将 rawComment 作为 Comment 字段（旧格式兼容）
   如果有第二个 "#"：
     - 第一部分作为 AlternateIP（验证 IP 格式或允许为空）
     - 第二部分作为 Description
```

### 代码变更位置

**文件**：`internal/hosts/service.go`

**现有 parseLine 方法**：
```go
func (s *Service) parseLine(line string) *HostEntry {
    // ... 现有 IP 和 DOMAIN 解析 ...
    
    // 分离注释
    var comment string
    if idx := strings.Index(line, "#"); idx != -1 {
        comment = strings.TrimSpace(line[idx+1:])
        line = strings.TrimSpace(line[:idx])
    }
    
    // ... 返回 HostEntry ...
}
```

**修改后的逻辑**：
```go
func (s *Service) parseLine(line string) *HostEntry {
    // ... 现有 IP 和 DOMAIN 解析 ...
    
    // 增强：分离注释并解析 AlternateIP 和 Description
    var comment, alternateIP, description string
    if idx := strings.Index(line, "#"); idx != -1 {
        rawComment := strings.TrimSpace(line[idx+1:])
        line = strings.TrimSpace(line[:idx])
        
        // 检查是否为新格式（包含第二个 #）
        if secondIdx := strings.Index(rawComment, "#"); secondIdx != -1 {
            // 新格式：ALT_IP#DESCRIPTION
            alternateIP = strings.TrimSpace(rawComment[:secondIdx])
            description = strings.TrimSpace(rawComment[secondIdx+1:])
            
            // 处理 "null" 字符串
            if alternateIP == "null" {
                alternateIP = ""
            }
            if description == "null" {
                description = ""
            }
        } else {
            // 旧格式：仅注释
            comment = rawComment
        }
    }
    
    // ... 返回扩展的 HostEntry ...
}
```

## 生成文件内容

### generateHostsContent() 增强逻辑

需要修改生成 Hosts 文件行的逻辑：

**当前逻辑**：
```
启用：IP DOMAIN # COMMENT
禁用：# IP DOMAIN # COMMENT
```

**修改后的逻辑**：
```
启用且有备选 IP：IP DOMAIN # ALT_IP#DESCRIPTION
启用但无备选 IP：IP DOMAIN # COMMENT（或 ## DESCRIPTION 如果有描述）
禁用：# IP DOMAIN # ALT_IP#DESCRIPTION
```

**优先级**：
1. 如果有新格式字段（AlternateIP 或 Description），生成新格式
2. 如果只有 Comment，生成旧格式
3. 如果都没有，仅输出 IP 和 DOMAIN

## 验证规则

### IP 验证

- 主 IP：已有验证逻辑（IPv4/IPv6）
- 备选 IP：使用相同的验证逻辑，但允许为空字符串
- 实现：复用现有的 `isValidIP()` 方法，添加空字符串检查

### 其他验证

- Description：无格式限制，允许任何文本
- AlternateIP 和 Description：都允许为空

## 前端交互设计

### 表格列显示

添加或修改表格列以显示：
- **备选 IP**：新增列，宽度约 120px
- **描述**：在现有备注列中或新增列
- **交换按钮**：在操作列中添加交换图标（⇅ 或 🔄）

### 编辑弹窗

修改编辑弹窗以包含：
- IP 地址输入框（现有）
- 域名输入框（现有）
- **备选 IP 输入框**（新增）
- **描述输入框**（新增）
- 备注输入框（现有，可选择保留或移除）

### 交换 IP 操作

实现一键交换主 IP 和备选 IP 的功能：
1. 点击交换按钮
2. 交换 `IP` 和 `AlternateIP` 的值
3. 标记为已修改
4. 无需立即保存，用户点击"应用"时一起保存

## 兼容性考虑

### 现有数据迁移

- 旧版本的 Hosts 数据（仅含 IP、Domain、Comment）在加载时正常工作
- 在编辑和保存时，新字段自动初始化为空
- 无需特殊的数据迁移脚本

### Hosts 文件格式

- 标准 Hosts 文件格式保持不变
- 新格式仅作为注释部分的扩展，不违反 Hosts 文件规范
- 其他 Hosts 管理工具仍然可以读取和编辑生成的文件

## 测试策略

### 后端单元测试

1. **解析测试**：
   - 测试各种格式的输入行
   - 验证新旧格式的正确解析
   - 测试边界情况（空值、特殊字符等）

2. **生成测试**：
   - 验证生成的 Hosts 行格式正确
   - 测试启用/禁用状态下的输出
   - 测试各种字段组合

### 前端集成测试

1. **编辑和交换**：验证交换 IP 操作的正确性
2. **表单验证**：验证备选 IP 格式检查
3. **保存和加载**：验证新格式数据正确持久化

## 性能考虑

- 解析性能：新增 `strings.Index()` 调用，但影响可忽略
- 内存：HostEntry 结构体增加两个字符串字段，影响可忽略
- 无需特殊优化

## 安全考虑

- 输入验证：备选 IP 需要格式验证（现有 isValidIP 方法已足够）
- 注入防护：无特殊风险，文本字段仅作为 Hosts 文件内容处理
- 权限：无新的权限要求
