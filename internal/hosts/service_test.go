package hosts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateIP(t *testing.T) {
	s, _ := NewService("test")

	tests := []struct {
		ip    string
		valid bool
	}{
		{"192.168.1.1", true},
		{"127.0.0.1", true},
		{"999.999.999.999", false},
		{"invalid", false},
		{"::1", true}, // IPv6
	}

	for _, tt := range tests {
		result := s.isValidIP(tt.ip)
		if result != tt.valid {
			t.Errorf("isValidIP(%s) = %v, want %v", tt.ip, result, tt.valid)
		}
	}
}

func TestValidateDomain(t *testing.T) {
	s, _ := NewService("test")

	tests := []struct {
		domain string
		valid  bool
	}{
		{"localhost", true},
		{"example.com", true},
		{"dev.example.com", true},
		{"test.local", true},
		{"", false},
	}

	for _, tt := range tests {
		result := s.isValidDomain(tt.domain)
		if result != tt.valid {
			t.Errorf("isValidDomain(%s) = %v, want %v", tt.domain, result, tt.valid)
		}
	}
}

func TestRenameProfile(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	s, err := NewService(tmpDir)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}

	// 保存一个初始方案
	entries := []HostEntry{
		{IP: "127.0.0.1", Domain: "localhost", Enabled: true},
	}
	if err := s.SaveProfile("开发环境", entries); err != nil {
		t.Fatalf("保存初始方案失败: %v", err)
	}

	// 测试成功重命名
	if err := s.RenameProfile("开发环境", "生产环境"); err != nil {
		t.Errorf("重命名方案失败: %v", err)
	}

	// 验证新方案文件存在
	newPath := filepath.Join(s.profileDir, "生产环境.json")
	if _, err := os.Stat(newPath); err != nil {
		t.Errorf("新方案文件不存在: %v", err)
	}

	// 验证旧方案文件不存在
	oldPath := filepath.Join(s.profileDir, "开发环境.json")
	if _, err := os.Stat(oldPath); err == nil {
		t.Error("旧方案文件仍然存在")
	}

	// 验证内容正确
	profile, err := s.LoadProfile("生产环境")
	if err != nil {
		t.Fatalf("加载重命名后的方案失败: %v", err)
	}
	if profile.Name != "生产环境" {
		t.Errorf("方案名称不正确，期望 '生产环境'，得到 '%s'", profile.Name)
	}
}

func TestRenameProfileNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := NewService(tmpDir)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}

	// 测试重命名不存在的方案
	if err := s.RenameProfile("不存在的方案", "新名称"); err == nil {
		t.Error("期望重命名不存在的方案时返回错误")
	}
}

func TestRenameProfileEmptyName(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := NewService(tmpDir)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}

	// 保存一个初始方案
	entries := []HostEntry{
		{IP: "127.0.0.1", Domain: "localhost", Enabled: true},
	}
	if err := s.SaveProfile("测试方案", entries); err != nil {
		t.Fatalf("保存初始方案失败: %v", err)
	}

	// 测试空名称
	if err := s.RenameProfile("测试方案", ""); err == nil {
		t.Error("期望空名称时返回错误")
	}

	// 测试仅空格
	if err := s.RenameProfile("测试方案", "   "); err == nil {
		t.Error("期望仅空格时返回错误")
	}
}

func TestRenameProfileSameName(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := NewService(tmpDir)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}

	// 保存一个初始方案
	entries := []HostEntry{
		{IP: "127.0.0.1", Domain: "localhost", Enabled: true},
	}
	if err := s.SaveProfile("测试方案", entries); err != nil {
		t.Fatalf("保存初始方案失败: %v", err)
	}

	// 测试相同名称（应返回 nil）
	if err := s.RenameProfile("测试方案", "测试方案"); err != nil {
		t.Errorf("重命名为相同名称应返回 nil，但得到: %v", err)
	}
}

// ==================== 扩展格式测试（AlternateIP + Description） ====================

func TestParseLineNewFormat(t *testing.T) {
	s, _ := NewService("test")

	tests := []struct {
		line           string
		expectedIP     string
		expectedDomain string
		expectedAltIP  string
		expectedDesc   string
	}{
		// 新格式：IP DOMAIN # ALT_IP#DESCRIPTION
		{"127.0.0.1 www.example.com # 192.168.1.100#生产环境", "127.0.0.1", "www.example.com", "192.168.1.100", "生产环境"},
		{"127.0.0.1 www.example.com # 192.168.1.100#", "127.0.0.1", "www.example.com", "192.168.1.100", ""},
		{"127.0.0.1 www.example.com # null#生产环境", "127.0.0.1", "www.example.com", "", "生产环境"},
		{"127.0.0.1 www.example.com # null#", "127.0.0.1", "www.example.com", "", ""},
	}

	for _, tt := range tests {
		entry := s.parseLine(tt.line)
		if entry == nil {
			t.Errorf("parseLine(%q) 返回 nil", tt.line)
			continue
		}

		if entry.IP != tt.expectedIP {
			t.Errorf("parseLine(%q): IP = %s, want %s", tt.line, entry.IP, tt.expectedIP)
		}
		if entry.Domain != tt.expectedDomain {
			t.Errorf("parseLine(%q): Domain = %s, want %s", tt.line, entry.Domain, tt.expectedDomain)
		}
		if entry.AlternateIP != tt.expectedAltIP {
			t.Errorf("parseLine(%q): AlternateIP = %s, want %s", tt.line, entry.AlternateIP, tt.expectedAltIP)
		}
		if entry.Description != tt.expectedDesc {
			t.Errorf("parseLine(%q): Description = %s, want %s", tt.line, entry.Description, tt.expectedDesc)
		}
	}
}

func TestParseLineOldFormat(t *testing.T) {
	s, _ := NewService("test")

	tests := []struct {
		line            string
		expectedIP      string
		expectedDomain  string
		expectedComment string
	}{
		// 旧格式：IP DOMAIN # COMMENT（向后兼容）
		{"127.0.0.1 www.example.com # 这是一条注释", "127.0.0.1", "www.example.com", "这是一条注释"},
		{"127.0.0.1 www.example.com # 备注", "127.0.0.1", "www.example.com", "备注"},
		{"127.0.0.1 www.example.com", "127.0.0.1", "www.example.com", ""},
	}

	for _, tt := range tests {
		entry := s.parseLine(tt.line)
		if entry == nil {
			t.Errorf("parseLine(%q) 返回 nil", tt.line)
			continue
		}

		if entry.IP != tt.expectedIP {
			t.Errorf("parseLine(%q): IP = %s, want %s", tt.line, entry.IP, tt.expectedIP)
		}
		if entry.Domain != tt.expectedDomain {
			t.Errorf("parseLine(%q): Domain = %s, want %s", tt.line, entry.Domain, tt.expectedDomain)
		}
		if entry.Comment != tt.expectedComment {
			t.Errorf("parseLine(%q): Comment = %s, want %s", tt.line, entry.Comment, tt.expectedComment)
		}
	}
}

func TestGenerateHostsContentNewFormat(t *testing.T) {
	s, _ := NewService("test")

	entries := []HostEntry{
		{
			IP:          "127.0.0.1",
			Domain:      "www.example.com",
			AlternateIP: "192.168.1.100",
			Description: "生产环境",
			Enabled:     true,
		},
		{
			IP:          "127.0.0.1",
			Domain:      "www.test.com",
			AlternateIP: "",
			Description: "开发环境",
			Enabled:     true,
		},
		{
			IP:      "127.0.0.1",
			Domain:  "www.old.com",
			Comment: "旧格式注释",
			Enabled: true,
		},
	}

	content := s.generateHostsContent(entries)

	// 验证内容包含所有条目和正确的格式
	if !contains(content, "127.0.0.1\twww.example.com\t# 192.168.1.100#生产环境") {
		t.Errorf("未找到预期的新格式条目（包含 AlternateIP 和 Description）")
	}
	if !contains(content, "127.0.0.1\twww.test.com\t# null#开发环境") {
		t.Errorf("未找到预期的新格式条目（null AlternateIP）")
	}
	if !contains(content, "127.0.0.1\twww.old.com\t# 旧格式注释") {
		t.Errorf("未找到预期的旧格式条目")
	}
}

func TestGenerateHostsContentDisabledEntries(t *testing.T) {
	s, _ := NewService("test")

	entries := []HostEntry{
		{
			IP:          "127.0.0.1",
			Domain:      "www.example.com",
			AlternateIP: "192.168.1.100",
			Description: "生产环境",
			Enabled:     false,
		},
	}

	content := s.generateHostsContent(entries)

	// 验证禁用条目以 # 开头
	if !contains(content, "# 127.0.0.1\twww.example.com\t# 192.168.1.100#生产环境") {
		t.Errorf("禁用条目的格式不正确")
	}
}

func TestParseAndGenerateRoundtrip(t *testing.T) {
	s, _ := NewService("test")

	// 原始行
	originalLine := "127.0.0.1 www.example.com # 192.168.1.100#生产环境"

	// 解析
	entry := s.parseLine(originalLine)
	if entry == nil {
		t.Fatalf("解析行失败")
	}

	// 生成
	content := s.generateHostsContent([]HostEntry{*entry})

	// 验证生成的内容包含相同的信息
	if !contains(content, "127.0.0.1\twww.example.com\t# 192.168.1.100#生产环境") {
		t.Errorf("往返转换（解析→生成）失败：原始=%q，生成=%q", originalLine, content)
	}
}

func TestSwapIPs(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := NewService(tmpDir)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}

	entries := []HostEntry{
		{IP: "127.0.0.1", Domain: "www.example.com", AlternateIP: "192.168.1.100", Enabled: true},
		{IP: "127.0.0.1", Domain: "www.test.com", AlternateIP: "", Enabled: true},
	}

	if err := s.WriteHostsFile(entries); err != nil {
		t.Fatalf("写入 hosts 文件失败: %v", err)
	}

	// 交换第一个条目的 IP
	result, err := s.SwapIPs(entries, 0)
	if err != nil {
		t.Errorf("交换 IP 失败: %v", err)
	}

	if result[0].IP != "192.168.1.100" || result[0].AlternateIP != "127.0.0.1" {
		t.Errorf("交换后 IP 不正确：IP=%s, AlternateIP=%s", result[0].IP, result[0].AlternateIP)
	}

	// 测试交换没有备用 IP 的条目
	_, err = s.SwapIPs(entries, 1)
	if err == nil {
		t.Error("应该返回错误：没有备用 IP")
	}

	// 测试无效索引
	_, err = s.SwapIPs(entries, 999)
	if err == nil {
		t.Error("应该返回错误：索引超出范围")
	}
}

// 辅助函数
func contains(text, substr string) bool {
	return len(text) > 0 && len(substr) > 0 && (text == substr || len(text) > len(substr) && (text[:len(substr)] == substr || text[len(text)-len(substr):] == substr || findSubstring(text, substr)))
}

func findSubstring(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
