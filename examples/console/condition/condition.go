package main

import (
	"github.com/MangioneAndrea/gonsole"
	"os"
)

// ShowIfNotEmpty user defined conditional for the console to show
var ShowIfNotEmpty gonsole.ShowIf = func(elem interface{}) bool {
	return elem != ""
}

func main() {
	// this file is in the same folder as the runnable
	data, err := os.ReadFile("./examples/console/condition/example.txt")
	gonsole.Error(err, "Error", gonsole.ShowIfNotNil)
	gonsole.Success(string(data), "Data", ShowIfNotEmpty)
	// This file does not exist
	data, err = os.ReadFile("./path/to/nowhere.go")
	gonsole.Error(err, "Error", gonsole.ShowIfNotNil)
	gonsole.Success(string(data), "Data", ShowIfNotEmpty)
}
