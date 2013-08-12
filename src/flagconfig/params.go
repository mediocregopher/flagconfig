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

// StrParam tells flagconfig to look for a param called name of type string in
// either the config file or on the command line, or use the given default
// instead
func StrParam(name, descr, def string) {
	newParam(name, descr, def, STRING)
}

// IntParam tells flagconfig to look for a param called name of type int in
// either the config file or on the command line, or use the given default
// instead
func IntParam(name, descr string, def int) {
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
