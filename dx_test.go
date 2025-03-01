package dx

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const testdataDir = "testdata"

func TestDx(t *testing.T) {
	files, err := os.ReadDir("./testdata")
	if err != nil {
		t.Error(err)
		return
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(testdataDir, file.Name()))
		if err != nil {
			t.Error(err)
			return
		}
		t.Run(file.Name(), func(t *testing.T) {
			helpTest(t, data, helpSlice, helpConvertert)
		})
	}
}

func TestMust(t *testing.T) {
	// ok
	func() {
		defer func() {
			const expect = ""
			var result string
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					result = err.Error()
				}
			}
			if result != expect {
				t.Error("\nexpect: " + expect + "\nresult: " + result)
			}
		}()
		_ = Must([]byte("{}"), nil, nil)
	}()
	// panic
	func() {
		defer func() {
			const expect = `json: cannot unmarshal number into Go value of type map[string]interface {}`
			var result string
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					result = err.Error()
				}
			}
			if result != expect {
				t.Error("\nexpect: " + expect + "\nresult: " + result)
			}
		}()
		_ = Must([]byte("1"), nil, nil)
	}()
}

func helpTest(t *testing.T, testdata []byte, slice SliceFunc, convert ConvertFunc) {
	t.Helper()
	var cases []struct {
		Name   string          `json:"name"`
		Exp    json.RawMessage `json:"exp"`
		Input  string          `json:"input"`
		Expect any             `json:"expect"`
	}
	err := json.Unmarshal(testdata, &cases)
	if err != nil {
		t.Error(err)
		return
	}
	for _, c := range cases {
		var result any
		exp, err := New(c.Exp, slice, convert)
		if err != nil {
			result = err.Error()
		} else {
			if c.Name == "item empty" {
				println("debug")
			}
			result, err = exp.Extract(c.Input)
			if err != nil {
				result = err.Error()
			} else {
				data, _ := json.Marshal(result)
				_ = json.Unmarshal(data, &result)
			}
		}
		if diff := cmp.Diff(c.Expect, result); diff != "" {
			t.Error(c.Name + "\n" + diff)
		}
	}
}

func helpSlice(input string, name string) (string, error) {
	switch name {
	case "found":
		return "found", nil
	case "error":
		return "", errors.New("func error")
	default:
		return "", errors.ErrUnsupported
	}
}

func helpConvertert(input string, name string) (any, error) {
	switch name {
	case "true":
		return true, nil
	case "error":
		return nil, errors.New("convert error")
	default:
		return nil, errors.ErrUnsupported
	}
}
