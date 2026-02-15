package cmd

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestExecute_VersionFlag(t *testing.T) {
	err := Execute([]string{"--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_HelpFlag(t *testing.T) {
	err := Execute([]string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_UnknownCommand(t *testing.T) {
	err := Execute([]string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestExecute_JSONAndPlainConflict(t *testing.T) {
	err := Execute([]string{"--json", "--plain", "version"})
	if err == nil {
		t.Fatal("expected error for conflicting output flags")
	}
}

func TestVersionCmd_Run(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &VersionCmd{}
	err := cmd.Run(context.TODO())

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "brightlocal-cli") {
		t.Errorf("expected output to contain 'brightlocal-cli', got: %s", output)
	}

	if !strings.Contains(output, "OS:") {
		t.Errorf("expected output to contain OS info, got: %s", output)
	}
}

func TestExitError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ExitError
		wantStr  string
		wantCode int
	}{
		{
			name:     "nil error",
			err:      nil,
			wantStr:  "",
			wantCode: 0,
		},
		{
			name:     "with wrapped error",
			err:      &ExitError{Code: 1, Err: os.ErrNotExist},
			wantStr:  "file does not exist",
			wantCode: 1,
		},
		{
			name:     "code only",
			err:      &ExitError{Code: 2},
			wantStr:  "exit code 2",
			wantCode: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantStr {
				t.Errorf("Error() = %q, want %q", got, tt.wantStr)
			}
		})
	}
}

func TestExitError_Unwrap(t *testing.T) {
	inner := os.ErrNotExist
	err := &ExitError{Code: 1, Err: inner}

	if !errors.Is(err.Unwrap(), inner) {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), inner)
	}

	var nilErr *ExitError
	if nilErr.Unwrap() != nil {
		t.Errorf("nil.Unwrap() should be nil")
	}
}

func TestNewParser(t *testing.T) {
	parser, cli, err := newParser("test description")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if parser == nil {
		t.Fatal("parser should not be nil")
	}

	if cli == nil {
		t.Fatal("cli should not be nil")
	}
}

func TestParseCLIStructure(t *testing.T) {
	_, cli, err := newParser("test")
	if err != nil {
		t.Fatalf("newParser: %v", err)
	}

	// Verify CLI struct has expected command groups
	_ = cli.Auth
	_ = cli.Locations
	_ = cli.Rankings
	_ = cli.Citations
	_ = cli.Reports
	_ = cli.VersionCmd
}

func TestEnvOr(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		envVal   string
		fallback string
		want     string
	}{
		{
			name:     "env set",
			key:      "BRIGHTLOCAL_CLI_TEST_VAR_1",
			envVal:   "from_env",
			fallback: "default",
			want:     "from_env",
		},
		{
			name:     "env empty",
			key:      "BRIGHTLOCAL_CLI_TEST_VAR_2",
			envVal:   "",
			fallback: "default",
			want:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				t.Setenv(tt.key, tt.envVal)
			}

			got := envOr(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("envOr(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestBoolString(t *testing.T) {
	if boolString(true) != "true" {
		t.Error("boolString(true) should return 'true'")
	}

	if boolString(false) != "false" {
		t.Error("boolString(false) should return 'false'")
	}
}
