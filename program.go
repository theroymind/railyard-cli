// Package cli provides the process plumbing every Railyard bin shares: argument parsing,
// subcommand dispatch, usage errors, exit codes and JSON in and out.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// helpCommand prints a program's usage on stdout, for every program that doesn't define its own.
const helpCommand = "help"

type Command func(arguments []string, stdin io.Reader, stdout io.Writer) error

type Program struct {
	Name     string
	Usage    string
	Commands map[string]Command
}

func (program Program) Run(arguments []string, stdin io.Reader, stdout, stderr io.Writer) ExitCode {
	message, exitCode := program.dispatch(arguments, stdin, stdout)
	if message != "" {
		_, _ = io.WriteString(stderr, message)
	}
	return exitCode
}

// dispatch runs the command the first argument names and returns what to print on stderr and the exit code. No
// arguments print the usage as a usage error; the built-in help command prints it on stdout and succeeds, unless
// the program defines its own; an unknown command says so in one line and points at help.
func (program Program) dispatch(arguments []string, stdin io.Reader, stdout io.Writer) (string, ExitCode) {
	if len(arguments) == 0 {
		return program.Usage, ExitCodeUsage
	}
	commandName, commandArguments := arguments[0], arguments[1:]
	command, exists := program.Commands[commandName]
	if !exists && commandName == helpCommand {
		_, _ = io.WriteString(stdout, program.Usage)
		return "", ExitCodeSuccess
	}
	if !exists {
		return fmt.Sprintf("%s: unknown command %s; run %q for usage\n", program.Name, strconv.Quote(commandName), program.Name+" "+helpCommand), ExitCodeUsage
	}
	err := command(commandArguments, stdin, stdout)
	if err == nil {
		return "", ExitCodeSuccess
	}
	message := fmt.Sprintf("%s %s: %v\n", program.Name, commandName, err)
	if usageError := (*UsageError)(nil); errors.As(err, &usageError) {
		return message + "\n" + program.Usage, ExitCodeUsage
	}
	return message, ExitCodeFailure
}
