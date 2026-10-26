// Package strictjson provides the strict JSON decoding every contract boundary uses: one value,
// no unknown fields and nothing after it.
package strictjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrNoValue      = errors.New("no JSON value")
	ErrTrailingData = errors.New("unexpected data after the JSON value")
)

func Decode(reader io.Reader, target any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			return ErrNoValue
		}
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return ErrTrailingData
	}
	return nil
}

func Unmarshal(text []byte, target any) error {
	return Decode(bytes.NewReader(text), target)
}
