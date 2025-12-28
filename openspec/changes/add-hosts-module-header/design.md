# 设计文档：HostsManager 模块 ModuleHeader 集成

## 设计目标

为 HostsManager 模块添加统一的 ModuleHeader 组件，实现以下目标：

1. **统一页头风格**: 与 NodeManager、Settings 等模块保持一致
2. **集成消息提示**: 将 toast 消息系统迁移到 ModuleHeader
3. **保持布局**: 不影响现有的两面板设计和功能
4. **改善用户体验**: 更清晰地标识模块身份

## HostsManager 现有结构分析

### 布局特点

```
HostsManager
├── 全局消息提示 (message-toast) - 顶部
├── module-layout (两面板)
│   ├── sidebar (左侧导航)
│   └── right-panel (右侧编辑区)
│       ├── editor-toolbar (工具栏)
│       ├── editor-container (内容区)
│       └── 表格等内容
```

### 消息提示系统

当前实现：
- **消息对象**: `message = { text: '', type: 'success|error' }`
- **更新方式**: 通过修改 `message` 状态
- **样式**: `.message-toast` 和 `.message-toast.success|error` 类
- **自动消退**: 通过 `showMessage()` 方法设置 2 秒超时

### 关键修改点

**当前代码结构**:
```vue
<transition name="message-fade">
    <div v-if="message.text" :class="['message-toast', message.type]">
        <!-- 消息图标和文本 -->
    </div>
</transition>

<div class="module-layout">
    <!-- 其他内容 -->
</div>
```

## 设计方案

### 集成策略

**替换方式**:
```vue
<ModuleHeader 
  :title="$t('hosts.title')" 
  :hint="messageHint"
/>

<div class="module-layout">
  <!-- 保持不变 -->
</div>
```

### 数据映射

当前的 `message` 对象需要转换为 ModuleHeader 的 `hint` 对象：

| 当前系统 | ModuleHeader |
|--------|-------------|
| `message.type: 'success'` | `hint.type: 'success'` |
| `message.type: 'error'` | `hint.type: 'error'` |
| `message.text` | `hint.message` |
| `message.text === ''` | `hint: null` |

### 消息提示逻辑迁移

**当前方法**:
```javascript
const showMessage = (text, type = 'success', duration = 2000) => {
  message.value = { text, type }
  setTimeout(() => {
    message.value = { text: '', type: 'success' }
  }, duration)
}
```

**新方法**:
```javascript
const showMessage = (text, type = 'success', duration = 2000) => {
  messageHint.value = { message: text, type }
  
  // 错误提示不自动消退，成功提示 2 秒后消退
  if (type === 'success') {
    setTimeout(() => {
      messageHint.value = null
    }, duration)
  }
}
```

### 自动消退行为

ModuleHeader 的自动消退设定：
- **成功提示**: 5 秒自动消退
- **错误提示**: 不自动消退（需用户手动关闭）

HostsManager 的需求：
- **成功提示**: 2 秒消退（更快的反馈）
- **错误提示**: 需要用户手动关闭

**解决方案**: 在 HostsManager 中管理消息生命周期，不依赖 ModuleHeader 的自动消退，通过手动设置 `messageHint.value = null` 控制

### 操作影响分析

当前消息提示的所有操作点：

1. **加载操作**:
   - `reloadHosts()`: 显示 "加载中..." 或 "重新加载成功"
   - 当前: 直接修改 `entries`，未显示消息
   - 新增: 可以显示 "加载成功" 提示

2. **保存操作**:
   - `applyChanges()`: 应用更改，显示成功或错误
   - 当前: 显示 "已保存" 或错误信息
   - 新增: 通过 ModuleHeader 显示

3. **删除操作**:
   - 删除条目后可显示 "已删除"
   - 当前: 无消息提示
   - 新增: 可集成到 ModuleHeader

4. **其他操作**:
   - 启用/禁用条目、交换 IP 等
   - 当前: 无消息提示
   - 新增: 可选择是否添加提示

## 样式一致性

### ModuleHeader 设计回顾

- 白色背景 (#ffffff)
- 下边框 1px #e1e4e8
- 内边距 24px
- 标题字号 24px，字重 600，颜色 #24292e

### HostsManager 容器

当前 `.hosts-manager`:
```css
display: flex;
flex-direction: column;
/* 其他样式 */
```

**调整**: ModuleHeader 应位于 `.hosts-manager` 顶部，作为第一个子元素

## 兼容性考虑

### 消息提示时序

注意 ModuleHeader 的提示框自动消退逻辑与 HostsManager 的消息系统的差异：

- ModuleHeader 成功提示: 5 秒自动消退
- HostsManager 历史: 2 秒自动消退
- **方案**: HostsManager 使用自己的定时器控制消息清除，不依赖 ModuleHeader 的自动消退

### CSS 冲突风险

需要验证：
- `.message-toast` 样式是否被完全移除
- 是否有其他样式依赖 `.message-fade` transition
- 两面板布局是否受到 ModuleHeader 插入的影响

## 实施步骤概览

1. 在 HostsManager 中导入 ModuleHeader
2. 添加 `messageHint` ref 来映射现有消息
3. 修改 `showMessage()` 方法以支持 ModuleHeader 的 hint
4. 将 `.message-toast` HTML 替换为 ModuleHeader 组件
5. 删除相关的 CSS 样式定义
6. 验证所有消息提示场景正常工作

## 扩展考虑

未来可以考虑：
1. 为 HostsManager 的其他操作添加提示（如删除、启用/禁用等）
2. 在 ModuleHeader 中添加副标题显示当前选中的配置文件
3. 优化消息提示的显示时长（可配置化）

## 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|---------|
| 消息提示逻辑错误 | 中 | 高 | 充分的手动测试 |
| 样式冲突 | 低 | 中 | 完整的 CSS 清理 |
| 布局破坏 | 低 | 高 | 验证两面板不受影响 |
| 性能问题 | 低 | 低 | 监测消息提示频率 |
