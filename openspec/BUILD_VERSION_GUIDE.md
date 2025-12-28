# 版本信息管理说明

## 概述

版本号、作者等全局信息定义在 `internal/version/version.go` 中，通过编译参数 (`-ldflags`) 动态设置，不依赖配置文件。

## 文件说明

### 1. `internal/version/version.go`
定义所有全局变量：
- `Version` - 应用版本号
- `Author` - 作者名称
- `AuthorEmail` - 作者邮箱
- `Repository` - 项目仓库地址
- `BuildTime` - 编译时间
- `BuildCommit` - Git 提交哈希
- `Description` - 应用描述

### 2. `build.env`
配置文件，存储编译参数的具体值。可根据需要修改。

### 3. `build.ps1`
Windows PowerShell 编译脚本，自动读取编译参数并执行编译。

### 4. `build.sh`
Linux/macOS 编译脚本，自动读取编译参数并执行编译。

## 使用方法

### 方法 1: 直接使用 go build 命令

**Windows**:
```powershell
$TIME = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$COMMIT = git rev-parse --short HEAD

go build -ldflags `
  "-X LiSteward/internal/version.Version=1.0.0 " + `
  "-X LiSteward/internal/version.Author='Your Name' " + `
  "-X LiSteward/internal/version.AuthorEmail='your.email@example.com' " + `
  "-X LiSteward/internal/version.Repository='https://github.com/yourname/listeward' " + `
  "-X LiSteward/internal/version.BuildTime='$TIME' " + `
  "-X LiSteward/internal/version.BuildCommit='$COMMIT'" `
  -o build\bin\LiSteward.exe .
```

**Linux/macOS**:
```bash
TIME=$(date '+%Y-%m-%d %H:%M:%S')
COMMIT=$(git rev-parse --short HEAD)

go build -ldflags \
  "-X LiSteward/internal/version.Version=1.0.0 \
   -X LiSteward/internal/version.Author='Your Name' \
   -X LiSteward/internal/version.AuthorEmail='your.email@example.com' \
   -X LiSteward/internal/version.Repository='https://github.com/yourname/listeward' \
   -X LiSteward/internal/version.BuildTime='$TIME' \
   -X LiSteward/internal/version.BuildCommit='$COMMIT'" \
  -o build/bin/LiSteward .
```

### 方法 2: 使用编译脚本

**Windows**:
```powershell
# 编辑 build.ps1，修改版本号和作者信息
.\build.ps1
```

**Linux/macOS**:
```bash
# 编辑 build.sh，修改版本号和作者信息
chmod +x build.sh
./build.sh
```

### 方法 3: 使用配置文件 + 脚本

编辑 `build.env` 文件设置参数，然后运行编译脚本会自动读取这些值。

## 在代码中访问版本信息

```go
package main

import (
    "fmt"
    "LiSteward/internal/version"
)

func main() {
    fmt.Printf("版本: %s\n", version.Version)
    fmt.Printf("作者: %s <%s>\n", version.Author, version.AuthorEmail)
    fmt.Printf("仓库: %s\n", version.Repository)
    fmt.Printf("编译时间: %s\n", version.BuildTime)
    fmt.Printf("提交: %s\n", version.BuildCommit)
    fmt.Printf("描述: %s\n", version.Description)
}
```

## 在前端访问版本信息

在 `app.go` 中添加 Wails 接口方法：

```go
import "LiSteward/internal/version"

// GetAppInfo 获取应用信息
func (a *App) GetAppInfo() map[string]string {
    return map[string]string{
        "version":      version.Version,
        "author":       version.Author,
        "email":        version.AuthorEmail,
        "repository":   version.Repository,
        "buildTime":    version.BuildTime,
        "buildCommit":  version.BuildCommit,
        "description":  version.Description,
    }
}
```

然后在前端 Vue 组件中可以调用：

```vue
<script setup>
import { GetAppInfo } from '../wailsjs/go/main/App'

const appInfo = ref(null)

onMounted(async () => {
  appInfo.value = await GetAppInfo()
})
</script>

<template>
  <div v-if="appInfo">
    <p>版本: {{ appInfo.version }}</p>
    <p>作者: {{ appInfo.author }} ({{ appInfo.email }})</p>
    <p>编译时间: {{ appInfo.buildTime }}</p>
  </div>
</template>
```

## 编译参数说明

### `-ldflags` 选项说明

- `-X` 标志用于在编译时设置包级别的字符串变量
- 格式: `-X 'package.Variable=value'`
- 注意：需要完整的包路径和变量名

### 特殊处理

- **BuildTime**: 自动捕获编译时的时间戳，无需手动指定（除非需要特定的时间）
- **BuildCommit**: 自动捕获当前 Git 提交哈希
- **Version**: 建议与 Git 标签同步（如 `v1.0.0`）

## 最佳实践

1. **版本号管理**
   - 遵循 [语义版本](https://semver.org/) 规范 (major.minor.patch)
   - 与 Git 标签同步：`git tag -a v1.0.0`

2. **编译脚本**
   - 在 CI/CD 流程中使用编译脚本自动设置版本信息
   - 不要在代码中硬编码版本号

3. **默认值**
   - 如果未通过 `-ldflags` 指定，使用 `version.go` 中的默认值
   - 默认值应该是 `"dev"` 或类似的开发标识

## 常见问题

### Q: 如何在发布时自动设置版本号？
A: 在 CI/CD 流程中（如 GitHub Actions），读取 Git 标签然后传递给编译脚本：
```yaml
- run: |
    VERSION=${GITHUB_REF#refs/tags/v}
    go build -ldflags "-X LiSteward/internal/version.Version=$VERSION" ...
```

### Q: 如何验证编译后的版本信息是否正确？
A: 可以在应用启动时打印版本信息，或在"关于"页面显示。

### Q: 能否在运行时修改版本信息？
A: 不能。版本信息在编译时链接到二进制文件中，运行时是只读的。

---

**更新日期**: 2025-12-26
**相关文件**: `internal/version/version.go`, `build.env`, `build.ps1`, `build.sh`
