package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "has": string | [string] | [string | [string]] )?
// ********************************************************************
type propHas []valueHasList

func parsePropHas(path string, ctn map[string]any) (propHas, *errloc) {
	if x, ok := ctn["has"]; ok {
		path += ".has"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []valueHasList{{parseValueHasItem(path, str)}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueHasList
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueHasList(ipath, it)
				if err != nil {
					return nil, err
				}
				if len(value) > 0 {
					result = append(result, value)
				}
			}
			return result, nil
		}
		return nil, newerr(msgMustStringArray, path)
	}
	return nil, nil
}

func (p propHas) process(input string) *errloc {
	for i := range p {
		if err := p[i].process(input); err != nil {
			return err
		}
	}
	return nil
}

type valueHasList []valueHas

func parseValueHasList(path string, x any) (valueHasList, *errloc) {
	if str, ok := x.(string); ok {
		if str == "" {
			return nil, nil
		}
		return []valueHas{parseValueHasItem(path, str)}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueHas
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if str, ok := it.(string); ok {
				if str != "" {
					result = append(result, parseValueHasItem(ipath, str))
				}
			} else {
				return nil, newerr(msgMustString, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustStringArray, path)
}

const (
	hasHas = "has:"
	hasNot = "not:"
)

func parseValueHasItem(path string, s string) valueHas {
	switch {
	case strings.HasPrefix(s, hasHas):
		return valueHas{path, false, strings.TrimPrefix(s, hasHas)}
	case strings.HasPrefix(s, hasNot):
		return valueHas{path, true, strings.TrimPrefix(s, hasNot)}
	default:
		return valueHas{path, false, s}
	}
}

func (p valueHasList) process(input string) *errloc {
	result := &errloc{}
	for i := range p {
		if err := p[i].process(input); err != nil {
			result.msg = append(result.msg, err.msg...)
		} else {
			return nil
		}
	}
	return result
}

type valueHas struct {
	path  string
	not   bool
	value string
}

func (p *valueHas) process(input string) *errloc {
	if p.not {
		if strings.Contains(input, p.value) {
			return newerr(msgMustNotHas, p.path)
		}
	} else {
		if !strings.Contains(input, p.value) {
			return newerr(msgNotFound, p.path)
		}
	}
	return nil
}
