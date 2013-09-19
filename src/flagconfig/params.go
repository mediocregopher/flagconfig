package flagconfig

import (
	"github.com/mediocregopher/flagconfig/src/flagconfig/params"
)

// StrParam tells flagconfig to look for a param called name of type string in
// either the config file or on the command line, or use the given default
// instead
func StrParam(name, descr, def string) {
	fullConfig[name] = params.NewString(name, descr, def, false)
}

// StrParams tells flagconfig to look for zero or more parameters of type string
// in either the config file or on the command line, or use the given default
// instead. If any are defined in one location they overwrite all from the
// other. For example, if there are three defined in the config file and one
// defined on the command-line, that one will be the only one in the returned
// value.
func StrParams(name, descr string, def ...string) {
	fullConfig[name] = params.NewStrings(name, descr, def, false)
}

// IntParam tells flagconfig to look for a param called name of type int in
// either the config file or on the command line, or use the given default
// instead
func IntParam(name, descr string, def int) {
	fullConfig[name] = params.NewInt(name, descr, def, false)
}
