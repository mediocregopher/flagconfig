package params

import (
	"errors"
	"strconv"
)

type intsParam struct {
	name, description           string
	defaultVal, cliVal, confVal []int
	finalVal                    []int
}

func NewInts(name, description string, def []int) *intsParam {
	if def == nil {
		def = []int{}
	}
	return &intsParam{
		name:        name,
		description: description,
		defaultVal:  def,
	}
}

func (p *intsParam) Name() string {
	return p.name
}

func (p *intsParam) Description() string {
	return p.description
}

func (p *intsParam) DefaultAsStrings() ([]string, bool) {
	if p.defaultVal != nil {
		defs := make([]string, len(p.defaultVal))
		for i := range p.defaultVal {
			defs[i] = strconv.Itoa(p.defaultVal[i])
		}
		return defs, true
	} else {
		return nil, false
	}
}

func (p *intsParam) CLAHasValue() bool {
	return true
}

func (p *intsParam) CLA(val string) {
	if vali, err := strconv.Atoi(val); err == nil {
		if p.cliVal == nil {
			p.cliVal = make([]int, 0, 4)
		}
		p.cliVal = append(p.cliVal, vali)
	}
}

func (p *intsParam) ConfFile(val string) {
	if vali, err := strconv.Atoi(val); err == nil {
		if p.confVal == nil {
			p.confVal = make([]int, 0, 4)
		}
		p.confVal = append(p.confVal, vali)
	}
}

func (p *intsParam) Post() error {
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

func (p *intsParam) Value() interface{} {
	return p.finalVal
}
