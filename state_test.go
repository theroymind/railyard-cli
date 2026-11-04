package cli

import (
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

func TestRepositoryKey(t *testing.T) {
	hash := "[0-9a-f]{" + strconv.Itoa(keyHashLength) + "}"

	tests := []struct {
		name       string
		repository string
		pattern    string
	}{
		{"Names the key after the repository directory", "/Users/someone/Development/railyard", "^railyard" + keyNameSeparator + hash + "$"},
		{"Replaces unsafe characters in the name", "/tmp/my repo:2", "^my" + keyNameSeparator + "repo" + keyNameSeparator + "2" + keyNameSeparator + hash + "$"},
		{"Falls back to the hash alone when no name is left", "/", "^" + hash + "$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RepositoryKey(tt.repository); !regexp.MustCompile(tt.pattern).MatchString(got) {
				t.Errorf("RepositoryKey(%q) = %q, want a match for %s", tt.repository, got, tt.pattern)
			}
		})
	}
}

func TestRepositoryKeySeparatesRepositoriesWithTheSameName(t *testing.T) {
	first, second := RepositoryKey("/work/railyard"), RepositoryKey("/personal/railyard")

	if first == second {
		t.Errorf("RepositoryKey() = %q for two repositories, want different keys", first)
	}
	if RepositoryKey("/work/railyard") != first {
		t.Errorf("RepositoryKey() is not stable for the same repository")
	}
}

func TestStateDirectory(t *testing.T) {
	home := t.TempDir()
	absolute := filepath.Join(t.TempDir(), "state")

	tests := []struct {
		name      string
		stateHome string
		expected  string
	}{
		{"Uses an absolute XDG_STATE_HOME", absolute, absolute},
		{"Falls back to the home directory when unset", "", filepath.Join(home, defaultStateHome)},
		{"Ignores a relative XDG_STATE_HOME", "relative" + string(filepath.Separator) + "state", filepath.Join(home, defaultStateHome)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOME", home)
			t.Setenv(StateHomeVariable, tt.stateHome)

			got, err := StateDirectory()
			if err != nil {
				t.Fatalf("StateDirectory() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("StateDirectory() = %q, want %q", got, tt.expected)
			}
		})
	}
}
