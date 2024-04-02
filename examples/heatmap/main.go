package main

import (
	"fmt"

	"github.com/amecky/table/heatmap"
)

func main() {
	hm := heatmap.New("Heatmap")
	hm.AddLine("Row 1", []int{0, 1, 2, 3, 4, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4})
	hm.AddLine("Next Row", []int{0, 1, 0, 1, 0, 1, 2, 2, 3, 3, 4, 4, 3, 2, 1, 0})
	hm.AddLine("Third one", []int{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4})
	fmt.Println(hm)
}
