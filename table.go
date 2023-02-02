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
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
			Marker:    0,
			Alignment: AlignLeft,
			Link:      url,
		})
	}
	return tr
}

func (tr *Row) AddDefaultText(txt string) *Row {
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
			Marker:    0,
			Alignment: AlignLeft,
		})
	}
	return tr
}

func (tr *Row) AddDate(txt string) *Row {
	tmp := txt
	if strings.Index(tmp, " ") != -1 {
		tmp = tmp[0:strings.Index(tmp, " ")]
	}
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      tmp,
			Marker:    0,
			Alignment: AlignRight,
		})
	}
	return tr
}

func (tr *Row) AddText(txt string, marker int) *Row {
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
			Marker:    marker,
			Alignment: AlignLeft,
		})
	}
	return tr
}

func (tr *Row) AddAlignedText(txt string, marker, alignment int) *Row {
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
			Marker:    marker,
			Alignment: TextAlign(alignment),
		})
	}
	return tr
}

func (tr *Row) AddTextRight(txt string, marker int) *Row {
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
			Marker:    marker,
			Alignment: AlignRight,
		})
	}
	return tr
}

func (tr *Row) AddCenteredText(txt string, marker int) *Row {
	if len(tr.Cells) < tr.Size {
		tr.Cells = append(tr.Cells, Cell{
			Text:      txt,
			Marker:    marker,
			Alignment: AlignCenter,
		})
	}
	return tr
}

func (tr *Row) AddBlock(positive bool) *Row {
	if len(tr.Cells) < tr.Size {
		marker := 6
		if !positive {
			marker = 2
		}
		tr.Cells = append(tr.Cells, Cell{
			Text:      "■",
			Marker:    marker,
			Alignment: AlignCenter,
		})
	}
	return tr
}

func (tr *Row) AddCategoryBlock(c float64) *Row {
	if len(tr.Cells) < tr.Size {
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
	}
	return tr
}

func (tr *Row) AddChangePercent(v float64) *Row {
	if len(tr.Cells) < tr.Size {
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
	}
	return tr
}

func (tr *Row) AddFloat(v float64, marker int) *Row {
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

func (tr *Row) AddMarkedFloat(v float64) *Row {
	marker := 0
	if v < 0.0 {
		marker = -1
	}
	if v > 0.0 {
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
<table class="table table-dark table-striped-columns">
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
			  	{{ $cls := ""}}
              	{{if eq .Marker 1 }}
			  		{{$cls = "class='text-success'"}}
				{{else if eq .Marker -1}}	  
					{{$cls = "class='text-danger'"}}
				{{end}}

				{{ $clr := ""}}
				{{if eq .Marker 2 }}
					{{$clr = "color:#ee4035;"}}
				{{end}}
				{{if eq .Marker 3 }}
					{{$clr = "color:#f37736;"}}
				{{end}}
				{{if eq .Marker 4 }}
					{{$clr = "color:#0392cf;"}}
				{{end}}
				{{if eq .Marker 5 }}
					{{$clr = "color:#fdf498;"}}
				{{end}}
				{{if eq .Marker 6 }}
					{{$clr = "color:#7bc043;"}}
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
					<td {{$cls}} {{$al}}><a href="{{.Link}}">{{.Text}}</a></td>
				{{else}}
					{{if eq .Marker 0 }}
                		<td {{$cls}} style='{{$al}}'>{{.Text}}</td>
					{{else}}
						<td {{$cls}} style='{{$al}};{{$clr}}'>{{.Text}}</td>
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
	Text           StyleFn
	Header         StyleFn
	PositiveMarker StyleFn
	NegativeMarker StyleFn
	ClassAMarker   StyleFn
	ClassBMarker   StyleFn
	ClassCMarker   StyleFn
	ClassDMarker   StyleFn
	ClassEMarker   StyleFn
	ClassFMarker   StyleFn
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

// colors := []string{"#ee4035", "#f37736", "#0392cf", "#fdf498", "#7bc043"}
func NewConsoleRenderer() *ConsoleRenderer {
	styles := Styles{
		Text: NewStyle(
			"#d0d0d0",
			"",
			false,
		),
		Header: NewStyle(
			"#81858d",
			"",
			true,
		),
		PositiveMarker: NewStyle(
			"#b2e539",
			"",
			true,
		),
		NegativeMarker: NewStyle(
			"#ff7940",
			"",
			true,
		),
		ClassAMarker: NewStyle(
			"#ee4035",
			"",
			true,
		),
		ClassBMarker: NewStyle(
			"#f37736",
			"",
			true,
		),
		ClassCMarker: NewStyle(
			"#0392cf",
			"",
			true,
		),
		ClassDMarker: NewStyle(
			"#fdf498",
			"",
			true,
		),
		ClassEMarker: NewStyle(
			"#7bc043",
			"",
			true,
		),
		ClassFMarker: NewStyle(
			"#7584d9",
			"",
			true,
		),
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
	cr.Append("┌", cr.Styles.Header)
	for j, i := range sizes {
		cr.Append(strings.Repeat("─", i), cr.Styles.Header)
		if j != len(sizes)-1 {
			cr.Append("┬", cr.Styles.Header)
		}
	}
	cr.Append("┐\n", cr.Styles.Header)

	cr.Append("│", cr.Styles.Header)
	for j, h := range rt.Headers {
		cr.Append(FormatString(h, sizes[j], AlignCenter), cr.Styles.Header)
		if j != len(sizes)-1 {
			cr.Append("│", cr.Styles.Header)
		}
	}
	cr.Append("│\n", cr.Styles.Header)

	cr.Append("├", cr.Styles.Header)
	for j, i := range sizes {
		cr.Append(strings.Repeat("─", i), cr.Styles.Header)
		if j != len(sizes)-1 {
			cr.Append("┼", cr.Styles.Header)
		}
	}
	cr.Append("┤\n", cr.Styles.Header)

	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			for i, c := range r.Cells {
				cr.Append("|", cr.Styles.Header)

				st := cr.Styles.Text
				switch c.Marker {
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
				str := FormatString(c.Text, sizes[i], c.Alignment)
				cr.Append(fmt.Sprintf("%s", str), st)
			}
			cr.Append("|\n", cr.Styles.Header)
		}
	}

	cr.Append("└", cr.Styles.Header)
	for j, i := range sizes {
		cr.Append(strings.Repeat("─", i), cr.Styles.Header)
		if j != len(sizes)-1 {
			cr.Append("┴", cr.Styles.Header)
		}
	}
	cr.Append("┘\n", cr.Styles.Header)
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
	cr.WriteString("│")
	for j, h := range rt.Headers {
		cr.WriteString(FormatString(h, sizes[j], AlignCenter))
		if j != len(sizes)-1 {
			cr.WriteString("│")
		}
	}
	cr.WriteString("│\n")

	for j, r := range rt.Rows {
		if rt.Limit == -1 || j < rt.Limit {
			for i, c := range r.Cells {
				cr.WriteString("|")
				str := FormatString(c.Text, sizes[i], c.Alignment)
				cr.WriteString(fmt.Sprintf("%s", str))
			}
			cr.WriteString("|\n")
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
