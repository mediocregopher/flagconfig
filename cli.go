package flagconfig

import (
	"bytes"
	"fmt"
	"github.com/mediocregopher/flagconfig/cla"
	"github.com/mediocregopher/flagconfig/params"
	"os"
	"strings"
	"github.com/mediocregopher/tablewriter"
)

type FlagConfig struct {
	projname      string
	projdescr     string
	projpostdescr string
	fullConfig    []params.Param
	pos           []string
	delim         string

	// if true, don't use configuration file at all
	forceNoConfig bool

	// if true and a config file is set but doesn't exist, continue and don't
	// error
	ignoreMissingConfig bool

	// the file to use as a config file by default
	defaultConfigFile string

	// if true we don't print usage on Help()
	dontPrintUsage bool

	cliOnly map[string]bool
}

// New returns a FlagConfig struct. Usage of this struct is:
//* Tell it what params to look for with StrParam, IntParam, etc...
//* Call Parse()
//* Get found parameters using GetStr, GetInt, etc...
func New(projname string) *FlagConfig {
	return &FlagConfig{
		projname:   projname,
		fullConfig: make([]params.Param,0,8),
		delim:      "--",
		cliOnly:    map[string]bool{
			"example": true,
			"config": true,
			"help": true,
		},
	}
}

// SetDescription sets the description which will be shown above the flag usage
// when --help is called on the command-line. This is optional.
func (f *FlagConfig) SetDescription(descr string) {
	f.projdescr = descr
}

// SetExtraHelp sets a string which will show up at the tail end of the --help
// message, after all flags are defined. This is optional.
func (f *FlagConfig) SetExtraHelp(help string) {
	f.projpostdescr = help
}

// Don't print a Usage line when printing help. Use this if you want to define
// your own usage line in SetDescription or SetExtraHelp.
func (f *FlagConfig) DontPrintUsage() {
	f.dontPrintUsage = true
}

// Set the argument delimiter which identifies and argument as being a flag.
// Default is -- (--username, for example).
func (f *FlagConfig) SetDelimiter(delim string) {
	f.delim = delim
}

// Sets the given field names to be cli-only. They will not be shown in the
// -example output. Can be called multiple times
func (f *FlagConfig) CLIOnly(fields ...string) {
	for i := range fields {
		f.cliOnly[fields[i]] = true
	}
}

func (f *FlagConfig) get(name string) params.Param {
	for _, param := range f.fullConfig {
		if name == param.Name() {
			return param
		}
	}
	return nil
}

func (f *FlagConfig) add(p params.Param) {
	f.fullConfig = append(f.fullConfig, p)
}

// IntParam tells flagconfig to look for a param called name of type int in
// either the config file or on the command line, or use the given default
// instead
func (f *FlagConfig) IntParam(name, descr string, def int) {
	f.add(params.NewInt(name, descr, def, false))
}

// RequiredIntParam tells flagconfig to look for a param called name of type int
// in either the config file or on the command line, or return an error from
// Parse if it's not specified anywhere
func (f *FlagConfig) RequiredIntParam(name, descr string) {
	f.add(params.NewInt(name, descr, 0, true))
}

// GetInt looks for a configuration parameter of the given name and
// returns its value (assuming the parameter is an integer)
func (f *FlagConfig) GetInt(name string) int {
	return f.get(name).Value().(int)
}

// IntParams tells flagconfig to look for zero or more parameters of type int in
// either the config file or on the command line, or use the given default
// instead. If any are defined in one location they overwrite all from the
// other. For example, if there are three defined in the config file and one
// defined on the command-line, that one will be the only one in the returned
// value.
func (f *FlagConfig) IntParams(name, descr string, defaults ...int) {
	f.add(params.NewInts(name, descr, defaults))
}

// GetInts looks for a configuration parameter of the given name and returns its
// value (assuming the parameter is a list of strings)
func (f *FlagConfig) GetInts(name string) []int {
	return f.get(name).Value().([]int)
}

// StrParam tells flagconfig to look for a param called name of type string in
// either the config file or on the command line, or use the given default
// instead
func (f *FlagConfig) StrParam(name, descr, def string) {
	f.add(params.NewString(name, descr, def, false))
}

// RequiredStrParam tells flagconfig to look for a param called name of type
// string in either the config file or on the command line, or return an error
// from Parse if it's not specified anywhere
func (f *FlagConfig) RequiredStrParam(name, descr string) {
	f.add(params.NewString(name, descr, "", true))
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
func (f *FlagConfig) StrParams(name, descr string, defaults ...string) {
	f.add(params.NewStrings(name, descr, defaults))
}

// GetStrs looks for a configuration parameter of the given name and returns its
// value (assuming the parameter is a list of strings)
func (f *FlagConfig) GetStrs(name string) []string {
	return f.get(name).Value().([]string)
}

// FlagParam tells flagconfig to look for a param called name on the command
// line or in the config file. Passing the flag on the command-line indicates a
// value of whatever the opposite of the default is (so if the default is false,
// passing it on the command-line would mean true). In the configuration file
// the value can be either "true" or "false".
func (f *FlagConfig) FlagParam(name, descr string, def bool) {
	f.add(params.NewFlag(name, descr, def))
}

// GetFlag looks for a configuration parameter of the given name and returns its
// value (assuming the parameter is a flag).
func (f *FlagConfig) GetFlag(name string) bool {
	return f.get(name).Value().(bool)
}

// Parse loads flagconfig's runtime configuration, using both command-line
// arguments and a possible configuration file
func (f *FlagConfig) Parse() error {
	f.FlagParam(
		"help",
		"Display help for parameters",
		false,
	)

	if !f.forceNoConfig {
		f.FlagParam(
			"example",
			"Dump example configuration to stdout and exit",
			false,
		)

		f.StrParam(
			"config",
			"Configuration file to load, empty means don't load any file and"+
			" only use command-line args",
			f.defaultConfigFile,
		)
	}

	claMap, pos := cla.Parse(f.delim, f.fullConfig)
	f.pos = pos
	_, printHelp := claMap["help"]
	_, printExample := claMap["example"]

	var configFile string
	if configFiles, ok := claMap["config"]; ok && len(configFiles) > 0 {
		configFile = configFiles[0]
	} else if f.defaultConfigFile != "" {
		configFile = f.defaultConfigFile
	}

	if printHelp {
		fmt.Println(f.Help())
		os.Exit(0)
	}

	//If the flag to dump example config is set to true, do that
	if printExample {
		fmt.Print(f.dumpExampleConfig(f.projname))
		os.Exit(0)
	}

	if configFile != "" {

		if configFile[:2] == "~/" {
			// We don't use os/user because it's not supported by osx or
			// windows, but this is (usually, I think)
			home := os.Getenv("HOME")
			if home == "" {
				return fmt.Errorf("Could not determine home directory. Please set the $HOME environment variable")
			}
			configFile = strings.Replace(configFile, "~", home, 1)
		}


		configFileMap, err := readConfig(configFile)
		if err != nil {
			// configFileMap isn't nil when the file couldn't be opened
			if !(configFileMap != nil && f.ignoreMissingConfig) {
				return err
			}
		}

		for _, param := range f.fullConfig {
			if vals, ok := configFileMap[param.Name()]; ok {
				for i := range vals {
					param.ConfFile(vals[i])
				}
			}
		}
	}

	for _, param := range f.fullConfig {
		if vals, ok := claMap[param.Name()]; ok {
			for i := range vals {
				param.CLA(vals[i])
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

// Returns a formatted help string for the parameters that have been given so
// far
func (f *FlagConfig) Help() string {
	buf := bytes.NewBuffer(make([]byte, 256))

	if !f.dontPrintUsage {
		fmt.Fprintf(buf, "Usage: %s [FLAGS]\n", os.Args[0])
	}

	if f.projdescr != "" {
		fmt.Fprintf(buf, "%s\n", f.projdescr)
	}

	fmt.Fprintf(buf, "\n")

	w := tablewriter.New(buf)
	w.SetTableWidth(150)
	w.SetBottomPadding(0)
	w.AddColumn(-1, 25)
	w.AddColumn(-1, 30)
	w.AddColumn(0, -1)

	fmtStr := "%s\t%s\t%s\n"
	fmt.Fprintf(w, fmtStr, "Flag", "Default(s)", "Description")
	fmt.Fprintf(w, fmtStr, "~~~~", "~~~~~~~~~~", "~~~~~~~~~~~")

	for _, param := range f.fullConfig {
		defj := "<required>"
		if defs, ok := param.DefaultAsStrings(); ok {
			defj = strings.Join(defs, ", ")
		}
		_, err := fmt.Fprintf(
			w, fmtStr,
			f.delim+param.Name(),
			defj,
			param.Description(),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR printing help: %s\n", err)
		}
	}

	if f.projpostdescr != "" {
		fmt.Fprintf(buf, "\n%s\n", f.projpostdescr)
	}

	return buf.String()
}

// If called then flagconfig won't have the ability to read in a configuration
// file. The --config and --example options won't be present in the --help menu.
func (f *FlagConfig) DisallowConfig() {
	f.forceNoConfig = true
}

// If set, this will become the default value for the --config parameter. Note
// that this will still generate an error if it's set but the file doesn't
// exist. Use IgnoreMissingConfigFile to have this not be the case.
func (f *FlagConfig) SetDefaultConfigFile(file string) {
	f.defaultConfigFile = file
}

// If called then flagconfig won't stop execution if it's given a configuration
// file to read but there isn't a file there. Instead it will act as if the file
// was empty.
func (f *FlagConfig) IgnoreMissingConfigFile() {
	f.ignoreMissingConfig = true
}

// GetPositionals returns all positional arguments found after Parse'ing. Will
// return empty slice if none were found, and nil if Parse hasn't been called
// yet
func (f *FlagConfig) GetPositionals() []string {
	return f.pos
}
