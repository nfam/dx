package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "after": string | [string] | [string | [string]] )?
// ********************************************************************
type propAfter []valueAfterList

func parsePropAfter(path string, ctn map[string]any) (propAfter, *errloc) {
	if x, ok := ctn["after"]; ok {
		path += ".after"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []valueAfterList{{{path, str}}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueAfterList
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueAfterList(ipath, it)
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

func (p propAfter) process(input string) (string, *errloc) {
	var err *errloc
	for i := range p {
		input, err = p[i].process(input)
		if err != nil {
			return "", err
		}
	}
	return input, nil
}

type valueAfterList []valueAfter

func parseValueAfterList(path string, x any) (valueAfterList, *errloc) {
	if str, ok := x.(string); ok {
		if str == "" {
			return nil, nil
		}
		return []valueAfter{{path, str}}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueAfter
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if str, ok := it.(string); ok {
				if str != "" {
					result = append(result, valueAfter{ipath, str})
				}
			} else {
				return nil, newerr(msgMustString, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustStringArray, path)
}

func (p valueAfterList) process(input string) (string, *errloc) {
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

type valueAfter struct {
	path  string
	value string
}

func (p *valueAfter) process(input string) (string, *errloc) {
	index := strings.Index(input, p.value)
	if index < 0 {
		return "", newerr(msgNotFound, p.path)
	}
	return input[index+len(p.value):], nil
}
