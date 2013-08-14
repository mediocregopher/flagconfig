package flagconfig

type paramType int

const (
	STRING paramType = iota
	STRINGS
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

// StrParams tells flagconfig to look for zero or more parameters of type string
// in either the config file or on the command line, or use the given default
// instead. If any are defined in one location they overwrite all from the
// other. For example, if there are three defined in the config file and one
// defined on the command-line, that one will be the only one in the returned
// value.
func StrParams(name, descr string, def ...string) {
	newParam(name, descr, def, STRINGS)
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
