package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/devsy-org/devsy/pkg/log"
	"github.com/spf13/cobra"
)

const (
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
	testProbe  = "probe"
)

func TestRootLoggingConfiguration(t *testing.T) {
	for _, tc := range []struct {
		level string
		want  string
	}{
		{"", infoLevel},
		{errorLevel, errorLevel},
		{warnLevel, warnLevel},
		{infoLevel, infoLevel},
		{"debug", "debug"},
		{"trace", "trace"},
		{"WaRn", warnLevel},
	} {
		t.Run("level_"+tc.level, func(t *testing.T) {
			t.Setenv("DEVSY_LOG_LEVEL", tc.level)
			t.Setenv("DEVSY_LOG_FORMAT", "")
			ran := false
			root := NewRootCmd()
			root.AddCommand(
				&cobra.Command{Use: testProbe, Run: func(*cobra.Command, []string) { ran = true }},
			)
			root.SetArgs([]string{testProbe})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if !ran || log.LevelString() != tc.want {
				t.Fatalf(
					"command ran = %v, effective level = %q, want %q",
					ran,
					log.LevelString(),
					tc.want,
				)
			}
		})
	}
}

func TestRootDiagnosticStream(t *testing.T) {
	if os.Getenv("GCLOUD_LOG_TEST_CHILD") == "1" {
		root := NewRootCmd()
		root.AddCommand(&cobra.Command{Use: testProbe, RunE: func(*cobra.Command, []string) error {
			log.Error("diagnostic-canary ", os.Getenv("GCLOUD_PROVIDER_TOKEN"))
			_, err := fmt.Fprint(os.Stdout, "protocol-canary")
			return err
		}})
		root.SetArgs([]string{testProbe})
		if err := root.Execute(); err != nil {
			t.Fatal(err)
		}
		return
	}

	for _, format := range []string{"", "auto", "json", "JsOn"} {
		t.Run("format_"+format, func(t *testing.T) {
			child := rootDiagnosticChild(t)
			child.Env = append(
				os.Environ(),
				"GCLOUD_LOG_TEST_CHILD=1",
				"DEVSY_LOG_LEVEL=warn",
				"DEVSY_LOG_FORMAT="+format,
				"GCLOUD_PROVIDER_TOKEN=known-private-token-canary",
			)
			var stdout, stderr bytes.Buffer
			child.Stdout = &stdout
			child.Stderr = &stderr
			if err := child.Run(); err != nil {
				t.Fatalf("child failed: %v, stderr: %s", err, stderr.String())
			}
			assertRootDiagnostic(t, stderr.Bytes())
			if !strings.Contains(stdout.String(), "protocol-canary") ||
				strings.Contains(stdout.String(), "diagnostic-canary") {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
		})
	}
}

func TestRootRejectsInvalidLoggingBeforeCommand(t *testing.T) {
	for _, name := range []string{"DEVSY_LOG_LEVEL", "DEVSY_LOG_FORMAT"} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("DEVSY_LOG_LEVEL", infoLevel)
			t.Setenv("DEVSY_LOG_FORMAT", "json")
			const secret = "secret-canary-invalid-value"
			t.Setenv(name, secret)
			ran := false
			root := NewRootCmd()
			root.AddCommand(
				&cobra.Command{Use: testProbe, Run: func(*cobra.Command, []string) { ran = true }},
			)
			root.SetArgs([]string{testProbe})
			err := root.Execute()
			if err == nil || !strings.Contains(err.Error(), name) ||
				strings.Contains(err.Error(), secret) ||
				ran {
				t.Fatalf("command ran = %v, error = %v", ran, err)
			}
		})
	}
}

func assertRootDiagnostic(t *testing.T, output []byte) {
	t.Helper()
	var record struct {
		Level string `json:"level"`
		Msg   string `json:"msg"`
		TS    string `json:"ts"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(output), &record); err != nil {
		t.Fatalf("invalid diagnostic JSON: %v, stderr: %s", err, output)
	}
	if record.Level != errorLevel || !strings.Contains(record.Msg, "diagnostic-canary") ||
		strings.Contains(record.Msg, "known-private-token-canary") ||
		record.TS == "" {
		t.Fatalf("unexpected diagnostic: %+v", record)
	}
}

func rootDiagnosticChild(t *testing.T) *exec.Cmd {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return &exec.Cmd{
		Path: executable,
		Args: []string{executable, "-test.run=^TestRootDiagnosticStream$"},
	}
}
