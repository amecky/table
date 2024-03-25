package table

import "strings"

type ConsoleRenderer struct {
	stylesCount      int
	builder          strings.Builder
	Styles           Styles
	additionalStyles []StyleFn
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
	/*
		styles := Styles{
			Text:                  NewStyle("#d0d0d0", "", false),
			Header:                NewStyle(HEADER_COLOR, "", true),
			PositiveMarker:        NewStyle(TEXT_COLOR, GREEN, true),
			NegativeMarker:        NewStyle(TEXT_COLOR, RED, true),
			ClassAMarker:          NewStyle(TEXT_COLOR, RED, true),
			ClassBMarker:          NewStyle(TEXT_COLOR, ORANGE, true),
			ClassCMarker:          NewStyle(TEXT_COLOR, BLUE, true),
			ClassDMarker:          NewStyle(TEXT_COLOR, LIGHT_GREEN, true),
			ClassEMarker:          NewStyle(TEXT_COLOR, GREEN, true),
			ClassFMarker:          NewStyle(TEXT_COLOR, "#209c05", true),
			HeaderStriped:         NewStyle(HEADER_COLOR, BG_COLOR, true),
			TextStriped:           NewStyle("#81858d", BG_COLOR, false),
			PositiveMarkerStriped: NewStyle(TEXT_COLOR, GREEN, true),
			NegativeMarkerStriped: NewStyle(TEXT_COLOR, RED, true),
			ClassAMarkerStriped:   NewStyle(TEXT_COLOR, RED, true),
			ClassBMarkerStriped:   NewStyle(TEXT_COLOR, ORANGE, true),
			ClassCMarkerStriped:   NewStyle(TEXT_COLOR, BLUE, true),
			ClassDMarkerStriped:   NewStyle(TEXT_COLOR, LIGHT_GREEN, true),
			ClassEMarkerStriped:   NewStyle(TEXT_COLOR, GREEN, true),
			ClassFMarkerStriped:   NewStyle(TEXT_COLOR, "#209c05", true),
		}
	*/
	styles := Styles{
		Text:                  NewStyle("#d0d0d0", "", false),
		Header:                NewStyle(HEADER_COLOR, "", true),
		PositiveMarker:        NewStyle(GREEN, "", true),
		NegativeMarker:        NewStyle(RED, "", true),
		ClassAMarker:          NewStyle(RED, "", true),
		ClassBMarker:          NewStyle(ORANGE, "", true),
		ClassCMarker:          NewStyle(BLUE, "", true),
		ClassDMarker:          NewStyle(LIGHT_GREEN, "", true),
		ClassEMarker:          NewStyle(GREEN, "", true),
		ClassFMarker:          NewStyle("#209c05", "", true),
		HeaderStriped:         NewStyle(HEADER_COLOR, BG_COLOR, true),
		TextStriped:           NewStyle("#81858d", BG_COLOR, false),
		PositiveMarkerStriped: NewStyle(GREEN, BG_COLOR, true),
		NegativeMarkerStriped: NewStyle(RED, BG_COLOR, true),
		ClassAMarkerStriped:   NewStyle(RED, BG_COLOR, true),
		ClassBMarkerStriped:   NewStyle(ORANGE, BG_COLOR, true),
		ClassCMarkerStriped:   NewStyle(BLUE, BG_COLOR, true),
		ClassDMarkerStriped:   NewStyle(LIGHT_GREEN, BG_COLOR, true),
		ClassEMarkerStriped:   NewStyle(GREEN, BG_COLOR, true),
		ClassFMarkerStriped:   NewStyle("#209c05", BG_COLOR, true),
	}
	return &ConsoleRenderer{
		stylesCount: 20,
		Styles:      styles,
	}
}

func (cr *ConsoleRenderer) AddStyle(style StyleFn) int {
	cr.additionalStyles = append(cr.additionalStyles, style)
	return len(cr.additionalStyles) - 1 + cr.stylesCount
}
func (cr *ConsoleRenderer) Append(txt string, style StyleFn) {
	cr.builder.WriteString(style(txt))
}

func (cr *ConsoleRenderer) String() string {
	return cr.builder.String()
}

func (cr *ConsoleRenderer) Marker(mk int, striped bool) StyleFn {
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
