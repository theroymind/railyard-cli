package cli

type ExitCode int

const (
	ExitCodeSuccess ExitCode = 0
	ExitCodeFailure ExitCode = 1
	ExitCodeUsage   ExitCode = 2
)
