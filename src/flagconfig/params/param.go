package params

// Param describes a type of param (no really). It could be a single string, and
// integer, multiple strings, or really anything that the CLA handler supports.
type Param interface {

	// Returns the name of the parameter. This will be both the name on the
	// command line and the key in the configuration file.
	Name() string

	// The description of the parameter. This will be displayed both on the
	// command-line and in the configuration file
	Description() string

	// The default value (or multiple values) of the parameter as a string, with
	// a boolean to indicate if it even has one
	DefaultAsStrings() ([]string,bool)

	// Make the appropriate call to the CLA handler for this param. The handler
	// will probably return a pointer of some sort which will be populated once
	// the handler's Parse method is called (which will be after this is called
	// on all defined parameters).  This pointer should be stored.
	CLA()

	// When parsing the conf file (if given) this is called with the value
	// attached to any key which matches Name(). This could be called multiple
	// times.
	ConfFile(string)

	// This is called after all CLA's and ConfFile's are done being called,
	// and is when the parameter should declare that its final value that it got
	// can't be properly formatted, is missing, or anything of that nature. Nil
	// means the final value is valid and can be used.
	Post() error

	// The final value of the parameter. Anything input by ConfFile (assuming
	// it's valid) should be overwritten by anything on the command-line from
	// the CLA handler (again, assuming it's valid).
	Value() interface{}
}
