package table

import "github.com/muesli/termenv"

type StyleFn func(string) string

type Styles struct {
	Text                  StyleFn
	Header                StyleFn
	PositiveMarker        StyleFn
	NegativeMarker        StyleFn
	ClassAMarker          StyleFn
	ClassBMarker          StyleFn
	ClassCMarker          StyleFn
	ClassDMarker          StyleFn
	ClassEMarker          StyleFn
	ClassFMarker          StyleFn
	HeaderStriped         StyleFn
	TextStriped           StyleFn
	PositiveMarkerStriped StyleFn
	NegativeMarkerStriped StyleFn
	ClassAMarkerStriped   StyleFn
	ClassBMarkerStriped   StyleFn
	ClassCMarkerStriped   StyleFn
	ClassDMarkerStriped   StyleFn
	ClassEMarkerStriped   StyleFn
	ClassFMarkerStriped   StyleFn
}

func NewStyle(fg string, bg string, bold bool) func(string) string {
	p := termenv.ColorProfile()
	s := termenv.Style{}.Foreground(p.Color(fg)).Background(p.Color(bg))
	if bold {
		s = s.Bold()
	}
	return s.Styled
}
