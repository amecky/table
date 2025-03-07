package heatmap

import (
	"strings"

	"github.com/amecky/table/term"
)

type ColorScheme struct {
	Text term.Style
	A    term.Style
	B    term.Style
	C    term.Style
	D    term.Style
	E    term.Style
}

func (cs ColorScheme) Get(idx int) term.Style {
	switch idx {
	case 0:
		return cs.E
	case 1:
		return cs.D
	case 2:
		return cs.C
	case 3:
		return cs.B
	case 4:
		return cs.A
	default:
		return cs.Text
	}
}

type symbols struct {
	symbol []string
}

// ■ ▢ ◆ ▲ ◻ ● ⬤
var BlockSymbols = symbols{
	[]string{"■", "■", "■", "■", "■"},
}

var FilledCircleSymbols = symbols{
	[]string{"⬤", "⬤", "⬤", "⬤", "⬤"},
}

var RoundSquareOutlinedSymbols = symbols{
	[]string{"▢", "▢", "▢", "▢", "▢"},
}

var ArrowSymbols = symbols{
	[]string{"⭣", "⭠", "■", "⭢", "⭡"},
}

var TriangleSymbols = symbols{
	[]string{"⯆", "⯇", "◼", "⯈", "⯅"},
}

var CharSymbols = symbols{
	[]string{"E", "D", "C", "B", "A"},
}

var WordSymbols = symbols{
	[]string{"Tiny", "Weak", "Medium", "Strong", "Huge"},
}

type HeatMapLine struct {
	Name      string
	Entries   []int
	Delimiter bool
}
type HeatMap struct {
	name      string
	offset    int
	recent    int
	Lines     []HeatMapLine
	Score     float64
	headers   []string
	scheme    ColorScheme
	oddScheme ColorScheme
	padding   int
	symbols   symbols
	delimiter int
	emptyChar string
}

func New(name string) *HeatMap {
	return &HeatMap{
		name:      name,
		Score:     0.0,
		offset:    0,
		padding:   1,
		delimiter: 0,
		symbols:   BlockSymbols,
		emptyChar: " ",
		scheme: ColorScheme{
			Text: term.NewStyle(term.WHITE, "", false),
			E:    term.NewStyle("#ff0000", "", false),
			D:    term.NewStyle("#ff6700", "", false),
			C:    term.NewStyle(term.BLUE, "", false),
			B:    term.NewStyle("#20600B", "", false),
			A:    term.NewStyle("#00ff00", "", false),
		},
		oddScheme: ColorScheme{
			Text: term.NewStyle(term.GRAY, term.BACKGROUND_ODD, false),
			E:    term.NewStyle("#ff0000", term.BACKGROUND_ODD, false),
			D:    term.NewStyle("#ff6700", term.BACKGROUND_ODD, false),
			C:    term.NewStyle(term.BLUE, term.BACKGROUND_ODD, false),
			B:    term.NewStyle("#20600B", term.BACKGROUND_ODD, false),
			A:    term.NewStyle("#00ff00", term.BACKGROUND_ODD, false),
		},
	}
}

func (h *HeatMap) Name(n string) *HeatMap {
	h.name = n
	return h
}

func (h *HeatMap) Headers(n []string) *HeatMap {
	h.headers = n
	return h
}

func (h *HeatMap) Offset(o int) *HeatMap {
	h.offset = o
	return h
}

func (h *HeatMap) Recent(r int) *HeatMap {
	h.recent = r
	return h
}

func (h *HeatMap) Symbols(s symbols) *HeatMap {
	h.symbols = s
	return h
}

func (h *HeatMap) Padding(p int) *HeatMap {
	h.padding = p
	return h
}

func (h *HeatMap) Delimiter(d int) *HeatMap {
	h.delimiter = d
	return h
}

func (h *HeatMap) EmptyChar(s string) *HeatMap {
	h.emptyChar = s
	return h
}

/*
	func NewHeatMapWithRows(name string, rowNames []string) *HeatMap {
		ret := &HeatMap{
			name:  name,
			Score: 0.0,
		}
		for _, c := range rowNames {
			ret.Lines = append(ret.Lines, HeatMapLine{
				Name: c,
			})
		}
		return ret
	}
*/
/*
func (hm *HeatMap) GetLine(index int) *HeatMapLine {
	return &hm.Lines[index]
}

func (l *HeatMapLine) Add(v int) {
	l.Entries = append(l.Entries, v)
}
*/
func (hm *HeatMap) AddValue(line, v int) {
	if v <= 4 {
		if v < 0 {
			v = -1
		}
		if line >= 0 && line < len(hm.Lines) {
			l := &hm.Lines[line]
			l.Entries = append(l.Entries, v)
		}
	}
}

func (hm *HeatMap) CreateLine(name string) int {
	hm.Lines = append(hm.Lines, HeatMapLine{
		Name: name,
	})
	return len(hm.Lines) - 1
}

func (hm *HeatMap) AddLine(name string, values []int) {
	hm.Lines = append(hm.Lines, HeatMapLine{
		Name:    name,
		Entries: values,
	})
}

func (hm *HeatMap) AddDelimiter() {
	hm.Lines = append(hm.Lines, HeatMapLine{
		Name:      "",
		Delimiter: true,
	})
}

func (hm *HeatMap) AddDelimiterText(txt string) {
	hm.Lines = append(hm.Lines, HeatMapLine{
		Name:      txt,
		Delimiter: true,
	})
}

func (hm *HeatMap) String() string {
	es := hm.emptyChar + strings.Repeat(" ", hm.padding)
	del := "|" + strings.Repeat(" ", hm.padding)
	dl := len(del)

	sb := strings.Builder{}
	if hm.name != "" {
		sb.WriteRune(' ')
		sb.WriteString(hm.scheme.Text.Convert(hm.name))
	}
	sb.WriteRune('\n')
	max := 0
	for _, l := range hm.Lines {
		if len(l.Name) > max && !l.Delimiter {
			max = len(l.Name)
		}
	}
	max += 2
	me := 0
	for _, l := range hm.Lines {
		if len(l.Entries) > me {
			me = len(l.Entries)
		}
	}
	q := hm.recent / hm.delimiter
	if hm.recent > 0 {
		me = hm.recent*(hm.padding+1) + q*(hm.padding+1)
	}
	total := me + max
	if hm.name != "" {
		sb.WriteString(hm.scheme.Text.Convert(strings.Repeat("-", me+max+1)))
		sb.WriteRune('\n')
	}
	if len(hm.headers) > 0 {
		sb.WriteString(hm.scheme.Text.Convert(strings.Repeat(" ", max-1)))
		for _, h := range hm.headers {
			sb.WriteString(hm.scheme.Text.Convert(h))
			d := me/q - 5
			sb.WriteString(hm.scheme.Text.Convert(strings.Repeat(" ", d)))
		}
		sb.WriteRune('\n')
	}
	for _, r := range hm.Lines {
		if r.Delimiter {
			if r.Name != "" {
				l := (total - len(r.Name) - 2) / 2
				sb.WriteString(hm.scheme.Text.Convert(strings.Repeat("-", l)))
				sb.WriteRune(' ')
				sb.WriteString(hm.scheme.Text.Convert(r.Name))
				sb.WriteRune(' ')
				l = total - l - len(r.Name) - 2
				sb.WriteString(hm.scheme.Text.Convert(strings.Repeat("-", l)))
			} else {
				sb.WriteString(hm.scheme.Text.Convert(strings.Repeat("-", total)))
			}
			sb.WriteRune('\n')

		} else {
			start := hm.offset
			if hm.recent > 0 {
				start = len(r.Entries) - hm.recent
			}
			st := hm.scheme.Text
			//if j%2 == 1 {
			//	st = hm.oddScheme.Text
			//}
			sb.WriteRune(' ')
			sb.WriteString(st.Convert(r.Name))
			d := max - len(r.Name)
			if d > 0 {
				sb.WriteString(st.Convert(strings.Repeat(" ", d)))
			}
			cnt := 0
			for i, v := range r.Entries {
				cv := v
				if i >= start {
					if hm.delimiter > 0 && i%hm.delimiter == 0 {
						sb.WriteString(hm.scheme.Text.Convert(del))
						cnt += dl
					}
					if cv > 4 {
						cv = 4
					}
					if v < 0 {
						cv = 0
					}
					s := hm.symbols.symbol[cv] + strings.Repeat(" ", hm.padding)
					sl := len(s)
					hst := hm.scheme.Get(cv)
					//if j%2 == 1 {
					//	hst = hm.oddScheme.Get(cv)
					//}
					if v < 0 {
						hst = hm.scheme.Text
						sb.WriteString(hst.Convert(es))
					} else {
						sb.WriteString(hst.Convert(s))
					}
					cnt += sl

				}
			}
			d = me - cnt
			if d > 0 {
				sb.WriteString(st.Convert(strings.Repeat(" ", d)))
			}
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

type HeatmapHtmlRenderer interface {
	Start() string
	StartRow() string
	EndRow() string
	RenderCell(v int) string
	End() string
}

type DefaultHeatmapHtmlRenderer struct {
	symbols symbols
}

func (dr DefaultHeatmapHtmlRenderer) Start() string {
	return "<table class=\"table table-bordered\">"
}
func (dr DefaultHeatmapHtmlRenderer) StartRow() string {
	return "<tr>\n"
}
func (dr DefaultHeatmapHtmlRenderer) EndRow() string {
	return "</tr>\n"
}
func (dr DefaultHeatmapHtmlRenderer) RenderCell(v int) string {
	ret := "<td style=\"color:"
	switch v {
	case 1:
		ret += "#ff0000"
	case 2:
		ret += "#ff6700"
	case 3:
		ret += term.BLUE
	case 4:
		ret += "#6fa287"
	case 5:
		ret += "#00ff00"
	default:
		ret += "#ffffff"
	}
	ret += "; text-align:center\">"
	ret += dr.symbols.symbol[v]
	ret += "</td>"
	return ret
}

func (dr DefaultHeatmapHtmlRenderer) End() string {
	return "</table>\n"
}

func (hm *HeatMap) RenderHtml(renderer HeatmapHtmlRenderer) string {
	sb := strings.Builder{}
	sb.WriteString(renderer.Start())
	if len(hm.headers) > 0 {
		sb.WriteString("<thead><tr>\n")
		sb.WriteString("<th>Value</th>")
		for _, h := range hm.headers {
			sb.WriteString("<th>" + h + "</th>")
		}
		sb.WriteString("</tr>")
		sb.WriteString("</thead>")
	}
	sb.WriteString("\n<tbody>\n")
	for _, r := range hm.Lines {
		sb.WriteString(renderer.StartRow())
		start := hm.offset
		if hm.recent > 0 {
			start = len(r.Entries) - hm.recent
		}
		sb.WriteString("<td>")
		sb.WriteString(r.Name)
		sb.WriteString("</td>")
		for i, v := range r.Entries {
			cv := v
			if i >= start {
				if cv > 4 {
					cv = 4
				}
				if v < 0 {
					cv = 0
				}
				sb.WriteString(renderer.RenderCell(cv))
			}
		}
		sb.WriteString(renderer.EndRow())
	}
	sb.WriteString("</tbody>\n")
	sb.WriteString(renderer.End())
	return sb.String()
}

func (hm *HeatMap) Html() string {
	return hm.RenderHtml(DefaultHeatmapHtmlRenderer{
		symbols: hm.symbols,
	})
}
