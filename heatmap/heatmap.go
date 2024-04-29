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

type HeatMapLine struct {
	Name    string
	Entries []int
}
type HeatMap struct {
	name      string
	offset    int
	recent    int
	Lines     []HeatMapLine
	Score     float64
	Headers   []string
	scheme    ColorScheme
	oddScheme ColorScheme
	padding   int
	symbols   symbols
	delimiter int
}

func New(name string) *HeatMap {
	return &HeatMap{
		name:      name,
		Score:     0.0,
		offset:    0,
		padding:   1,
		delimiter: 0,
		symbols:   BlockSymbols,
		scheme: ColorScheme{
			Text: term.NewStyle(term.WHITE, term.BACKGROUND, false),
			E:    term.NewStyle("#ff0000", term.BACKGROUND, false),
			D:    term.NewStyle("#ff6700", term.BACKGROUND, false),
			C:    term.NewStyle(term.BLUE, term.BACKGROUND, false),
			B:    term.NewStyle("#6fa287", term.BACKGROUND, false),
			A:    term.NewStyle("#00ff00", term.BACKGROUND, false),
		},
		oddScheme: ColorScheme{
			Text: term.NewStyle(term.GRAY, term.BACKGROUND_ODD, false),
			E:    term.NewStyle("#ff0000", term.BACKGROUND_ODD, false),
			D:    term.NewStyle("#ff6700", term.BACKGROUND_ODD, false),
			C:    term.NewStyle(term.BLUE, term.BACKGROUND_ODD, false),
			B:    term.NewStyle("#6fa287", term.BACKGROUND_ODD, false),
			A:    term.NewStyle("#00ff00", term.BACKGROUND_ODD, false),
		},
	}
}

func (h *HeatMap) Name(n string) *HeatMap {
	h.name = n
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

func (hm *HeatMap) String() string {

	es := " " + strings.Repeat(" ", hm.padding)
	del := "|" + strings.Repeat(" ", hm.padding)
	dl := len(del)

	sb := strings.Builder{}
	if hm.name != "" {
		sb.WriteString(hm.scheme.Text.Convert(hm.name))
	}
	sb.WriteRune('\n')
	max := 0
	for _, l := range hm.Lines {
		if len(l.Name) > max {
			max = len(l.Name)
		}
	}
	max += 2
	me := 0
	for _, l := range hm.Lines {
		if len(l.Entries)*2 > me {
			me = len(l.Entries) * 2
		}
	}
	if hm.recent > 0 {
		me = hm.recent * 2
	}
	sb.WriteString(hm.scheme.Text.Convert(strings.Repeat("-", me+max)))
	sb.WriteRune('\n')

	for j, r := range hm.Lines {
		start := hm.offset
		if hm.recent > 0 {
			start = len(r.Entries) - hm.recent
		}
		st := hm.scheme.Text
		if j%2 == 1 {
			st = hm.oddScheme.Text
		}
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
				if j%2 == 1 {
					hst = hm.oddScheme.Get(cv)
				}
				if v < 0 {
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
	return sb.String()
}

func (hm *HeatMap) Html() string {

	del := "|" + strings.Repeat(" ", hm.padding)
	dl := len(del)

	sb := strings.Builder{}
	sb.WriteString("<table class=\"table table-bordered\">\n<tbody>\n")
	for _, r := range hm.Lines {
		sb.WriteString("<tr>")
		start := hm.offset
		if hm.recent > 0 {
			start = len(r.Entries) - hm.recent
		}
		sb.WriteString("<td>")
		sb.WriteString(r.Name)
		sb.WriteString("</td>")
		cnt := 0

		for i, v := range r.Entries {
			cv := v
			if i >= start {
				if hm.delimiter > 0 && i%hm.delimiter == 0 {
					sb.WriteString(del)
					cnt += dl
				}
				if cv > 4 {
					cv = 4
				}
				if v < 0 {
					cv = 0
				}
				sb.WriteString("<td style=\"color:")
				switch cv {
				case 0:
					sb.WriteString("#ff0000")
				case 1:
					sb.WriteString("#ff6700")
				case 2:
					sb.WriteString(term.BLUE)
				case 3:
					sb.WriteString("#6fa287")
				case 4:
					sb.WriteString("#00ff00")
				}
				sb.WriteString("\">■</td>")
			}
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</tbody></table>\n")
	return sb.String()
}

/*
func (hm *HeatMap) BuildHtml() string {
	reportTemplate, err := template.New("report").Parse(HeatMapTemplate)
	if err != nil {
		fmt.Println(err)
		return "<h6>" + err.Error() + "</h6>"
	} else {
		var doc bytes.Buffer
		err := reportTemplate.Execute(&doc, hm)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return doc.String()
	}
}

func (hm *HeatMap) BuildInlineHtml() string {
	reportTemplate, err := template.New("report").Parse(InlineHeatMapTemplate)
	if err != nil {
		fmt.Println(err)
		return "<h6>" + err.Error() + "</h6>"
	} else {
		var doc bytes.Buffer
		err := reportTemplate.Execute(&doc, hm)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return doc.String()
	}
}
*/
