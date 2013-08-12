package main

import (
	"flagconfig"
	"fmt"
)

func main() {
	//Specify the parameters we want to fetch
	flagconfig.StrParam("foo","Some foo","foofoofoo")
	flagconfig.IntParam("bar","Some bar",64)

	//Parse command line and possibly config file
	flagconfig.Parse("flagconfigtest")

	//Display the values that have been parsed
	fmt.Println(flagconfig.GetStr("foo"))
	fmt.Println(flagconfig.GetInt("bar"))
}
