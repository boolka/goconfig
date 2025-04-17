package main_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	vaultMock "github.com/boolka/goconfig/pkg/vault"
)

func TmpConfigDir(t *testing.T) string {
	d, err := os.MkdirTemp(os.TempDir(), "config")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(d)
	})

	return d
}

func CreateConfigFile(dir, name, content string) error {
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}

	f.Write([]byte(content))
	f.Close()

	return nil
}

func TestGoconfigGet(t *testing.T) {
	t.Parallel()

	d := TmpConfigDir(t)

	CreateConfigFile(d, "default.toml", `field="value"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=field"},
		{"go", "run", "./goconfig.go", "--config", d, "--get", "field"},
		{"go", "run", "./goconfig.go", "-c", d, "-g", "field"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigGet(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != "value" {
				t.Fatal(stdout.String(), stderr.String(), "value")
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigDeployment(t *testing.T) {
	t.Parallel()

	d := TmpConfigDir(t)

	CreateConfigFile(d, "default.toml", `field="value"`)
	CreateConfigFile(d, "production.toml", `field="production_value"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--deployment=production", "--config=" + d, "--get=field"},
		{"go", "run", "./goconfig.go", "--deployment", "production", "--config", d, "--get", "field"},
		{"go", "run", "./goconfig.go", "-d", "production", "-c", d, "-g", "field"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigDeployment(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != "production_value" {
				t.Fatal(stdout.String(), stderr.String(), "production_value")
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigInstance(t *testing.T) {
	t.Parallel()

	d := TmpConfigDir(t)

	CreateConfigFile(d, "default.toml", `field="value"`)
	CreateConfigFile(d, "default-1.toml", `field="value-instance-1"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--instance=1", "--config=" + d, "--get=field"},
		{"go", "run", "./goconfig.go", "--instance", "1", "--config", d, "--get", "field"},
		{"go", "run", "./goconfig.go", "-i", "1", "-c", d, "-g", "field"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigInstance(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != "value-instance-1" {
				t.Fatal(stdout.String(), stderr.String(), "value-instance-1")
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigHostname(t *testing.T) {
	t.Parallel()

	d := TmpConfigDir(t)

	CreateConfigFile(d, "default.toml", `field="value"`)
	CreateConfigFile(d, "local-machine.toml", `field="local-machine"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--hostname=local-machine", "-c", d, "-g", "field"},
		{"go", "run", "./goconfig.go", "--hostname", "local-machine", "-c", d, "-g", "field"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigHostname(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != "local-machine" {
				t.Fatal(stdout.String(), stderr.String(), "local-machine")
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigVaultConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(vaultServer.Close)
	vaultClient := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	err := vaultClient.WriteSecret(ctx, "secret", "goconfig_cmd_secret", map[string]any{
		"password": "abc123",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = vaultClient.DeleteSecret(ctx, "secret", "goconfig_cmd_secret")
		if err != nil {
			t.Fatal(err)
		}
	})

	d := TmpConfigDir(t)

	vaultCfgFile := fmt.Sprintf("[goconfig.vault]\naddress=\"%s\"\n[goconfig.vault.auth]\ntoken=\"root\"\n", vaultServer.URL)
	CreateConfigFile(d, "default.toml", vaultCfgFile)
	CreateConfigFile(d, "vault.toml", `password="secret,goconfig_cmd_secret"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=password"},
		{"go", "run", "./goconfig.go", "--config", d, "--get", "password"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigVaultConfig(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != "abc123" {
				t.Fatal(stdout.String(), stderr.String(), "abc123")
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigVaultDirectToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(vaultServer.Close)
	vaultClient := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	err := vaultClient.WriteSecret(ctx, "secret", "goconfig_cmd_secret", map[string]any{
		"password": "abc123",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = vaultClient.DeleteSecret(ctx, "secret", "goconfig_cmd_secret")
		if err != nil {
			t.Fatal(err)
		}
	})

	d := TmpConfigDir(t)

	vaultCfgFile := fmt.Sprintf("[goconfig.vault]\naddress=\"%s\"", vaultServer.URL)
	CreateConfigFile(d, "default.toml", vaultCfgFile)
	CreateConfigFile(d, "vault.toml", `password="secret,goconfig_cmd_secret"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=password", "--token", "root"},
		{"go", "run", "./goconfig.go", "--config", d, "--get=password", "--token=root"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigVaultDirectToken(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != "abc123" {
				t.Fatal(stdout.String(), stderr.String(), "abc123")
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigVerbose(t *testing.T) {
	t.Parallel()

	d := TmpConfigDir(t)

	CreateConfigFile(d, "default.toml", `field="value"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=field", "-v"},
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=field", "--verbose"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigVerbose(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if !strings.Contains(stdout.String(), "module=github.com/boolka/goconfig") {
				t.Fatal(stdout.String())
			}

			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigVerboseError(t *testing.T) {
	t.Parallel()

	d := TmpConfigDir(t)

	CreateConfigFile(d, "default.toml", `field="value"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=empty", "-v"},
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=empty", "--verbose"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigVerboseError(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if !strings.Contains(stdout.String(), "module=github.com/boolka/goconfig") {
				t.Fatal(stdout.String())
			}

			if !strings.Contains(stderr.String(), "key not found") {
				t.Fatal(stderr.String())
			}
		})
	}
}

func TestGoconfigVaultTokenError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vaultServer := vaultMock.NewServer("root")
	t.Cleanup(vaultServer.Close)
	vaultClient := vaultMock.NewClient(vaultServer.URL, "root", vaultServer.Client())

	err := vaultClient.WriteSecret(ctx, "secret", "goconfig_cmd_secret", map[string]any{
		"password": "abc123",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = vaultClient.DeleteSecret(ctx, "secret", "goconfig_cmd_secret")
		if err != nil {
			t.Fatal(err)
		}
	})

	d := TmpConfigDir(t)

	vaultCfgFile := fmt.Sprintf("password=\"qwerty123456\"\n[goconfig.vault]\naddress=\"%s\"", vaultServer.URL)
	CreateConfigFile(d, "default.toml", vaultCfgFile)
	CreateConfigFile(d, "vault.toml", `password="secret,goconfig_cmd_secret"`)

	testCases := [][]string{
		{"go", "run", "./goconfig.go", "--config=" + d, "--get=password", "--token", "root1"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestGoconfigVaultTokenError(%d)", i), func(t *testing.T) {
			cmd := exec.Command(testCase[0], testCase[1:]...)

			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatal(err, stderr.String())
			}

			if stdout.String() != `qwerty123456` {
				t.Fatal(stdout.String())
			}
			if stderr.String() != "" {
				t.Fatal(stderr.String())
			}
		})
	}
}
