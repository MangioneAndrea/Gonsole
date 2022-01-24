package main

import "github.com/MangioneAndrea/gonsole"

func main() {
	var name string
	var choco bool
	var color string
	var languages []string

	gonsole.Cli().
		Input("What's your name?", &name).
		Confirm("Do you like chocolate?", &choco).
		SelectOne("What's your favourite colour?", []string{"red", "blue", "yellow", "cyan"}, &color).
		SelectMany("What languages do you speak?", []string{"English", "German", "Italian", "French"}, &languages)
}
