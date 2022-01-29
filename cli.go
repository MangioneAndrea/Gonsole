package gonsole

import (
	"github.com/MangioneAndrea/GoUtils/structures"
	"github.com/gdamore/tcell"
	term "github.com/nsf/termbox-go"
	"strings"
	"time"
)

type cli struct {
	interrupted bool
	screen      *screen
}

func Cli() *cli {
	c := &cli{
		interrupted: false,
		screen:      newScreen(),
	}
	term.Init()
	// Fixes a graphic glitch
	time.Sleep(150 * time.Millisecond)

	return c
}

var QuestionColor tcell.Color = tcell.ColorViolet
var ActiveQuestionColor tcell.Color = tcell.ColorForestGreen
var SelectedLineColor tcell.Color = tcell.ColorCadetBlue

func (c *cli) Confirm(message string, confirmed *bool) *cli {
	if c.interrupted {
		return c
	}
	c.screen.ColorLine(ActiveQuestionColor).WriteF("%s (y / n) : \n", message).ShowCursor(true).Draw()
	*confirmed = c.screen.pollYN()
	if *confirmed {
		c.screen.Write("y")
	} else {
		c.screen.Write("n")
	}
	c.screen.ColorLine(QuestionColor).ShowCursor(false).Newline().Draw()
	return c
}
func (c *cli) Clear() *cli {
	c.screen.Clear()
	return c
}

func (c *cli) Input(message string, input *string) *cli {
	if c.interrupted {
		return c
	}
	c.screen.ColorLine(ActiveQuestionColor).WriteF("%s: ", message).ShowCursor(true).Draw()
	*input = c.screen.pollText()
	c.screen.ColorLine(QuestionColor).Newline().Draw()
	return c
}

func (c *cli) SelectOne(question string, answers []string, selection *string) *cli {
	if c.interrupted {
		return c
	}
	c.screen.ColorLine(ActiveQuestionColor).Write(question).Newline()

	for i, message := range answers {
		c.screen.WriteF("- %s", message)
		if i < len(answers)-1 {
			c.screen.Newline()
		}
	}
	c.screen.Top().Start().Draw()
	idx := c.screen.ShowCursor(true).Draw().pollVerticalSelect(len(answers))

	if idx != -1 {
		*selection = answers[idx]
	} else {
		*selection = ""
	}

	c.screen.DeleteLastLines(len(answers)).End().WriteF(": %s", *selection).ColorLine(QuestionColor).Newline().Draw()

	return c
}

func (c *cli) KillIf(condition bool) *cli {
	if condition {
		c.interrupted = true
	}
	return c
}

func (c *cli) SelectMany(question string, answers []string, selection *[]string) *cli {
	if c.interrupted {
		return c
	}
	c.screen.ColorLine(ActiveQuestionColor).Write(question).Newline()

	for i, message := range answers {
		c.screen.WriteF("[ ] %s", message)
		if i < len(answers)-1 {
			c.screen.Newline()
		}
	}
	c.screen.Top().Start().Right(1).Draw()
	idxs := c.screen.ShowCursor(true).Draw().pollVerticalManySelect(len(answers), func(selected bool, i int) {
		sign := " "
		if selected {
			sign = "x"
		}
		c.screen.Start()
		c.screen.Right(2)
		c.screen.DeleteOne()
		c.screen.Write(sign)
		c.screen.Left(1).Draw()
	})

	*selection = structures.Filter(answers, func(_ string, i int) bool {
		return idxs[i]
	})

	formatted := strings.Join(*selection, ", ")

	c.screen.DeleteLastLines(len(answers)).End().WriteF(": %s", formatted).ColorLine(QuestionColor).Newline().Draw()

	return c
}
