package gonsole

import (
	"fmt"
	"io"
	"os"
)

var out io.Writer = os.Stdout
var in io.Writer = os.Stdin

type Color string

var (
	Default = Color("")
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

// Color Format text with custom color
func setColor(colorString Color) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		if colorString == Default {
			return fmt.Sprint(args...)
		}
		return fmt.Sprintf(string(colorString),
			fmt.Sprint(args...))
	}
	return sprint
}

type ShowIf func(elem interface{}) bool

var ShowAlways ShowIf = func(elem interface{}) bool { return true }

var ShowIfNotNil ShowIf = func(elem interface{}) bool { return elem != nil }

func resolveProps(props []interface{}) (interface{}, ShowIf) {
	switch len(props) {
	case 1:
		return props[0], ShowAlways
	case 2:
		return props[0], props[1].(ShowIf)
	}
	return nil, ShowAlways
}

func Print(a interface{}, desc interface{}, showIf ShowIf, color Color) {
	if showIf(a) {
		if desc != nil {
			fmt.Fprint(out, setColor(color)(desc), setColor(color)(": "))
		}
		fmt.Fprintln(out, setColor(color)(a))
	}
}

// Log Print a text
func Log(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, Default)
}

// Success Print a text green
func Success(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, Green)
}

// Error Print a text red
func Error(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, Red)
}

// Warn Print a text yellow
func Warn(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, Yellow)
}
