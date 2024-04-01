package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/amecky/table/table"
	term "github.com/amecky/table/term"
)

type TerminalCell struct {
	Char  rune
	Style int
}

type TerminalText struct {
	txt   string
	style int
}

type TerminalLine struct {
	Row int

	Cells []TerminalCell
}

type TerminalMatrix struct {
	out    io.Writer
	Height int
	Width  int
	Lines  []TerminalLine
	Styles []term.Style
	First  bool
}

func NewStyle(fg string, bg string, bold bool) term.Style {
	return term.NewStyle(term.Hex(fg), term.Hex(bg), bold)
}

const (
	BG_COLOR = "#0c0c0c"

	MARKED_BG_COLOR = "#1c1c1c"
)

const (
	ST_TEXT = iota
	ST_HEADER
	ST_TWO
	ST_THREE
	ST_MARKER_RED
	ST_MARKER_ORANGE
	ST_MARKER_BLUE
	ST_MARKER_LIGHT_GREEN
	ST_MARKER_GREEN
	ST_MARKER_BG_RED
	ST_MARKER_BG_ORANGE
	ST_MARKER_BG_BLUE
	ST_MARKER_BG_LIGHT_GREEN
	ST_MARKER_BG_GREEN
	ST_HEADER_BG
	ST_INFO
	ST_INFO_BG
	ST_HEADER_2
	ST_DARK
	ST_ALERT
	ST_INFO_MESSAGE
)

func NewTerminalMatrix(width, height int) *TerminalMatrix {
	ret := TerminalMatrix{
		Height: height,
		Width:  width,
		out:    os.Stdout,
		First:  true,
	}
	for h := 0; h < height; h++ {
		line := TerminalLine{
			Row: height - h - 1,
		}
		line.Cells = make([]TerminalCell, width)
		ret.Lines = append(ret.Lines, line)
	}
	// ST_TEXT
	ret.Styles = append(ret.Styles, NewStyle("#d0d0d0", BG_COLOR, false))
	// ST_HEADER
	ret.Styles = append(ret.Styles, NewStyle("#81858d", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#b2e539", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#ff7940", BG_COLOR, true))
	// ST_MARKER_RED
	ret.Styles = append(ret.Styles, NewStyle("#ee4035", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#f37736", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#0392cf", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#fdf498", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#7bc043", BG_COLOR, true))
	// ST_MARKER_BG_RED
	ret.Styles = append(ret.Styles, NewStyle("#ee4035", MARKED_BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#f37736", MARKED_BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#0392cf", MARKED_BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#fdf498", MARKED_BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#7bc043", MARKED_BG_COLOR, true))
	// ST_HEADER_BG
	ret.Styles = append(ret.Styles, NewStyle("#81858d", MARKED_BG_COLOR, false))
	// ST_INFO
	ret.Styles = append(ret.Styles, NewStyle("#7584d9", BG_COLOR, true))
	// ST_INFO_BG
	ret.Styles = append(ret.Styles, NewStyle("#7584d9", MARKED_BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#d0d0d0", BG_COLOR, true))
	ret.Styles = append(ret.Styles, NewStyle("#81858d", MARKED_BG_COLOR, true))

	ret.Styles = append(ret.Styles, NewStyle("#efefef", "#ee4035", true))
	ret.Styles = append(ret.Styles, NewStyle("#efefef", "#7584d9", true))
	return &ret

}

func (t *TerminalMatrix) AddStyle(fg, bg string, bold bool) int {
	t.Styles = append(t.Styles, NewStyle(fg, bg, bold))
	return len(t.Styles) - 1
}

func TT(txt string, style int) TerminalText {
	return TerminalText{
		txt:   txt,
		style: style,
	}
}

func TV(v float64, style int) TerminalText {
	return TerminalText{
		txt:   fmt.Sprintf("%.2f", v),
		style: style,
	}
}

func TMV(v float64) TerminalText {
	style := 4
	if v >= 0.0 {
		style = 8
	}
	return TerminalText{
		txt:   fmt.Sprintf("%.2f", v),
		style: style,
	}
}

func (tm *TerminalMatrix) Set(x, y int, txt rune, style int) {
	if x < tm.Width && x >= 0 && y >= 0 && y < tm.Height {
		r := &tm.Lines[y]
		r.Cells[x].Char = txt
		r.Cells[x].Style = style
	}
}

func (tm *TerminalMatrix) writeLines(x, y int, txt string, style int) int {
	ret := 0
	if x < tm.Width && x >= 0 && y >= 0 && y < tm.Height && x+len(txt) < tm.Width {
		r := &tm.Lines[y]
		ln := len(txt)
		if x+ln > tm.Width {
			ln = tm.Width - x
		}
		ra := []rune(txt)
		for i, cr := range ra {
			if i < ln {
				r.Cells[x+i].Char = cr
				r.Cells[x+i].Style = style
				ret++
			}
		}
	}
	return ret
}

func (tm *TerminalMatrix) Write(x, y int, txt string, style int) int {
	if strings.Contains("\n", txt) {
		ret := 0
		lines := strings.Split(txt, "\n")
		for i, l := range lines {
			tmp := tm.writeLines(x, y+i, l, style)
			if tmp > ret {
				ret = tmp
			}
		}
		return ret
	} else {
		return tm.writeLines(x, y, txt, style)
	}
}

func (tm *TerminalMatrix) WriteMarkedFloat(x, y int, txt string, v float64) int {
	ln := tm.Write(x, y, txt, 0)
	nt := fmt.Sprintf("%.2f", v)
	if v >= 0.0 {
		ln += tm.Write(x+ln, y, nt, 6)
	} else {
		ln += tm.Write(x+ln, y, nt, 2)
	}
	return ln
}

func (tm *TerminalMatrix) WriteText(x, y int, TT ...TerminalText) int {
	xp := x
	ln := 0
	for _, t := range TT {
		// FIXME: space between?
		l := tm.Write(xp, y, t.txt, t.style)
		xp += l + 1
		ln += l + 1
	}
	return ln
}

func (tm *TerminalMatrix) Clear() {
	for y := 0; y < tm.Height; y++ {
		r := &tm.Lines[y]
		for x := 0; x < tm.Width; x++ {
			r.Cells[x].Char = 0
			r.Cells[x].Style = 0
		}
	}
}

func (tm *TerminalMatrix) ClearBox(x, y, w, h int) {
	if x+w > tm.Width {
		w = tm.Width - x
	}
	if y+h > tm.Height {
		h = tm.Height - y
	}
	for i := y; i < y+h; i++ {
		l := tm.Lines[i]
		for j := x; j < x+w; j++ {
			r := &l.Cells[j]
			r.Char = 0
			r.Style = 0
		}
	}
}

func (tm *TerminalMatrix) FillBox(x, y, w, h, style int) {
	if x+w > tm.Width {
		w = tm.Width - x
	}
	if y+h > tm.Height {
		h = tm.Height - y
	}
	for i := y; i < y+h; i++ {
		l := tm.Lines[i]
		for j := x; j < x+w; j++ {
			r := &l.Cells[j]
			r.Char = ' '
			r.Style = style
		}
	}
}

func (tm *TerminalMatrix) Flush() {
	out := new(bytes.Buffer)
	//hideCursor(out)
	//showCursor(out)
	if !tm.First {
		for i := 0; i < tm.Height; i++ {
			clearLine(out)
			cursorUp(out)
		}
	}
	if tm.First {
		tm.First = false
	}
	//var ret strings.Builder
	for _, l := range tm.Lines {
		io.WriteString(out, fmt.Sprintf("%2d: ", l.Row))
		for _, r := range l.Cells {

			if r.Char == 0 {
				io.WriteString(out, " ")
			} else {
				io.WriteString(out, tm.Styles[r.Style].Convert(fmt.Sprintf("%c", r.Char)))
			}
		}
		if l.Row != 0 {
			io.WriteString(out, "\n")
		}
	}
	tm.out.Write(out.Bytes())
}

func (tm *TerminalMatrix) String() string {
	var ret strings.Builder
	for _, l := range tm.Lines {
		//ret.WriteString(fmt.Sprintf("%2d: ", l.Row))
		for _, r := range l.Cells {

			if r.Char == 0 {
				ret.WriteString(" ")
			} else {
				ret.WriteString(tm.Styles[r.Style].Convert(fmt.Sprintf("%c", r.Char)))
			}
		}
		if l.Row != 0 {
			ret.WriteString("\n")
		}
	}
	return ret.String()
}

const (
	H_LINE     = '│'
	V_LINE     = '─'
	TL_CORNER  = '┌'
	TR_CORNER  = '┐'
	BR_CORNER  = '┘'
	BL_CORNER  = '└'
	CROSS      = '┼'
	TOP_DEL    = '┬'
	BOT_DEL    = '┴'
	LEFT_DEL   = '├'
	RIGHT_DEL  = '┤'
	SM_H_LINE  = '|'
	BLOCK      = '■'
	FULL_BLOCK = '█'
	SPACE      = ' '
)

func (tm *TerminalMatrix) BoxWithHeader(x, y, width, height int, header string, style int) {
	tm.Box(x, y, width, height, style)
	txt := fmt.Sprintf(" %s ", header)
	px := (width - len(txt)) / 2
	tm.WriteText(x+px, y, TT(txt, style))
}

func (tm *TerminalMatrix) Box(x, y, width, height, style int) {
	if x < tm.Width && x >= 0 && y >= 0 && y < tm.Height {
		tm.Set(x, y, TL_CORNER, style)
		tm.Set(x, y+height, BL_CORNER, style)
		tm.Set(x+width, y, TR_CORNER, style)
		tm.Set(x+width, y+height, BR_CORNER, style)
		for i := x + 1; i < x+width; i++ {
			tm.Set(i, y, V_LINE, style)
			tm.Set(i, y+height, V_LINE, style)
		}
		for i := y + 1; i < y+height; i++ {
			tm.Set(x, i, H_LINE, style)
			tm.Set(x+width, i, H_LINE, style)
		}

	}
}

func (tm *TerminalMatrix) VLine(x, y, width, style int) {
	if x < tm.Width && x >= 0 && y >= 0 && y < tm.Height {
		for i := x; i < x+width; i++ {
			tm.Set(i, y, V_LINE, style)
		}
	}
}

func (tm *TerminalMatrix) HLine(x, y, height, style int) {
	if x < tm.Width && x >= 0 && y >= 0 && y < tm.Height {
		for i := y; i < y+height; i++ {
			tm.Set(x, i, H_LINE, style)
		}

	}
}

type TextAlign int

const (
	// AlignLeft align text within a cell
	AlignLeft TextAlign = iota
	// AlignRight align text within a cell
	AlignRight
	// AlignCenter align
	AlignCenter
)

func FormatString(txt string, length int, align TextAlign) string {
	var ret string
	d := length - internalLen(txt)
	// left
	if align == AlignLeft {
		ret = " " + txt + strings.Repeat(" ", d-1)
	}
	// right
	if align == AlignRight {
		ret = strings.Repeat(" ", d-1) + txt + " "
	}
	// center
	if align == AlignCenter {
		d /= 2
		if d > 0 {
			for i := 0; i < d; i++ {
				ret += " "
			}
		}
		ret += txt
		d = length - internalLen(txt) - d
		if d > 0 {
			for i := 0; i < d; i++ {
				ret += " "
			}
		}
	}
	return ret
}

func internalLen(txt string) int {
	return utf8.RuneCountInString(txt)
}

func (tm *TerminalMatrix) GetMarker(mk int, striped bool) int {
	st := 0
	if striped {
		st = 14
		switch mk {
		case -1:
			st = ST_MARKER_BG_RED
		case 1:
			st = ST_MARKER_BG_GREEN
		case 2:
			st = ST_MARKER_BG_RED
		case 3:
			st = 10
		case 4:
			st = 11
		case 5:
			st = 12
		case 6:
			st = ST_MARKER_BG_GREEN
		case 7:
			st = 14
		}
	} else {
		switch mk {
		case -1:
			st = ST_MARKER_RED
		case 1:
			st = ST_MARKER_GREEN
		case 2:
			st = ST_MARKER_RED
		case 3:
			st = 5
		case 4:
			st = 6
		case 5:
			st = 7
		case 6:
			st = ST_MARKER_GREEN
		case 7:
			st = 9
		}
	}
	return st
}

func (tm *TerminalMatrix) ConvertTable(x, y int, rt *table.Table) {
	var sizes = make([]int, 0)
	var indices = make([]int, 0)
	total := 0
	ind := 0
	for _, th := range rt.TableHeaders {
		sizes = append(sizes, internalLen(th)+2)
	}
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > sizes[j] {
				sizes[j] = internalLen(c.Text) + 1
			}

		}
	}
	for _, s := range sizes {
		indices = append(indices, ind)
		ind += s + 1
		total += s + 1
	}
	// Name
	tm.Write(x, y, rt.Description, 0)
	yp := y + 1
	tm.VLine(x, yp, total, 1)
	yp++
	xp := x
	// Headers
	for j, h := range rt.TableHeaders {
		xp = x + indices[j] + 1
		tm.Write(xp, yp, FormatString(h, sizes[j], AlignCenter), 0)
	}
	yp++
	tm.VLine(x, yp, total, 1)
	yp++
	xp = x
	for j, r := range rt.Rows {
		//bs := 1
		//if j%2 == 0 {
		//	bs = 18
		//}
		for i, c := range r.Cells {
			st := tm.GetMarker(c.Marker, j%2 == 0)
			//tm.Set(x+indices[i], yp, H_LINE, bs)
			al := AlignCenter
			if c.Alignment == table.AlignLeft {
				al = AlignLeft
			}
			if c.Alignment == table.AlignRight {
				al = AlignRight
			}
			str := FormatString(c.Text, sizes[i], al)
			tm.Write(x+indices[i], yp, str, st)
		}
		//tm.Set(x+total, yp, H_LINE, bs)
		yp++
	}
}

func (tm *TerminalMatrix) ConvertHeatMap(x, y int, rt *table.HeatMap) {
	tm.Write(x, y, rt.Name, ST_TEXT)
	sz := 0
	for j, r := range rt.Lines {
		st := tm.GetMarker(0, j%2 == 0)
		tmp := tm.Write(x, y+j+3, r.Name, st)
		if tmp > sz {
			sz = tmp
		}

	}
	sz += 2
	tm.VLine(x, y+1, sz+rt.Columns*2, ST_TEXT)
	start := len(rt.Headers) - rt.Columns
	if start < 0 {
		start = 0
	}
	xp := x + sz - 1
	for i, h := range rt.Headers {
		if i >= start {
			if i%5 == 0 {
				tm.Write(xp, y+2, h, ST_TEXT)
				xp += 10
			}
		}
	}
	yp := y + 3
	for j, r := range rt.Lines {
		start := len(r.Entries) - rt.Columns
		if start < 0 {
			start = 0
		}
		for i, c := range r.Entries {
			if i >= start {
				st := tm.GetMarker(c+2, j%2 == 0)
				// ●•
				tm.Write(x+(i-start)*2+sz, y+j+3, " ◼", st)
			}
		}
		yp++
	}
}

func (tm *TerminalMatrix) ConvertFullTable(x, y int, rt *table.Table) {
	borderStyle := ST_HEADER
	var sizes = make([]int, 0)
	var indices = make([]int, 0)
	total := 0
	ind := 0
	for _, th := range rt.TableHeaders {
		sizes = append(sizes, internalLen(th)+2)
	}
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > sizes[j] {
				sizes[j] = internalLen(c.Text) + 2
			}

		}
	}
	for _, s := range sizes {
		indices = append(indices, ind)
		ind += s + 1
		total += s + 1
	}
	tm.Write(x, y, rt.Description, 0)
	yp := y + 1
	xp := x
	tm.Set(xp, yp, TL_CORNER, borderStyle)
	tm.VLine(xp+1, yp, total, borderStyle)
	for i := 1; i < len(indices); i++ {
		tm.Set(x+indices[i], yp, TOP_DEL, borderStyle)
	}
	tm.Set(xp+total, yp, TR_CORNER, borderStyle)
	yp++
	xp = x
	tm.Set(xp, yp, H_LINE, borderStyle)
	for j, h := range rt.TableHeaders {
		xp = x + indices[j] + 1
		tm.Write(xp, yp, FormatString(h, sizes[j], AlignCenter), 0)
	}
	for i := 1; i < len(indices); i++ {
		tm.Set(x+indices[i], yp, H_LINE, borderStyle)
	}
	tm.Set(x+total, yp, H_LINE, borderStyle)
	yp++
	tm.Set(x, yp, LEFT_DEL, borderStyle)
	tm.VLine(x+1, yp, total, borderStyle)
	for i := 1; i < len(indices); i++ {
		tm.Set(x+indices[i], yp, CROSS, borderStyle)
	}
	tm.Set(x+total, yp, RIGHT_DEL, borderStyle)
	yp++
	xp = x
	for j, r := range rt.Rows {
		bs := ST_HEADER
		if j%2 == 0 {
			bs = ST_HEADER_BG
		}
		for i, c := range r.Cells {
			st := tm.GetMarker(c.Marker, j%2 == 0)
			tm.Set(x+indices[i], yp, H_LINE, bs)
			al := AlignCenter
			if c.Alignment == table.AlignLeft {
				al = AlignLeft
			}
			if c.Alignment == table.AlignRight {
				al = AlignRight
			}
			str := FormatString(c.Text, sizes[i], al)
			tm.Write(x+indices[i]+1, yp, str, st)
		}
		tm.Set(x+total, yp, H_LINE, bs)
		yp++
	}
	tm.Set(xp, yp, BL_CORNER, borderStyle)
	tm.VLine(xp+1, yp, total, borderStyle)
	for i := 1; i < len(indices); i++ {
		tm.Set(x+indices[i], yp, BOT_DEL, borderStyle)
	}
	tm.Set(xp+total, yp, BR_CORNER, borderStyle)
}
