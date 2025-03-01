package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "before": string | [string] | [string | [string]] )?
// ********************************************************************
type propBefore []valueBeforeList

func parsePropBefore(path string, ctn map[string]any) (propBefore, *errloc) {
	if x, ok := ctn["before"]; ok {
		path += ".before"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []valueBeforeList{{{path, str}}}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueBeforeList
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				value, err := parseValueBeforeList(ipath, it)
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

func (p propBefore) process(input string) (string, *errloc) {
	var err *errloc
	for i := range p {
		input, err = p[i].process(input)
		if err != nil {
			return "", err
		}
	}
	return input, nil
}

type valueBeforeList []valueBefore

func parseValueBeforeList(path string, x any) (valueBeforeList, *errloc) {
	if str, ok := x.(string); ok {
		if str == "" {
			return nil, nil
		}
		return []valueBefore{{path, str}}, nil
	} else if array, ok := x.([]any); ok {
		var result []valueBefore
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if str, ok := it.(string); ok {
				if str != "" {
					result = append(result, valueBefore{ipath, str})
				}
			} else {
				return nil, newerr(msgMustString, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustStringArray, path)
}

func (p valueBeforeList) process(input string) (string, *errloc) {
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

type valueBefore struct {
	path  string
	value string
}

func (p *valueBefore) process(input string) (string, *errloc) {
	index := strings.Index(input, p.value)
	if index < 0 {
		return "", newerr(msgNotFound, p.path)
	}
	return input[:index], nil
}
