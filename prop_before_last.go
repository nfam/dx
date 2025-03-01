package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "beforeLast": string | [string] | [string | [string]] )?
// ********************************************************************
type propBeforeLast []valueBeforeLastList

func parsePropBeforeLast(path string, ctn map[string]any) (propBeforeLast, *errloc) {
	if x, ok := ctn["beforeLast"]; ok {
		path += ".beforeLast"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []valueBeforeLastList{{{path, str}}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueBeforeLastList
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueBeforeLastList(ipath, it)
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

func (p propBeforeLast) process(input string) (string, *errloc) {
	var err *errloc
	for i := range p {
		input, err = p[i].process(input)
		if err != nil {
			return "", err
		}
	}
	return input, nil
}

type valueBeforeLastList []valueBeforeLast

func parseValueBeforeLastList(path string, x any) (valueBeforeLastList, *errloc) {
	if str, ok := x.(string); ok {
		if str == "" {
			return nil, nil
		}
		return []valueBeforeLast{{path, str}}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueBeforeLast
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if str, ok := it.(string); ok {
				if str != "" {
					result = append(result, valueBeforeLast{ipath, str})
				}
			} else {
				return nil, newerr(msgMustString, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustStringArray, path)
}

func (p valueBeforeLastList) process(input string) (string, *errloc) {
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

type valueBeforeLast struct {
	path  string
	value string
}

func (p *valueBeforeLast) process(input string) (string, *errloc) {
	index := strings.LastIndex(input, p.value)
	if index < 0 {
		return "", newerr(msgNotFound, p.path)
	}
	return input[:index], nil
}
