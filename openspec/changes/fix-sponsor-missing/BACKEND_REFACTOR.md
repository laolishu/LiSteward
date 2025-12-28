# 赞赏二维码来源重构 - 从 app.json 迁移到 version.go

## 改动说明

现在赞赏二维码由后端的 `version.go` 进行定义和管理，而不是从 `app.json` 中读取。

## 修改内容

### 1. 后端修改

#### internal/version/version.go
添加新变量：
```go
// SponsorQrCode 赞赏二维码 URL（微信收款码）
var SponsorQrCode = "https://coffee.laolishu.com/wechat-qr.png"
```

#### config/config.go
- 添加导入：`"LiSteward/internal/version"`
- 修改 `Load()` 函数，使其从 `version.go` 读取 `SponsorQrCode`：
```go
func Load() (*AppConfig, error) {
    // ... 加载 app.json ...
    
    // 从 version.go 中读取赞赏二维码，覆盖 app.json 中的值
    cfg.SponsorQrCode = version.SponsorQrCode
    
    // ...
}
```

#### config/app.json
移除 `sponsorQrCode` 字段（现在由 version.go 管理）

### 2. 数据流

**修改前**:
```
Settings.vue 
    ↓
GetAppConfig() 
    ↓
config.Load() 
    ↓
config/app.json (读取 sponsorQrCode)
```

**修改后**:
```
Settings.vue 
    ↓
GetAppConfig() 
    ↓
config.Load() 
    ↓
config/app.json (不包含 sponsorQrCode)
    ↓
version.SponsorQrCode (从 version.go 读取)
```

## 好处

1. ✅ **集中管理** - 所有版本相关的信息都在 `version.go` 中
2. ✅ **编译时管理** - 可以通过 `-ldflags` 在编译时指定（将来改进）
3. ✅ **简化配置** - `app.json` 不需要维护 `sponsorQrCode`
4. ✅ **代码一致性** - 与其他版本信息（Version, Author 等）集中管理

## 修改列表

| 文件 | 修改 |
|------|------|
| `internal/version/version.go` | 添加 `SponsorQrCode` 变量 |
| `config/config.go` | 导入 `version` 包，修改 `Load()` 函数 |
| `config/app.json` | 移除 `sponsorQrCode` 字段 |

## 编译检查

```
✅ 无编译错误
✅ 所有 import 正确
✅ 所有变量已定义
```

## 前端无需改动

Settings.vue 仍然通过 `GetAppConfig()` 获取 `appConfig.sponsorQrCode`，逻辑不变。

## 后续改进建议

### 1. 支持编译时参数
使用 `-ldflags` 在编译时指定二维码：
```bash
go build -ldflags "-X LiSteward/internal/version.SponsorQrCode=https://..."
```

### 2. 支持本地环境变量
```go
if qrCode := os.Getenv("SPONSOR_QR_CODE"); qrCode != "" {
    cfg.SponsorQrCode = qrCode
}
```

### 3. 支持本地文件
```go
if data, err := os.ReadFile("local-qr-code.txt"); err == nil {
    cfg.SponsorQrCode = strings.TrimSpace(string(data))
}
```

---

**状态**: ✅ 完成  
**编译验证**: ✅ 无错误  
**兼容性**: ✅ 完全向后兼容
