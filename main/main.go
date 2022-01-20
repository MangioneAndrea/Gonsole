package main

import (
	"github.com/MangioneAndrea/gonsole"
)

func main() {
	gonsole.Cli().Confirm("Confirm").Input("Enter name")
}
