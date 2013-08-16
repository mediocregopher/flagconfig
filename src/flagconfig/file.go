package flagconfig

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
)

// dumpExampleConfig returns the string representation of a configuration file
// with all parameters filled in with their default values
func dumpExampleConfig(projname string) string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("#####################################\n")
	buffer.WriteString("### " + projname + " configuration\n")
	buffer.WriteString("#####################################\n")
	buffer.WriteString("\n")

	for name, param := range params {
		buffer.WriteString("# " + param.Description + "\n")

		var def string
		if param.Type == STRING {
			def = param.Default.(string)
			buffer.WriteString(name + ": " + def + "\n")
			buffer.WriteString("\n")
		} else if param.Type == INT {
			def = strconv.Itoa(param.Default.(int))
			buffer.WriteString(name + ": " + def + "\n")
			buffer.WriteString("\n")
		} else {
			for _, def := range param.Default.([]string) {
				buffer.WriteString(name + ": " + def + "\n")
			}
			buffer.WriteString("\n")
		}

	}
	return buffer.String()
}

func setOrAppend(m map[string][]string, name, val string) {
	_, ok := m[name]
	if !ok {
		m[name] = make([]string, 0, 8)
	}
	m[name] = append(m[name], val)
}

// readConfig returns a map of the key/values found in a given configuration file.
// Note: these key/values don't have to actually correspond to expected parameters,
// that parsing is done elsewhere
func readConfig(file string) (map[string][]string, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(fi)
	ret := map[string][]string{}
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\n")
		if len(line) > 0 && line[0] != '#' {
			idx := strings.Index(line, ":")
			if idx < 0 {
				continue
			}
			name := strings.Trim(line[0:idx], " \t")
			val := strings.Trim(line[idx+1:], " \t")
			setOrAppend(ret, name, val)
		}
	}

	return ret, nil

}
