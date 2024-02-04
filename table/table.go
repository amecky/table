package table

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/muesli/termenv"
)

const (
	ARROW_DOWN    = "▼"
	ARROW_UP      = "▲"
	DIAMOND       = "◆"
	SQUARE        = "■"
	MEDIUM_SQUARE = "◼"
	CIRCLE        = "●"
	WHITE_CIRCLE  = "○"
)

// ◼■
type FilterFunc func(r *Row) bool

type FilterDef struct {
	Header     string
	Comparator int
	Value      string
}

func BuildFilterDef(txt string) FilterDef {
	ret := FilterDef{}
	entries := strings.Split(txt, " ")
	if len(entries) == 3 {
		ret.Header = entries[0]
		ret.Value = entries[2]
		if entries[1] == "==" {
			ret.Comparator = 1
		}
		if entries[1] == "!=" {
			ret.Comparator = 2
		}
		if entries[1] == ">=" {
			ret.Comparator = 3
		}
		if entries[1] == ">" {
			ret.Comparator = 5
		}
		if entries[1] == "<=" {
			ret.Comparator = 4
		}
		if entries[1] == "<" {
			ret.Comparator = 6
		}
	}
	return ret
}

type ReportOptions struct {
	StartDate   string
	EndDate     string
	DetailsLink string
	Isin        string
	Mode        string
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

type Table struct {
	Name        string
	Created     string
	Headers     []string
	Rows        []Row
	Count       int
	HeaderSizes []int
	Limit       int
	Formatters  Formatters
}

type KeyValueTable struct {
	Table *Table
}

type Row struct {
	Size  int
	Cells []Cell
}

type Cell struct {
	Text      string
	Value     float64
	Marker    int
	Alignment TextAlign
	Link      string
}

type FormatterFn func(values []float64, index int) (string, int, int)

type Formatters struct {
	Float           FormatterFn
	Int             FormatterFn
	HistoInt        FormatterFn
	MarkedFloat     FormatterFn
	Histo           FormatterFn
	Block           FormatterFn
	HistoBlock      FormatterFn
	ColorBlock      FormatterFn
	Percentage      FormatterFn
	Relation        FormatterFn
	Categorized     FormatterFn
	CategorizedNorm FormatterFn
	BuySell         FormatterFn
}

type MarkedText struct {
	Text   string
	Marker int
}

type MarkedTextList struct {
	Items []MarkedText
}

func NewMarkedTextList() *MarkedTextList {
	ret := &MarkedTextList{}
	ret.Items = make([]MarkedText, 0)
	return ret
}

func (m *MarkedTextList) AddEmpty() {
	m.Items = append(m.Items, MarkedText{})
}

func (m *MarkedTextList) Add(txt string, marker int) {
	m.Items = append(m.Items, MarkedText{
		Text:   txt,
		Marker: marker,
	})
}

func NewKeyValueTable(name string) *KeyValueTable {
	tbl := Table{
		Name:    name,
		Headers: []string{"Name", "Value"},
		Count:   0,
		Limit:   -1,
		Created: time.Now().Format("2006-01-02 15:04"),
	}
	return &KeyValueTable{
		Table: &tbl,
	}
}

func (kvt *KeyValueTable) AddRow(key, value string) {
	r := kvt.Table.CreateRow()
	r.AddDefaultText(key)
	r.AddDefaultText(value)
}

func (kvt *KeyValueTable) AddFloat(key string, value float64, marker int) {
	r := kvt.Table.CreateRow()
	r.AddDefaultText(key)
	r.AddFloat(value, marker)
}

func (kvt *KeyValueTable) String() string {
	return kvt.Table.String()
}

func NewTable(name string, headers []string) *Table {
	tbl := Table{
		Name:    name,
		Headers: headers,
		Count:   0,
		Limit:   -1,
		Created: time.Now().Format("2006-01-02 15:04"),
	}
	tbl.Formatters = Formatters{
		Percentage: func(values []float64, index int) (string, int, int) {
			marker := 4
			vc := math.Round(values[index]*100) / 100
			if vc < 0.0 {
				marker = 2
			}
			if vc > 0.0 {
				marker = 6
			}
			return fmt.Sprintf("%.2f%%", vc), marker, 1
		},
		Relation: func(values []float64, index int) (string, int, int) {
			marker := 4
			vc := math.Round(values[index]*100) / 100
			if vc < 1.0 {
				marker = 2
			}
			if vc > 1.0 {
				marker = 6
			}
			return fmt.Sprintf("%.2f", vc), marker, 1
		},
		Categorized: func(values []float64, index int) (string, int, int) {
			v := values[index]
			marker := 4
			//vc := math.Round(values[index]*100) / 100
			if v < 20.0 {
				marker = 2
			} else if v < 40.0 {
				marker = 3
			} else if v < 60.0 {
				marker = 4
			} else if v < 80.0 {
				marker = 5
			} else {
				marker = 6
			}
			return fmt.Sprintf("%.2f", v), marker, 1
		},
		CategorizedNorm: func(values []float64, index int) (string, int, int) {
			v := values[index]
			marker := 4
			if v < 0.2 {
				marker = 2
			} else if v < 0.4 {
				marker = 3
			} else if v < 0.6 {
				marker = 4
			} else if v < 0.8 {
				marker = 5
			} else {
				marker = 6
			}
			return fmt.Sprintf("%.2f", v), marker, 1
		},
		BuySell: func(values []float64, index int) (string, int, int) {
			v := values[index]
			txt := ""
			marker := 4
			if v == 1.0 {
				txt = "BUY"
				marker = 6
			}
			if v == -1.0 {
				txt = "SELL"
				marker = 2
			}
			return txt, marker, 1
		},
		Block: func(values []float64, index int) (string, int, int) {
			marker := 4
			if values[index] < 0.0 {
				marker = 2
			} else if values[index] > 0.0 {
				marker = 6
			}
			return "■", marker, 2
		},
		HistoBlock: func(values []float64, index int) (string, int, int) {
			txt := "■"
			marker := 6
			prev := 0.0
			if index > 0 {
				prev = values[index-1]
			}
			cur := values[index]
			if prev > cur {
				txt = ARROW_DOWN
				//txt = SQUARE
			} else if prev < cur {
				txt = ARROW_UP
				//txt = CIRCLE
			} else {
				txt = DIAMOND
			}
			if cur < 0.0 {
				marker = 2
				if prev < cur {
					marker = 3
				}
			} else if cur > 0.0 {
				marker = 6
				if prev > cur {
					marker = 5
				}
			} else {
				marker = 4
			}
			return txt, marker, 2
		},

		ColorBlock: func(values []float64, index int) (string, int, int) {
			v := values[index]
			if v < 0.0 {
				v = 0.0
			}
			if v > 4.0 {
				v = 4.0
			}
			marker := int(v) + 2
			return "■", marker, 2
		},
		Float: func(values []float64, index int) (string, int, int) {
			v := values[index]
			return fmt.Sprintf("%.2f", v), 0, 1
		},
		Int: func(values []float64, index int) (string, int, int) {
			v := values[index]
			return fmt.Sprintf("%d", int(v)), 0, 1
		},
		MarkedFloat: func(values []float64, index int) (string, int, int) {
			v := values[index]
			mk := 4
			if v < 0.0 {
				mk = 2
			}
			if v > 0.0 {
				mk = 6
			}
			return fmt.Sprintf("%.2f", v), mk, 1
		},
		Histo: func(values []float64, index int) (string, int, int) {
			c := values[index]
			p := 0.0
			mk := 4
			if index > 0 {
				p = values[index-1]
			}
			if c > 0.0 && p < 0.0 {
				mk = 6
			}
			if c < 0.0 && p > 0.0 {
				mk = 2
			}
			if c < 0.0 {
				if p > c {
					mk = 2
				} else {
					mk = 3
				}
			} else {
				if p > c {
					mk = 5
				} else {
					mk = 6
				}
			}
			v := values[index]
			return fmt.Sprintf("%.2f", v), mk, 1
		},
		HistoInt: func(values []float64, index int) (string, int, int) {
			c := values[index]
			p := 0.0
			mk := 4
			if index > 0 {
				p = values[index-1]
			}
			if c > 0.0 && p < 0.0 {
				mk = 6
			}
			if c < 0.0 && p > 0.0 {
				mk = 2
			}
			if c < 0.0 {
				if p > c {
					mk = 2
				} else {
					mk = 3
				}
			} else {
				if p > c {
					mk = 5
				} else {
					mk = 6
				}
			}
			v := values[index]
			return fmt.Sprintf("%d", int(v)), mk, 1
		},
	}
	for _, header := range headers {
		tbl.HeaderSizes = append(tbl.HeaderSizes, len(header))
	}
	return &tbl
}

func (rt *Table) SetHeaderName(idx int, name string) *Table {
	rt.Headers[idx] = name
	return rt
}

func (rt *Table) FindColumnIndex(name string) int {
	for i, h := range rt.Headers {
		if h == name {
			return i
		}
	}
	return -1
}

func (rt *Table) SetText(row, col int, txt string) {
	if row < len(rt.Rows) {
		tr := rt.Rows[row]
		if col < len(tr.Cells) {
			tr.Cells[col].Text = txt
		}
	}
}

func (rt *Table) Sort(name string) {
	idx := -1
	for i, n := range rt.Headers {
		if n == name {
			idx = i
		}
	}
	if idx != -1 {
		sort.Slice(rt.Rows, func(i, j int) bool {
			return rt.Rows[i].Cells[idx].Value > rt.Rows[j].Cells[idx].Value
		})
	}
}

func (rt *Table) SortReverse(name string) {
	idx := -1
	for i, n := range rt.Headers {
		if n == name {
			idx = i
		}
	}
	if idx != -1 {
		sort.Slice(rt.Rows, func(i, j int) bool {
			return rt.Rows[i].Cells[idx].Value < rt.Rows[j].Cells[idx].Value
		})
	}
}

func (rt *Table) CreateRow() *Row {
	r := Row{
		Size: len(rt.Headers),
	}
	rt.Rows = append(rt.Rows, r)
	rt.Count++
	return &rt.Rows[rt.Count-1]
}

func (rt *Table) DelimiterLine(txt string) *Row {
	row := rt.CreateRow()
	for i := 0; i < len(rt.Headers); i++ {
		row.AddCenteredText(txt, 0)
	}
	return row
}

func (rt *Table) AddColumn(name string, values []float64, fn FormatterFn) {
	rt.Headers = append(rt.Headers, name)
	if len(rt.Rows) == 0 {
		for i := 0; i < len(values); i++ {
			r := rt.CreateRow()
			txt, mk, al := fn(values, i)
			r.AddAlignedText(txt, mk, al)
		}
	} else {
		for i := 0; i < len(rt.Rows); i++ {
			r := &rt.Rows[i]
			if i < len(values) {
				txt, mk, al := fn(values, i)
				r.AddAlignedText(txt, mk, al)
			}
		}
	}
}

func (rt *Table) AddIntColumn(name string, values []int) {
	rt.Headers = append(rt.Headers, name)
	if len(rt.Rows) == 0 {
		for i := 0; i < len(values); i++ {
			r := rt.CreateRow()
			r.AddAlignedText(fmt.Sprintf("%d", values[i]), 0, int(AlignRight))
		}
	} else {
		for i := 0; i < len(rt.Rows); i++ {
			r := &rt.Rows[i]
			if i < len(values) {
				r.AddAlignedText(fmt.Sprintf("%d", values[i]), 0, int(AlignRight))
			}
		}
	}
}

func (rt *Table) AddStringColumn(name string, values []string) {
	rt.Headers = append(rt.Headers, name)
	if len(rt.Rows) == 0 {
		for _, s := range values {
			r := rt.CreateRow()
			r.AddDefaultText(s)
		}
	} else {
		for i := 0; i < len(rt.Rows); i++ {
			r := &rt.Rows[i]
			if i < len(values) {
				r.AddDefaultText(values[i])
			}
		}
	}
}

func (rt *Table) AddMarkedTextColumn(name string, values []MarkedText) {
	rt.Headers = append(rt.Headers, name)
	if len(rt.Rows) == 0 {
		for _, s := range values {
			r := rt.CreateRow()
			r.AddCenteredText(s.Text, s.Marker)
		}
	} else {
		for i := 0; i < len(rt.Rows); i++ {
			r := &rt.Rows[i]
			if i < len(values) {
				r.AddCenteredText(values[i].Text, values[i].Marker)
			}
		}
	}
}

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

func (tr *Row) AddLink(txt, url string) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    0,
		Alignment: AlignLeft,
		Link:      url,
	})
	//}
	return tr
}

func (tr *Row) AddDefaultText(txt string) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    0,
		Alignment: AlignLeft,
	})
	//}
	return tr
}

func (tr *Row) AddEmpty() *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      "",
		Marker:    0,
		Alignment: AlignLeft,
	})
	//}
	return tr
}

func (tr *Row) AddDate(txt string) *Row {
	tmp := txt
	if strings.Index(tmp, " ") != -1 {
		tmp = tmp[0:strings.Index(tmp, " ")]
	}
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      tmp,
		Marker:    0,
		Alignment: AlignRight,
	})
	//}
	return tr
}

func (tr *Row) AddTime(txt string) *Row {
	tmp := txt
	if strings.Index(tmp, " ") != -1 {
		tmp = tmp[strings.Index(tmp, " ")+1:]
	}
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      tmp,
		Marker:    0,
		Alignment: AlignRight,
	})
	//}
	return tr
}

func (tr *Row) AddText(txt string, marker int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignLeft,
	})
	//}
	return tr
}

func (tr *Row) AddAlignedText(txt string, marker, alignment int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: TextAlign(alignment),
	})
	//}
	return tr
}

func (tr *Row) AddTextRight(txt string, marker int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignRight,
	})
	//}
	return tr
}

func (tr *Row) AddCenteredText(txt string, marker int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignCenter,
	})
	//}
	return tr
}

func (tr *Row) AddBlock(positive bool) *Row {
	//if len(tr.Cells) < tr.Size {
	marker := 6
	if !positive {
		marker = 2
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      "■",
		Marker:    marker,
		Alignment: AlignCenter,
	})
	//}
	return tr
}

func (tr *Row) AddHistoBlock(prev, cur float64) *Row {
	txt := "■"
	marker := 6
	/*
		if prev > cur {
			txt += " ↓"
		} else if prev < cur {
			txt += " ↑"
		} else {
			txt += " →"
		}
	*/
	if prev > cur {
		txt = ARROW_DOWN
		//txt = SQUARE
	} else if prev < cur {
		txt = ARROW_UP
		//txt = CIRCLE
	} else {
		txt = DIAMOND
	}
	if cur < 0.0 {
		marker = 2
		if prev < cur {
			marker = 3
		}
	} else if cur > 0.0 {
		marker = 6
		if prev > cur {
			marker = 5
		}
	} else {
		marker = 4
	}
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignCenter,
	})
	//}
	return tr
}

func (tr *Row) AddMarkedBlock(marker int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      "■",
		Marker:    marker,
		Alignment: AlignCenter,
	})
	//}
	return tr
}

func (tr *Row) AddCategoryBlock(c float64) *Row {
	//if len(tr.Cells) < tr.Size {
	txt := "■"
	if c == 1.0 || c == 5.0 {
		txt = "■ ■"
	}
	marker := 6
	if c < 3.0 {
		marker = 2
	}
	if c == 3.0 {
		marker = 4
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignLeft,
	})
	//}
	return tr
}

func (tr *Row) AddChangePercent(v float64) *Row {
	//if len(tr.Cells) < tr.Size {
	marker := 0
	vc := math.Round(v*100) / 100
	if vc < 0.0 {
		marker = -1
	}
	if vc > 0.0 {
		marker = 1
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f%%", vc),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     vc,
	})
	//}
	return tr
}

func (tr *Row) AddCategorizedFloat(v, l1, l2, l3 float64) *Row {
	mk := 4
	if v <= l1 {
		mk = 2
	} else if v <= l2 {
		mk = 3
	} else if v <= l3 {
		mk = 5
	} else if v > l3 {
		mk = 6
	}
	return tr.AddFloat(v, mk)

}

func (tr *Row) AddFloat(v float64, marker int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f", v),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     v,
	})
	//}
	return tr
}

func (tr *Row) AddExtendedFloat(v float64, pattern string, marker int) *Row {
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf(pattern, v),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     v,
	})
	//}
	return tr
}

func (tr *Row) AddFlaggedFloat(v float64, flag bool) *Row {
	marker := 1
	if !flag {
		marker = -1
	}
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f", v),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     v,
	})
	//}
	return tr
}

func (tr *Row) AddHisto(p, c float64) *Row {
	marker := 0
	dv := 0.01
	rc := math.Round(c/dv) * dv
	if rc == 0.0 {
		c = 0.0
		marker = 4
	} else if p < 0.0 && rc > 0.0 {
		marker = 6
	} else if p > 0.0 && rc < 0.0 {
		marker = 2
	} else if rc > 0.0 {
		if c > p {
			marker = 6
		} else {
			marker = 5
		}
	} else {
		if p < c {
			marker = 3
		} else {
			marker = 2
		}
	}
	//if len(tr.Cells) < tr.Size {
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f", c),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     c,
	})
	//}
	return tr
}

var CAT_DESCRIPTORS = [5]string{"Tiny", "Small", "Medium", "Big", "Huge"}

func (tr *Row) AddCategorizedPercentage(v, steps float64) *Row {
	s := 0.0
	marker := 1
	vp := v * 100.0
	for i := 0; i < 5; i++ {
		if vp >= s {
			marker++
		}
		s += steps
	}
	text := "-"
	if marker >= 2 && marker <= 6 {
		text = CAT_DESCRIPTORS[marker-2]
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      text,
			Marker:    marker,
			Alignment: AlignLeft,
			Value:     v,
		})
	}
	return tr
}

func (tr *Row) AddColouredPercentage(v, steps float64) *Row {
	s := 0.0
	marker := 1
	for i := 0; i < 5; i++ {
		if v >= s {
			marker++
		}
		s += steps
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      fmt.Sprintf("%.2f%%", v),
			Marker:    marker,
			Alignment: AlignLeft,
			Value:     v,
		})
	}
	return tr
}

func (tr *Row) AddMarkedPercentage(v, steps float64) *Row {
	s := 0.0
	marker := 1
	vp := v * 100.0
	for i := 0; i < 5; i++ {
		if vp >= s {
			marker++
		}
		s += steps
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      fmt.Sprintf("%.2f%%", vp),
			Marker:    marker,
			Alignment: AlignRight,
			Value:     v,
		})
	}
	return tr
}

func (tr *Row) AddMarkedFloat(v float64) *Row {
	dv := 0.01
	rp := math.Round(v/dv) * dv
	marker := 0
	if rp < 0.0 {
		marker = -1
	} else if rp > 0.0 {
		marker = 1
	} else {
		rp = 0.0
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      fmt.Sprintf("%.2f", rp),
			Marker:    marker,
			Alignment: AlignRight,
			Value:     v,
		})
	}
	return tr
}

func (tr *Row) AddMarkedInt(v int) *Row {
	marker := 0
	if v < 0 {
		marker = -1
	}
	if v > 0 {
		marker = 1
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      fmt.Sprintf("%d", v),
			Marker:    marker,
			Alignment: AlignRight,
			Value:     float64(v),
		})
	}
	return tr
}

func (tr *Row) AddMarkedFloatThreshold(v, t float64) *Row {
	marker := 0
	if v < t {
		marker = -1
	}
	if v >= t {
		marker = 1
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      fmt.Sprintf("%.2f", v),
			Marker:    marker,
			Alignment: AlignRight,
			Value:     v,
		})
	}
	return tr
}

func (tr *Row) AddInt(v int, marker int) *Row {
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      fmt.Sprintf("%d", v),
			Marker:    marker,
			Alignment: AlignRight,
			Value:     float64(v),
		})
	}
	return tr
}

func (tr *Row) AddChange(change, changePercent float64) *Row {
	marker := 0
	if changePercent > 0.0 {
		marker = 1
	} else if changePercent < 0.0 {
		marker = -1
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f (%.2f%%)", change, changePercent),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     change,
	})
	return tr
}

func (tr *Row) AddRelationPercentage(first, second float64) *Row {
	if second == 0.0 {
		tr.Cells = append(tr.Cells, Cell{
			Text:      "-",
			Marker:    0,
			Alignment: AlignRight,
			Value:     0.0,
		})
		return tr
	}
	changePercent := (first/second - 1.0) * 100.0
	marker := 0
	if changePercent > 0.0 {
		marker = 1
	} else if changePercent < 0.0 {
		marker = -1
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f%%", changePercent),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     changePercent,
	})
	return tr
}

func (tr *Row) AddRelation(first, second float64) *Row {
	change := first - second
	changePercent := (first/second - 1.0) * 100.0
	return tr.AddChange(change, changePercent)
}

const TableTemplate = `
<table class="table table-dark table-bordered">
        <thead>
          <tr>
            {{ range .Headers }}
            <th scope="col" style='text-align:center'>{{.}}</th>
            {{ end }}
          </tr>
        </thead>
        <tbody>
          {{ range .Rows }}
            <tr>
              {{range .Cells}}
				{{ $clr := ""}}
				{{if eq .Marker -1 }}
					{{$clr = "color:#be0000;"}}
				{{end}}
				{{if eq .Marker 2 }}
					{{$clr = "color:#be0000;"}}
				{{end}}
				{{if eq .Marker 3 }}
					{{$clr = "color:#c0a102;"}}
				{{end}}
				{{if eq .Marker 4 }}
					{{$clr = "color:#1a7091;"}}
				{{end}}
				{{if eq .Marker 5 }}
					{{$clr = "color:#166a03;"}}
				{{end}}
				{{if eq .Marker 6 }}
					{{$clr = "color:#6cc717;"}}
				{{end}}
				{{if eq .Marker 1 }}
					{{$clr = "color:#6cc717;"}}
				{{end}}

				{{$al := ""}}				
				{{if eq .Alignment 0}}
					{{$al = "text-align: left"}}
				{{else if eq .Alignment 1}}
					{{$al = "text-align: right"}}
				{{else if eq .Alignment 2}}
					{{$al = "text-align: center"}}
				{{end}}

				{{if gt .Link ""}}
					<td {{$al}}><a href="{{.Link}}">{{.Text}}</a></td>
				{{else}}
					{{if eq .Marker 0 }}
                		<td style='{{$al}}'>{{.Text}}</td>
					{{else}}
						<td style='{{$al}};{{$clr}}'>{{.Text}}</td>
					{{end}}
				{{end}}
              {{end}}
          </tr>
          {{ end }}
        </tbody>
      </table>
`

const HeadlessTableTemplate = `
<table class="table table-dark table-bordered">
        <tbody>
          {{ range .Rows }}
            <tr>
              {{range .Cells}}
				{{ $clr := ""}}
				{{if eq .Marker -1 }}
					{{$clr = "color:#be0000;"}}
				{{end}}
				{{if eq .Marker 2 }}
					{{$clr = "color:#be0000;"}}
				{{end}}
				{{if eq .Marker 3 }}
					{{$clr = "color:#c0a102;"}}
				{{end}}
				{{if eq .Marker 4 }}
					{{$clr = "color:#1a7091;"}}
				{{end}}
				{{if eq .Marker 5 }}
					{{$clr = "color:#166a03;"}}
				{{end}}
				{{if eq .Marker 6 }}
					{{$clr = "color:#6cc717;"}}
				{{end}}
				{{if eq .Marker 1 }}
					{{$clr = "color:#6cc717;"}}
				{{end}}

				{{$al := ""}}				
				{{if eq .Alignment 0}}
					{{$al = "text-align: left"}}
				{{else if eq .Alignment 1}}
					{{$al = "text-align: right"}}
				{{else if eq .Alignment 2}}
					{{$al = "text-align: center"}}
				{{end}}

				{{if gt .Link ""}}
					<td {{$al}}><a href="{{.Link}}">{{.Text}}</a></td>
				{{else}}
					{{if eq .Marker 0 }}
                		<td style='{{$al}}'>{{.Text}}</td>
					{{else}}
						<td style='{{$al}};{{$clr}}'>{{.Text}}</td>
					{{end}}
				{{end}}
              {{end}}
          </tr>
          {{ end }}
        </tbody>
      </table>
`

func (rt *Table) BuildHtml() string {
	reportTemplate, err := template.New("report").Parse(TableTemplate)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		var doc bytes.Buffer
		err := reportTemplate.Execute(&doc, rt)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		txt := "<h4>" + rt.Name + "</h4>"
		txt += doc.String()
		return txt
	}
}

func (rt *Table) BuildHeadlessHtml() string {
	reportTemplate, err := template.New("report").Parse(HeadlessTableTemplate)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		var doc bytes.Buffer
		err := reportTemplate.Execute(&doc, rt)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return doc.String()
	}
}

func (rt *Table) BuildPlainHtml() string {
	reportTemplate, err := template.New("report").Parse(TableTemplate)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		var doc bytes.Buffer
		err := reportTemplate.Execute(&doc, rt)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return doc.String()
	}
}

func (rt *Table) RebuildSizes() {
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > rt.HeaderSizes[j] {
				rt.HeaderSizes[j] = internalLen(c.Text) + 2
			}
		}
	}
}

func (tr *Table) Sub(start, end int) *Table {
	ret := NewTable(tr.Name, tr.Headers)
	if end > len(tr.Rows) {
		end = len(tr.Rows)
	}
	for i := start; i < end; i++ {
		ret.Rows = append(ret.Rows, tr.Rows[i])
	}
	return ret
}

func (tr *Table) Filter(f string) *Table {
	def := BuildFilterDef(f)
	ret := NewTable(tr.Name, tr.Headers)
	rc := tr.FindColumnIndex(def.Header)
	if rc != -1 {
		for _, r := range tr.Rows {
			add := 0
			n := r.Cells[rc].Text
			if def.Comparator == 1 && n == def.Value {
				add = 1
			}
			if def.Comparator == 2 && n != def.Value {
				add = 1
			}
			if def.Comparator == 3 && n >= def.Value {
				add = 1
			}
			if def.Comparator == 4 && n <= def.Value {
				add = 1
			}
			if def.Comparator == 5 && n > def.Value {
				add = 1
			}
			if def.Comparator == 6 && n < def.Value {
				add = 1
			}
			if add == 1 {
				ret.Rows = append(ret.Rows, r)
			}
		}
	}
	return ret
}

func (tr *Table) FilterRecent(num int) *Table {
	if num == -1 {
		return tr
	}
	ret := NewTable(tr.Name, tr.Headers)
	start := len(tr.Rows) - num
	if start < 0 {
		start = 0
	}
	for i := start; i < len(tr.Rows); i++ {
		ret.Rows = append(ret.Rows, tr.Rows[i])
	}
	return ret
}

func internalLen(txt string) int {
	return utf8.RuneCountInString(txt)
}

type StyleFn func(string) string

type Styles struct {
	Text                  StyleFn
	Header                StyleFn
	PositiveMarker        StyleFn
	NegativeMarker        StyleFn
	ClassAMarker          StyleFn
	ClassBMarker          StyleFn
	ClassCMarker          StyleFn
	ClassDMarker          StyleFn
	ClassEMarker          StyleFn
	ClassFMarker          StyleFn
	HeaderStriped         StyleFn
	TextStriped           StyleFn
	PositiveMarkerStriped StyleFn
	NegativeMarkerStriped StyleFn
	ClassAMarkerStriped   StyleFn
	ClassBMarkerStriped   StyleFn
	ClassCMarkerStriped   StyleFn
	ClassDMarkerStriped   StyleFn
	ClassEMarkerStriped   StyleFn
	ClassFMarkerStriped   StyleFn
}

func NewStyle(fg string, bg string, bold bool) func(string) string {
	p := termenv.ColorProfile()
	s := termenv.Style{}.Foreground(p.Color(fg)).Background(p.Color(bg))
	if bold {
		s = s.Bold()
	}
	return s.Styled
}

type ConsoleRenderer struct {
	builder strings.Builder
	Styles  Styles
}

// #094A25, #0C6B37, #F8B324, #EB442C, #BC2023
const (
	BG_COLOR = "#131313"
	//RED         = "#be0000"
	//ORANGE      = "#c0a102"
	//BLUE        = "#1a7091"
	//GREEN       = "#6cc717"
	//LIGHT_GREEN = "#166a03"

	RED         = "#cc0000"
	ORANGE      = "#ff7800"
	BLUE        = "#1a7091"
	LIGHT_GREEN = "#008a33"
	GREEN       = "#82cc00"

	//TEXT_COLOR  = "#efefef"
	TEXT_COLOR   = "#fff"
	HEADER_COLOR = "#81858d"
)

func NewConsoleRenderer() *ConsoleRenderer {
	/*
		styles := Styles{
			Text:                  NewStyle("#d0d0d0", "", false),
			Header:                NewStyle(HEADER_COLOR, "", true),
			PositiveMarker:        NewStyle(TEXT_COLOR, GREEN, true),
			NegativeMarker:        NewStyle(TEXT_COLOR, RED, true),
			ClassAMarker:          NewStyle(TEXT_COLOR, RED, true),
			ClassBMarker:          NewStyle(TEXT_COLOR, ORANGE, true),
			ClassCMarker:          NewStyle(TEXT_COLOR, BLUE, true),
			ClassDMarker:          NewStyle(TEXT_COLOR, LIGHT_GREEN, true),
			ClassEMarker:          NewStyle(TEXT_COLOR, GREEN, true),
			ClassFMarker:          NewStyle(TEXT_COLOR, "#209c05", true),
			HeaderStriped:         NewStyle(HEADER_COLOR, BG_COLOR, true),
			TextStriped:           NewStyle("#81858d", BG_COLOR, false),
			PositiveMarkerStriped: NewStyle(TEXT_COLOR, GREEN, true),
			NegativeMarkerStriped: NewStyle(TEXT_COLOR, RED, true),
			ClassAMarkerStriped:   NewStyle(TEXT_COLOR, RED, true),
			ClassBMarkerStriped:   NewStyle(TEXT_COLOR, ORANGE, true),
			ClassCMarkerStriped:   NewStyle(TEXT_COLOR, BLUE, true),
			ClassDMarkerStriped:   NewStyle(TEXT_COLOR, LIGHT_GREEN, true),
			ClassEMarkerStriped:   NewStyle(TEXT_COLOR, GREEN, true),
			ClassFMarkerStriped:   NewStyle(TEXT_COLOR, "#209c05", true),
		}
	*/
	styles := Styles{
		Text:                  NewStyle("#d0d0d0", "", false),
		Header:                NewStyle(HEADER_COLOR, "", true),
		PositiveMarker:        NewStyle(GREEN, "", true),
		NegativeMarker:        NewStyle(RED, "", true),
		ClassAMarker:          NewStyle(RED, "", true),
		ClassBMarker:          NewStyle(ORANGE, "", true),
		ClassCMarker:          NewStyle(BLUE, "", true),
		ClassDMarker:          NewStyle(LIGHT_GREEN, "", true),
		ClassEMarker:          NewStyle(GREEN, "", true),
		ClassFMarker:          NewStyle("#209c05", "", true),
		HeaderStriped:         NewStyle(HEADER_COLOR, BG_COLOR, true),
		TextStriped:           NewStyle("#81858d", BG_COLOR, false),
		PositiveMarkerStriped: NewStyle(GREEN, BG_COLOR, true),
		NegativeMarkerStriped: NewStyle(RED, BG_COLOR, true),
		ClassAMarkerStriped:   NewStyle(RED, BG_COLOR, true),
		ClassBMarkerStriped:   NewStyle(ORANGE, BG_COLOR, true),
		ClassCMarkerStriped:   NewStyle(BLUE, BG_COLOR, true),
		ClassDMarkerStriped:   NewStyle(LIGHT_GREEN, BG_COLOR, true),
		ClassEMarkerStriped:   NewStyle(GREEN, BG_COLOR, true),
		ClassFMarkerStriped:   NewStyle("#209c05", BG_COLOR, true),
	}
	return &ConsoleRenderer{
		Styles: styles,
	}
}

func (cr *ConsoleRenderer) Append(txt string, style StyleFn) {
	cr.builder.WriteString(style(txt))
}

func (cr *ConsoleRenderer) String() string {
	return cr.builder.String()
}

func (cr *ConsoleRenderer) BorderLine(left, right string, sizes []int) {
	cr.Append(left, cr.Styles.Header)
	for _, i := range sizes {
		cr.Append(strings.Repeat(V_LINE, i+1), cr.Styles.Header)
	}
	cr.Append(right, cr.Styles.Header)
	cr.Append("\n", cr.Styles.Header)
}

// https://en.wikipedia.org/wiki/Box-drawing_character
const (
	H_LINE    = "│"
	V_LINE    = "─"
	TL_CORNER = "┌"
	TR_CORNER = "┐"
	BR_CORNER = "┘"
	BL_CORNER = "└"
	CROSS     = "┼"
	TOP_DEL   = "┬"
	BOT_DEL   = "┴"
	LEFT_DEL  = "├"
	RIGHT_DEL = "┤"
)

func (rt *Table) Width() int {
	ret := 0
	for _, th := range rt.Headers {
		ret += internalLen(th) + 3
	}
	for _, r := range rt.Rows {
		cr := 0
		for _, c := range r.Cells {
			cr += internalLen(c.Text) + 3
		}
		if cr > ret {
			ret = cr
		}
	}
	ret += len(rt.Headers) + 4
	return ret
}

func (rt *Table) getMarker(mk int, cr *ConsoleRenderer, striped bool) StyleFn {
	st := cr.Styles.Text

	if striped {
		st = cr.Styles.TextStriped
		switch mk {
		case -1:
			st = cr.Styles.NegativeMarkerStriped
		case 1:
			st = cr.Styles.PositiveMarkerStriped
		case 2:
			st = cr.Styles.ClassAMarkerStriped
		case 3:
			st = cr.Styles.ClassBMarkerStriped
		case 4:
			st = cr.Styles.ClassCMarkerStriped
		case 5:
			st = cr.Styles.ClassDMarkerStriped
		case 6:
			st = cr.Styles.ClassEMarkerStriped
		case 7:
			st = cr.Styles.ClassFMarkerStriped
		}
	} else {
		switch mk {
		case -1:
			st = cr.Styles.NegativeMarker
		case 1:
			st = cr.Styles.PositiveMarker
		case 2:
			st = cr.Styles.ClassAMarker
		case 3:
			st = cr.Styles.ClassBMarker
		case 4:
			st = cr.Styles.ClassCMarker
		case 5:
			st = cr.Styles.ClassDMarker
		case 6:
			st = cr.Styles.ClassEMarker
		case 7:
			st = cr.Styles.ClassFMarker
		}
	}
	return st
}

/*
	func (rt *Table) String() string {
		var sizes = make([]int, 0)
		for _, th := range rt.Headers {
			sizes = append(sizes, internalLen(th)+2)
		}
		for _, r := range rt.Rows {
			for j, c := range r.Cells {
				if internalLen(c.Text)+2 > sizes[j] {
					sizes[j] = internalLen(c.Text) + 2
				}
			}
		}

		cr := NewConsoleRenderer()
		cr.Append(rt.Name, cr.Styles.Text)
		cr.Append("\n", cr.Styles.Text)

		cr.BorderLine(TL_CORNER, TR_CORNER, sizes)

		cr.Append(H_LINE, cr.Styles.Header)
		for j, h := range rt.Headers {
			cr.Append(FormatString(h, sizes[j], AlignCenter), cr.Styles.Header)
			cr.Append(" ", cr.Styles.Header)
		}
		cr.Append(H_LINE+"\n", cr.Styles.Header)

		cr.BorderLine(LEFT_DEL, RIGHT_DEL, sizes)

		for j, r := range rt.Rows {
			if j%2 == 0 {
				cr.Append(H_LINE, cr.Styles.HeaderStriped)
			} else {
				cr.Append(H_LINE, cr.Styles.Header)
			}
			if rt.Limit == -1 || j < rt.Limit {
				for i, c := range r.Cells {
					if j%2 == 0 {
						cr.Append(" ", cr.Styles.HeaderStriped)
					} else {
						cr.Append(" ", cr.Styles.Header)
					}
					st := rt.getMarker(c.Marker, cr, j%2 == 0)
					str := FormatString(c.Text, sizes[i], c.Alignment)
					cr.Append(str, st)
				}
				if j%2 == 0 {
					cr.Append(H_LINE, cr.Styles.HeaderStriped)
				} else {
					cr.Append(H_LINE, cr.Styles.Header)
				}
				cr.Append("\n", cr.Styles.Header)
			}
		}
		cr.BorderLine(BL_CORNER, BR_CORNER, sizes)
		return cr.String()
	}
*/
func (rt *Table) String() string {
	var sizes = make([]int, 0)
	for _, th := range rt.Headers {
		sizes = append(sizes, internalLen(th)+2)
	}
	total := 0
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > sizes[j] {
				sizes[j] = internalLen(c.Text) + 2
			}
		}
	}
	for _, s := range sizes {
		total += s + 1
	}
	cr := NewConsoleRenderer()
	cr.Append(rt.Name, cr.Styles.Text)
	cr.Append("\n", cr.Styles.Text)
	for j, h := range rt.Headers {
		cr.Append(FormatString(h, sizes[j], AlignCenter), cr.Styles.Header)
		cr.Append(" ", cr.Styles.Header)
	}
	cr.Append("\n", cr.Styles.Header)

	cr.Append(strings.Repeat(V_LINE, total), cr.Styles.Header)
	cr.Append("\n", cr.Styles.Header)
	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			for i, c := range r.Cells {
				if j%2 == 0 {
					cr.Append(" ", cr.Styles.HeaderStriped)
				} else {
					cr.Append(" ", cr.Styles.Header)
				}
				st := rt.getMarker(c.Marker, cr, j%2 == 0)
				str := FormatString(c.Text, sizes[i], c.Alignment)
				cr.Append(str, st)
			}
			cr.Append("\n", cr.Styles.Header)
		}
	}
	return cr.String()
}

func (rt *Table) BorderlessString() string {
	var sizes = make([]int, 0)
	for _, th := range rt.Headers {
		sizes = append(sizes, internalLen(th)+2)
	}
	total := 0
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > sizes[j] {
				sizes[j] = internalLen(c.Text) + 2
			}
		}
	}
	for _, s := range sizes {
		total += s
	}
	cr := NewConsoleRenderer()
	cr.Append(rt.Name, cr.Styles.Text)
	cr.Append("\n", cr.Styles.Text)
	for j, h := range rt.Headers {
		cr.Append(FormatString(h, sizes[j], AlignCenter), cr.Styles.Header)
		cr.Append(" ", cr.Styles.Header)
	}
	cr.Append("\n", cr.Styles.Header)

	cr.Append(strings.Repeat(V_LINE, total+1), cr.Styles.Header)
	cr.Append("\n", cr.Styles.Header)
	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			for i, c := range r.Cells {
				if j%2 == 0 {
					cr.Append(" ", cr.Styles.HeaderStriped)
				} else {
					cr.Append(" ", cr.Styles.Header)
				}
				st := rt.getMarker(c.Marker, cr, j%2 == 0)
				str := FormatString(c.Text, sizes[i], c.Alignment)
				cr.Append(str, st)
			}
			cr.Append("\n", cr.Styles.Header)
		}
	}
	return cr.String()
}

func (rt *Table) HeadlessString() string {
	var sizes = make([]int, 0)
	for _, th := range rt.Headers {
		sizes = append(sizes, internalLen(th)+2)
	}
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > sizes[j] {
				sizes[j] = internalLen(c.Text) + 2
			}
		}
	}
	cr := NewConsoleRenderer()
	// Name
	cr.Append(rt.Name, cr.Styles.Header)
	cr.Append("\n", cr.Styles.Header)
	cr.BorderLine(TL_CORNER, TR_CORNER, sizes)
	// table
	for j, r := range rt.Rows {
		cr.Append(H_LINE, cr.Styles.Header)
		if rt.Limit == -1 || j < rt.Limit {
			for i, c := range r.Cells {
				st := rt.getMarker(c.Marker, cr, j%2 == 0)
				str := FormatString(c.Text, sizes[i], c.Alignment)
				cr.Append(str, st)
				cr.Append(" ", cr.Styles.Text)
			}
		}
		cr.Append(H_LINE+"\n", cr.Styles.Header)
	}
	cr.BorderLine(BL_CORNER, BR_CORNER, sizes)
	return cr.String()
}

func (rt *Table) PlainString() string {
	var sizes = make([]int, 0)
	for _, th := range rt.Headers {
		sizes = append(sizes, internalLen(th)+2)
	}
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text)+2 > sizes[j] {
				sizes[j] = internalLen(c.Text) + 2
			}
		}
	}
	var cr strings.Builder
	cr.WriteString(rt.Name)
	cr.WriteString("\n")
	cr.WriteString("┃")
	for j, h := range rt.Headers {
		cr.WriteString(FormatString(h, sizes[j], AlignCenter))
		if j != len(sizes)-1 {
			cr.WriteString("┃")
		}
	}
	cr.WriteString("┃\n")

	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			for i, c := range r.Cells {
				cr.WriteString("┃")
				str := FormatString(c.Text, sizes[i], c.Alignment)
				cr.WriteString(fmt.Sprintf("%s", str))
			}
			cr.WriteString("┃\n")
		}
	}
	return cr.String()
}

type TableCell struct {
	Value  string `json:"value"`
	Marker int    `json:"marker"`
}
type TableRow struct {
	Cells []TableCell `json:"row"`
}

func (rt *Table) JSON(w io.Writer) error {
	var rows = make([]TableRow, 0)
	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			tr := TableRow{}
			for _, c := range r.Cells {
				tr.Cells = append(tr.Cells, TableCell{
					Value:  c.Text,
					Marker: c.Marker,
				})
			}
			rows = append(rows, tr)
		}
	}
	reply := map[string]interface{}{
		"headers": rt.Headers,
		"rows":    rows,
	}

	return json.NewEncoder(w).Encode(reply)
}

const ReportMailTemplate = `
{{ range .Sections}}
<h2>{{.Name}}</h2>
<table class="table table-bordered">
        <thead>
          <tr>
            {{ range .Headers }}
            <th scope="col">{{.}}</th>
            {{ end }}
          </tr>
        </thead>
        <tbody>
          {{ range .Lines }}
            <tr>
              {{range .Entries}}
			  {{if gt .Link ""}}
			  	<td><a href="{{.Link}}">{{ .Text}}</a></td>
			  {{else}}
              {{if eq .Marker 1 }}
                <td class="text-success" style="text-align: right;">{{.Text}}</td>
              {{else if eq .Marker -1}}
                <td class="text-danger" style="text-align: right;">{{ .Text}}</td>
              {{else}}
                <td>{{.Text}}</td>
              {{end}}
              {{end}}
			  {{end}}
          </tr>
          {{ end }}
        </tbody>
      </table>
{{end}}
`

// -------------------------------------------------------
//
// Heatmap
//
// -------------------------------------------------------

const HeatMapTemplate = `
<h4>{{.Name}}</h4>
<p>Score: {{printf "%.2f" .Score}}</p>
<table class="table table-bordered">
	<tbody>
		{{ range .Lines}}
		<tr>
			{{range .Entries}}
			<td class="hm-{{.Category}}"><b>{{.Name}}</b> <br> {{.Price}} 
			{{if gt .Change ""}}
			<br/> {{.Change}}
			{{end}}
			</td>
			{{end}}
		</tr>
		{{end}}
	</tbody>
</table>
`

const InlineHeatMapTemplate = `
{{ range .Lines}}	
<div class="row">	
	<div class="col--lg-12">
		<div class="panel panel-default">
		  	<ul class="hmul">
				{{range .Entries}}
					<il class="ilhm hmc-{{.Category}}">{{.Name}} | {{.Price}} | <b>{{.Change}}</b></il>
				{{end}}	
			</ul>
		</div>
	</div>
</div>			
{{end}}
`
