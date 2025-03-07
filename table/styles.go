package table

import "github.com/amecky/table/term"

type Styles struct {
	Text                  term.Style
	Header                term.Style
	PositiveMarker        term.Style
	NegativeMarker        term.Style
	ClassAMarker          term.Style
	ClassBMarker          term.Style
	ClassCMarker          term.Style
	ClassDMarker          term.Style
	ClassEMarker          term.Style
	ClassFMarker          term.Style
	HeaderStriped         term.Style
	TextStriped           term.Style
	PositiveMarkerStriped term.Style
	NegativeMarkerStriped term.Style
	ClassAMarkerStriped   term.Style
	ClassBMarkerStriped   term.Style
	ClassCMarkerStriped   term.Style
	ClassDMarkerStriped   term.Style
	ClassEMarkerStriped   term.Style
	ClassFMarkerStriped   term.Style
	HeaderPositive        term.Style
	HeaderNegative        term.Style
	HeaderClassA          term.Style
	HeaderClassB          term.Style
	HeaderClassC          term.Style
	HeaderClassD          term.Style
	HeaderClassE          term.Style
	HeaderClassF          term.Style
}

var DEFAULT_STYLE = Styles{
	Text:                  term.NewStyle("#d0d0d0", "", false),
	Header:                term.NewStyle(HEADER_COLOR, "", true),
	PositiveMarker:        term.NewStyle(LIGHT_GREEN, "", true),
	NegativeMarker:        term.NewStyle(RED, "", true),
	ClassAMarker:          term.NewStyle(RED, "", true),
	ClassBMarker:          term.NewStyle(ORANGE, "", true),
	ClassCMarker:          term.NewStyle(BLUE, "", true),
	ClassDMarker:          term.NewStyle(GREEN, "", true),
	ClassEMarker:          term.NewStyle(LIGHT_GREEN, "", true),
	ClassFMarker:          term.NewStyle("#209c05", "", true),
	HeaderStriped:         term.NewStyle(HEADER_COLOR, BG_COLOR_ODD, true),
	TextStriped:           term.NewStyle("#81858d", BG_COLOR_ODD, false),
	PositiveMarkerStriped: term.NewStyle(LIGHT_GREEN, BG_COLOR_ODD, true),
	NegativeMarkerStriped: term.NewStyle(RED, BG_COLOR_ODD, true),
	ClassAMarkerStriped:   term.NewStyle(RED, BG_COLOR_ODD, true),
	ClassBMarkerStriped:   term.NewStyle(ORANGE, BG_COLOR_ODD, true),
	ClassCMarkerStriped:   term.NewStyle(BLUE, BG_COLOR_ODD, true),
	ClassDMarkerStriped:   term.NewStyle(GREEN, BG_COLOR_ODD, true),
	ClassEMarkerStriped:   term.NewStyle(LIGHT_GREEN, BG_COLOR_ODD, true),
	ClassFMarkerStriped:   term.NewStyle("#209c05", BG_COLOR_ODD, true),
	HeaderPositive:        term.NewStyle("#0C0C0C", GREEN, true),
	HeaderNegative:        term.NewStyle("#ffffff", RED, true),
	HeaderClassA:          term.NewStyle("#ffffff", RED, true),
	HeaderClassB:          term.NewStyle("#ffffff", ORANGE, true),
	HeaderClassC:          term.NewStyle("#ffffff", BLUE, true),
	HeaderClassD:          term.NewStyle("#0C0C0C", LIGHT_GREEN, true),
	HeaderClassE:          term.NewStyle("#0C0C0C", GREEN, true),
	HeaderClassF:          term.NewStyle("#0C0C0C", "#209c05", true),
}

var GUV_DARK_STYLE = Styles{
	Text:                  term.NewStyle("#d0d0d0", "", false),
	Header:                term.NewStyle(HEADER_COLOR, "", true),
	PositiveMarker:        term.NewStyle(LIGHT_GREEN, "", true),
	NegativeMarker:        term.NewStyle(RED, "", true),
	ClassAMarker:          term.NewStyle(RED, "", true),
	ClassBMarker:          term.NewStyle(ORANGE, "", true),
	ClassCMarker:          term.NewStyle(BLUE, "", true),
	ClassDMarker:          term.NewStyle(LIGHT_GREEN, "", true),
	ClassEMarker:          term.NewStyle(GREEN, "", true),
	ClassFMarker:          term.NewStyle("#209c05", "", true),
	HeaderStriped:         term.NewStyle(HEADER_COLOR, BG_COLOR_ODD, true),
	TextStriped:           term.NewStyle("#81858d", BG_COLOR_ODD, false),
	PositiveMarkerStriped: term.NewStyle(LIGHT_GREEN, BG_COLOR_ODD, true),
	NegativeMarkerStriped: term.NewStyle(RED, BG_COLOR_ODD, true),
	ClassAMarkerStriped:   term.NewStyle(RED, BG_COLOR_ODD, true),
	ClassBMarkerStriped:   term.NewStyle(ORANGE, BG_COLOR_ODD, true),
	ClassCMarkerStriped:   term.NewStyle(BLUE, BG_COLOR_ODD, true),
	ClassDMarkerStriped:   term.NewStyle(LIGHT_GREEN, BG_COLOR_ODD, true),
	ClassEMarkerStriped:   term.NewStyle(GREEN, BG_COLOR_ODD, true),
	ClassFMarkerStriped:   term.NewStyle("#209c05", BG_COLOR_ODD, true),
	HeaderPositive:        term.NewStyle("#0C0C0C", GREEN, true),
	HeaderNegative:        term.NewStyle("#ffffff", RED, true),
	HeaderClassA:          term.NewStyle("#ffffff", RED, true),
	HeaderClassB:          term.NewStyle("#ffffff", ORANGE, true),
	HeaderClassC:          term.NewStyle("#ffffff", BLUE, true),
	HeaderClassD:          term.NewStyle("#0C0C0C", LIGHT_GREEN, true),
	HeaderClassE:          term.NewStyle("#0C0C0C", GREEN, true),
	HeaderClassF:          term.NewStyle("#0C0C0C", "#209c05", true),
}

var BG_STYLE = Styles{
	Text:                  term.NewStyle("#d0d0d0", "", false),
	Header:                term.NewStyle(HEADER_COLOR, "", true),
	PositiveMarker:        term.NewStyle("#ffffff", LIGHT_GREEN, true),
	NegativeMarker:        term.NewStyle("#ffffff", RED, true),
	ClassAMarker:          term.NewStyle("#ffffff", RED, true),
	ClassBMarker:          term.NewStyle("#ffffff", ORANGE, true),
	ClassCMarker:          term.NewStyle("#ffffff", BLUE, true),
	ClassDMarker:          term.NewStyle("#ffffff", GREEN, true),
	ClassEMarker:          term.NewStyle("#ffffff", LIGHT_GREEN, true),
	ClassFMarker:          term.NewStyle("#209c05", "", true),
	HeaderStriped:         term.NewStyle(HEADER_COLOR, BG_COLOR_ODD, true),
	TextStriped:           term.NewStyle("#81858d", BG_COLOR_ODD, false),
	PositiveMarkerStriped: term.NewStyle("#ffffff", LIGHT_GREEN, true),
	NegativeMarkerStriped: term.NewStyle("#ffffff", RED, true),
	ClassAMarkerStriped:   term.NewStyle("#ffffff", RED, true),
	ClassBMarkerStriped:   term.NewStyle("#ffffff", ORANGE, true),
	ClassCMarkerStriped:   term.NewStyle("#ffffff", BLUE, true),
	ClassDMarkerStriped:   term.NewStyle("#ffffff", GREEN, true),
	ClassEMarkerStriped:   term.NewStyle("#ffffff", LIGHT_GREEN, true),
	ClassFMarkerStriped:   term.NewStyle("#209c05", BG_COLOR_ODD, true),
	HeaderPositive:        term.NewStyle("#0C0C0C", GREEN, true),
	HeaderNegative:        term.NewStyle("#ffffff", RED, true),
	HeaderClassA:          term.NewStyle("#ffffff", RED, true),
	HeaderClassB:          term.NewStyle("#ffffff", ORANGE, true),
	HeaderClassC:          term.NewStyle("#ffffff", BLUE, true),
	HeaderClassD:          term.NewStyle("#0C0C0C", GREEN, true),
	HeaderClassE:          term.NewStyle("#0C0C0C", LIGHT_GREEN, true),
	HeaderClassF:          term.NewStyle("#0C0C0C", "#209c05", true),
}
