package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var ErrChildFailed = errors.New("child process failed")

func RunChild(name string, arguments []string, stdin io.Reader) ([]byte, error) {
	var stdout bytes.Buffer
	child := exec.Command(name, arguments...)
	child.Stdin = stdin
	child.Stdout = &stdout
	child.Stderr = os.Stderr
	if err := child.Run(); err != nil {
		return stdout.Bytes(), fmt.Errorf("%w: %s: %w", ErrChildFailed, name, err)
	}
	return stdout.Bytes(), nil
}
