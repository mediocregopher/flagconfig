package cla

import (
	"os"
	"github.com/mediocregopher/flagconfig/src/flagconfig/params"
	"strings"
)

func looksLikeOpt(s string) bool {
	return s[0] == '-' && s[1] == '-'
}

func equalSignSplit(s string) (string,string,bool) {
	sub := strings.SplitN(s,"=",2)
	if len(sub) == 2 {
		return sub[0],sub[1],true
	} else {
		return "","",false
	}
}

func setAppend(r map[string][]string, k, v string) {
	if _,ok := r[k]; !ok {
		r[k] = make([]string,0,2)
	}
	r[k] = append(r[k],v)
}

func Parse(fullConfig map[string]params.Param) map[string][]string {
	r := map[string][]string{}
	args := os.Args[1:]

	for i:=0; i<len(args); i++ {
		if looksLikeOpt(args[i]) {
			fullname := args[i][2:]

			// If there's an equal-sign in there we split on that and use the
			// rest as the value (if param is wanted)
			if potName,potVal,ok := equalSignSplit(fullname); ok {
				if param,ok := fullConfig[potName]; ok {
					if !param.CLAHasValue() {
						potVal = ""
					}
					setAppend(r,potName,potVal)
				}

			// Otherwise check if this full key is wanted. If so process it.
			} else if param,ok := fullConfig[fullname]; ok {
				val := ""
				if param.CLAHasValue() {
					if len(args) > i+1 {
						//Skip this one next time-round the for loop, since
						//we're using it as a value it can't be a key
						val = args[i+1]
						i++ 
					}
				}
				setAppend(r,fullname,val)
			}
		}
	}

	return r
}
