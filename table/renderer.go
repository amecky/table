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
	BG_COLOR_ODD = "#1C1C1C"
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
	styles := Styles{
		Text:                  term.NewStyle("#d0d0d0", "", false),
		Header:                term.NewStyle(HEADER_COLOR, "", true),
		PositiveMarker:        term.NewStyle(GREEN, "", true),
		NegativeMarker:        term.NewStyle(RED, "", true),
		ClassAMarker:          term.NewStyle(RED, "", true),
		ClassBMarker:          term.NewStyle(ORANGE, "", true),
		ClassCMarker:          term.NewStyle(BLUE, "", true),
		ClassDMarker:          term.NewStyle(LIGHT_GREEN, "", true),
		ClassEMarker:          term.NewStyle(GREEN, "", true),
		ClassFMarker:          term.NewStyle("#209c05", "", true),
		HeaderStriped:         term.NewStyle(HEADER_COLOR, BG_COLOR_ODD, true),
		TextStriped:           term.NewStyle("#81858d", BG_COLOR_ODD, false),
		PositiveMarkerStriped: term.NewStyle(GREEN, BG_COLOR_ODD, true),
		NegativeMarkerStriped: term.NewStyle(RED, BG_COLOR_ODD, true),
		ClassAMarkerStriped:   term.NewStyle(RED, BG_COLOR_ODD, true),
		ClassBMarkerStriped:   term.NewStyle(ORANGE, BG_COLOR_ODD, true),
		ClassCMarkerStriped:   term.NewStyle(BLUE, BG_COLOR_ODD, true),
		ClassDMarkerStriped:   term.NewStyle(LIGHT_GREEN, BG_COLOR_ODD, true),
		ClassEMarkerStriped:   term.NewStyle(GREEN, BG_COLOR_ODD, true),
		ClassFMarkerStriped:   term.NewStyle("#209c05", BG_COLOR_ODD, true),
	}
	return &ConsoleRenderer{
		stylesCount: 20,
		Styles:      styles,
	}
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
