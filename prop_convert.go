package dx

import (
	"errors"
	"strconv"
	"time"
)

type converter interface {
	extract(input string) (any, *errloc)
}

// ********************************************************************
// ( "convert":  string )?
// ********************************************************************
type propConvert struct {
	path      string
	converter converter
	convert   ConvertFunc
	name      string
}

func parsePropConvert(path string, ctn map[string]any) (*propConvert, *errloc) {
	if x, ok := ctn["convert"]; ok {
		path += ".convert"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, nil
			}
			return &propConvert{path, nil, nil, str}, nil
		} else if obj, ok := x.(map[string]any); ok {
			if name, ok := obj["name"]; ok {
				kpath := path + ".name"
				if name, ok := name.(string); ok {
					switch name {
					case "":
						return nil, newerr(msgMustNotEmpty, kpath)
					case timeUnix, timeUnixMilli:
						if conv, err := parseValueConvertTime(path, obj, name); err != nil {
							return nil, err
						} else {
							return &propConvert{path, conv, nil, ""}, nil
						}
					}
					return nil, newerr("unsupported convert "+name, kpath)
				}
				return nil, newerr(msgMustString, kpath)
			}
			return nil, newerr("missing property name", path)
		}
		return nil, newerr(msgMustStringObject, path)
	}
	return nil, nil
}

func (p *propConvert) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	if p.converter == nil {
		p.convert = convert
		if _, err := p.extract(""); err != nil {
			if err.msg[0] == "unsupported convert "+p.name {
				return err
			}
		}
	}
	return nil
}

func (p *propConvert) extract(input string) (any, *errloc) {
	if p.converter != nil {
		if v, err := p.converter.extract(input); err != nil {
			return nil, newerr(err.Error(), p.path)
		} else {
			return v, nil
		}
	} else if p.convert != nil {
		if v, err := p.convert(input, p.name); err != nil {
			if err != errors.ErrUnsupported {
				return nil, newerr(err.Error(), p.path)
			}
		} else {
			return v, nil
		}
	}
	switch p.name {
	case "atoi":
		if v, err := strconv.Atoi(input); err != nil {
			return nil, newerr(err.Error(), p.path)
		} else {
			return v, nil
		}
	case "atof":
		if v, err := strconv.ParseFloat(input, 64); err != nil {
			return nil, newerr(err.Error(), p.path)
		} else {
			return v, nil
		}
	default:
		return nil, newerr("unsupported convert "+p.name, p.path)
	}
}

const (
	timeUnix      = "timeUnix"
	timeUnixMilli = "timeUnixMilli"
)

type valueConvertTime struct {
	path     string
	name     string
	layout   []string
	location *time.Location
}

func parseValueConvertTime(path string, ctn map[string]any, name string) (*valueConvertTime, *errloc) {
	result := &valueConvertTime{name: name}
	if x, ok := ctn["layout"]; ok {
		kpath := path + ".layout"
		if str, ok := x.(string); ok {
			if str == "" {
				return nil, newerr(msgMustNotEmpty, kpath)
			}
			result.layout = []string{str}
		} else if array, ok := x.([]any); ok {
			for i, it := range array {
				ipath := kpath + "[" + strconv.Itoa(i) + "]"
				if str, ok := it.(string); ok {
					if str == "" {
						return nil, newerr(msgMustNotEmpty, ipath)
					}
					result.layout = append(result.layout, str)
				} else {
					return nil, newerr(msgMustString, ipath)
				}
			}
			if len(result.layout) == 0 {
				return nil, newerr(msgMustNotEmpty, kpath)
			}
		} else {
			return nil, newerr(msgMustStringArray, kpath)
		}
	} else {
		return nil, newerr("missing property layout", path)
	}
	if x, ok := ctn["location"]; ok {
		kpath := path + ".location"
		if x == nil {
			result.location = time.UTC
		} else if str, ok := x.(string); ok {
			if loc, err := time.LoadLocation(str); err != nil {
				return nil, newerr(err.Error(), kpath)
			} else {
				result.location = loc
			}
		} else if fnum, ok := x.(float64); ok {
			offset := int64(fnum)
			if fnum != float64(offset) {
				return nil, newerr(msgMustStringIntegerNull, kpath)
			}
			result.location = time.FixedZone("UTC"+strconv.FormatInt(offset, 10), int(offset))
		} else {
			return nil, newerr(msgMustStringIntegerNull, kpath)
		}
	} else {
		result.location = time.UTC
	}
	return result, nil
}

func (p *valueConvertTime) extract(input string) (any, *errloc) {
	for _, l := range p.layout {
		if t, err := time.ParseInLocation(l, input, p.location); err == nil {
			switch p.name {
			case timeUnix:
				return t.Unix(), nil
			case timeUnixMilli:
				return t.UnixMilli(), nil
			}
		}
	}
	return nil, newerr("unable parsing time \""+input+"\"", p.path)
}
