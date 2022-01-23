package main

import (
	"github.com/MangioneAndrea/gonsole"
)

func main() {
	var multi []string
	var conf bool
	var name string
	gonsole.Cli().
		SelectMany("Select greeting", []string{"hey", "hi", "", "hola"}, &multi).
		Confirm("Confirm", &conf).
		Input("Enter name", &name)
}
