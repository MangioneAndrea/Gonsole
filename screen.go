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
	x      int
	y      int
	cursor bool
	s      tcell.Screen
	style  tcell.Style
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
		lines: []string{""},
		x:     0,
		y:     0,
		style: tcell.StyleDefault,
		s:     s,
	}
	scr.s.SetStyle(scr.style)
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
		drawText(s.s, 0, y, len(line), y, s.style, line)
	}
	fmt.Println(s.lines)
	if s.cursor {
		s.s.ShowCursor(s.x, s.y)
	}
	s.s.Show()
	return s
}

func (s *screen) Up() *screen {
	if s.y == 0 {
		return s
	}
	s.y--
	s.x = min(s.x, len(s.lines[s.y]))
	return s
}
func (s *screen) Down() *screen {
	if s.y+1 < len(s.lines) {
		s.y++
	}
	s.x = min(s.x, len(s.lines[s.y]))
	return s
}
func (s *screen) Left() *screen {
	if s.x == 0 {
		return s
	}
	s.x--
	return s
}

func (s *screen) Right(amount int) *screen {
	s.x = min(len(s.lines[s.y]), amount+s.x)
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

func (s *screen) DeleteOne() *screen {
	p := s.lines[s.y]
	s.lines[s.y] = p[:s.x-1] + p[s.x:]
	s.Left()
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
					s.Left().Draw()
				}
			case tcell.KeyRune:
				s.Write(string(ev.Rune())).Draw()
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
