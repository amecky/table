package table

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/amecky/table/term"
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
	Description  string
	Created      string
	TableHeaders []string
	Rows         []Row
	Count        int
	HeaderSizes  []int
	Limit        int
	Formatters   Formatters
	BorderStyle  Border
	PaddingSize  int
	cr           *ConsoleRenderer
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

/*
	func NewTable(name string, headers []string) *Table {
		tbl := Table{
			Description:  name,
			TableHeaders: headers,
			Count:        0,
			Limit:        -1,
			Created:      time.Now().Format("2006-01-02 15:04"),
			BorderStyle:  DefaultBorder,
			PaddingSize:  1,
		}
		tbl.Formatters = DefaultFormatters()
		for _, header := range headers {
			tbl.HeaderSizes = append(tbl.HeaderSizes, len(header))
		}
		return &tbl
	}
*/
func New() *Table {
	tbl := Table{
		Count:       0,
		Limit:       -1,
		Created:     time.Now().Format("2006-01-02 15:04"),
		BorderStyle: DefaultBorder,
		PaddingSize: 1,
		cr:          NewConsoleRenderer(),
	}
	tbl.Formatters = DefaultFormatters()
	return &tbl
}

func (rt *Table) AddStyle(fg, bg string, bold bool) int {
	return rt.cr.AddStyle(term.NewStyle(fg, bg, bold))
}

func (rt *Table) Name(name string) *Table {
	rt.Description = name
	return rt
}

func (rt *Table) Headers(headers ...string) *Table {
	rt.TableHeaders = headers
	return rt
}

func (rt *Table) Recent(count int) *Table {
	rt.Limit = count
	return rt
}

func (rt *Table) Border(border Border) *Table {
	rt.BorderStyle = border
	return rt
}

func (rt *Table) Padding(padding int) *Table {
	rt.PaddingSize = padding
	return rt
}

func (rt *Table) TableHeader(idx int, name string) *Table {
	rt.TableHeaders[idx] = name
	return rt
}

func (rt *Table) FindColumnIndex(name string) int {
	for i, h := range rt.TableHeaders {
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
	for i, n := range rt.TableHeaders {
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
	for i, n := range rt.TableHeaders {
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
		Size: len(rt.TableHeaders),
	}
	rt.Rows = append(rt.Rows, r)
	rt.Count++
	return &rt.Rows[rt.Count-1]
}

func (rt *Table) DelimiterLine(txt string) *Row {
	row := rt.CreateRow()
	for i := 0; i < len(rt.TableHeaders); i++ {
		row.AddCenteredText(txt, 0)
	}
	return row
}

func (rt *Table) AddTableHeader(name string) {
	rt.TableHeaders = append(rt.TableHeaders, name)
}

func (rt *Table) AddColumn(name string, values []float64, fn FormatterFn) {
	rt.TableHeaders = append(rt.TableHeaders, name)
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
	rt.TableHeaders = append(rt.TableHeaders, name)
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
	rt.TableHeaders = append(rt.TableHeaders, name)
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
	rt.TableHeaders = append(rt.TableHeaders, name)
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
	if d < 0 {
		d = 0
	}
	// left
	if align == AlignLeft {
		ret = txt
		if d > 0 {
			ret += strings.Repeat(" ", d)
		}
	}
	// right
	if align == AlignRight {
		if d > 0 {
			ret = strings.Repeat(" ", d)
		}
		ret += txt
	}
	// center
	if align == AlignCenter {
		d /= 2
		if d > 0 {
			ret = strings.Repeat(" ", d)
		}
		ret += txt
		d = length - internalLen(txt) - d
		if d > 0 {
			ret += strings.Repeat(" ", d)
		}
	}
	return ret
}

func (tr *Row) AddLink(txt, url string) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    0,
		Alignment: AlignLeft,
		Link:      url,
	})
	return tr
}

func (tr *Row) AddDefaultText(txt string) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    0,
		Alignment: AlignLeft,
	})
	return tr
}

func (tr *Row) AddEmpty() *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      "",
		Marker:    0,
		Alignment: AlignLeft,
	})
	return tr
}

func (tr *Row) AddDate(txt string) *Row {
	tmp := txt
	if strings.Contains(tmp, " ") {
		tmp = tmp[0:strings.Index(tmp, " ")]
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      tmp,
		Marker:    0,
		Alignment: AlignRight,
	})
	return tr
}

func (tr *Row) AddTime(txt string) *Row {
	tmp := txt
	if strings.Contains(tmp, " ") {
		tmp = tmp[strings.Index(tmp, " ")+1:]
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      tmp,
		Marker:    0,
		Alignment: AlignRight,
	})
	return tr
}

func (tr *Row) AddText(txt string, marker int) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignLeft,
	})
	return tr
}

func (tr *Row) AddAlignedText(txt string, marker, alignment int) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: TextAlign(alignment),
	})
	return tr
}

func (tr *Row) AddTextRight(txt string, marker int) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignRight,
	})
	return tr
}

func (tr *Row) AddCenteredText(txt string, marker int) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignCenter,
	})
	return tr
}

func (tr *Row) AddBlock(positive bool) *Row {
	marker := 6
	if !positive {
		marker = 2
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      "■",
		Marker:    marker,
		Alignment: AlignCenter,
	})
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
	tr.Cells = append(tr.Cells, Cell{
		Text:      txt,
		Marker:    marker,
		Alignment: AlignCenter,
	})
	return tr
}

func (tr *Row) AddMarkedBlock(marker int) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      "■",
		Marker:    marker,
		Alignment: AlignCenter,
	})
	return tr
}

func (tr *Row) AddCategoryBlock(c float64) *Row {
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
	return tr
}

func (tr *Row) AddChangePercent(v float64) *Row {
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
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f", v),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     v,
	})
	return tr
}

func (tr *Row) AddExtendedFloat(v float64, pattern string, marker int) *Row {
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf(pattern, v),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     v,
	})
	return tr
}

func (tr *Row) AddFlaggedFloat(v float64, flag bool) *Row {
	marker := 1
	if !flag {
		marker = -1
	}
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f", v),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     v,
	})
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
	tr.Cells = append(tr.Cells, Cell{
		Text:      fmt.Sprintf("%.2f", c),
		Marker:    marker,
		Alignment: AlignRight,
		Value:     c,
	})
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
			Alignment: AlignRight,
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
			Alignment: AlignRight,
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
	txt := fmt.Sprintf("%.2f", rp)
	if txt == "0.00" {
		marker = 0
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
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
	ret := New().Name(tr.Description).Headers(tr.TableHeaders...)
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
	ret := New().Name(tr.Description).Headers(tr.TableHeaders...)
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
	ret := New().Name(tr.Description).Headers(tr.TableHeaders...)
	start := len(tr.Rows) - num
	if start < 0 {
		start = 0
	}
	for i := start; i < len(tr.Rows); i++ {
		ret.Rows = append(ret.Rows, tr.Rows[i])
	}
	return ret
}

func (tr *Table) Top(num int) *Table {
	if num == -1 {
		return tr
	}
	if num > len(tr.Rows) {
		num = len(tr.Rows)
	}
	ret := New().Name(tr.Description).Headers(tr.TableHeaders...)
	for i := 0; i < num; i++ {
		ret.Rows = append(ret.Rows, tr.Rows[i])
	}
	return ret
}

func internalLen(txt string) int {
	return utf8.RuneCountInString(txt)
}

func (rt *Table) Width() int {
	ret := 0
	for _, th := range rt.TableHeaders {
		ret += internalLen(th) + rt.PaddingSize*2
	}
	for _, r := range rt.Rows {
		cr := 0
		for _, c := range r.Cells {
			cr += internalLen(c.Text) + rt.PaddingSize*2
		}
		if cr > ret {
			ret = cr
		}
	}
	ret += len(rt.TableHeaders) + 2
	return ret
}

func (rt *Table) String() string {
	var sizes = make([]int, 0)
	for _, th := range rt.TableHeaders {
		sizes = append(sizes, internalLen(th))
	}
	total := 0
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text) > sizes[j] {
				sizes[j] = internalLen(c.Text)
			}
		}
	}
	for _, s := range sizes {
		total += s + rt.PaddingSize*2
	}

	if rt.Description != "" {
		rt.cr.Append(rt.Description, rt.cr.Styles.Text)
		rt.cr.Append("\n", rt.cr.Styles.Text)
	}
	// top line
	if rt.BorderStyle.Size > 0 {
		rt.cr.Append(rt.BorderStyle.TL_CORNER, rt.cr.Styles.Header)
		for i, s := range sizes {
			rt.cr.Append(strings.Repeat(rt.BorderStyle.V_LINE, s+rt.PaddingSize*2), rt.cr.Styles.Header)
			if i < len(sizes)-1 {
				rt.cr.Append(rt.BorderStyle.TOP_DEL, rt.cr.Styles.Header)
			}
		}
		rt.cr.Append(rt.BorderStyle.TR_CORNER, rt.cr.Styles.Header)
		rt.cr.Append("\n", rt.cr.Styles.Text)
	}

	// headers
	for j, h := range rt.TableHeaders {
		if rt.BorderStyle.Size > 0 {
			rt.cr.Append(rt.BorderStyle.H_LINE, rt.cr.Styles.Header)
		}
		rt.cr.Append(strings.Repeat(" ", rt.PaddingSize), rt.cr.Styles.Header)
		rt.cr.Append(FormatString(h, sizes[j], AlignCenter), rt.cr.Styles.Header)
		rt.cr.Append(strings.Repeat(" ", rt.PaddingSize), rt.cr.Styles.Header)
	}
	if rt.BorderStyle.Size > 0 {
		rt.cr.Append(rt.BorderStyle.H_LINE, rt.cr.Styles.Header)
	}
	rt.cr.Append("\n", rt.cr.Styles.Header)

	// header delimiter line
	if rt.BorderStyle.Size > 0 {
		rt.cr.Append(rt.BorderStyle.LEFT_DEL, rt.cr.Styles.Header)
		for i, s := range sizes {
			rt.cr.Append(strings.Repeat(rt.BorderStyle.V_LINE, s+rt.PaddingSize*2), rt.cr.Styles.Header)
			if i < len(sizes)-1 {
				rt.cr.Append(rt.BorderStyle.CROSS, rt.cr.Styles.Header)
			}
		}
		rt.cr.Append(rt.BorderStyle.RIGHT_DEL, rt.cr.Styles.Header)
		rt.cr.Append("\n", rt.cr.Styles.Text)
	} else {
		rt.cr.Append(strings.Repeat(rt.BorderStyle.V_LINE, total), rt.cr.Styles.Header)
		rt.cr.Append("\n", rt.cr.Styles.Text)
	}

	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			even := j % 2
			bst := rt.cr.Styles.Header
			if even == 0 {
				bst = rt.cr.Styles.HeaderStriped
			}
			for i, c := range r.Cells {
				if rt.BorderStyle.Size > 0 {
					rt.cr.Append(rt.BorderStyle.H_LINE, bst)
				}
				rt.cr.Append(strings.Repeat(" ", rt.PaddingSize), bst)

				st := rt.cr.Marker(c.Marker, j%2 == 0)
				str := FormatString(c.Text, sizes[i], c.Alignment)
				rt.cr.Append(str, st)

				rt.cr.Append(strings.Repeat(" ", rt.PaddingSize), bst)
			}
			if rt.BorderStyle.Size > 0 {
				rt.cr.Append(rt.BorderStyle.H_LINE, bst)
			}
			rt.cr.Append("\n", rt.cr.Styles.Header)
		}
	}

	// bottom line
	if rt.BorderStyle.Size > 0 {
		rt.cr.Append(rt.BorderStyle.BL_CORNER, rt.cr.Styles.Header)
		for i, s := range sizes {
			rt.cr.Append(strings.Repeat(rt.BorderStyle.V_LINE, s+rt.PaddingSize*2), rt.cr.Styles.Header)
			if i < len(sizes)-1 {
				rt.cr.Append(rt.BorderStyle.BOT_DEL, rt.cr.Styles.Header)
			}
		}
		rt.cr.Append(rt.BorderStyle.BR_CORNER, rt.cr.Styles.Header)
		//rt.cr.Append("\n", rt.cr.Styles.Text)
	}
	return rt.cr.String()
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
