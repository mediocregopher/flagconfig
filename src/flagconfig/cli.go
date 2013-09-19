package flagconfig

import (
	"fmt"
	"github.com/droundy/goopt"
	"os"
	"github.com/mediocregopher/flagconfig/src/flagconfig/params"
)

type FlagConfig struct {
	projname string
	fullConfig map[string]params.Param
}

// New returns a FlagConfig struct. Usage of this struct is:
//	* Tell it what params to look for with StrParam, IntParam, etc...
//	* Call Parse()
//	* Get found parameters using GetStr, GetInt, etc...
func New(projname string) *FlagConfig {
	return &FlagConfig{
		projname: projname,
		fullConfig: map[string]params.Param{},
	}
}

func (f *FlagConfig) get(name string) params.Param {
	val, ok := f.fullConfig[name]
	if !ok {
		panic("attempted to access non-parameter " + name)
	}
	return val
}

// IntParam tells flagconfig to look for a param called name of type int in
// either the config file or on the command line, or use the given default
// instead
func (f *FlagConfig) IntParam(name, descr string, def int) {
	f.fullConfig[name] = params.NewInt(name, descr, def, false)
}

// GetInt looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is an integer)
func (f *FlagConfig) GetInt(name string) int {
	return f.get(name).Value().(int)
}

// StrParam tells flagconfig to look for a param called name of type string in
// either the config file or on the command line, or use the given default
// instead
func (f *FlagConfig) StrParam(name, descr, def string) {
	f.fullConfig[name] = params.NewString(name, descr, def, false)
}

// GetStr looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is a string)
func (f *FlagConfig) GetStr(name string) string {
	return f.get(name).Value().(string)
}

// StrParams tells flagconfig to look for zero or more parameters of type string
// in either the config file or on the command line, or use the given default
// instead. If any are defined in one location they overwrite all from the
// other. For example, if there are three defined in the config file and one
// defined on the command-line, that one will be the only one in the returned
// value.
func (f *FlagConfig) StrParams(name, descr string, def ...string) {
	f.fullConfig[name] = params.NewStrings(name, descr, def, false)
}

// GetStrs looks for a configuration parameter of the given name and returns its
// value (assuming the parameter is a list of strings)
func (f *FlagConfig) GetStrs(name string) []string {
	return f.get(name).Value().([]string)
}

// Parse loads flagconfig's runtime configuration, using both command-line
// arguments and a possible configuration file
func (f *FlagConfig) Parse() error {

	for _,param := range f.fullConfig {
		param.CLA()
	}

	//Some extra cli args
	dumpExample := goopt.Flag(
		[]string{"--example"},
		[]string{},
		"Dump example configuration to stdout and exit",
		"",
	)
	configFile := goopt.String(
		[]string{"--config"},
		"",
		"Configuration file to load, empty means don't load any file and only"+
			" use command-line args",
	)

	goopt.Parse(nil)

	//If the flag to dump example config is set to true, do that
	if *dumpExample {
		fmt.Print(f.dumpExampleConfig(f.projname))
		os.Exit(0)
	}

	if *configFile != "" {
		configFileMap, err := readConfig(*configFile)
		if err != nil {
			return err
		}

		for _,param := range f.fullConfig {
			if vals, ok := configFileMap[param.Name()]; ok {
				for i := range vals {
					param.ConfFile(vals[i])
				}
			}
		}
	}

	for _, param := range f.fullConfig {
		if err := param.Post(); err != nil {
			return fmt.Errorf(
				"error in param %s: %s",
				param.Name(),
				err.Error(),
			)
		}
	}

	return nil
}
