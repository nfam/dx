package dx

import (
	"encoding/json"
)

// SliceFunc represents custom slicer.
// If the operation is not supported, return errors.ErrUnsupported.
type SliceFunc func(input string, name string) (string, error)

// Convertert represents custom coverter to process final value.
// If the operation is not supported, return errors.ErrUnsupported.
type ConvertFunc func(input string, name string) (any, error)

// Expression represents a data extraction handler.
type Expression struct {
	json.Unmarshaler
	exp valueExp
}

// New returns the parsed `Expression` from given json data.
func New(data []byte, slice SliceFunc, convert ConvertFunc) (*Expression, error) {
	e := &Expression{}
	if err := e.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	if err := e.exp.setFunc(slice, convert); err != nil {
		return nil, err
	}
	return e, nil
}

// Must returns the parsed `Expression` from given json data, panic if fail.
func Must(data []byte, slice SliceFunc, convert ConvertFunc) *Expression {
	e, err := New(data, slice, convert)
	if err != nil {
		panic(err)
	}
	return e
}

// UnmarshalJSON implements `json.Unmarshaler` interface.
func (e *Expression) UnmarshalJSON(data []byte) error {
	ctn := make(map[string]any)
	if err := json.Unmarshal(data, &ctn); err != nil {
		return err
	}
	if ctn == nil {
		return newerr(msgMustObject, "")
	}
	if value, err := parseValueEpx("", ctn); err != nil {
		return err
	} else {
		e.exp = *value
	}
	return nil
}

// Extract returns result from extracting given input.
func (e *Expression) Extract(input string) (any, error) {
	result, err := e.exp.extract(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}
