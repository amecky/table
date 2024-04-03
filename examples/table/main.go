package main

import (
	"fmt"

	"github.com/amecky/table/table"
)

func main() {
	tbl := table.New().Name("Example").Headers("H1", "H2", "H3", "H4", "H5").Border(table.RoundedBorder)
	specialStyle := tbl.AddStyle("#FF00FF", "#00FF00", false)
	for i := 1; i < 6; i++ {
		marker := 1
		if i%2 == 0 {
			marker = -1
		}
		r1 := tbl.CreateRow()
		r1.AddDefaultText(fmt.Sprintf("Test %d", i))
		r1.AddInt(i, marker)
		r1.AddFloat(float64(i), marker)
		r1.AddMarkedFloat(3.0 - float64(i))
		r1.AddText("For testing", specialStyle)
	}
	fmt.Println(tbl)
}
