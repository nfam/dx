package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "slice": slice | [slice] | [slice | [slice]] )?
// ********************************************************************
type propSlice []valueSliceList

func parsePropSlice(path string, ctn map[string]any) (propSlice, *errloc) {
	if x, ok := ctn["slice"]; ok {
		path += ".slice"
		if obj, ok := x.(map[string]any); ok {
			value, err := parseValueSlice(path, obj)
			if err != nil {
				return nil, err
			}
			return []valueSliceList{{*value}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueSliceList
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueSliceList(ipath, it)
				if err != nil {
					return nil, err
				}
				if len(value) > 0 {
					result = append(result, value)
				}
			}
			return result, nil
		}
		return nil, newerr(msgMustObjectArray, path)
	}
	return nil, nil
}

func (p propSlice) setFunc(slice SliceFunc) *errloc {
	for i := range p {
		if err := p[i].setFunc(slice); err != nil {
			return err
		}
	}
	return nil
}

func (p propSlice) process(input string) (string, *errloc) {
	var err *errloc
	for i := range p {
		if input, err = p[i].process(input); err != nil {
			return "", err
		}
	}
	return input, nil
}

type valueSliceList []valueSlice

func parseValueSliceList(path string, x any) (valueSliceList, *errloc) {
	if obj, ok := x.(map[string]any); ok {
		value, err := parseValueSlice(path, obj)
		if err != nil {
			return nil, err
		}
		return []valueSlice{*value}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueSlice
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if obj, ok := it.(map[string]any); ok {
				value, err := parseValueSlice(ipath, obj)
				if err != nil {
					return nil, err
				}
				result = append(result, *value)
			} else {
				return nil, newerr(msgMustObject, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustObjectArray, path)
}

func (p valueSliceList) setFunc(slice SliceFunc) *errloc {
	for i := range p {
		if err := p[i].setFunc(slice); err != nil {
			return err
		}
	}
	return nil
}

func (p valueSliceList) process(input string) (string, *errloc) {
	var (
		merr   = &errloc{}
		result string
		err    *errloc
	)
	for i := range p {
		if result, err = p[i].process(input); err != nil {
			merr.append(err)
		} else {
			return result, nil
		}
	}
	return "", merr
}

type valueSlice struct {
	after      propAfter
	before     propBefore
	beforeLast propBeforeLast
	afterLast  propAfterLast
	join       propJoin
	_func      propFunc
}

func parseValueSlice(path string, ctn map[string]any) (*valueSlice, *errloc) {
	slice := &valueSlice{}
	var err *errloc
	if slice.after, err = parsePropAfter(path, ctn); err != nil {
		return nil, err
	}
	if slice.before, err = parsePropBefore(path, ctn); err != nil {
		return nil, err
	}
	if slice.beforeLast, err = parsePropBeforeLast(path, ctn); err != nil {
		return nil, err
	}
	if slice.afterLast, err = parsePropAfterLast(path, ctn); err != nil {
		return nil, err
	}
	if slice.join, err = parsePropJoin(path, ctn); err != nil {
		return nil, err
	}
	if slice._func, err = parsePropFunc(path, ctn); err != nil {
		return nil, err
	}
	for name := range ctn {
		if !strings.HasPrefix(name, "//") {
			switch name {
			case "after", "before", "beforeLast", "afterLast", "join", "func":
			default:
				return nil, newerr("unknown property "+name, path)
			}
		}
	}
	return slice, nil
}

func (p *valueSlice) setFunc(slice SliceFunc) *errloc {
	if p.join != nil {
		if err := p.join.setFunc(slice); err != nil {
			return err
		}
	}
	if p._func != nil {
		if err := p._func.setFunc(slice); err != nil {
			return err
		}
	}
	return nil
}

func (p *valueSlice) process(input string) (string, *errloc) {
	var err *errloc
	if p.after != nil {
		if input, err = p.after.process(input); err != nil {
			return "", err
		}
	}
	if p.before != nil {
		if input, err = p.before.process(input); err != nil {
			return "", err
		}
	}
	if p.beforeLast != nil {
		if input, err = p.beforeLast.process(input); err != nil {
			return "", err
		}
	}
	if p.afterLast != nil {
		if input, err = p.afterLast.process(input); err != nil {
			return "", err
		}
	}
	if p.join != nil {
		if input, err = p.join.process(input); err != nil {
			return "", err
		}
	}
	if p._func != nil {
		if input, err = p._func.process(input); err != nil {
			return "", err
		}
	}
	return input, nil
}
