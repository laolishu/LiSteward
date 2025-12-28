package nvm

import (
	"errors"
)
func TestGetNvmVersion_Mock(t *testing.T) {
	// 备份原始 ExecCommand
	orig := ExecCommand
	defer func() { ExecCommand = orig }()

	tests := []struct {
		name      string
		mockOut   string
		mockErr   error
		expectVer string
		expectErr bool
	}{
		{
			name:      "正常返回版本号",
			mockOut:   "v0.39.5\n",
			mockErr:   nil,
			expectVer: "v0.39.5",
			expectErr: false,
		},
		{
			name:      "无输出",
			mockOut:   "",
			mockErr:   nil,
			expectVer: "",
			expectErr: true,
		},
		{
			name:      "命令执行失败",
			mockOut:   "",
			mockErr:   errors.New("not found"),
			expectVer: "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExecCommand = func(name string, args ...string) *exec.Cmd {
				return mockCmdGetNvmVersion(tt.mockOut, tt.mockErr)
			}
			ver, err := GetNvmVersion()
			if (err != nil) != tt.expectErr {
				t.Errorf("期望错误: %v, 实际: %v", tt.expectErr, err)
			}
			if ver != tt.expectVer {
				t.Errorf("期望版本: %q, 实际: %q", tt.expectVer, ver)
			}
		})
	}
}

// mockCmdGetNvmVersion 返回一个模拟 exec.Cmd，CombinedOutput返回指定内容
func mockCmdGetNvmVersion(output string, err error) *exec.Cmd {
	return &exec.Cmd{
		// 利用Run/CombinedOutput的hook
	}
}

// 重写 CombinedOutput 方法（仅用于测试）
func (c *exec.Cmd) CombinedOutput() ([]byte, error) {
	// 仅用于TestGetNvmVersion_Mock
	if c.ProcessState == nil {
		// 通过环境变量传递输出和错误
		out := os.Getenv("MOCK_CMD_OUT")
		errStr := os.Getenv("MOCK_CMD_ERR")
		if errStr != "" {
			return []byte(out), errors.New(errStr)
		}
		return []byte(out), nil
	}
	return nil, errors.New("not mocked")
}
package nvm

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDetect_WithNvmDirEnv(t *testing.T) {
	orig := os.Getenv("NVM_DIR")
	defer os.Setenv("NVM_DIR", orig)

	tmp := t.TempDir()
	if err := os.MkdirAll(tmp, 0755); err != nil {
		t.Fatalf("failed to create tmp dir: %v", err)
	}
	if err := os.Setenv("NVM_DIR", tmp); err != nil {
		t.Fatalf("failed to set NVM_DIR: %v", err)
	}

	ok, dir := Detect()
	if !ok {
		t.Fatalf("expected Detect to return true when NVM_DIR is set")
	}
	if dir != tmp {
		t.Fatalf("expected Detect to return dir %s, got %s", tmp, dir)
	}
}

func TestListVersions_FromNvmDir(t *testing.T) {
	tmp := t.TempDir()
	v1 := filepath.Join(tmp, "v14.17.0")
	v2 := filepath.Join(tmp, "v16.0.0")
	if err := os.MkdirAll(v1, 0755); err != nil {
		t.Fatalf("failed to create v1: %v", err)
	}
	if err := os.MkdirAll(v2, 0755); err != nil {
		t.Fatalf("failed to create v2: %v", err)
	}

	list, err := ListVersions(tmp)
	if err != nil {
		t.Fatalf("ListVersions returned error: %v", err)
	}

	// expect at least 2 entries with these versions
	foundV1, foundV2 := false, false
	for _, v := range list {
		if v.Version == "v14.17.0" || v.Version == "14.17.0" {
			foundV1 = true
		}
		if v.Version == "v16.0.0" || v.Version == "16.0.0" {
			foundV2 = true
		}
	}
	if !foundV1 || !foundV2 {
		t.Fatalf("expected to find both versions in list: %#v", list)
	}
}

func TestGetNvmDir_FallbackHome(t *testing.T) {
	// backup envs
	origNvm := os.Getenv("NVM_DIR")
	defer os.Setenv("NVM_DIR", origNvm)

	// ensure NVM_DIR is unset
	_ = os.Unsetenv("NVM_DIR")

	tmpHome := t.TempDir()
	nvmDir := filepath.Join(tmpHome, ".nvm")
	if err := os.MkdirAll(nvmDir, 0755); err != nil {
		t.Fatalf("failed to create nvm dir: %v", err)
	}

	// set appropriate home env based on platform
	if runtime.GOOS == "windows" {
		origProfile := os.Getenv("USERPROFILE")
		defer os.Setenv("USERPROFILE", origProfile)
		if err := os.Setenv("USERPROFILE", tmpHome); err != nil {
			t.Fatalf("failed to set USERPROFILE: %v", err)
		}
	} else {
		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)
		if err := os.Setenv("HOME", tmpHome); err != nil {
			t.Fatalf("failed to set HOME: %v", err)
		}
	}

	got := GetNvmDir("")
	if got == "" {
		t.Fatalf("expected GetNvmDir to find fallback .nvm in home, got empty string")
	}
	if filepath.Clean(got) != filepath.Clean(nvmDir) {
		t.Fatalf("expected GetNvmDir to return %s, got %s", nvmDir, got)
	}
}

func TestInstallUseUninstall_MockedExec(t *testing.T) {
	// backup and restore ExecCommand
	orig := ExecCommand
	defer func() { ExecCommand = orig }()

	// mock ExecCommand to echo arguments
	ExecCommand = func(name string, arg ...string) *exec.Cmd {
		// produce a simple echo of joined args
		out := strings.Join(append([]string{name}, arg...), " ")
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", "echo", out)
		}
		return exec.Command("sh", "-c", "echo '"+out+"'")
	}

	// InstallVersion
	out, err := InstallVersion("14.17.0")
	if err != nil {
		t.Fatalf("InstallVersion returned error: %v", err)
	}
	if !strings.Contains(out, "nvm install 14.17.0") && !strings.Contains(out, "nvm 14.17.0") {
		t.Fatalf("unexpected output from InstallVersion: %s", out)
	}

	// UseVersion
	out, err = UseVersion("14.17.0")
	if err != nil {
		t.Fatalf("UseVersion returned error: %v", err)
	}
	if !strings.Contains(out, "nvm use 14.17.0") {
		t.Fatalf("unexpected output from UseVersion: %s", out)
	}

	// UninstallVersion
	out, err = UninstallVersion("14.17.0")
	if err != nil {
		t.Fatalf("UninstallVersion returned error: %v", err)
	}
	if !strings.Contains(out, "nvm uninstall 14.17.0") {
		t.Fatalf("unexpected output from UninstallVersion: %s", out)
	}
}
func TestParseNvmLsOutput_WithAsteriskAndParentheses(t *testing.T) {
	// Test output similar to actual `nvm ls` with asterisk and parentheses
	nvmOutput := `v14.17.0
v16.0.0
* 20.19.4 (Currently using 64-bit executable)
v18.20.5 (some note here)
->  v20.4.0
system`

	result := parseNvmLsOutput(nvmOutput)

	expectedVersions := map[string]bool{
		"v14.17.0": false,
		"v16.0.0":  false,
		"20.19.4":  true, // 标记为当前版本
		"v18.20.5": false,
		"v20.4.0":  false,
	}

	// Check that we have all expected versions (allow for v20.4.0 or variations)
	if len(result) < 5 {
		t.Fatalf("expected at least 5 versions, got %d: %v", len(result), result)
	}

	// Verify all expected versions were found and current flag is correct
	foundVersions := make(map[string]bool)
	for _, v := range result {
		foundVersions[v.Version] = v.IsCurrent
	}

	// Verify all expected versions were found with correct IsCurrent status
	for ver, shouldBeCurrent := range expectedVersions {
		if isCurrent, found := foundVersions[ver]; !found {
			t.Errorf("expected version %s not found in parsed result %v", ver, foundVersions)
		} else if isCurrent != shouldBeCurrent {
			t.Errorf("version %s has IsCurrent=%v, expected %v", ver, isCurrent, shouldBeCurrent)
		}
	}
}
