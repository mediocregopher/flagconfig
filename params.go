package flagconfig

type paramType int

const (
	STRING paramType = iota
	INT
)

type param struct {
	Description string
	Type        paramType
	Default     interface{}
}

func StrParam(name, descr, def string) {
	newParam(name, descr, def, STRING)
}

func IntParam(name, descr string, def int64) {
	newParam(name, descr, def, INT)
}

func newParam(name, descr string, def interface{}, t paramType) {
	params[name] = param{
		Description: descr,
		Type:        t,
		Default:     def,
	}
}

var params = map[string]param{}
