package cli

import (
	"errors"
	"slices"
	"testing"
)

const (
	stringOption   = "type"
	multipleOption = "tag"
	booleanOption  = "draft"
)

var testSpecification = Specification{
	Options: []Option{
		{Name: stringOption, Type: OptionTypeString},
		{Name: multipleOption, Type: OptionTypeString, IsMultiple: true},
		{Name: booleanOption, Type: OptionTypeBoolean},
	},
	AllowsPositionals: true,
}

func TestParse(t *testing.T) {
	tests := []struct {
		name                string
		arguments           []string
		expectedString      string
		expectedStrings     []string
		expectedIsEnabled   bool
		expectedPositionals []string
	}{
		{"Options and positionals interleave", []string{"first", "--type", "task", "second"}, "task", nil, false, []string{"first", "second"}},
		{"An inline value may be empty or start with a dash", []string{"--type=-x"}, "-x", nil, false, nil},
		{"The last value of a single option wins", []string{"--type=a", "--type", "b"}, "b", nil, false, nil},
		{"A multiple option collects every value", []string{"--tag", "a", "--tag="}, "", []string{"a", ""}, false, nil},
		{"A boolean option is enabled by its presence", []string{"--draft", "--draft"}, "", nil, true, nil},
		{"A lone dash is a value, not an option", []string{"--type", "-", "-"}, "-", nil, false, []string{"-"}},
		{"Everything after the terminator is positional", []string{"--", "--type", "-x"}, "", nil, false, []string{"--type", "-x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := Parse(tt.arguments, testSpecification)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.arguments, err)
			}

			if got, _ := parsed.String(stringOption); got != tt.expectedString {
				t.Errorf("String(%q) = %q, want %q", stringOption, got, tt.expectedString)
			}
			if got := parsed.Strings(multipleOption); !slices.Equal(got, tt.expectedStrings) {
				t.Errorf("Strings(%q) = %q, want %q", multipleOption, got, tt.expectedStrings)
			}
			if got := parsed.IsEnabled(booleanOption); got != tt.expectedIsEnabled {
				t.Errorf("IsEnabled(%q) = %v, want %v", booleanOption, got, tt.expectedIsEnabled)
			}
			if !slices.Equal(parsed.Positionals, tt.expectedPositionals) {
				t.Errorf("Positionals = %q, want %q", parsed.Positionals, tt.expectedPositionals)
			}
		})
	}
}

func TestParseRejects(t *testing.T) {
	withoutPositionals := testSpecification
	withoutPositionals.AllowsPositionals = false

	tests := []struct {
		name          string
		arguments     []string
		specification Specification
		expected      error
	}{
		{"An unknown long option", []string{"--bogus"}, testSpecification, ErrUnknownOption},
		{"A short option", []string{"-x"}, testSpecification, ErrUnknownOption},
		{"A negative number as a positional", []string{"-1"}, testSpecification, ErrUnknownOption},
		{"A positional when none are allowed", []string{"extra"}, withoutPositionals, ErrUnexpectedArgument},
		{"A positional after the terminator when none are allowed", []string{"--", "extra"}, withoutPositionals, ErrUnexpectedArgument},
		{"A string option at the end", []string{"--type"}, testSpecification, ErrArgumentMissing},
		{"A string option followed by an option", []string{"--type", "--draft"}, testSpecification, ErrArgumentAmbiguous},
		{"A string option followed by a negative number", []string{"--type", "-1"}, testSpecification, ErrArgumentAmbiguous},
		{"A boolean option with a value", []string{"--draft=false"}, testSpecification, ErrArgumentNotAccepted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.arguments, tt.specification)

			if !errors.Is(err, tt.expected) {
				t.Errorf("Parse(%q) error = %v, want %v", tt.arguments, err, tt.expected)
			}
			var usageError *UsageError
			if !errors.As(err, &usageError) {
				t.Errorf("Parse(%q) error = %v, want a usage error", tt.arguments, err)
			}
		})
	}
}

func TestRequiredString(t *testing.T) {
	parsed, err := Parse(nil, testSpecification)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	_, err = parsed.RequiredString(stringOption)

	var usageError *UsageError
	if !errors.Is(err, ErrRequiredOption) || !errors.As(err, &usageError) {
		t.Errorf("RequiredString(%q) error = %v, want a usage error wrapping %v", stringOption, err, ErrRequiredOption)
	}
}
