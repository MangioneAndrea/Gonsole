package main

import (
	"fmt"
	"github.com/MangioneAndrea/gonsole"
)

func main() {
	var res []string
	gonsole.Cli().
		SelectMany([]string{"hey", "hi", "", "hola"}, &res) /*.
	Confirm("Confirm").
	Input("Enter name")*/
	fmt.Println(res)
}
