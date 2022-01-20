package gonsole

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"os"
)

var out = os.Stdout
var in = os.Stdin

type cli struct {
	interrupted bool
	screen      *screen
}

func Cli() *cli {
	c := &cli{
		interrupted: false,
		screen:      NewScreen(),
	}
	term.Init()
	return c
}

func (c *cli) Confirm(message string, onAbort ...func(*cli)) *cli {
	if c.interrupted {
		return c
	}
	c.screen.WriteF("- %s? (y / n) : \n", message).ShowCursor(true).Draw()
	confirmed := c.screen.pollYN()
	if !confirmed && onAbort != nil && len(onAbort) == 1 {
		onAbort[0](c)
	}
	if confirmed {
		c.screen.Write("y")
	} else {
		c.screen.Write("n")
	}
	c.screen.ShowCursor(false).Newline().Draw()
	return c
}

func (c *cli) Input(message string, onEnter ...func(str string)) *cli {
	if c.interrupted {
		return c
	}
	c.screen.WriteF("- %s: ", message).ShowCursor(true).Draw()
	input := c.screen.pollText()
	c.screen.Newline().Draw()
	fmt.Println(input)
	if onEnter != nil && len(onEnter) > 0 {
		onEnter[0](input)
	}

	return c
}
