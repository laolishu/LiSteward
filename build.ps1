# build.ps1 - PowerShell 编译脚本示例

# 获取当前时间
$BuildTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

# 获取 Git 提交哈希（如果是 Git 仓库）
try {
    $BuildCommit = & git rev-parse --short HEAD 2>$null
} catch {
    $BuildCommit = "unknown"
}

# 编译参数
$LDFlags = @(
    "-X", "LiSteward/internal/version.Version=1.0.0",
    "-X", "LiSteward/internal/version.Author=Your Name",
    "-X", "LiSteward/internal/version.AuthorEmail=your.email@example.com",
    "-X", "LiSteward/internal/version.Repository=https://github.com/yourname/listeward",
    "-X", "LiSteward/internal/version.BuildTime=$BuildTime",
    "-X", "LiSteward/internal/version.BuildCommit=$BuildCommit"
) -join ' '

# 编译
Write-Host "开始编译..."
go build -ldflags $LDFlags -o build\bin\LiSteward.exe .

if ($LASTEXITCODE -eq 0) {
    Write-Host "编译完成！" -ForegroundColor Green
    Write-Host "版本: 1.0.0"
    Write-Host "作者: Your Name"
    Write-Host "编译时间: $BuildTime"
    Write-Host "提交哈希: $BuildCommit"
} else {
    Write-Host "编译失败！" -ForegroundColor Red
}
