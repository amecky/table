package table

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type HeatMapLine struct {
	Name    string
	Entries []int
}
type HeatMap struct {
	Name    string
	Columns int
	Offset  int
	Lines   []HeatMapLine
	Score   float64
	Headers []string
}

func NewHeatMap(name string, cols int) *HeatMap {
	return &HeatMap{
		Columns: cols,
		Name:    name,
		Score:   0.0,
		Offset:  0,
	}
}

func NewHeatMapWithRows(name string, cols int, rowNames []string) *HeatMap {
	ret := &HeatMap{
		Columns: cols,
		Name:    name,
		Score:   0.0,
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
	if line >= 0 && line < len(hm.Lines) {
		l := &hm.Lines[line]
		l.Entries = append(l.Entries, v)
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
	size := 0
	for _, th := range hm.Lines {
		cs := len(th.Name) + 2
		if cs > size {
			size = cs
		}
	}
	cr := NewConsoleRenderer()
	cr.Append(hm.Name, cr.Styles.Text)
	cr.Append("\n", cr.Styles.Text)
	cr.Append(strings.Repeat(DefaultBorder.V_LINE, size+hm.Columns*2+5), cr.Styles.Header)
	cr.Append("\n", cr.Styles.Header)
	if len(hm.Headers) > 0 {
		cr.Append(strings.Repeat(" ", size-1), cr.Styles.Header)
		start := len(hm.Headers) - hm.Columns
		if hm.Offset > 0 {
			start -= hm.Offset
		}
		if start < 0 {
			start = 0
		}
		end := start + hm.Columns
		for i, h := range hm.Headers {
			if i >= start && i < end {
				if i%5 == 0 {
					cr.Append(h, cr.Styles.Header)
					cr.Append(strings.Repeat(" ", 5), cr.Styles.Header)
				}
			}
		}
		cr.Append("\n", cr.Styles.Header)
	}
	for j, r := range hm.Lines {
		str := FormatString(r.Name, size, 0)
		if j%2 == 0 {
			cr.Append(str, cr.Styles.Text)
		} else {
			cr.Append(str, cr.Styles.TextStriped)
		}
		start := len(r.Entries) - hm.Columns
		if hm.Offset > 0 {
			start -= hm.Offset
		}
		if start < 0 {
			start = 0
		}
		end := start + hm.Columns
		for i, v := range r.Entries {
			if i >= start && i < end {

				st := cr.Styles.Text
				if v == -1 {
					cr.Append("  ", st)
				} else {
					if j%2 != 0 {
						switch v {
						case 0:
							st = cr.Styles.ClassAMarkerStriped
						case 1:
							st = cr.Styles.ClassBMarkerStriped
						case 2:
							st = cr.Styles.ClassCMarkerStriped
						case 3:
							st = cr.Styles.ClassDMarkerStriped
						case 4:
							st = cr.Styles.ClassEMarkerStriped
						}
					} else {
						switch v {
						case 0:
							st = cr.Styles.ClassAMarker
						case 1:
							st = cr.Styles.ClassBMarker
						case 2:
							st = cr.Styles.ClassCMarker
						case 3:
							st = cr.Styles.ClassDMarker
						case 4:
							st = cr.Styles.ClassEMarker
						}
					}
					cr.Append(" ◼", st)
				}
			}
		}
		cr.Append("\n", cr.Styles.Text)
	}
	return cr.String()
}

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
