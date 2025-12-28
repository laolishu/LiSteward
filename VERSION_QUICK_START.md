# 版本信息快速参考

## 文件结构

```
LiSteward/
├── internal/version/
│   └── version.go          # 版本变量定义
├── build.env               # 编译配置（编辑这个文件）
├── build.ps1               # Windows 编译脚本
├── build.sh                # Linux/Mac 编译脚本
└── app.go                  # GetAppInfo() 接口
```

## 编辑版本信息

编辑 `build.env` 文件：

```env
VERSION=1.0.0
AUTHOR=Your Name
AUTHOR_EMAIL=your.email@example.com
REPOSITORY=https://github.com/yourname/listeward
```

## 编译命令

### Windows (快速)

```powershell
go build -ldflags `
  "-X LiSteward/internal/version.Version=1.0.0 " + `
  "-X LiSteward/internal/version.Author='Your Name' " + `
  "-X LiSteward/internal/version.AuthorEmail='your.email@example.com' " + `
  "-X LiSteward/internal/version.Repository='https://github.com/yourname/listeward'" `
  -o build\bin\LiSteward.exe .
```

或使用脚本（编辑 build.ps1 后）：

```powershell
.\build.ps1
```

### Linux/macOS (快速)

```bash
go build -ldflags \
  "-X LiSteward/internal/version.Version=1.0.0 \
   -X LiSteward/internal/version.Author='Your Name' \
   -X LiSteward/internal/version.AuthorEmail='your.email@example.com' \
   -X LiSteward/internal/version.Repository='https://github.com/yourname/listeward'" \
  -o build/bin/LiSteward .
```

或使用脚本（编辑 build.sh 后）：

```bash
./build.sh
```

## 前端访问版本信息

在前端 Vue 组件中：

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
    <h3>关于</h3>
    <p><strong>版本:</strong> {{ appInfo.version }}</p>
    <p><strong>作者:</strong> {{ appInfo.author }} &lt;{{ appInfo.email }}&gt;</p>
    <p><strong>仓库:</strong> 
      <a :href="appInfo.repository" target="_blank">{{ appInfo.repository }}</a>
    </p>
    <p><strong>编译时间:</strong> {{ appInfo.buildTime }}</p>
    <p><strong>提交:</strong> {{ appInfo.buildCommit }}</p>
    <p><strong>Go版本:</strong> {{ appInfo.goVersion }}</p>
    <p><strong>平台:</strong> {{ appInfo.os }}/{{ appInfo.arch }}</p>
  </div>
</template>
```

## 版本号规范

遵循 [语义版本](https://semver.org/)：

- `MAJOR.MINOR.PATCH`
- 示例: `1.0.0`, `1.2.3`, `2.1.0`

## 版本控制建议

```bash
# 创建版本标签
git tag -a v1.0.0 -m "Release 1.0.0"

# 推送标签
git push origin v1.0.0

# 编译带版本号的二进制
go build -ldflags "-X LiSteward/internal/version.Version=1.0.0" ...
```

## 自动化 CI/CD 示例

### GitHub Actions

```yaml
name: Build Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set version from tag
        id: version
        run: |
          $TAG = $env:GITHUB_REF -replace '^refs/tags/v', ''
          echo "VERSION=$TAG" >> $env:GITHUB_OUTPUT
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.23
      
      - name: Build
        run: |
          $BuildTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
          $BuildCommit = git rev-parse --short HEAD
          
          go build -ldflags `
            "-X LiSteward/internal/version.Version=${{ steps.version.outputs.VERSION }} " + `
            "-X LiSteward/internal/version.BuildTime='$BuildTime' " + `
            "-X LiSteward/internal/version.BuildCommit='$BuildCommit' " + `
            "-X LiSteward/internal/version.Author='LiSteward Contributors' " + `
            "-X LiSteward/internal/version.Repository='https://github.com/yourname/listeward'" `
            -o build\bin\LiSteward.exe .
```

---

**更新日期**: 2025-12-26
**状态**: ✅ 已实现和集成
