package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	StateHomeVariable    = "XDG_STATE_HOME"
	StateApplication     = "railyard"
	defaultStateHome     = ".local/state"
	keyNameSeparator     = "-"
	keyHashLength        = 16
	unsafeKeyNamePattern = `[^A-Za-z0-9._-]+`
)

var unsafeKeyNameExpression = regexp.MustCompile(unsafeKeyNamePattern)

// RepositoryKey names a repository's file under the state directory: its directory name, with unsafe characters
// replaced, and a hash of its absolute path, so that two repositories with the same name never collide.
func RepositoryKey(repository string) string {
	digest := sha256.Sum256([]byte(repository))
	hash := hex.EncodeToString(digest[:])[:keyHashLength]
	name := unsafeKeyNameExpression.ReplaceAllString(filepath.Base(repository), keyNameSeparator)
	name = strings.Trim(name, keyNameSeparator)
	if name == "" {
		return hash
	}
	return name + keyNameSeparator + hash
}

// StateDirectory is $XDG_STATE_HOME, or ~/.local/state when it is unset or relative, as the XDG Base Directory
// spec requires.
func StateDirectory() (string, error) {
	if stateHome := os.Getenv(StateHomeVariable); filepath.IsAbs(stateHome) {
		return stateHome, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultStateHome), nil
}
