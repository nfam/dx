package dx

import (
	"strconv"
	"strings"
)

// ********************************************************************
// ( "sep":  string | [string]  )?
// ( "item": exp | [exp | null] )?
// ********************************************************************
type propList struct {
	valueList
}

func parsePropList(path string, ctn map[string]any) (*propList, *errloc) {
	if x, ok := ctn["list"]; ok {
		path += ".list"
		if obj, ok := x.(map[string]any); ok {
			if value, err := parseValueList(path, obj); err != nil {
				return nil, err
			} else {
				return &propList{*value}, nil
			}
		}
		return nil, newerr(msgMustObject, path)
	}
	return nil, nil
}

type valueList struct {
	sep  propSep
	item propItem
}

func parseValueList(path string, ctn map[string]any) (*valueList, *errloc) {
	var (
		list = &valueList{}
		err  *errloc
	)
	if list.sep, err = parsePropSep(path, ctn); err != nil {
		return nil, err
	}
	if list.item, err = parsePropItem(path, ctn); err != nil {
		return nil, err
	}
	for name := range ctn {
		if !strings.HasPrefix(name, "//") {
			switch name {
			case "sep", "item":
			default:
				return nil, newerr("unknown property "+name, path)
			}
		}
	}
	return list, nil
}

func (p *valueList) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	if p.item != nil {
		return p.item.setFunc(slice, convert)
	}
	return nil
}

func (p *valueList) extract(input string) (any, *errloc) {
	var subs []string
	if p.sep != nil {
		subs = p.sep.process(input)
	} else {
		subs = []string{input}
	}
	list := make([]any, 0, len(subs))
	for _, sub := range subs {
		if p.item != nil {
			if r, err := p.item.extract(sub); err != nil {
				err.msg = append(err.msg, "@list("+strconv.Itoa(len(list))+") "+sub)
				return nil, err
			} else {
				list = append(list, r)
			}
		} else {
			list = append(list, sub)
		}
	}
	return list, nil
}
