package table

import (
	"bytes"
	"fmt"
	"text/template"
)

type HeatMapEntry struct {
	Name     string
	ISIN     string
	Price    string
	Category int
	Change   string
}

type HeatMapLine struct {
	Entries []HeatMapEntry
}
type HeatMap struct {
	Name    string
	Columns int
	Current int
	Row     int
	Lines   []HeatMapLine
	Score   float64
}

func NewHeatMap(name string, cols int) *HeatMap {
	return &HeatMap{
		Columns: cols,
		Current: 0,
		Row:     -1,
		Name:    name,
		Score:   0.0,
	}
}

func (hm *HeatMap) Append(name, value, additional string, category int) {
	if hm.Current%hm.Columns == 0 {
		hm.Row++
		hm.Lines = append(hm.Lines, HeatMapLine{})
	}
	entry := HeatMapEntry{
		Name:     name,
		ISIN:     "",
		Change:   additional,
		Category: category,
		Price:    value,
	}

	hm.Lines[hm.Row].Entries = append(hm.Lines[hm.Row].Entries, entry)
	hm.Current++
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
