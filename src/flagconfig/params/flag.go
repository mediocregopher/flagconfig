package params

import (
	"errors"
	"strings"
)

type flagParam struct {
	name, description           string
	defaultVal, cliVal, confVal *bool
	finalVal                    bool
}

func NewFlag(name, description string, def bool) *flagParam {
	defVal := new(bool)
	*defVal = def
	return &flagParam{
		name:        name,
		description: description,
		defaultVal:  defVal,
	}
}

func (p *flagParam) Name() string {
	return p.name
}

func (p *flagParam) Description() string {
	return p.description
}

func (p *flagParam) DefaultAsStrings() ([]string, bool) {
	if p.defaultVal != nil {
		if *p.defaultVal {
			return []string{"true"}, true
		} else {
			return []string{"false"}, true
		}
	} else {
		return nil, false
	}
}

func (p *flagParam) CLAHasValue() bool {
	return false
}

func (p *flagParam) CLA(v string) {
	p.cliVal = new(bool)
	*p.cliVal = !*p.defaultVal
}

func (p *flagParam) ConfFile(val string) {
	valL := strings.ToLower(val)
	p.confVal = new(bool)
	if valL == "true" {
		*p.confVal = true
	} else {
		*p.confVal = false
	}
}

func (p *flagParam) Post() error {
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

func (p *flagParam) Value() interface{} {
	return p.finalVal
}

