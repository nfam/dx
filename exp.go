package dx

import "strings"

type extractor interface {
	setFunc(slice SliceFunc, convert ConvertFunc) *errloc
	extract(input string) (any, *errloc)
}

type valueExp struct {
	has       propHas
	slice     propSlice
	extractor extractor
}

func parseValueEpx(path string, ctn map[string]any) (*valueExp, *errloc) {
	var (
		exp = &valueExp{}
		err *errloc
	)
	if exp.has, err = parsePropHas(path, ctn); err != nil {
		return nil, err
	}
	if exp.slice, err = parsePropSlice(path, ctn); err != nil {
		return nil, err
	}
	{
		var next propNext
		if next, err = parsePropNext(path, ctn); err != nil {
			return nil, err
		}
		if next != nil {
			exp.extractor = next
		}
	}
	if exp.extractor == nil {
		if value, ok := parsePropValue(path, ctn); ok {
			exp.extractor = value
		}
	}
	if exp.extractor == nil {
		if list, err := parsePropList(path, ctn); err != nil {
			return nil, err
		} else if list != nil {
			exp.extractor = list
		}
	}
	if exp.extractor == nil {
		if dict, err := parsePropDict(path, ctn); err != nil {
			return nil, err
		} else if dict != nil {
			exp.extractor = dict
		}
	}
	if exp.extractor == nil {
		if convert, err := parsePropConvert(path, ctn); err != nil {
			return nil, err
		} else if convert != nil {
			exp.extractor = convert
		}
	}
	for name := range ctn {
		if !strings.HasPrefix(name, "//") {
			switch name {
			case "has", "slice", "next":
			case "convert":
				if _, ok := ctn["dict"]; ok {
					return nil, newerr("unwanted property "+name+" as dict provided", path)
				}
				fallthrough
			case "dict":
				if _, ok := ctn["list"]; ok {
					return nil, newerr("unwanted property "+name+" as list provided", path)
				}
				fallthrough
			case "list":
				if _, ok := ctn["value"]; ok {
					return nil, newerr("unwanted property "+name+" as value provided", path)
				}
				fallthrough
			case "value":
				if _, ok := ctn["next"]; ok {
					return nil, newerr("unwanted property "+name+" as next provided", path)
				}
			default:
				return nil, newerr("unknown property "+name, path)
			}
		}
	}
	return exp, nil
}

func (p *valueExp) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	var err *errloc
	if p.slice != nil {
		if err = p.slice.setFunc(slice); err != nil {
			return err
		}
	}
	if p.extractor != nil {
		if err = p.extractor.setFunc(slice, convert); err != nil {
			return err
		}
	}
	return nil
}

func (p *valueExp) extract(input string) (any, *errloc) {
	var err *errloc
	if p.has != nil {
		if err = p.has.process(input); err != nil {
			return nil, err
		}
	}
	if p.slice != nil {
		if input, err = p.slice.process(input); err != nil {
			return nil, err
		}
	}
	if p.extractor != nil {
		return p.extractor.extract(input)
	}
	return input, nil
}
