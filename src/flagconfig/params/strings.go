package params

import (
	"github.com/droundy/goopt"
	"errors"
)

type stringsParam struct {
	name, description string
	defaultVal, confVal []string
	cliVal *[]string
	finalVal []string
}

func NewStrings(name, description string, def []string, required bool) *stringsParam {
	var defVal []string = nil
	if !required {
		defVal = def
	}
	return &stringsParam{
		name: name,
		description: description,
		defaultVal: defVal,
	}
}

func (p *stringsParam) Name() string {
	return p.name
}

func (p *stringsParam) Description() string {
	return p.description
}

func (p *stringsParam) DefaultAsStrings() ([]string,bool) {
	return p.defaultVal, p.defaultVal != nil
}

func (p *stringsParam) CLA() {
	p.cliVal = goopt.Strings(
		[]string{"--" + p.name},
		"<string>",
		p.description,
	)
}

func (p *stringsParam) ConfFile(val string) {
	if p.confVal == nil {
		p.confVal = make([]string,0,4)
	}
	p.confVal = append(p.confVal,val)
}

func (p *stringsParam) Post() error {
	cliSet := p.cliVal != nil && *p.cliVal != nil && len(*p.cliVal) > 0
	if cliSet && !stringSlicesEqual(*p.cliVal,p.defaultVal) {
		p.finalVal = *p.cliVal
	} else if p.confVal != nil {
		p.finalVal = p.confVal
	} else if p.defaultVal != nil {
		p.finalVal = p.defaultVal
	} else {
		return errors.New("parameter required but no valid value given")
	}

	return nil
}

func stringSlicesEqual (a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (p *stringsParam) Value() interface{} {
	return p.finalVal
}
