package params

import (
	"errors"
	"github.com/droundy/goopt"
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

func (p *stringParam) CLA() {
	d := ""
	if p.defaultVal != nil {
		d = *p.defaultVal
	}
	p.cliVal = goopt.String(
		[]string{"--" + p.name},
		d,
		p.description,
	)
}

func (p *stringParam) ConfFile(val string) {
	p.confVal = new(string)
	*p.confVal = val
}

func (p *stringParam) Post() error {
	if p.cliVal != nil && *p.cliVal != *p.defaultVal {
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
