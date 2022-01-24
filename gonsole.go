package gonsole

import (
	"fmt"
	"github.com/gdamore/tcell"
	"io"
	"os"
)

var out io.Writer = os.Stdout
var in io.Writer = os.Stdin

// Color Format text with custom color
func setColor(colorString tcell.Color) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		if colorString == tcell.ColorDefault {
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

func Print(a interface{}, desc interface{}, showIf ShowIf, color tcell.Color) {
	if showIf(a) {
		if desc != nil {
			fmt.Fprint(out, setColor(color)(desc), setColor(color)(": "))
		}

		fmt.Printf("set'%s'", setColor(color)(a))
		fmt.Fprintln(out, setColor(color)(a))
	}
}

// Log Print a text
func Log(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, tcell.ColorDefault)
}

// Success Print a text green
func Success(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, tcell.ColorGreen)
}

// Error Print a text red
func Error(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, tcell.ColorRed)
}

// Warn Print a text yellow
func Warn(a interface{}, props ...interface{}) {
	desc, showIf := resolveProps(props)
	Print(a, desc, showIf, tcell.ColorYellow)
}
