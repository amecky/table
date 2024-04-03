package main

import (
	"fmt"

	"github.com/amecky/table/table"
	term "github.com/amecky/table/term"
)

func main() {
	tbl := table.New().Border(table.RoundedBorder).Headers("H1", "H2", "H3")
	for i := 0; i < 5; i++ {
		r := tbl.CreateRow()
		r.AddDefaultText(fmt.Sprintf("Row %d", i))
		r.AddInt(i, 0)
		r.AddDefaultText("Just testing")
	}
	tbl2 := table.New().Border(table.RoundedBorder).Headers("H1", "H2")
	for i := 0; i < 8; i++ {
		r := tbl2.CreateRow()
		mk := -1
		if i%2 == 0 {
			mk = 1
		}
		r.AddDefaultText(fmt.Sprintf("Row %d", i))
		r.AddInt(i*2+200, mk)
	}
	st := term.NewStyle("#ff0000", "", true)
	st2 := term.NewStyle("#ababab", "", false)
	g := term.Grid{
		Rows: []term.GridRow{
			{
				Padding: 2,
				Cells: []term.GridCell{
					{Style: term.TEXT_STYLE, Width: 20, Text: "Hello\nworld"},
					{Style: term.TEXT_STYLE_ODD, Width: 20, Text: "Second text"},
					{Style: term.TEXT_STYLE, Width: 20, Text: "Third one"},
				},
			},
			{
				Padding: 0,
				Cells: []term.GridCell{
					{Style: st2, Width: 20, Text: "--------------------"},
					{Style: st2, Width: 20, Text: "--------------------"},
					{Style: st2, Width: 20, Text: "--------------------"},
				},
			},
			{
				Padding: 2,
				Cells: []term.GridCell{
					{Style: st, Width: 20, Text: "Here\nis\nmuch\nmore"},
					{Style: st, Width: 20, Text: "Simple\n\nTest"},
					{Style: st, Width: 20, Text: "1\n2\n3\n4\n5"},
				},
			},
			{
				Padding: 0,
				Cells: []term.GridCell{
					{Style: st2, Width: 100, Text: "|123456789|123456789|123456789|123456789|123456789|123456789|123456789"},
				},
			},
			{
				Padding: 2,
				Cells: []term.GridCell{
					{Plain: true, Width: 40, Text: tbl.String()},
					{Plain: true, Width: 30, Text: tbl2.String()},
				},
			},
			{
				Padding: 0,
				Cells: []term.GridCell{
					{Style: st2, Width: 100, Text: "|123456789|123456789|123456789|123456789|123456789|123456789|123456789"},
				},
			},
		},
	}
	fmt.Println(g)
}
