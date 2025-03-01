package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "sep": string | [string] )?
// ********************************************************************
type propSep []string

func parsePropSep(path string, ctn map[string]any) (propSep, *errloc) {
	if x, ok := ctn["sep"]; ok {
		path += ".sep"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return []string{str}, nil
		} else if array, ok := x.([]any); ok {
			var result []string
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				if str, ok := it.(string); ok {
					if str != "" {
						result = append(result, str)
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

func (p propSep) process(input string) []string {
	result := []string{input}
	for _, sep := range p {
		var iresult []string
		for _, r := range result {
			iresult = append(iresult, strings.Split(r, sep)...)
		}
		result = iresult
	}
	return result
}
