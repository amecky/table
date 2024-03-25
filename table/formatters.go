package table

import (
	"fmt"
	"math"
)

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

func DefaultFormatters() Formatters {
	return Formatters{
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
}
