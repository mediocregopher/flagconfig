package params

import (
	"github.com/droundy/goopt"
	"strconv"
	"errors"
)

type intParam struct {
	name, description string
	defaultVal, cliVal, confVal *int
	finalVal int
}

func NewInt(name, description string, def int, required bool) *intParam {
	var defVal *int = nil
	if !required {
		defVal = new(int)
		*defVal = def
	}
	return &intParam{
		name: name,
		description: description,
		defaultVal: defVal,
	}
}

func (p *intParam) Name() string {
	return p.name
}

func (p *intParam) Description() string {
	return p.description
}

func (p *intParam) DefaultAsStrings() ([]string,bool) {
	if p.defaultVal != nil {
		return []string{strconv.Itoa(*p.defaultVal)},true
	} else {
		return nil,false
	}
}

func (p *intParam) CLA() {
	d := 0
	if p.defaultVal != nil {
		d = *p.defaultVal
	}
	p.cliVal = goopt.Int(
		[]string{"--" + p.name},
		d,
		p.description,
	)
}

func (p *intParam) ConfFile(val string) {
	if valint, err := strconv.Atoi(val); err == nil {
		p.confVal = new(int)
		*p.confVal = valint
	}
}

func (p *intParam) Post() error {
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

func (p *intParam) Value() interface{} {
	return p.finalVal
}
