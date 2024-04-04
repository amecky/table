package term

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	ResetSeq = "0"
)

const (
	// Escape character
	ESC = '\x1b'
	// Control Sequence Introducer
	CSI = string(ESC) + "["
)

type Style struct {
	foreground Color
	background Color
	bold       bool
	flags      uint8
}

const (
	BLACK                = "#0c0c0c"
	RED                  = "#cc0000"
	GREEN                = "#4e9a06"
	YELLOW               = "#c4a000"
	BLUE                 = "#3465a4"
	PURPLE               = "#75507b"
	CYAN                 = "#06989a"
	WHITE                = "#d3d7cf"
	GRAY                 = "#81858d"
	BRIGHT_BLACK         = "#555753"
	BRIGHT_RED           = "#ef2929"
	BRIGHT_GREEN         = "#8ae234"
	BRIGHT_YELLOW        = "#fce94f"
	BRIGHT_BLUE          = "#729fcf"
	BRIGHT_PURPLE        = "#ad7fa8"
	BRIGHT_CYAN          = "#34e2e2"
	BRIGHT_WHITE         = "#eeeeec"
	BACKGROUND           = "#0c0c0c"
	BACKGROUND_ODD       = "#1c1c1c"
	FOREGROUND           = "#eeeeec"
	CURSOR               = "#bbbbbb"
	SELECTION_BACKGROUND = "#b5d5ff"
)

var TEXT_STYLE = NewStyle(WHITE, "", false)
var TEXT_STYLE_ODD = NewStyle(GRAY, "", false)

func NewStyle(f, b string, bld bool) Style {
	s := Style{}
	if f != "" {
		s.foreground = Hex(f)
		s.flags = 1
	}
	if b != "" {
		s.background = Hex(b)
		s.flags = s.flags | 2
	}
	if bld {
		s.flags = s.flags | 4
	}
	return s
}

func (s Style) Convert(t string) string {
	b := newBuffer()
	if s.flags&4 != 0 {
		b.bold()
	}
	if s.flags&1 != 0 {
		b.forground(s.foreground)
	}
	if s.flags&2 != 0 {
		b.background(s.background)
	}
	return b.text(t).String()
}

func (s Style) Debug() string {
	b := newBuffer()
	if s.bold {
		b.bold()
	}
	b.forground(s.foreground).background(s.background)
	txt := ""
	for i := 2; i < b.index; i++ {
		r := b.runes[i]
		if r == ESC {
			txt += "ESC"
		} else {
			txt += fmt.Sprintf("%c", b.runes[i])
		}
	}
	txt += "m"
	return txt
}

type Color struct {
	r byte
	g byte
	b byte
}

func Hex(scol string) Color {
	format := "#%02x%02x%02x"
	if len(scol) == 4 {
		format = "#%1x%1x%1x"
	}

	var r, g, b byte
	n, err := fmt.Sscanf(scol, format, &r, &g, &b)
	if err != nil {
		return Color{}
	}
	if n != 3 {
		return Color{}
	}

	return Color{r: r, g: g, b: b}
}

type styleBuffer struct {
	runes []rune
	index int
}

func newBuffer() *styleBuffer {
	r := make([]rune, 200)
	r[0] = ESC
	r[1] = '['
	return &styleBuffer{
		runes: r,
		index: 2,
	}
}

func (b *styleBuffer) append(r rune) {
	if b.index < 200 {
		b.runes[b.index] = r
		b.index++
	}
}

func (b *styleBuffer) byte(bt byte) {
	t := bt
	force := false
	if t >= 100 {
		tt := t / 100
		b.append(rune(tt + 48))
		t = t - tt*100
		force = true
	}
	if t >= 10 {
		tt := t / 10
		b.append(rune(tt + 48))
		t = t - tt*10
	} else if force {
		b.append(rune(48))
	}
	b.append(rune(t + 48))
}

func (b *styleBuffer) bold() *styleBuffer {
	if b.index > 2 {
		b.append(';')
	}
	b.append('1')
	return b
}

func (b *styleBuffer) forground(c Color) *styleBuffer {
	if b.index > 2 {
		b.append(';')
	}
	b.sequence("38;2;")
	b.byte(c.r)
	b.append(';')
	b.byte(c.g)
	b.append(';')
	b.byte(c.b)
	return b
}

func (b *styleBuffer) background(c Color) *styleBuffer {
	if b.index > 2 {
		b.append(';')
	}
	b.sequence("48;2;")
	b.byte(c.r)
	b.append(';')
	b.byte(c.g)
	b.append(';')
	b.byte(c.b)
	return b
}

func (b *styleBuffer) text(s string) *styleBuffer {
	b.append('m')
	for _, r := range s {
		b.append(r)
	}
	return b
}

func (b *styleBuffer) sequence(s string) *styleBuffer {
	for _, r := range s {
		b.append(r)
	}
	return b
}

func (b *styleBuffer) String() string {
	b.append(ESC)
	b.append('[')
	b.append(0)
	b.append('m')
	ret := string(b.runes[0:b.index])
	b.runes[0] = ESC
	b.runes[1] = '['
	b.index = 2
	return ret
}

type Grid struct {
	Rows []GridRow
}

func (g Grid) String() string {
	sb := strings.Builder{}
	for _, r := range g.Rows {
		sb.WriteString(r.String())

	}
	return sb.String()
}

type GridRow struct {
	Cells   []GridCell
	Padding int
}

type GridCell struct {
	Style Style
	Width int
	Text  string
	Align int
	Plain bool
}

func isSequenceChar(r rune) bool {
	if r == ESC {
		return true
	}
	if r == 0 {
		return true
	}
	if r >= 48 && r <= 57 {
		return true
	}
	if r == 59 {
		return true
	}
	return false
}

func internalLen(txt string) int {
	return utf8.RuneCountInString(txt)
}

func internalLength(s string) int {
	if strings.Contains(s, CSI) {
		pos := make([]int, 0)
		for i := 0; i < len(s)-2; i++ {
			if s[i] == ESC && s[i+1] == '[' && s[i+2] != '0' {
				pos = append(pos, i)
			}
		}
		l := 0
		for i, p := range pos {
			if i < len(pos)-1 {
				t := s[p+3 : pos[i+1]]
				start := 0
				for _, c := range t {
					if isSequenceChar(c) {
						start++
					} else {
						break
					}
				}
				t = t[start+1:]
				l += internalLen(t)
			}
		}
		return l
	}
	return len(s)
}

func (r GridRow) String() string {

	lines := make([][]string, 0)
	for _, c := range r.Cells {
		entries := strings.Split(c.Text, "\n")
		if len(entries) > 0 {
			lines = append(lines, entries)
		} else {
			lines = append(lines, []string{c.Text})
		}
	}
	ll := make([]int, 0)
	cnt := 0
	for i := 0; i < len(lines); i++ {
		ll = append(ll, len(lines[i]))
		if len(lines[i]) > cnt {
			cnt = len(lines[i])

		}
	}
	sb := strings.Builder{}
	for i := 0; i < cnt; i++ {
		for j, c := range r.Cells {
			if ll[j] > i {
				txt := lines[j][i]
				if r.Padding > 0 {
					sb.WriteString(strings.Repeat(" ", r.Padding))
				}
				d := c.Width - internalLength(txt) - 2*r.Padding
				if c.Align == 1 && d > 0 {
					if c.Plain {
						sb.WriteString(strings.Repeat(" ", d))
					} else {
						sb.WriteString(c.Style.Convert(strings.Repeat(" ", d)))
					}
				}
				if c.Plain {
					sb.WriteString(txt)
				} else {
					sb.WriteString(c.Style.Convert(txt))
				}
				if c.Align == 0 && d > 0 {
					if c.Plain {
						sb.WriteString(strings.Repeat(" ", d))
					} else {
						sb.WriteString(c.Style.Convert(strings.Repeat(" ", d)))
					}
				}
				if r.Padding > 0 {
					sb.WriteString(strings.Repeat(" ", r.Padding))
				}
			} else {
				if c.Plain {
					sb.WriteString(strings.Repeat(" ", c.Width))
				} else {
					sb.WriteString(strings.Repeat(" ", c.Width))
				}
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
