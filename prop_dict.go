package dx

import "strconv"

// ********************************************************************
// ( $name: exp | [exp | null] )*
// ********************************************************************
type propDict map[string]valueDict

func parsePropDict(path string, ctn map[string]any) (propDict, *errloc) {
	if x, ok := ctn["dict"]; ok {
		path += ".dict"
		if obj, ok := x.(map[string]any); ok {
			result := make(map[string]valueDict, len(obj))
			for k, v := range obj {
				kpath := path + "[\"" + k + "\"]"
				value, err := parseValueDict(kpath, k, v)
				if err != nil {
					return nil, err
				}
				if value != nil {
					result[k] = value
				}
			}
			return result, nil
		}
		return nil, newerr(msgMustObject, path)
	}
	return nil, nil
}

func (p propDict) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	var err *errloc
	for _, v := range p {
		if err = v.setFunc(slice, convert); err != nil {
			return err
		}
	}
	return nil
}

func (p propDict) extract(input string) (any, *errloc) {
	result := make(map[string]any, len(p))
	for k, v := range p {
		r, err := v.extract(input)
		if err != nil {
			return nil, err
		}
		if r != nil {
			result[k] = r
		}
	}
	return result, nil
}

type valueDict []*valueExp

func parseValueDict(path string, name string, x any) (valueDict, *errloc) {
	if obj, ok := x.(map[string]any); ok {
		value, err := parseValueEpx(path, obj)
		if err != nil {
			return nil, err
		}
		return []*valueExp{value}, nil
	} else if array, ok := x.([]any); ok {
		var result []*valueExp
		for i, it := range array {
			ipath := path + "[" + strconv.Itoa(i) + "]"
			if it == nil {
				result = append(result, nil)
			} else if obj, ok := it.(map[string]any); ok {
				value, err := parseValueEpx(ipath, obj)
				if err != nil {
					return nil, err
				}
				result = append(result, value)
			} else {
				return nil, newerr(msgMustObjectNull, ipath)
			}
		}
		return result, nil
	}
	return nil, newerr(msgMustObjectArray, path)
}

func (p valueDict) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	var err *errloc
	for _, v := range p {
		if v != nil {
			if err = v.setFunc(slice, convert); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p valueDict) extract(input string) (any, *errloc) {
	var (
		merr   = &errloc{}
		result any
		err    *errloc
	)
	for _, v := range p {
		if v == nil {
			return nil, nil
		}
		if result, err = v.extract(input); err != nil {
			merr.append(err)
		} else {
			return result, nil
		}
	}
	return nil, merr
}
