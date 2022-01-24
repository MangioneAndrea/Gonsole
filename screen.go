package gonsole

import (
	"fmt"
	"github.com/gdamore/tcell"
	"log"
)

type numbers interface {
	int
}

func min[T numbers](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T numbers](a, b T) T {
	if a > b {
		return a
	}
	return b
}

type screen struct {
	lines  []string
	colors map[int]tcell.Color
	x      int
	y      int
	cursor bool
	s      tcell.Screen
}

func NewScreen() *screen {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	scr := &screen{
		lines:  []string{""},
		colors: make(map[int]tcell.Color),
		x:      0,
		y:      0,
		s:      s,
	}
	scr.s.Clear()
	scr.s.Show()
	return scr
}

func (s *screen) reset() *screen {
	s.s.Clear()
	s.lines = []string{}

	return s
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (s *screen) Draw() *screen {
	s.s.Clear()

	for y, line := range s.lines {
		style := tcell.StyleDefault
		if val, ok := s.colors[y]; ok {
			style = style.Foreground(val)
		}
		drawText(s.s, 0, y, len(line), y, style, line)
	}
	if s.cursor {
		s.s.ShowCursor(s.x, s.y)
	} else {
		s.s.HideCursor()
	}
	s.s.Show()
	return s
}
func (s *screen) GoTo(x, y int) *screen {
	s.x = x
	s.y = y
	return s
}

func (s *screen) Up(amount int) *screen {
	s.y = max(0, s.y-amount)
	s.x = min(s.x, len(s.lines[s.y]))
	return s
}
func (s *screen) Down(amount int) *screen {
	s.y = min(len(s.lines)-1, s.y+amount)
	s.x = min(s.x, len(s.lines[s.y]))
	return s
}
func (s *screen) Left(amount int) *screen {
	s.x = max(0, s.x-amount)
	return s
}

func (s *screen) Right(amount int) *screen {
	if len(s.lines[s.y]) != 0 {
		s.x = min(len(s.lines[s.y]), amount+s.x)
	}
	return s
}

func (s *screen) Top() *screen {
	return s.Up(s.y)
}
func (s *screen) Bottom() *screen {
	return s.Down(len(s.lines) - s.y)
}
func (s *screen) Start() *screen {
	return s.Left(s.x)
}
func (s *screen) End() *screen {
	return s.Right(len(s.lines[s.y]) - s.y + 1)
}

func (s *screen) ClearColor() *screen {
	delete(s.colors, s.y)
	return s
}

func (s *screen) ColorLine(color tcell.Color) *screen {
	s.colors[s.y] = color
	return s
}

func (s *screen) Newline() *screen {
	s.y++
	s.x = 0
	s.lines = append(append(s.lines[:s.y], ""), s.lines[s.y:]...)
	s.lines[s.y] = ""
	return s
}

func (s *screen) WriteF(format string, args ...interface{}) *screen {
	str := fmt.Sprintf(format, args...)

	p := s.lines[s.y]
	s.lines[s.y] = p[:s.x] + str + p[s.x:]
	s.Right(len(str))
	return s
}
func (s *screen) Write(str string) *screen {
	return s.WriteF(str)
}

func (s *screen) DeleteLastLines(amount int) *screen {
	s.lines = s.lines[:len(s.lines)-amount]
	s.Down(0)
	return s
}

func (s *screen) DeleteOne() *screen {
	p := s.lines[s.y]
	s.lines[s.y] = p[:s.x-1] + p[s.x:]
	s.Left(1)
	return s
}

func (s *screen) ShowCursor(show bool) *screen {
	s.cursor = show
	return s
}

func (s *screen) pollYN() bool {
	for {
		ev := s.s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'y', 'Y':
					return true
				case 'n', 'N':
					return false
				}
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return false
			}
		}
	}
}
func (s *screen) pollText() string {
	minx := s.x
	for {
		ev := s.s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRight:
				s.Right(1).Draw()
			case tcell.KeyLeft:
				if s.x > minx {
					s.Left(1).Draw()
				}
			case tcell.KeyRune:
				s.Write(string(ev.Rune())).Draw()
			case tcell.KeyDelete:
				if s.x < len(s.lines[s.y]) {
					s.Right(1).DeleteOne().Draw()
				}
			case tcell.KeyBackspace:
				if s.x > minx {
					s.DeleteOne().Draw()
				}
			case tcell.KeyEnter:
				return s.lines[s.y][minx:]
			}

		}
	}
}

func (s *screen) pollVerticalSelect(count int) int {
	minY := len(s.lines) - count
	s.y = minY
	s.ColorLine(SelectedLineColor)
	s.Draw()
	for {
		ev := s.s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				}
			case tcell.KeyUp:
				if s.y > minY {
					s.ClearColor().Up(1).ColorLine(SelectedLineColor).Draw()
				}
			case tcell.KeyDown:
				if s.y < len(s.lines)-1 {
					s.ClearColor().Down(1).ColorLine(SelectedLineColor).Draw()
				}
			case tcell.KeyEnter:
				s.Bottom().Draw()
				return s.y - minY
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return -1
			}
		}
	}
}

func (s *screen) pollVerticalManySelect(count int, onUpdate func(bool, int)) []bool {
	minY := len(s.lines) - count
	selection := make([]bool, len(s.lines)-minY)
	s.y = minY
	s.ColorLine(SelectedLineColor)
	s.Draw()
	for {
		ev := s.s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case ' ':
					selection[s.y-minY] = !selection[s.y-minY]
					onUpdate(selection[s.y-minY], s.y-minY)
				}
			case tcell.KeyUp:
				if s.y > minY {
					s.ClearColor().Up(1).ColorLine(SelectedLineColor).Draw()
				}
			case tcell.KeyDown:
				if s.y < len(s.lines)-1 {
					s.ClearColor().Down(1).ColorLine(SelectedLineColor).Draw()
				}
			case tcell.KeyEnter:
				s.Bottom().Draw()
				return selection

			case tcell.KeyEscape, tcell.KeyCtrlC:
				return selection
			}
		}
	}
}

func (s *screen) Clear() *screen {
	s.lines = []string{""}
	s.x = 0
	s.y = 0
	return s
}
