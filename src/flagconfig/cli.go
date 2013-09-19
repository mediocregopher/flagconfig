package flagconfig

import (
	"fmt"
	"github.com/droundy/goopt"
	"os"
	"github.com/mediocregopher/flagconfig/src/flagconfig/params"
)

var fullConfig = map[string]params.Param{}

func get(name string) params.Param {
	val, ok := fullConfig[name]
	if !ok {
		panic("attempted to access non-parameter " + name)
	}
	return val
}

// GetInt looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is an integer)
func GetInt(name string) int {
	return get(name).Value().(int)
}

// GetStr looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is a string)
func GetStr(name string) string {
	return get(name).Value().(string)
}

// GetStrs looks for a configuration parameter of the given name and returns its
// value (assuming the parameter is a list of strings)
func GetStrs(name string) []string {
	return get(name).Value().([]string)
}

// Parse loads flagconfig's runtime configuration, using both command-line
// arguments and a possible configuration file
func Parse(projname string) error {
	fullConfig = map[string]params.Param{}

	for _,param := range fullConfig {
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
		fmt.Print(dumpExampleConfig(projname))
		os.Exit(0)
	}

	if *configFile != "" {
		configFileMap, err := readConfig(*configFile)
		if err != nil {
			return err
		}

		for _,param := range fullConfig {
			if vals, ok := configFileMap[param.Name()]; ok {
				for i := range vals {
					param.ConfFile(vals[i])
				}
			}
		}
	}

	for _, param := range fullConfig {
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
