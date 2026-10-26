package cli

import (
	"encoding/json"
	"io"
)

const jsonIndent = "  "

func WriteJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", jsonIndent)
	return encoder.Encode(value)
}
