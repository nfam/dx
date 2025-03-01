package dx

import "strings"

// errloc represents an error at expression path.
type errloc struct {
	msg []string
}

// Error implements error interface.
func (e *errloc) Error() string {
	return strings.Join(e.msg, "\n")
}

func (e *errloc) append(n *errloc) {
	e.msg = append(e.msg, n.msg...)
}

func newerr(msg string, path string) *errloc {
	if path == "" {
		return &errloc{[]string{msg}}
	}
	return &errloc{[]string{msg, path}}
}

const (
	msgMustArray             = "must be an array"
	msgMustObject            = "must be an object"
	msgMustObjectArray       = "must be an object or array"
	msgMustObjectNull        = "must be an object or null"
	msgMustString            = "must be a string"
	msgMustStringArray       = "must be a string or array"
	msgMustStringObject      = "must be a string or object"
	msgMustStringObjectArray = "must be a string, object, or array"
	msgMustStringIntegerNull = "must be a string, integer, or null"
	msgMustNotEmpty          = "must not be empty"
	msgMustNotHas            = "must not has"
	msgNotFound              = "not found"
)
