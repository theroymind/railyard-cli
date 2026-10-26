package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type OptionType string

const (
	OptionTypeString  OptionType = "string"
	OptionTypeBoolean OptionType = "boolean"
)

const (
	optionMarker    = "-"
	longOptionStart = "--"
	terminator      = "--"
	valueSeparator  = "="
)

var (
	ErrUnknownOption       = errors.New("unknown option")
	ErrUnexpectedArgument  = errors.New("unexpected argument")
	ErrArgumentMissing     = errors.New("argument missing")
	ErrArgumentAmbiguous   = errors.New("argument is ambiguous")
	ErrArgumentNotAccepted = errors.New("does not take an argument")
	ErrRequiredOption      = errors.New("is required")
)

type Option struct {
	Name       string
	Type       OptionType
	IsMultiple bool
}

type Specification struct {
	Options           []Option
	AllowsPositionals bool
}

type Arguments struct {
	Positionals []string
	strings     map[string][]string
	booleans    map[string]bool
}

func (arguments Arguments) String(name string) (string, bool) {
	values := arguments.strings[name]
	if len(values) == 0 {
		return "", false
	}
	return values[len(values)-1], true
}

func (arguments Arguments) RequiredString(name string) (string, error) {
	value, exists := arguments.String(name)
	if !exists {
		return "", NewUsageError(fmt.Errorf("%s%s %w", longOptionStart, name, ErrRequiredOption))
	}
	return value, nil
}

func (arguments Arguments) Strings(name string) []string {
	return arguments.strings[name]
}

func (arguments Arguments) IsEnabled(name string) bool {
	return arguments.booleans[name]
}

func Parse(arguments []string, specification Specification) (Arguments, error) {
	parsed := Arguments{strings: map[string][]string{}, booleans: map[string]bool{}}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == terminator {
			for _, positional := range arguments[index+1:] {
				if err := parsed.addPositional(positional, specification); err != nil {
					return Arguments{}, err
				}
			}
			break
		}
		if !isOptionLike(argument) {
			if err := parsed.addPositional(argument, specification); err != nil {
				return Arguments{}, err
			}
			continue
		}
		rawName, value, hasValue := strings.Cut(argument, valueSeparator)
		option, isKnown := specification.find(rawName)
		if !isKnown {
			return Arguments{}, NewUsageError(unknownOption(rawName, specification))
		}
		switch option.Type {
		case OptionTypeBoolean:
			if hasValue {
				return Arguments{}, NewUsageError(fmt.Errorf("option '%s' %w", rawName, ErrArgumentNotAccepted))
			}
			parsed.booleans[option.Name] = true
		case OptionTypeString:
			if !hasValue {
				if index+1 == len(arguments) {
					return Arguments{}, NewUsageError(fmt.Errorf("option '%s <value>' %w", rawName, ErrArgumentMissing))
				}
				index++
				value = arguments[index]
				if isOptionLike(value) {
					return Arguments{}, NewUsageError(fmt.Errorf(
						"option '%s' %w.\nDid you forget to specify the option argument for '%s'?\n"+
							"To specify an option argument starting with a dash use '%s=-XYZ'",
						rawName, ErrArgumentAmbiguous, rawName, rawName,
					))
				}
			}
			if option.IsMultiple {
				parsed.strings[option.Name] = append(parsed.strings[option.Name], value)
			} else {
				parsed.strings[option.Name] = []string{value}
			}
		}
	}
	return parsed, nil
}

func (arguments *Arguments) addPositional(positional string, specification Specification) error {
	if !specification.AllowsPositionals {
		return NewUsageError(fmt.Errorf(
			"%w '%s'. This command does not take positional arguments", ErrUnexpectedArgument, positional,
		))
	}
	arguments.Positionals = append(arguments.Positionals, positional)
	return nil
}

func (specification Specification) find(rawName string) (Option, bool) {
	name, isLong := strings.CutPrefix(rawName, longOptionStart)
	if !isLong {
		return Option{}, false
	}
	index := slices.IndexFunc(specification.Options, func(option Option) bool { return option.Name == name })
	if index < 0 {
		return Option{}, false
	}
	return specification.Options[index], true
}

func unknownOption(rawName string, specification Specification) error {
	if !specification.AllowsPositionals {
		return fmt.Errorf("%w '%s'", ErrUnknownOption, rawName)
	}
	return fmt.Errorf(
		"%w '%s'. To specify a positional argument starting with a '-', "+
			"place it at the end of the command after '--', as in '-- \"%s\"",
		ErrUnknownOption, rawName, rawName,
	)
}

func isOptionLike(argument string) bool {
	return len(argument) > len(optionMarker) && strings.HasPrefix(argument, optionMarker)
}
