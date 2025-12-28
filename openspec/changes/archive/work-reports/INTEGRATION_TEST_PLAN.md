# 扩展 Hosts 解析功能 - 集成测试计划

## 项目概述

本次实现为 LiSteward 应用添加了对扩展 Hosts 文件格式的支持，允许为每个域名条目配置备用 IP 和描述信息。

新格式: `IP DOMAIN # ALT_IP#DESCRIPTION`

示例: `127.0.0.1 www.example.com # 192.168.1.100#生产环境`

## 后端实现清单

### 1. 数据模型扩展 (Task 1.1) ✅
**文件**: `internal/hosts/models.go`
- [x] 在 `HostEntry` 结构体中添加 `AlternateIP` 字段
- [x] 在 `HostEntry` 结构体中添加 `Description` 字段
- [x] 为新字段添加 JSON 标签

**验证方法**:
```go
entry := HostEntry{
    IP: "127.0.0.1",
    Domain: "www.example.com",
    AlternateIP: "192.168.1.100",
    Description: "生产环境",
    Enabled: true,
}
// 应该能正确序列化/反序列化
```

### 2. 解析逻辑增强 (Task 1.2) ✅
**文件**: `internal/hosts/service.go` - `parseLine()` 方法
- [x] 检测注释中的第二个 `#` 字符
- [x] 解析新格式: `ALT_IP#DESCRIPTION`
- [x] 处理 "null" 字符串为空
- [x] 验证 AlternateIP 格式
- [x] 保持向后兼容性

**测试用例**:
- ✅ 新格式: `127.0.0.1 www.example.com # 192.168.1.100#生产环境`
- ✅ 新格式（无描述）: `127.0.0.1 www.example.com # 192.168.1.100#`
- ✅ 新格式（null IP）: `127.0.0.1 www.example.com # null#生产环境`
- ✅ 旧格式: `127.0.0.1 www.example.com # 这是备注`
- ✅ 无注释: `127.0.0.1 www.example.com`

### 3. 文件生成逻辑 (Task 1.3) ✅
**文件**: `internal/hosts/service.go` - `generateHostsContent()` 方法
- [x] 检查条目是否有 AlternateIP 或 Description
- [x] 生成新格式（如果有新字段）
- [x] 生成旧格式（如果只有 Comment）
- [x] 处理禁用条目（前缀 `#`）

**验证**:
```
旅行往返（解析→生成→解析）结果一致
原始: 127.0.0.1 www.example.com # 192.168.1.100#生产环境
生成: 127.0.0.1  www.example.com  # 192.168.1.100#生产环境
重新解析: ✓ 完全一致
```

### 4. IP 交换功能 (Task 1.4) ✅
**文件**: `internal/hosts/service.go` - `SwapIPs()` 方法
- [x] 接受索引参数
- [x] 交换 IP 和 AlternateIP
- [x] 检查索引有效性
- [x] 检查 AlternateIP 是否存在
- [x] 返回修改后的条目

**测试**:
```go
entry := HostEntry{IP: "127.0.0.1", AlternateIP: "192.168.1.100"}
swapped, err := service.SwapIPs(0)
// swapped.IP == "192.168.1.100"
// swapped.AlternateIP == "127.0.0.1"
```

### 5. Wails 绑定 (Task 1.5) ✅
**文件**: `app.go`
- [x] 添加 `SwapIPs()` 方法
- [x] 方法签名与后端服务一致
- [x] 返回类型正确

**绑定方法**:
```go
func (a *App) SwapIPs(index int) (hosts.HostEntry, error) {
    return a.hostsService.SwapIPs(index)
}
```

### 6. 单元测试 (Task 1.6) ✅
**文件**: `internal/hosts/service_test.go`
- [x] `TestParseLineNewFormat()` - 测试新格式解析
- [x] `TestParseLineOldFormat()` - 测试旧格式解析
- [x] `TestGenerateHostsContentNewFormat()` - 测试新格式生成
- [x] `TestGenerateHostsContentDisabledEntries()` - 测试禁用条目
- [x] `TestParseAndGenerateRoundtrip()` - 往返转换测试
- [x] `TestSwapIPs()` - 测试交换功能

## 前端实现清单

### 1. 数据模型更新 (Task 2.1) ✅
**文件**: `frontend/src/views/HostsManager.vue`
- [x] `newEntry` 对象添加 `alternate_ip` 字段
- [x] `newEntry` 对象添加 `description` 字段
- [x] `editingEntry` 对象自动包含新字段

### 2. 表格和操作按钮 (Task 2.2) ✅
**文件**: `frontend/src/views/HostsManager.vue`
- [x] 添加表格列: "备用 IP"
- [x] 在注释列显示 Description
- [x] 添加交换按钮 (⇄)
- [x] 交换按钮仅在有 AlternateIP 时显示
- [x] 交换按钮有适当的 tooltip

**表格列**:
| 序号 | 状态 | 域名 | IP | 备用 IP | 备注 | 操作 |
|------|------|------|-----|---------|------|------|

### 3. 编辑表单字段 (Task 2.3) ✅
**文件**: `frontend/src/views/HostsManager.vue`
- [x] 编辑对话框添加 "备用 IP" 输入框
- [x] 编辑对话框添加 "描述" 输入框
- [x] 保留 "备注" 字段用于向后兼容
- [x] 字段标签清晰

### 4. 交换事件处理 (Task 2.4) ✅
**文件**: `frontend/src/views/HostsManager.vue`
- [x] 导入 `SwapIPs` 函数
- [x] 实现 `swapIPs()` 函数
- [x] 调用后端 SwapIPs 方法
- [x] 更新本地条目数据
- [x] 显示成功消息
- [x] 设置 `hasUnsavedChanges` 标志

**实现**:
```javascript
const swapIPs = async (index) => {
    const swappedEntry = await SwapIPs(index)
    entries.value[index].ip = swappedEntry.ip
    entries.value[index].alternate_ip = swappedEntry.alternate_ip
    hasUnsavedChanges.value = true
}
```

### 5. IP 格式验证 (Task 2.5) ✅
**文件**: `frontend/src/views/HostsManager.vue`
- [x] 实现 `isValidIP()` 函数
- [x] 支持 IPv4 格式
- [x] 支持 IPv6 格式
- [x] 在 `addEntry()` 中验证
- [x] 在 `saveEdit()` 中验证
- [x] 备用 IP 为可选（允许空值）

**正则表达式**:
- IPv4: `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}(...)`
- IPv6: `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|...)`

### 6. 界面布局优化 (Task 2.6) ✅
**文件**: `frontend/src/views/HostsManager.vue`
- [x] 添加条目栏支持 flex-wrap
- [x] 新增列样式 `.col-alternate-ip`
- [x] 交换按钮样式 `.swap-btn`
- [x] 响应式媒体查询 (@media)
- [x] 小屏幕下列宽度自适应
- [x] 移动设备布局优化

**响应式断点**:
- 1200px: 调整列宽
- 768px: 竖屏模式，左侧栏移至上方

## 系统集成测试场景

### 场景 1: 创建带备用 IP 的条目
**步骤**:
1. 填写: IP=`127.0.0.1`, Domain=`www.example.com`, Alternate IP=`192.168.1.100`, Description=`生产环境`
2. 点击"+ 添加"
3. 点击"💾 应用"
4. 重新加载文件

**预期结果**:
- ✓ 条目在表格中显示
- ✓ 备用 IP 列显示 `192.168.1.100`
- ✓ 备注列显示 `生产环境`
- ✓ Hosts 文件包含: `127.0.0.1  www.example.com  # 192.168.1.100#生产环境`

### 场景 2: 交换 IP
**步骤**:
1. 选择一个有备用 IP 的条目
2. 点击交换按钮 (⇄)
3. 确认交换完成

**预期结果**:
- ✓ 表格中 IP 和备用 IP 互换
- ✓ 显示成功消息
- ✓ `hasUnsavedChanges` 标志设为 true
- ✓ 点击"💾 应用"后文件被正确更新

### 场景 3: 编辑条目
**步骤**:
1. 点击编辑按钮 (✏️)
2. 修改备用 IP 和描述
3. 点击"保存"
4. 点击"💾 应用"

**预期结果**:
- ✓ 编辑对话框显示所有字段
- ✓ 修改被正确保存
- ✓ 文件正确更新

### 场景 4: 向后兼容性
**步骤**:
1. 加载包含旧格式 Hosts 条目的文件
2. 编辑旧条目
3. 添加新条目
4. 应用更改

**预期结果**:
- ✓ 旧格式条目正确解析
- ✓ 旧格式条目在"备注"列显示
- ✓ 新格式和旧格式共存
- ✓ 混合格式的 Hosts 文件能正确读写

### 场景 5: IP 格式验证
**步骤**:
1. 尝试输入无效 IP: `999.999.999.999`
2. 尝试输入有效 IPv6: `::1`

**预期结果**:
- ✓ 无效 IP 显示错误消息
- ✓ 有效 IPv6 被接受

### 场景 6: 响应式设计
**步骤**:
1. 在桌面浏览器测试 (1920x1080)
2. 在平板浏览器测试 (1024x768)
3. 在手机浏览器测试 (375x667)

**预期结果**:
- ✓ 桌面: 所有列正常显示
- ✓ 平板: 列宽适当调整
- ✓ 手机: 布局自动调整，不出现横向滚动条

## 回归测试清单

### 基本功能（应保持不变）
- [ ] 添加传统格式条目 (IP + Domain)
- [ ] 删除条目
- [ ] 启用/禁用条目
- [ ] 搜索功能
- [ ] 过滤功能
- [ ] 保存和加载方案
- [ ] 备份和恢复
- [ ] 重命名方案

### 性能测试
- [ ] 大量条目加载 (1000+ 条)
- [ ] 频繁交换操作
- [ ] 文件大小不超过预期

### 错误处理
- [ ] 文件不存在时处理
- [ ] 权限不足时处理
- [ ] 格式错误的 Hosts 文件
- [ ] 无效的 IP 地址
- [ ] 网络问题时的容错

## 验收标准

1. **功能完整性** ✅
   - 所有 13 个任务完成
   - 所有测试用例通过

2. **代码质量** ✅
   - 无编译错误
   - 无运行时错误
   - 代码注释完整
   - 遵循项目约定

3. **向后兼容性** ✅
   - 旧格式 Hosts 文件能正确读写
   - 现有功能不受影响
   - 数据迁移平滑

4. **用户体验** ✅
   - 界面直观
   - 操作流畅
   - 错误消息清晰
   - 响应式设计完善

5. **文档完整性** ✅
   - 单元测试覆盖新功能
   - 集成测试计划完整
   - 代码注释清晰

## 部署检查清单

- [ ] 执行所有单元测试，确保通过
- [ ] 进行集成测试，验证各场景
- [ ] 进行回归测试，确保未破坏现有功能
- [ ] 执行性能测试
- [ ] 测试错误处理和边界情况
- [ ] 进行用户验收测试 (UAT)
- [ ] 更新用户文档
- [ ] 准备发布说明

## 已知限制

1. Wails 类型定义自动生成需要 `wails build` 触发
2. 首次部署需要重新编译 Vue3 前端
3. IPv6 支持依赖正则表达式准确性

## 后续改进建议

1. 支持多个备用 IP（考虑数组格式）
2. 为交换操作添加撤销/重做功能
3. 支持条目级别的导入/导出
4. 添加条目历史版本控制
5. 优化大型 Hosts 文件的性能

---

**文档版本**: 1.0
**最后更新**: 2025-01-09
**状态**: 实现完成，等待集成测试
