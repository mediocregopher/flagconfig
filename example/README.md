# flagconfig example

This example shows how you can use flagconfig to easily parse options off either
the command line or a configuration file. If you call `go run main.go` (assuming
flagconfig is installed on your command-line) then you will see the default
values for three parameters, `--foo` and `--bar` and `--baz`. You can use `-h`
to see all options.

## --example

flagconfig can read the configuration from a configuration file. These files
have a very simple syntax that you'll be able to easily pick up at a glance. You
can use the `--example` flag to dump your application's example configuration to
stdout (which can then be piped to a file).

## --config

You can specify a configuration file for flagconfig to read in and attempt to
populate the defined paramters with. If a parameter is defined on the
command-line AND in a configuration file then the command-line value takes
precedence. If it's not defined in either then the default value specified is
used.
