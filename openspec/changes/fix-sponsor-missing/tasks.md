# Tasks — fix-sponsor-missing

- [x] 诊断问题并记录根因（见 `DIAGNOSIS_AND_FIX.md`）
- [x] 在 `frontend/src/views/Settings.vue` 中添加加载状态（`qrCodeLoaded`）与重试计数（`qrCodeErrorRetried`）
- [x] 改进模板渲染逻辑：赞赏区域始终显示，根据加载状态展示“加载中/成功/失败”三种内容
- [x] 添加 `onQrCodeLoad` / `onQrCodeLoadError` / `retryLoadQrCode` 事件处理逻辑
- [x] 添加样式（`.qr-loading`、`.sponsor-error-content` 等）并验证样式效果
- [x] 在构建和发布文档中建议使用本地资源或 `#file:` 引用以避免网络依赖（已在仓库 README 与 `internal/version/version.go` 中更新）
- [x] 完成本地验证：`wails dev` 与 `wails build` 环境下行为符合预期

Notes:
- 所有变更均已应用并在 `frontend/src/views/Settings.vue` 中实现。
- `internal/version/version.go` 中的 `SponsorQrCode` 已更新为 `#file:wechat-qr.png`，并已同步到 `README.md` 与 `README_CN.md`。
