package cli

import (
	"errors"
	"strings"
	"testing"
)

const shell = "sh"

func TestRunChildInReportsTheLastStderrLine(t *testing.T) {
	tests := []struct {
		name          string
		script        string
		expectedCause string
	}{
		{"The last of several lines", `echo first >&2; echo "real cause" >&2; exit 3`, ": real cause"},
		{"Trailing blank lines and spaces are skipped", `echo "  real cause  " >&2; printf '\n\n' >&2; exit 3`, ": real cause"},
		{"No stderr leaves the exit status last", `exit 3`, "exit status 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RunChildIn("", shell, []string{"-c", tt.script}, nil)

			if !errors.Is(err, ErrChildFailed) {
				t.Fatalf("RunChildIn() error = %v, want %v", err, ErrChildFailed)
			}
			if !strings.HasSuffix(err.Error(), tt.expectedCause) {
				t.Errorf("RunChildIn() error = %q, want it to end with %q", err, tt.expectedCause)
			}
		})
	}
}
