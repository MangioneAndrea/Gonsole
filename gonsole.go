package gonsole

import (
	"github.com/MangioneAndrea/GoUtils/structures"
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

func (c *cli) Confirm(message string, confirmed *bool) *cli {
	if c.interrupted {
		return c
	}
	c.screen.WriteF("- %s? (y / n) : \n", message).ShowCursor(true).Draw()
	*confirmed = c.screen.pollYN()
	if *confirmed {
		c.screen.Write("y")
	} else {
		c.screen.Write("n")
	}
	c.screen.ShowCursor(false).Newline().Draw()
	return c
}

func (c *cli) Input(message string, input *string) *cli {
	if c.interrupted {
		return c
	}
	c.screen.WriteF("- %s: ", message).ShowCursor(true).Draw()
	*input = c.screen.pollText()
	c.screen.Newline().Draw()
	return c
}

func (c *cli) SelectOne(messages []string, selection *string) *cli {
	if c.interrupted {
		return c
	}

	for i, message := range messages {
		c.screen.WriteF("- %s", message)
		if i < len(messages)-1 {
			c.screen.Newline()
		}
	}
	c.screen.Top().Start().Draw()
	idx := c.screen.ShowCursor(true).Draw().pollVerticalSelect()

	if idx != -1 {
		*selection = messages[idx]
	}

	return c
}

func (c *cli) SelectMany(messages []string, selection *[]string) *cli {
	if c.interrupted {
		return c
	}

	for i, message := range messages {
		c.screen.WriteF("[ ] %s", message)
		if i < len(messages)-1 {
			c.screen.Newline()
		}
	}
	c.screen.Top().Start().Right(1).Draw()
	idxs := c.screen.ShowCursor(true).Draw().pollVerticalManySelect(func(selected bool, i int) {
		sign := " "
		if selected {
			sign = "x"
		}
		c.screen.Start().Right(2).DeleteOne().Write(sign).Left(1).Draw()
	})

	*selection = structures.Filter(messages, func(_ string, i int) bool {
		return idxs[i]
	})

	return c
}
