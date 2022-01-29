package main

import "github.com/MangioneAndrea/gonsole"

func main() {
	var stop bool
	var cont bool

	gonsole.Cli().
		Confirm("Do you want to continue?", &cont).
		KillIf(!cont).
		Confirm("You decided to continue. Do you want to stop?", &stop).
		KillIf(stop).
		Confirm("You decided not to stop! No matter what you do now, I'm going to end this", &stop).
		KillIf(true).
		Confirm("I'm never going to ask this anyway...", nil)

}
