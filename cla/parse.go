package cla

import (
	"github.com/mediocregopher/flagconfig/params"
	"os"
	"strings"
)

func looksLikeOpt(s, delim string) bool {
	return strings.Index(s, delim) == 0
}

func equalSignSplit(s string) (string, string, bool) {
	sub := strings.SplitN(s, "=", 2)
	if len(sub) == 2 {
		return sub[0], sub[1], true
	} else {
		return "", "", false
	}
}

func setAppend(r map[string][]string, k, v string) {
	if _, ok := r[k]; !ok {
		r[k] = make([]string, 0, 2)
	}
	r[k] = append(r[k], v)
}

func getParam(fullConfig []params.Param, name string) (params.Param,bool) {
	for _, p := range fullConfig {
		if p.Name() == name {
			return p,true
		}
	}
	return nil,false
}

func Parse(delim string, fullConfig []params.Param) (map[string][]string, []string) {
	r := map[string][]string{}
	pos := make([]string, 0, 2)
	delimL := len(delim)

	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		if looksLikeOpt(args[i], delim) {
			fullname := args[i][delimL:]

			// If there's an equal-sign in there we split on that and use the
			// rest as the value (if param is wanted)
			if potName, potVal, ok := equalSignSplit(fullname); ok {
				if param, ok := getParam(fullConfig, potName); ok {
					if !param.CLAHasValue() {
						potVal = ""
					}
					setAppend(r, potName, potVal)
				}

			// Otherwise check if this full key is wanted. If so process it.
			} else if param, ok := getParam(fullConfig, fullname); ok {
				val := ""
				if param.CLAHasValue() {
					if len(args) > i+1 {
						//Skip this one next time-round the for loop, since
						//we're using it as a value it can't be a key
						val = args[i+1]
						i++
					}
				}
				setAppend(r, fullname, val)
			} else {
				pos = append(pos, args[i])
			}
		} else {
			pos = append(pos, args[i])
		}
	}

	return r, pos
}
