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
	styles := Styles{
		Text:                  term.NewStyle(term.Hex("#d0d0d0"), term.Hex(""), false),
		Header:                term.NewStyle(term.Hex(HEADER_COLOR), term.Hex(""), true),
		PositiveMarker:        term.NewStyle(term.Hex(GREEN), term.Hex(""), true),
		NegativeMarker:        term.NewStyle(term.Hex(RED), term.Hex(""), true),
		ClassAMarker:          term.NewStyle(term.Hex(RED), term.Hex(""), true),
		ClassBMarker:          term.NewStyle(term.Hex(ORANGE), term.Hex(""), true),
		ClassCMarker:          term.NewStyle(term.Hex(BLUE), term.Hex(""), true),
		ClassDMarker:          term.NewStyle(term.Hex(LIGHT_GREEN), term.Hex(""), true),
		ClassEMarker:          term.NewStyle(term.Hex(GREEN), term.Hex(""), true),
		ClassFMarker:          term.NewStyle(term.Hex("#209c05"), term.Hex(""), true),
		HeaderStriped:         term.NewStyle(term.Hex(HEADER_COLOR), term.Hex(BG_COLOR), true),
		TextStriped:           term.NewStyle(term.Hex("#81858d"), term.Hex(BG_COLOR), false),
		PositiveMarkerStriped: term.NewStyle(term.Hex(GREEN), term.Hex(BG_COLOR), true),
		NegativeMarkerStriped: term.NewStyle(term.Hex(RED), term.Hex(BG_COLOR), true),
		ClassAMarkerStriped:   term.NewStyle(term.Hex(RED), term.Hex(BG_COLOR), true),
		ClassBMarkerStriped:   term.NewStyle(term.Hex(ORANGE), term.Hex(BG_COLOR), true),
		ClassCMarkerStriped:   term.NewStyle(term.Hex(BLUE), term.Hex(BG_COLOR), true),
		ClassDMarkerStriped:   term.NewStyle(term.Hex(LIGHT_GREEN), term.Hex(BG_COLOR), true),
		ClassEMarkerStriped:   term.NewStyle(term.Hex(GREEN), term.Hex(BG_COLOR), true),
		ClassFMarkerStriped:   term.NewStyle(term.Hex("#209c05"), term.Hex(BG_COLOR), true),
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
