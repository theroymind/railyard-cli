// Package strictjson provides the strict JSON decoding every contract boundary uses: one value,
// no unknown or duplicate fields, exact-case field names and nothing after it.
package strictjson

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
)

var (
	ErrNoValue      = errors.New("no JSON value")
	ErrTrailingData = errors.New("unexpected data after the JSON value")
)

func Decode(reader io.Reader, target any) error {
	decoder := jsontext.NewDecoder(reader)
	value, err := decoder.ReadValue()
	if errors.Is(err, io.EOF) {
		return ErrNoValue
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(value, target, json.RejectUnknownMembers(true)); err != nil {
		return err
	}
	if _, err := decoder.ReadToken(); !errors.Is(err, io.EOF) {
		return ErrTrailingData
	}
	return nil
}

func Unmarshal(text []byte, target any) error {
	return Decode(bytes.NewReader(text), target)
}
