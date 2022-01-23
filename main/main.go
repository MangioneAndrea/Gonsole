package main

import (
	"github.com/MangioneAndrea/gonsole"
)

func main() {
	var nat string
	var color []string
	var conf bool
	var name string

	gonsole.Cli().
		SelectOne("Select nationality", []string{"en", "de", "it"}, &nat).
		SelectMany("Select colors", []string{"red", "blue", "green", "violet"}, &color).
		Confirm("Confirm", &conf).
		Input("Enter name", &name)
}
