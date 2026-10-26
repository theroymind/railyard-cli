package cli

type UsageError struct {
	Err error
}

func (usageError *UsageError) Error() string {
	return usageError.Err.Error()
}

func (usageError *UsageError) Unwrap() error {
	return usageError.Err
}

func NewUsageError(err error) error {
	return &UsageError{Err: err}
}
