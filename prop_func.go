package dx

import (
	"errors"
	"strconv"
	"strings"
)

// ********************************************************************
// ( "func": string | [string] )?
// ********************************************************************
type propFunc []valueFunc

func parsePropFunc(path string, ctn map[string]any) (propFunc, *errloc) {
	if x, ok := ctn["func"]; ok {
		path += ".func"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []valueFunc{{path, nil, str}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueFunc
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				if str, ok := it.(string); ok {
					if str != "" {
						result = append(result, valueFunc{ipath, nil, str})
					}
				} else {
					return nil, newerr(msgMustString, ipath)
				}
			}
			return result, nil
		}
		return nil, newerr(msgMustStringArray, path)
	}
	return nil, nil
}

func (p propFunc) setFunc(slice SliceFunc) *errloc {
	var err *errloc
	for i := range p {
		if err = p[i].setFunc(slice); err != nil {
			return err
		}
	}
	return nil
}

func (p propFunc) process(input string) (string, *errloc) {
	var err *errloc
	for i := range p {
		input, err = p[i].process(input)
		if err != nil {
			return "", err
		}
	}
	return input, nil
}

type valueFunc struct {
	path  string
	_func SliceFunc
	name  string
}

func (p *valueFunc) setFunc(slice SliceFunc) *errloc {
	p._func = slice
	if _, err := p.process(""); err != nil {
		if err.msg[0] == "unsupported func "+p.name {
			return err
		}
	}
	return nil
}

func (p *valueFunc) process(input string) (string, *errloc) {
	if p._func != nil {
		if d, err := p._func(input, p.name); err != nil {
			if err != errors.ErrUnsupported {
				return "", newerr(err.Error(), p.path)
			}
		} else {
			return d, nil
		}
	}
	switch p.name {
	case "trim":
		return strings.TrimSpace(input), nil
	default:
		if index := strings.Index(p.name, ":"); index > 0 {
			name := p.name[:index]
			arg := p.name[index+1:]
			switch name {
			case "trimPrefix":
				return strings.TrimPrefix(input, arg), nil
			case "trimSuffix":
				return strings.TrimSuffix(input, arg), nil
			case "trimLeft":
				return strings.TrimLeft(input, arg), nil
			case "trimRight":
				return strings.TrimRight(input, arg), nil
			}
		}
		return "", newerr("unsupported func "+p.name, p.path)
	}
}
