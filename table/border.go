package table

type Border struct {
	Size      int
	H_LINE    string
	V_LINE    string
	TL_CORNER string
	TR_CORNER string
	BR_CORNER string
	BL_CORNER string
	CROSS     string
	TOP_DEL   string
	BOT_DEL   string
	LEFT_DEL  string
	RIGHT_DEL string
}

var DefaultBorder = Border{
	Size:      1,
	H_LINE:    "│",
	V_LINE:    "─",
	TL_CORNER: "┌",
	TR_CORNER: "┐",
	BR_CORNER: "┘",
	BL_CORNER: "└",
	CROSS:     "┼",
	TOP_DEL:   "┬",
	BOT_DEL:   "┴",
	LEFT_DEL:  "├",
	RIGHT_DEL: "┤",
}

var HiddenBorder = Border{
	Size:      0,
	H_LINE:    "",
	V_LINE:    "─",
	TL_CORNER: "",
	TR_CORNER: "",
	BR_CORNER: "",
	BL_CORNER: "",
	CROSS:     "",
	TOP_DEL:   "",
	BOT_DEL:   "",
	LEFT_DEL:  "",
	RIGHT_DEL: "",
}

var RoundedBorder = Border{
	Size:      1,
	V_LINE:    "─",
	H_LINE:    "│",
	TL_CORNER: "╭",
	TR_CORNER: "╮",
	BL_CORNER: "╰",
	BR_CORNER: "╯",
	LEFT_DEL:  "├",
	RIGHT_DEL: "┤",
	CROSS:     "┼",
	TOP_DEL:   "┬",
	BOT_DEL:   "┴",
}

var ThickBorder = Border{
	Size:      1,
	V_LINE:    "━",
	H_LINE:    "┃",
	TL_CORNER: "┏",
	TR_CORNER: "┓",
	BL_CORNER: "┗",
	BR_CORNER: "┛",
	LEFT_DEL:  "┣",
	RIGHT_DEL: "┫",
	CROSS:     "╋",
	TOP_DEL:   "┳",
	BOT_DEL:   "┻",
}

var DoubleBorder = Border{
	Size:      1,
	V_LINE:    "═",
	H_LINE:    "║",
	TL_CORNER: "╔",
	TR_CORNER: "╗",
	BL_CORNER: "╚",
	BR_CORNER: "╝",
	LEFT_DEL:  "╠",
	RIGHT_DEL: "╣",
	CROSS:     "╬",
	TOP_DEL:   "╦",
	BOT_DEL:   "╩",
}

/*
blockBorder = Border{
	Top:         "█",
	Bottom:      "█",
	Left:        "█",
	Right:       "█",
	TopLeft:     "█",
	TopRight:    "█",
	BottomLeft:  "█",
	BottomRight: "█",
}

outerHalfBlockBorder = Border{
	Top:         "▀",
	Bottom:      "▄",
	Left:        "▌",
	Right:       "▐",
	TopLeft:     "▛",
	TopRight:    "▜",
	BottomLeft:  "▙",
	BottomRight: "▟",
}

innerHalfBlockBorder = Border{
	Top:         "▄",
	Bottom:      "▀",
	Left:        "▐",
	Right:       "▌",
	TopLeft:     "▗",
	TopRight:    "▖",
	BottomLeft:  "▝",
	BottomRight: "▘",
}
*/
