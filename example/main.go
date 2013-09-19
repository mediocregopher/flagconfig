package main

import (
	"github.com/mediocregopher/flagconfig/src/flagconfig"
	"fmt"
)

func main() {

	//Create new flagconfig object
	FC := flagconfig.New("flagconfigtest")

	//Specify the parameters we want to fetch
	FC.StrParam("foo","Some foo","foofoofoo")
	FC.IntParam("bar","Some bar",64)
	FC.StrParams("baz","Some baz","a","b","c")

	//Parse command line and possibly config file
	err := FC.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Display the values that have been parsed
	fmt.Println(FC.GetStr("foo"))
	fmt.Println(FC.GetInt("bar"))
	fmt.Println(FC.GetStrs("baz"))
}
