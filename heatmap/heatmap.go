package heatmap

import (
	"strings"

	"github.com/amecky/table/term"
)

var STYLES = []term.Style{
	term.NewStyle(term.Hex(term.WHITE), term.Hex(term.BACKGROUND), false),
	term.NewStyle(term.Hex(term.GRAY), term.Hex(term.BACKGROUND_ODD), false),
	term.NewStyle(term.Hex(term.RED), term.Hex(term.BACKGROUND), false),
	term.NewStyle(term.Hex(term.YELLOW), term.Hex(term.BACKGROUND), false),
	term.NewStyle(term.Hex(term.BLUE), term.Hex(term.BACKGROUND), false),
	term.NewStyle(term.Hex("#008a33"), term.Hex(term.BACKGROUND), false),
	term.NewStyle(term.Hex("#82cc00"), term.Hex(term.BACKGROUND), false),
	term.NewStyle(term.Hex(term.RED), term.Hex(term.BACKGROUND_ODD), false),
	term.NewStyle(term.Hex(term.YELLOW), term.Hex(term.BACKGROUND_ODD), false),
	term.NewStyle(term.Hex(term.BLUE), term.Hex(term.BACKGROUND_ODD), false),
	term.NewStyle(term.Hex("#008a33"), term.Hex(term.BACKGROUND_ODD), false),
	term.NewStyle(term.Hex("#82cc00"), term.Hex(term.BACKGROUND_ODD), false),
}

type HeatMapLine struct {
	Name    string
	Entries []int
}
type HeatMap struct {
	name    string
	Offset  int
	Lines   []HeatMapLine
	Score   float64
	Headers []string
}

func New(name string) *HeatMap {
	return &HeatMap{
		name:   name,
		Score:  0.0,
		Offset: 0,
	}
}

func (h *HeatMap) Name(n string) *HeatMap {
	h.name = n
	return h
}

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

func (hm *HeatMap) GetLine(index int) *HeatMapLine {
	return &hm.Lines[index]
}

func (l *HeatMapLine) Add(v int) {
	l.Entries = append(l.Entries, v)
}

func (hm *HeatMap) AddValue(line, v int) {
	if v >= 0 && v <= 4 {
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
	sb := strings.Builder{}
	if hm.name != "" {
		sb.WriteString(STYLES[0].Convert(hm.name))
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
	sb.WriteString(STYLES[0].Convert(strings.Repeat("-", me+max)))
	sb.WriteRune('\n')
	for j, r := range hm.Lines {
		sb.WriteString(STYLES[j%2].Convert(r.Name))
		d := max - len(r.Name)
		if d > 0 {
			sb.WriteString(STYLES[j%2].Convert(strings.Repeat(" ", d)))
		}
		cnt := 0
		for _, v := range r.Entries {
			if v < 0 || v > 4 {
				v = 0
			}
			st := STYLES[v+2+(j%2)*5]
			sb.WriteString(st.Convert("◼ "))
			cnt += 2
		}
		d = me - cnt
		if d > 0 {
			sb.WriteString(STYLES[j%2].Convert(strings.Repeat(" ", d)))
		}
		sb.WriteRune('\n')
	}
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
