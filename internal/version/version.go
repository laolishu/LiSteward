/*
 * @Descripttion:
 * @version:
 * @Author: lfzxs@qq.com
 * @Date: 2025-12-26 17:28:43
 * @LastEditors: lfzxs@qq.com
 * @LastEditTime: 2025-12-28 19:59:10
 */
package version

// Version 应用版本号，通过编译参数 -ldflags 指定
var Version = "0.2.0"

// Author 作者名称
var Author = "laolishu"

// AuthorEmail 作者邮箱
var AuthorEmail = "lfzxs@hotmail.com"

// Repository 项目仓库地址
var Repository = "https://github.com/laolishu/listeward"

// SponsorQrCode 赞赏二维码 路径或 URL（微信收款码）
// 使用本地资源引用时以 "#file:filename" 形式表示仓库内文件
var SponsorQrCode = "wechat-qr.png"

// BuildTime 编译时间，格式: 2025-12-26 12:00:00
var BuildTime = "unknown"

// BuildCommit 编译时的 Git 提交哈希
var BuildCommit = "unknown"

// Description 应用描述
var Description = "LiSteward - A lightweight hosts file manager with system tray and global hotkey support"
