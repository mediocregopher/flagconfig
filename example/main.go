package main

import (
	"flagconfig"
	"fmt"
)

func main() {
	flagconfig.StrParam("foo","Some foo","foofoofoo")
	flagconfig.Int64Param("bar","Some bar",64)
	flagconfig.Parse("flagconfigtest")

	fmt.Println(flagconfig.GetStr("foo"))
	fmt.Println(flagconfig.GetInt64("bar"))
}
