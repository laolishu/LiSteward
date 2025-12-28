package hosts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	// Windows Hosts 文件路径
	HostsFilePath = `C:\Windows\System32\drivers\etc\hosts`

	// 备份文件名格式
	BackupTimeFormat = "20060102_150405"

	// 最大备份数量
	MaxBackups = 30
)

// Service 提供 Hosts 文件管理服务
type Service struct {
	mu         sync.Mutex
	appDataDir string // 应用数据目录
	backupDir  string // 备份目录
	profileDir string // 方案目录
}

// NewService 创建新的 Hosts 服务实例
func NewService(appDataDir string) (*Service, error) {
	s := &Service{
		appDataDir: appDataDir,
		backupDir:  filepath.Join(appDataDir, "backups", "hosts"),
		profileDir: filepath.Join(appDataDir, "profiles", "hosts"),
	}

	// 创建必要的目录
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("创建备份目录失败: %w", err)
	}
	if err := os.MkdirAll(s.profileDir, 0755); err != nil {
		return nil, fmt.Errorf("创建方案目录失败: %w", err)
	}

	return s, nil
}

// ReadHostsFile 读取并解析 Hosts 文件
func (s *Service) ReadHostsFile() ([]HostEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(HostsFilePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 Hosts 文件: %w (请确保以管理员权限运行)", err)
	}
	defer file.Close()

	var entries []HostEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		entry := s.parseLine(line)
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取 Hosts 文件失败: %w", err)
	}

	return entries, nil
}

// parseLine 解析单行 Hosts 配置
func (s *Service) parseLine(line string) *HostEntry {
	line = strings.TrimSpace(line)

	// 跳过空行
	if line == "" {
		return nil
	}

	enabled := true
	// 检查是否是注释或禁用的条目
	if strings.HasPrefix(line, "#") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
		// 如果是纯注释（不包含有效的 IP 和域名），则跳过
		// 注意：Hosts 文件中 IP 与域名可能用空格或制表符分隔，
		// 不应仅依赖单一空格字符判断。改为使用 strings.Fields 判断字段数。
		if line == "" || len(strings.Fields(line)) < 2 {
			return nil
		}
		enabled = false
	}

	// 增强：分离注释并解析 AlternateIP 和 Description
	var comment, alternateIP, description string
	if idx := strings.Index(line, "#"); idx != -1 {
		rawComment := strings.TrimSpace(line[idx+1:])
		line = strings.TrimSpace(line[:idx])

		// 检查是否为新格式（包含第二个 #）
		if secondIdx := strings.Index(rawComment, "#"); secondIdx != -1 {
			// 新格式：ALT_IP#DESCRIPTION
			alternateIP = strings.TrimSpace(rawComment[:secondIdx])
			description = strings.TrimSpace(rawComment[secondIdx+1:])

			// 处理 "null" 字符串
			if alternateIP == "null" {
				alternateIP = ""
			}
			if description == "null" {
				description = ""
			}
		} else {
			// 旧格式：仅注释
			comment = rawComment
		}
	}

	// 解析 IP 和域名
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}

	ip := fields[0]
	domain := fields[1]

	// 基本验证
	if !s.isValidIP(ip) || !s.isValidDomain(domain) {
		return nil
	}

	// 验证备选 IP（如果提供了）
	if alternateIP != "" && !s.isValidIP(alternateIP) {
		return nil
	}

	return &HostEntry{
		IP:          ip,
		Domain:      domain,
		Comment:     comment,
		AlternateIP: alternateIP,
		Description: description,
		Enabled:     enabled,
	}
}

// WriteHostsFile 写入 Hosts 文件
func (s *Service) WriteHostsFile(entries []HostEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 先备份当前文件
	if err := s.backupHostsFileInternal(); err != nil {
		return fmt.Errorf("备份失败: %w", err)
	}

	// 2. 生成新的 Hosts 文件内容
	content := s.generateHostsContent(entries)

	// 3. 使用 WriteFile 直接覆盖原始文件（更兼容 Windows）
	// 注意：这要求有 Hosts 文件的写入权限
	if err := os.WriteFile(HostsFilePath, []byte(content), 0644); err != nil {
		// 检查是否是权限问题
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足: 无法写入 Hosts 文件。请确保: 1.以管理员身份运行程序 2.Hosts 文件未被其他程序独占")
		}
		return fmt.Errorf("写入 Hosts 文件失败: %w", err)
	}

	return nil
}

// WriteHostsFileWithReadOnlyRestore 写入 Hosts 文件，自动处理只读属性
// 如果文件为只读，先临时取消只读标志，保存后恢复原有属性
func (s *Service) WriteHostsFileWithReadOnlyRestore(entries []HostEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取文件当前属性
	fileInfo, err := os.Stat(HostsFilePath)
	if err != nil {
		return fmt.Errorf("无法获取 Hosts 文件属性: %w", err)
	}

	// 检查文件是否为只读
	isReadOnly := fileInfo.Mode()&0200 == 0

	// 如果文件为只读，临时取消只读标志
	if isReadOnly {
		if err := s.setReadOnly(HostsFilePath, false); err != nil {
			return fmt.Errorf("无法取消只读属性: %w", err)
		}
		// 确保在函数结束时恢复只读属性
		defer func() {
			if restoreErr := s.setReadOnly(HostsFilePath, true); restoreErr != nil {
				fmt.Printf("警告：无法恢复只读属性: %v\n", restoreErr)
			}
		}()
	}

	// 执行备份和写入
	if err := s.backupHostsFileInternal(); err != nil {
		return fmt.Errorf("备份失败: %w", err)
	}

	content := s.generateHostsContent(entries)

	if err := os.WriteFile(HostsFilePath, []byte(content), 0644); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足: 无法写入 Hosts 文件。请确保: 1.以管理员身份运行程序 2.Hosts 文件未被其他程序独占")
		}
		return fmt.Errorf("写入 Hosts 文件失败: %w", err)
	}

	return nil
}

// setReadOnly 设置 Windows 文件的只读属性
func (s *Service) setReadOnly(filePath string, readonly bool) error {
	// 将路径转换为 UTF-16 指针（Windows API 要求）
	pathPtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	// 获取文件属性
	attrs, err := syscall.GetFileAttributes(pathPtr)
	if err != nil {
		return err
	}

	// 设置或取消只读属性
	if readonly {
		attrs |= syscall.FILE_ATTRIBUTE_READONLY
	} else {
		attrs &^= syscall.FILE_ATTRIBUTE_READONLY
	}

	// 应用新属性
	if err := syscall.SetFileAttributes(pathPtr, attrs); err != nil {
		return err
	}

	return nil
}

// generateHostsContent 生成 Hosts 文件内容
func (s *Service) generateHostsContent(entries []HostEntry) string {
	var builder strings.Builder

	builder.WriteString("# Hosts 文件由 LiSteward 管理\n")
	builder.WriteString(fmt.Sprintf("# 最后更新: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	for _, entry := range entries {
		if !entry.Enabled {
			builder.WriteString("# ")
		}

		builder.WriteString(entry.IP)
		builder.WriteString("\t")
		builder.WriteString(entry.Domain)

		// 使用扩展格式（如果有 AlternateIP 或 Description）或向后兼容格式
		if entry.AlternateIP != "" || entry.Description != "" {
			// 扩展格式: IP DOMAIN # ALT_IP#DESCRIPTION
			builder.WriteString("\t# ")
			if entry.AlternateIP != "" {
				builder.WriteString(entry.AlternateIP)
			} else {
				builder.WriteString("null")
			}
			builder.WriteString("#")
			if entry.Description != "" {
				builder.WriteString(entry.Description)
			}
		} else if entry.Comment != "" {
			// 旧格式: IP DOMAIN # COMMENT（向后兼容）
			builder.WriteString("\t# ")
			builder.WriteString(entry.Comment)
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// SwapIPs 交换条目的主 IP 和备用 IP
func (s *Service) SwapIPs(entries []HostEntry, index int) ([]HostEntry, error) {
	if index < 0 || index >= len(entries) {
		return nil, fmt.Errorf("索引超出范围: %d", index)
	}

	if entries[index].AlternateIP == "" {
		return nil, fmt.Errorf("条目 %d 没有备用 IP，无法交换", index)
	}

	// 交换 IP 和 AlternateIP
	entries[index].IP, entries[index].AlternateIP = entries[index].AlternateIP, entries[index].IP

	return entries, nil
}

// isValidIP 验证 IP 地址格式
func (s *Service) isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// isValidDomain 验证域名格式
func (s *Service) isValidDomain(domain string) bool {
	// 允许 localhost
	if domain == "localhost" {
		return true
	}

	// 允许 .local 域名
	if strings.HasSuffix(domain, ".local") {
		return true
	}

	// 标准域名验证
	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, domain)
	return matched
}

// ValidateEntry 验证条目的有效性
func (s *Service) ValidateEntry(entry HostEntry) error {
	if entry.IP == "" {
		return fmt.Errorf("IP 地址不能为空")
	}
	if !s.isValidIP(entry.IP) {
		return fmt.Errorf("无效的 IP 地址: %s", entry.IP)
	}
	if entry.Domain == "" {
		return fmt.Errorf("域名不能为空")
	}
	if !s.isValidDomain(entry.Domain) {
		return fmt.Errorf("无效的域名: %s", entry.Domain)
	}
	return nil
}

// backupHostsFileInternal 内部备份方法（不加锁）
func (s *Service) backupHostsFileInternal() error {
	timestamp := time.Now().Format(BackupTimeFormat)
	backupPath := filepath.Join(s.backupDir, fmt.Sprintf("hosts_backup_%s.txt", timestamp))

	// 复制文件
	src, err := os.Open(HostsFilePath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	// 清理旧备份
	go s.cleanOldBackups()

	return nil
}

// BackupHosts 手动创建备份
func (s *Service) BackupHosts() (*BackupInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.backupHostsFileInternal(); err != nil {
		return nil, err
	}

	// 获取最新的备份信息
	backups, err := s.listBackupsInternal()
	if err != nil {
		return nil, err
	}

	if len(backups) > 0 {
		return &backups[0], nil
	}

	return nil, fmt.Errorf("备份创建后未找到备份文件")
}

// ListBackups 列出所有备份
func (s *Service) ListBackups() ([]BackupInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.listBackupsInternal()
}

// listBackupsInternal 内部列出备份方法（不加锁）
func (s *Service) listBackupsInternal() ([]BackupInfo, error) {
	files, err := os.ReadDir(s.backupDir)
	if err != nil {
		return nil, err
	}

	var backups []BackupInfo

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), "hosts_backup_") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// 解析时间戳
		timeStr := strings.TrimPrefix(file.Name(), "hosts_backup_")
		timeStr = strings.TrimSuffix(timeStr, ".txt")
		timestamp, err := time.Parse(BackupTimeFormat, timeStr)
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			ID:        timeStr,
			Timestamp: timestamp,
			Size:      info.Size(),
			Path:      filepath.Join(s.backupDir, file.Name()),
		})
	}

	// 按时间倒序排列
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	return backups, nil
}

// RestoreFromBackup 从备份恢复
func (s *Service) RestoreFromBackup(backupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupPath := filepath.Join(s.backupDir, fmt.Sprintf("hosts_backup_%s.txt", backupID))

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupID)
	}

	// 先备份当前状态
	if err := s.backupHostsFileInternal(); err != nil {
		return fmt.Errorf("备份当前状态失败: %w", err)
	}

	// 读取备份文件内容
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %w", err)
	}

	// 直接覆盖原始 Hosts 文件
	if err := os.WriteFile(HostsFilePath, content, 0644); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足: 无法恢复 Hosts 文件。请确保以管理员身份运行程序")
		}
		return fmt.Errorf("恢复失败: %w", err)
	}

	return nil
}

// DeleteBackup 删除指定备份
func (s *Service) DeleteBackup(backupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupPath := filepath.Join(s.backupDir, fmt.Sprintf("hosts_backup_%s.txt", backupID))

	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("删除备份失败: %w", err)
	}

	return nil
}

// cleanOldBackups 清理超过限制的旧备份
func (s *Service) cleanOldBackups() {
	backups, err := s.listBackupsInternal()
	if err != nil {
		return
	}

	if len(backups) <= MaxBackups {
		return
	}

	// 删除最旧的备份
	for i := MaxBackups; i < len(backups); i++ {
		os.Remove(backups[i].Path)
	}
}

// SaveProfile 保存当前配置为方案
func (s *Service) SaveProfile(name string, entries []HostEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	profile := HostsProfile{
		Name:      name,
		Entries:   entries,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 如果方案已存在，保留创建时间
	existingPath := filepath.Join(s.profileDir, name+".json")
	if data, err := os.ReadFile(existingPath); err == nil {
		var existing HostsProfile
		if err := json.Unmarshal(data, &existing); err == nil {
			profile.CreatedAt = existing.CreatedAt
		}
	}

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化方案失败: %w", err)
	}

	profilePath := filepath.Join(s.profileDir, name+".json")
	if err := os.WriteFile(profilePath, data, 0644); err != nil {
		return fmt.Errorf("保存方案失败: %w", err)
	}

	return nil
}

// LoadProfile 加载方案
func (s *Service) LoadProfile(name string) (*HostsProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	profilePath := filepath.Join(s.profileDir, name+".json")

	data, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("读取方案失败: %w", err)
	}

	var profile HostsProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("解析方案失败: %w", err)
	}

	return &profile, nil
}

// ListProfiles 列出所有方案
func (s *Service) ListProfiles() ([]HostsProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := os.ReadDir(s.profileDir)
	if err != nil {
		return nil, err
	}

	var profiles []HostsProfile

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.profileDir, file.Name()))
		if err != nil {
			continue
		}

		var profile HostsProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			continue
		}

		profiles = append(profiles, profile)
	}

	// 按更新时间倒序排列
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].UpdatedAt.After(profiles[j].UpdatedAt)
	})

	return profiles, nil
}

// DeleteProfile 删除方案
func (s *Service) DeleteProfile(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	profilePath := filepath.Join(s.profileDir, name+".json")

	if err := os.Remove(profilePath); err != nil {
		return fmt.Errorf("删除方案失败: %w", err)
	}

	return nil
}

// RenameProfile 重命名方案
func (s *Service) RenameProfile(oldName, newName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 验证新名称
	if strings.TrimSpace(newName) == "" {
		return fmt.Errorf("方案名称不能为空")
	}

	// 检查非法字符
	if strings.ContainsAny(newName, `\/:*?"<>|`) {
		return fmt.Errorf("方案名称包含非法字符")
	}

	// 检查原方案是否存在
	oldPath := filepath.Join(s.profileDir, oldName+".json")
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return fmt.Errorf("方案不存在或已删除")
	}

	// 如果新名称与原名称相同，直接返回
	if oldName == newName {
		return nil
	}

	// 新方案路径
	newPath := filepath.Join(s.profileDir, newName+".json")

	// 如果新名称已存在，删除旧文件
	if _, err := os.Stat(newPath); err == nil {
		if err := os.Remove(newPath); err != nil {
			return fmt.Errorf("删除目标方案失败: %w", err)
		}
	}

	// 重命名文件
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("重命名方案失败: %w", err)
	}

	// 更新 JSON 中的名称和时间戳
	var profile HostsProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		// 如果解析失败，尝试恢复原文件
		os.Rename(newPath, oldPath)
		return fmt.Errorf("解析方案失败: %w", err)
	}

	profile.Name = newName
	profile.UpdatedAt = time.Now()

	updatedData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		// 如果序列化失败，尝试恢复原文件
		os.Rename(newPath, oldPath)
		return fmt.Errorf("序列化方案失败: %w", err)
	}

	if err := os.WriteFile(newPath, updatedData, 0644); err != nil {
		// 如果写入失败，尝试恢复原文件
		os.Rename(newPath, oldPath)
		return fmt.Errorf("保存重命名后的方案失败: %w", err)
	}

	return nil
}

// ApplyProfile 应用方案到 Hosts 文件
func (s *Service) ApplyProfile(name string) error {
	profile, err := s.LoadProfile(name)
	if err != nil {
		return err
	}

	return s.WriteHostsFile(profile.Entries)
}
