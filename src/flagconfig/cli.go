package flagconfig

import (
	"github.com/droundy/goopt"
	"fmt"
	"os"
	"strconv"
)

var fullConfig = map[string]interface{}{}

// GetInt looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is an integer)
func GetInt(name string) int {
	val, ok := fullConfig[name]
	if !ok {
		panic("attempted to access non-int-parameter " + name)
	}
	return val.(int)
}

// GetStr looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is a string)
func GetStr(name string) string {
	val, ok := fullConfig[name]
	if !ok {
		panic("attempted to access non-str-parameter " + name)
	}
	return val.(string)
}

// GetStr looks for a configuration parameter of the given name and returns its
// value (assuming the parameter is a list of strings)
func GetStrs(name string) []string {
	val, ok := fullConfig[name]
	if !ok {
		panic("attempted to access non-str-slice-parameter " + name)
	}
	return val.([]string)
}

// Pars loads flagconfig's runtime configuration, using both command-line arguments
// and a possible configuration file
func Parse(projname string) {

	//Load cli into its own set of config maps
	cliConfig := map[string]interface{}{}
	for name, param := range params {
		if param.Type == INT {
			cliConfig[name] = goopt.Int(
				[]string{"--"+name},
				param.Default.(int),
				param.Description,
			)
		} else if param.Type == STRING {
			cliConfig[name] = goopt.String(
				[]string{"--"+name},
				param.Default.(string),
				param.Description,
			)
		} else {
			cliConfig[name] = goopt.Strings(
				[]string{"--"+name},
				"<string>",
				param.Description,
			)
		}
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

	// Check if any of the multiple arg params didn't have any cli's set, and
	// therefore should have the default set
	for name, param := range params {
		if param.Type == STRINGS && len(*cliConfig[name].(*[]string)) < 1 {
			*cliConfig[name].(*[]string) = param.Default.([]string)
		}
	}

	//If the flag to dump example config is set to true, do that
	if *dumpExample {
		fmt.Print(dumpExampleConfig(projname))
		os.Exit(0)
	}

	//If config file is specified, load the string map from it and load the values into
	//global config
	if *configFile != "" {
		configFileMap, err := readConfig(*configFile)
		if err != nil {
			panic(err)
		}

		for name, val := range configFileMap {
			if param, ok := params[name]; ok {
				if param.Type != STRINGS && len(val) > 1 {
					panic("field "+name+" in "+*configFile+
					" cannot have multiple values")
				}

				if param.Type == INT {
					valint, err := strconv.Atoi(val[0])
					if err != nil {
						panic("field "+name+" in "+*configFile+
						" cannot be read as a number")
					}
					fullConfig[name] = valint
				} else if param.Type == STRING {
					fullConfig[name] = val[0]
				} else {
					fullConfig[name] = val
				}
			}
		}
	}

	//Now we look through each param. If it's set on the command-line (not to the default)
	//we set that in the global config maps. If it's also not set in the conf we set it
	//to the param's default. If it is set in the conf then it's already set in the global
	//configs by the previous section
	for name, param := range params {
		if param.Type == INT {
			cliVal := *cliConfig[name].(*int)
			_, confSet := fullConfig[name]
			if cliVal != param.Default {
				fullConfig[name] = cliVal
			} else if !confSet {
				fullConfig[name] = param.Default.(int)
			}
		} else if param.Type == STRING {
			cliVal := *cliConfig[name].(*string)
			_, confSet := fullConfig[name]
			if cliVal != param.Default {
				fullConfig[name] = cliVal
			} else if !confSet {
				fullConfig[name] = param.Default.(string)
			}
		} else {
			cliVal := *cliConfig[name].(*[]string)
			_, confSet := fullConfig[name]
			if !stringSlicesEqual(cliVal,param.Default.([]string)) {
				fullConfig[name] = cliVal
			} else if !confSet {
				fullConfig[name] = param.Default.([]string)
			}
		}
	}

}

func stringSlicesEqual(a, b []string) bool {
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
