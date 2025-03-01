package dx

// ********************************************************************
// ( "value": any )?
// ********************************************************************
type propValue struct {
	value any
}

func parsePropValue(path string, ctn map[string]any) (*propValue, bool) {
	if x, ok := ctn["value"]; ok {
		return &propValue{x}, true
	}
	return nil, false
}

func (p *propValue) setFunc(slice SliceFunc, convert ConvertFunc) *errloc {
	return nil
}

func (p *propValue) extract(input string) (any, *errloc) {
	return p.value, nil
}
