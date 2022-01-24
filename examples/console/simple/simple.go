package main

import (
	"github.com/MangioneAndrea/gonsole"
)

func main() {
	gonsole.Error("File not found", "Error")
	gonsole.Log("Simple message", "Log")
	gonsole.Warn("Wrong password", "Warning")
	gonsole.Success("Action successful", "Success")
}
