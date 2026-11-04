package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/theroymind/railyard-cli/strictjson"
)

const jsonIndent = "  "

var ErrInvalidInput = errors.New("invalid input")

type validatable interface {
	Validate() error
}

func WriteJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", jsonIndent)
	return encoder.Encode(value)
}

func ReadJSON[Value validatable](reader io.Reader) (Value, error) {
	var value Value
	if err := strictjson.Decode(reader, &value); err != nil {
		return value, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}
	if err := value.Validate(); err != nil {
		return value, fmt.Errorf("%w:\n%w", ErrInvalidInput, err)
	}
	return value, nil
}
