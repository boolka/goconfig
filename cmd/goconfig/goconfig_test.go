package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

	for _, testCase := range testCases {
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

	for _, testCase := range testCases {
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

	for _, testCase := range testCases {
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

	for _, testCase := range testCases {
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
	}
}
