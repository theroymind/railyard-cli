package strictjson

import (
	"errors"
	"testing"
)

type target struct {
	Name string `json:"name"`
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedName  string
		expectedError error
	}{
		{"Decodes one value with surrounding whitespace", " {\"name\":\"a\"}\n", "a", nil},
		{"Rejects empty input", "  \n", "", ErrNoValue},
		{"Rejects a second value", `{"name":"a"} {"name":"b"}`, "a", ErrTrailingData},
		{"Rejects trailing garbage", `{"name":"a"} x`, "a", ErrTrailingData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded target

			err := Unmarshal([]byte(tt.input), &decoded)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Unmarshal(%q) error = %v, want %v", tt.input, err, tt.expectedError)
			}
			if decoded.Name != tt.expectedName {
				t.Errorf("Unmarshal(%q) name = %q, want %q", tt.input, decoded.Name, tt.expectedName)
			}
		})
	}
}

func TestUnmarshalRejectsLooseFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"An unknown field", `{"name":"a","extra":1}`},
		{"A field name in another case", `{"Name":"a"}`},
		{"A duplicate field", `{"name":"a","name":"b"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded target

			if err := Unmarshal([]byte(tt.input), &decoded); err == nil {
				t.Errorf("Unmarshal(%q) error = nil, want an error", tt.input)
			}
		})
	}
}
