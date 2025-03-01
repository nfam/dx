package dx

import (
	"strconv"
)

// ********************************************************************
// ( "item": exp | [exp] )?
// ********************************************************************
type propItem []valueExp

func parsePropItem(path string, ctn map[string]any) (propItem, *errloc) {
	if x, ok := ctn["item"]; ok {
		path += ".item"
		if obj, ok := x.(map[string]any); ok {
			value, err := parseValueEpx(path, obj)
			if err != nil {
				return nil, err
			}
			return []valueExp{*value}, nil
		} else if array, ok := x.([]any); ok {
			var result []valueExp
			for i, it := range array {
				ipath := path + "[" + strconv.Itoa(i) + "]"
				if obj, ok := it.(map[string]any); ok {
					value, err := parseValueEpx(ipath, obj)
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
	return nil, nil
}

func (p propItem) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	var err *errloc
	for i := range p {
		if err = p[i].setFunc(slice, convert); err != nil {
			return err
		}
	}
	return nil
}

func (p propItem) extract(input string) (any, *errloc) {
	var (
		merr   = &errloc{}
		result any
		err    *errloc
	)
	for i := range p {
		if result, err = p[i].extract(input); err != nil {
			merr.append(err)
		} else {
			return result, nil
		}
	}
	return nil, merr
}
