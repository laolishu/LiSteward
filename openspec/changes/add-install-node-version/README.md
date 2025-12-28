# 文档导读指南

## 推荐阅读顺序

### 快速了解 (5 分钟)
1. **[SOLUTION_SUMMARY.md](SOLUTION_SUMMARY.md)** ⭐ 从这里开始
   - 快速了解问题、解决方案、修改清单
   - 了解 wails dev vs build 差异
   - 查看关键代码片段

### 深入理解 (10 分钟)
2. **[WAILS_BUILD_FIX.md](WAILS_BUILD_FIX.md)**
   - Windows/macOS/Linux 平台的详细实现
   - 为什么只有 Build 版本出现问题
   - 静默模式的技术细节
   - 后续优化建议

3. **[IMPLEMENTATION_VERIFICATION.md](IMPLEMENTATION_VERIFICATION.md)**
   - 实现状态汇总
   - 代码验证细节
   - 验收清单
   - 相关文件修改历史

### 实施测试 (15 分钟)
4. **[QUICK_TEST.md](QUICK_TEST.md)**
   - 开发模式测试步骤
   - 编译版本测试步骤
   - 故障排除指南
   - 验收标准

### 项目管理 (可选)
5. **[tasks.md](tasks.md)**
   - 项目任务追踪
   - 各个功能点的实现状态
   - 完成度评估

6. **[proposal.md](proposal.md)** & **[design.md](design.md)**
   - 原始需求和设计文档
   - 了解功能的来源和目标

---

## 场景导航

### "我是项目经理，想快速了解进度"
→ 阅读 [SOLUTION_SUMMARY.md](SOLUTION_SUMMARY.md) 中的"验收清单"部分

### "我是开发者，想理解为什么会出现多窗口问题"
→ 阅读 [WAILS_BUILD_FIX.md](WAILS_BUILD_FIX.md) 中的"根本原因分析"

### "我需要测试这个功能"
→ 按照 [QUICK_TEST.md](QUICK_TEST.md) 中的步骤操作

### "我要修复编译错误"
→ 参考 [IMPLEMENTATION_VERIFICATION.md](IMPLEMENTATION_VERIFICATION.md) 中的"编译验证"部分

### "我想学习 Windows 进程管理"
→ 阅读 [WAILS_BUILD_FIX.md](WAILS_BUILD_FIX.md) 中的"技术细节"部分

### "我需要调试窗口问题"
→ 参考 [QUICK_TEST.md](QUICK_TEST.md) 中的"开发者调试"部分

---

## 文档结构概览

```
OpenSpec 文档结构:
├─ SOLUTION_SUMMARY.md           ⭐ 总览 - 必读
├─ WAILS_BUILD_FIX.md            📋 详细分析
├─ IMPLEMENTATION_VERIFICATION.md 🔍 实现验证
├─ QUICK_TEST.md                 🧪 测试指南
│
├─ proposal.md                    📝 原始提案
├─ design.md                      🎨 设计文档
├─ tasks.md                       ✅ 任务清单
│
└─ specs/                         📚 技术规范
```

---

## 快速问答

**Q: 多窗口问题已经解决了吗？**
A: ✅ 是的。已在 `internal/nvm/nvm.go` 中添加 `syscall.SysProcAttr{CreationFlags: 0x08000000}` 隐藏中间窗口。

**Q: 需要我做什么？**
A: 按照 [QUICK_TEST.md](QUICK_TEST.md) 进行测试验证。

**Q: 代码能编译吗？**
A: ✅ 是的。已验证无编译错误。

**Q: 这个修复会影响其他功能吗？**
A: ❌ 不会。仅修改 NVM 相关代码，不影响其他模块。

**Q: macOS 和 Linux 也支持吗？**
A: ✅ 是的。已实现完整的跨平台支持。

**Q: 窗口会一直打开吗？**
A: ✅ 是的。使用 `cmd.Start()` 确保窗口独立存在，用户可自由关闭。

---

## 关键术语说明

| 术语 | 解释 |
|------|------|
| **wails dev** | 开发模式，运行热重载开发服务器 |
| **wails build** | 编译模式，生成独立可执行程序 |
| **CREATE_NO_WINDOW** | Windows API 标志，隐藏进程窗口 |
| **SysProcAttr** | Go syscall 包中的进程属性结构体 |
| **cmd.Start()** | 启动进程但不等待其完成 |
| **中间窗口** | `cmd /c start` 这个启动进程的窗口 |

---

## 联系与反馈

如果文档有问题或需要补充，请参考：
- 代码问题 → 检查 `internal/nvm/nvm.go`
- 测试问题 → 参考 [QUICK_TEST.md](QUICK_TEST.md)
- 设计问题 → 参考 [design.md](design.md)

---

**最后更新**: 2025-01-13
**状态**: ✅ 完成
**版本**: 1.0.0
