package hosts

import "time"

// HostEntry 表示单个 Hosts 文件条目
type HostEntry struct {
	IP          string `json:"ip"`           // IP 地址（IPv4 或 IPv6）
	Domain      string `json:"domain"`       // 域名
	Comment     string `json:"comment"`      // 注释（可选）
	AlternateIP string `json:"alternate_ip"` // 备选 IP 地址（可选）
	Description string `json:"description"`  // 详细描述（可选）
	Enabled     bool   `json:"enabled"`      // 是否启用
}

// HostsProfile 表示一个 Hosts 配置方案
type HostsProfile struct {
	Name      string      `json:"name"`       // 方案名称
	Entries   []HostEntry `json:"entries"`    // Hosts 条目列表
	CreatedAt time.Time   `json:"created_at"` // 创建时间
	UpdatedAt time.Time   `json:"updated_at"` // 更新时间
}

// BackupInfo 表示备份信息
type BackupInfo struct {
	ID        string    `json:"id"`        // 备份ID（时间戳格式）
	Timestamp time.Time `json:"timestamp"` // 备份时间
	Size      int64     `json:"size"`      // 文件大小（字节）
	Path      string    `json:"path"`      // 备份文件路径
}
