// Package cli provides the process plumbing every Railyard bin shares: argument parsing,
// subcommand dispatch, usage errors, exit codes and JSON output.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Command func(arguments []string, stdout io.Writer) error

type Program struct {
	Name     string
	Usage    string
	Commands map[string]Command
}

func (program Program) Run(arguments []string, stdout, stderr io.Writer) ExitCode {
	message, exitCode := program.dispatch(arguments, stdout)
	if message != "" {
		_, _ = io.WriteString(stderr, message)
	}
	return exitCode
}

func (program Program) dispatch(arguments []string, stdout io.Writer) (string, ExitCode) {
	if len(arguments) == 0 {
		return program.Usage, ExitCodeUsage
	}
	commandName, commandArguments := arguments[0], arguments[1:]
	command, exists := program.Commands[commandName]
	if !exists {
		return fmt.Sprintf("%s: unknown command %s\n\n%s", program.Name, strconv.Quote(commandName), program.Usage), ExitCodeUsage
	}
	err := command(commandArguments, stdout)
	if err == nil {
		return "", ExitCodeSuccess
	}
	message := fmt.Sprintf("%s %s: %v\n", program.Name, commandName, err)
	if usageError := (*UsageError)(nil); errors.As(err, &usageError) {
		return message + "\n" + program.Usage, ExitCodeUsage
	}
	return message, ExitCodeFailure
}
