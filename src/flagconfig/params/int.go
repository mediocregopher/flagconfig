package params

import (
	"errors"
	"strconv"
)

type intParam struct {
	name, description           string
	defaultVal, cliVal, confVal *int
	finalVal                    int
}

func NewInt(name, description string, def int, required bool) *intParam {
	var defVal *int = nil
	if !required {
		defVal = new(int)
		*defVal = def
	}
	return &intParam{
		name:        name,
		description: description,
		defaultVal:  defVal,
	}
}

func (p *intParam) Name() string {
	return p.name
}

func (p *intParam) Description() string {
	return p.description
}

func (p *intParam) DefaultAsStrings() ([]string, bool) {
	if p.defaultVal != nil {
		return []string{strconv.Itoa(*p.defaultVal)}, true
	} else {
		return nil, false
	}
}

func (p *intParam) CLAHasValue() bool {
	return true
}

func (p *intParam) CLA(s string) {
	if si, err := strconv.Atoi(s); err == nil {
		p.cliVal = new(int)
		*p.cliVal = si
	}
}

func (p *intParam) ConfFile(val string) {
	if valint, err := strconv.Atoi(val); err == nil {
		p.confVal = new(int)
		*p.confVal = valint
	}
}

func (p *intParam) Post() error {
	if p.cliVal != nil {
		p.finalVal = *p.cliVal
	} else if p.confVal != nil {
		p.finalVal = *p.confVal
	} else if p.defaultVal != nil {
		p.finalVal = *p.defaultVal
	} else {
		return errors.New("parameter required but no valid value given")
	}

	return nil
}

func (p *intParam) Value() interface{} {
	return p.finalVal
}
