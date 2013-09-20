package params

import (
	"errors"
)

type stringsParam struct {
	name, description   string
	defaultVal, cliVal, confVal []string
	finalVal            []string
}

func NewStrings(name, description string, def []string, required bool) *stringsParam {
	var defVal []string = nil
	if !required {
		defVal = def
	}
	return &stringsParam{
		name:        name,
		description: description,
		defaultVal:  defVal,
	}
}

func (p *stringsParam) Name() string {
	return p.name
}

func (p *stringsParam) Description() string {
	return p.description
}

func (p *stringsParam) DefaultAsStrings() ([]string, bool) {
	if p.defaultVal != nil {
		defs := make([]string,len(p.defaultVal))
		for i := range p.defaultVal {
			defs[i] = "\""+p.defaultVal[i]+"\""
		}
		return defs,true
	} else {
		return nil,false
	}
}

func (p *stringsParam) CLAHasValue() bool {
	return true
}

func (p *stringsParam) CLA(val string) {
	if p.cliVal == nil {
		p.cliVal = make([]string, 0, 4)
	}
	p.cliVal = append(p.cliVal, val)
}

func (p *stringsParam) ConfFile(val string) {
	if p.confVal == nil {
		p.confVal = make([]string, 0, 4)
	}
	p.confVal = append(p.confVal, val)
}

func (p *stringsParam) Post() error {
	if p.cliVal != nil {
		p.finalVal = p.cliVal
	} else if p.confVal != nil {
		p.finalVal = p.confVal
	} else if p.defaultVal != nil {
		p.finalVal = p.defaultVal
	} else {
		return errors.New("parameter required but no valid value given")
	}

	return nil
}

func (p *stringsParam) Value() interface{} {
	return p.finalVal
}
