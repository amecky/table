package table

import (
	"bytes"
	"fmt"
	"text/template"
)

const TableTemplate = `
<table class="table table-bordered">
        <thead>
          <tr>
            {{ range .TableHeaders }}
            <th scope="col" style='text-align:center'>{{.Text}}</th>
            {{ end }}
          </tr>
        </thead>
        <tbody>
          {{ range .Rows }}
            <tr>
              {{range .Cells}}
				{{ $clr := ""}}
				{{if eq .Marker -1 }}
					{{$clr = "background-color:#ff2222;"}}
				{{end}}
				{{if eq .Marker 2 }}
					{{$clr = "background-color:#ff2222;"}}
				{{end}}
				{{if eq .Marker 3 }}
					{{$clr = "background-color:#c0a102;"}}
				{{end}}
				{{if eq .Marker 4 }}
					{{$clr = "background-color:#1a7091;"}}
				{{end}}
				{{if eq .Marker 5 }}
					{{$clr = "background-color:#21870a;"}}
				{{end}}
				{{if eq .Marker 6 }}
					{{$clr = "background-color:#00ff00;"}}
				{{end}}
				{{if eq .Marker 1 }}
					{{$clr = "background-color:#00ff00;"}}
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
					{{$clr = "color:#ff2222;"}}
				{{end}}
				{{if eq .Marker 2 }}
					{{$clr = "color:#ff2222;"}}
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

const ReportMailTemplate = `
{{ range .Sections}}
<h2>{{.Name}}</h2>
<table class="table table-bordered">
        <thead>
          <tr>
            {{ range .TableHeaders }}
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
		txt := "<h4>" + rt.Description + "</h4>"
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
