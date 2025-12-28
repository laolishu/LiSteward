#!/bin/bash
# build.sh - 编译脚本示例

# 获取当前时间
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')

# 获取 Git 提交哈希（如果是 Git 仓库）
BUILD_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 编译参数
LDFLAGS="-X 'LiSteward/internal/version.Version=1.0.0' \
         -X 'LiSteward/internal/version.Author=Your Name' \
         -X 'LiSteward/internal/version.AuthorEmail=your.email@example.com' \
         -X 'LiSteward/internal/version.Repository=https://github.com/yourname/listeward' \
         -X 'LiSteward/internal/version.BuildTime=${BUILD_TIME}' \
         -X 'LiSteward/internal/version.BuildCommit=${BUILD_COMMIT}'"

# 编译
go build -ldflags "$LDFLAGS" -o build/bin/LiSteward.exe .

echo "编译完成！"
echo "版本: 1.0.0"
echo "作者: Your Name"
echo "编译时间: $BUILD_TIME"
echo "提交哈希: $BUILD_COMMIT"
