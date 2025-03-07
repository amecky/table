package table

import (
	"strings"

	"github.com/amecky/table/term"
)

type ConsoleRenderer struct {
	stylesCount      int
	builder          strings.Builder
	Styles           Styles
	additionalStyles []term.Style
}

// #094A25, #0C6B37, #F8B324, #EB442C, #BC2023
const (
	BG_COLOR     = "#0C0C0C"
	BG_COLOR_ODD = "#070707"
	//RED         = "#be0000"
	//ORANGE      = "#c0a102"
	//BLUE        = "#1a7091"
	//GREEN       = "#6cc717"
	//LIGHT_GREEN = "#166a03"

	RED         = "#a21a1a"
	ORANGE      = "#d56f1a"
	BLUE        = "#1a7091"
	LIGHT_GREEN = "#389a1d"
	GREEN       = "#287114"

	//TEXT_COLOR  = "#efefef"
	TEXT_COLOR   = "#fff"
	HEADER_COLOR = "#81858d"
)

func NewConsoleRenderer() *ConsoleRenderer {
	styles := DEFAULT_STYLE
	return &ConsoleRenderer{
		stylesCount: 20,
		Styles:      styles,
	}
}

func (cr *ConsoleRenderer) SetStyle(styles Styles) {
	cr.Styles = styles
}

func (cr *ConsoleRenderer) AddStyle(style term.Style) int {
	cr.additionalStyles = append(cr.additionalStyles, style)
	return len(cr.additionalStyles) - 1 + cr.stylesCount
}
func (cr *ConsoleRenderer) Append(txt string, style term.Style) {
	cr.builder.WriteString(style.Convert(txt))
}

func (cr *ConsoleRenderer) String() string {
	return cr.builder.String()
}

func (cr *ConsoleRenderer) HeaderMarker(mk int) term.Style {
	if mk >= cr.stylesCount {
		return cr.additionalStyles[mk-cr.stylesCount]
	}
	st := cr.Styles.Text
	switch mk {
	case -1:
		st = cr.Styles.HeaderNegative
	case 1:
		st = cr.Styles.HeaderPositive
	case 2:
		st = cr.Styles.HeaderClassA
	case 3:
		st = cr.Styles.HeaderClassB
	case 4:
		st = cr.Styles.HeaderClassC
	case 5:
		st = cr.Styles.HeaderClassD
	case 6:
		st = cr.Styles.HeaderClassE
	case 7:
		st = cr.Styles.HeaderClassF
	}
	return st
}

func (cr *ConsoleRenderer) Marker(mk int, striped bool) term.Style {
	if mk >= cr.stylesCount {
		return cr.additionalStyles[mk-cr.stylesCount]
	}
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
