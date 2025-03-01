package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "afterLast": string | [string] | [string | [string]] )?
// ********************************************************************
type propAfterLast []valueAfterLastList

func parsePropAfterLast(path string, ctn map[string]any) (propAfterLast, *errloc) {
	if x, ok := ctn["afterLast"]; ok {
		path += ".afterLast"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []valueAfterLastList{{{path, str}}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueAfterLastList
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueAfterLastList(ipath, it)
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

func (p propAfterLast) process(input string) (string, *errloc) {
	var err *errloc
	for i := range p {
		input, err = p[i].process(input)
		if err != nil {
			return "", err
		}
	}
	return input, nil
}

type valueAfterLastList []valueAfterLast

func parseValueAfterLastList(path string, x any) (valueAfterLastList, *errloc) {
	if str, ok := x.(string); ok {
		if str == "" {
			return nil, nil
		}
		return []valueAfterLast{{path, str}}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueAfterLast
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if str, ok := it.(string); ok {
				if str != "" {
					result = append(result, valueAfterLast{ipath, str})
				}
			} else {
				return nil, newerr(msgMustString, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustStringArray, path)
}

func (p valueAfterLastList) process(input string) (string, *errloc) {
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

type valueAfterLast struct {
	path  string
	value string
}

func (p *valueAfterLast) process(input string) (string, *errloc) {
	index := strings.LastIndex(input, p.value)
	if index < 0 {
		return "", newerr(msgNotFound, p.path)
	}
	return input[index+len(p.value):], nil
}
