package main

import (
	"github.com/mediocregopher/flagconfig/src/flagconfig"
	"fmt"
)

func main() {

	//Create new flagconfig object
	FC := flagconfig.New("flagconfigtest")

	//Optionally set a description that'll show up in the command-line help
	FC.SetDescription("A super cool test program")

	//Specify the parameters we want to fetch
	FC.StrParam("foo","Some foo","foofoofoo")
	FC.IntParam("bar","Some bar",64)
	FC.StrParams("baz","Some baz (can set multiple times)","a","b","c")
	FC.FlagParam("bax", "Some bax", false)
	FC.RequiredIntParam("baw","Some baw")
	FC.IntParams("baq","Some baq (can set multiple times)",1,2,3)

	//Optionally set a message to show up at the end of the --help message
	FC.SetExtraHelp("Thanks for reading the --help message!")

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
	fmt.Println(FC.GetFlag("bax"))
	fmt.Println(FC.GetInt("baw"))
	fmt.Println(FC.GetInts("baq"))
}
