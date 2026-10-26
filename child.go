package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const lineSeparator = "\n"

var ErrChildFailed = errors.New("child process failed")

func RunChild(name string, arguments []string, stdin io.Reader) ([]byte, error) {
	return RunChildIn("", name, arguments, stdin)
}

func RunChildIn(directory, name string, arguments []string, stdin io.Reader) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	child := exec.Command(name, arguments...)
	child.Dir = directory
	child.Stdin = stdin
	child.Stdout = &stdout
	child.Stderr = io.MultiWriter(os.Stderr, &stderr)
	if err := child.Run(); err != nil {
		if cause := lastLine(stderr.String()); cause != "" {
			return stdout.Bytes(), fmt.Errorf("%w: %s: %w: %s", ErrChildFailed, name, err, cause)
		}
		return stdout.Bytes(), fmt.Errorf("%w: %s: %w", ErrChildFailed, name, err)
	}
	return stdout.Bytes(), nil
}

func ReadChildJSON[Value validatable](directory, name string, arguments []string, stdin io.Reader) (Value, error) {
	output, err := RunChildIn(directory, name, arguments, stdin)
	if err != nil {
		var empty Value
		return empty, err
	}
	return ReadJSON[Value](bytes.NewReader(output))
}

func lastLine(text string) string {
	lines := strings.Split(strings.TrimSpace(text), lineSeparator)
	return strings.TrimSpace(lines[len(lines)-1])
}
