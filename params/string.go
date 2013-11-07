package params

import (
	"errors"
)

type stringParam struct {
	name, description           string
	defaultVal, cliVal, confVal *string
	finalVal                    string
}

func NewString(name, description string, def string, required bool) *stringParam {
	var defVal *string = nil
	if !required {
		defVal = new(string)
		*defVal = def
	}
	return &stringParam{
		name:        name,
		description: description,
		defaultVal:  defVal,
	}
}

func (p *stringParam) Name() string {
	return p.name
}

func (p *stringParam) Description() string {
	return p.description
}

func (p *stringParam) DefaultAsStrings() ([]string, bool) {
	if p.defaultVal != nil {
		return []string{*p.defaultVal}, true
	} else {
		return nil, false
	}
}

func (p *stringParam) CLAHasValue() bool {
	return true
}

func (p *stringParam) CLA(val string) {
	p.cliVal = new(string)
	*p.cliVal = val
}

func (p *stringParam) ConfFile(val string) {
	p.confVal = new(string)
	*p.confVal = val
}

func (p *stringParam) Post() error {
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

func (p *stringParam) Value() interface{} {
	return p.finalVal
}
