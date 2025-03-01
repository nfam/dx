package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "join": [string | join | [join]] )?
// ********************************************************************
type propJoin []valueJoinList

func parsePropJoin(path string, ctn map[string]any) (propJoin, *errloc) {
	if x, ok := ctn["join"]; ok {
		path += ".join"
		if array, ok := x.([]any); ok {
			result := []valueJoinList{}
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueJoinList(ipath, it)
				if err != nil {
					return nil, err
				}
				if len(value) > 0 {
					result = append(result, value)
				}
			}
			return result, nil
		}
		return nil, newerr(msgMustArray, path)
	}
	return nil, nil
}

func (p propJoin) setFunc(slice SliceFunc) *errloc {
	var err *errloc
	for i := range p {
		if err = p[i].setFunc(slice); err != nil {
			return err
		}
	}
	return nil
}

func (p propJoin) process(input string) (string, *errloc) {
	var (
		b   strings.Builder
		s   string
		err *errloc
	)
	for i := range p {
		if s, err = p[i].process(input); err != nil {
			return "", err
		}
		b.WriteString(s)
	}
	return b.String(), nil
}

type valueJoinList []valueJoin

func parseValueJoinList(path string, x any) (valueJoinList, *errloc) {
	if str, ok := x.(string); ok {
		return []valueJoin{{path, nil, nil, &str}}, nil
	} else if obj, ok := x.(map[string]any); ok {
		value, err := parseValueJoin(path, obj)
		if err != nil {
			return nil, err
		}
		return []valueJoin{*value}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueJoin
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if obj, ok := it.(map[string]any); ok {
				value, err := parseValueJoin(ipath, obj)
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
	return nil, newerr(msgMustStringObjectArray, path)
}

func (p valueJoinList) setFunc(slice SliceFunc) *errloc {
	var err *errloc
	for i := range p {
		if err = p[i].setFunc(slice); err != nil {
			return err
		}
	}
	return nil
}

func (p valueJoinList) process(input string) (string, *errloc) {
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

type valueJoin struct {
	path  string
	has   propHas
	slice propSlice
	value *string
}

func parseValueJoin(path string, ctn map[string]any) (*valueJoin, *errloc) {
	var (
		join = &valueJoin{}
		err  *errloc
	)
	if join.has, err = parsePropHas(path, ctn); err != nil {
		return nil, err
	}
	if join.slice, err = parsePropSlice(path, ctn); err != nil {
		return nil, err
	}
	if x, ok := ctn["value"]; ok {
		if str, ok := x.(string); ok {
			join.value = &str
		} else {
			return nil, newerr(msgMustString, path+".value")
		}
	}
	for name := range ctn {
		if !strings.HasPrefix(name, "//") {
			switch name {
			case "has", "slice", "value":
			default:
				return nil, newerr("unknown property "+name, path)
			}
		}
	}
	return join, nil
}

func (p *valueJoin) setFunc(slice SliceFunc) *errloc {
	if p.slice != nil {
		return p.slice.setFunc(slice)
	}
	return nil
}

func (p *valueJoin) process(input string) (string, *errloc) {
	var err *errloc
	if p.has != nil {
		if err = p.has.process(input); err != nil {
			return "", err
		}
	}
	if p.slice != nil {
		if input, err = p.slice.process(input); err != nil {
			return "", err
		}
	}
	if p.value != nil {
		return *p.value, nil
	}
	return input, nil
}
