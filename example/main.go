package main

import (
	"flagconfig"
	"fmt"
)

func main() {
	//Specify the parameters we want to fetch
	flagconfig.StrParam("foo","Some foo","foofoofoo")
	flagconfig.IntParam("bar","Some bar",64)
	flagconfig.StrParams("baz","Some baz","a","b","c")

	//Parse command line and possibly config file
	err := flagconfig.Parse("flagconfigtest")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Display the values that have been parsed
	fmt.Println(flagconfig.GetStr("foo"))
	fmt.Println(flagconfig.GetInt("bar"))
	fmt.Println(flagconfig.GetStrs("baz"))
}
