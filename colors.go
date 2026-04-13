package main

type Palette struct {
	BG, Surface, Overlay     string
	FG, FGDim                string
	Cursor, CursorText       string
	SelectionBG, SelectionFG string
	ButtonBG, ButtonFG       string
	DebugFG                  string
	DiffAddBG, DiffDeleteBG  string
	DiffChangeBG, SearchBG   string
	VisualBG                 string

	Black, Red, Green, Yellow          string
	Blue, Magenta, Cyan, White         string
	Orange                             string
	BrightBlack, BrightRed, BrightGreen string
	BrightYellow, BrightBlue           string
	BrightMagenta, BrightCyan          string
	BrightWhite                        string
}

var Dark = Palette{
	BG:            "#282a36",
	Surface:       "#44475a",
	Overlay:       "#6272a4",
	FG:            "#f8f8f2",
	FGDim:         "#8894c4",
	Cursor:        "#f8f8f2",
	CursorText:    "#282a36",
	SelectionBG:   "#44475a",
	SelectionFG:   "#f8f8f2",
	ButtonBG:      "#7260b5",
	ButtonFG:      "#f8f8f2",
	DebugFG:       "#21222c",
	DiffAddBG:     "#2a3a2e",
	DiffDeleteBG:  "#3a2a2c",
	DiffChangeBG:  "#2a2d3a",
	SearchBG:      "#3a3520",
	VisualBG:      "#44475a",
	Black:         "#21222c",
	Red:           "#ff5555",
	Green:         "#50fa7b",
	Yellow:        "#f1fa8c",
	Blue:          "#bd93f9",
	Magenta:       "#ff79c6",
	Cyan:          "#8be9fd",
	White:         "#f8f8f2",
	Orange:        "#ffb86c",
	BrightBlack:   "#8894c4",
	BrightRed:     "#ff6e6e",
	BrightGreen:   "#69ff94",
	BrightYellow:  "#ffffa5",
	BrightBlue:    "#d6acff",
	BrightMagenta: "#ff92df",
	BrightCyan:    "#a4ffff",
	BrightWhite:   "#ffffff",
}

var Light = Palette{
	BG:            "#d4d6e2",
	Surface:       "#c8cbda",
	Overlay:       "#5b5c6b",
	FG:            "#282a36",
	FGDim:         "#49557b",
	Cursor:        "#282a36",
	CursorText:    "#d4d6e2",
	SelectionBG:   "#aeb2cb",
	SelectionFG:   "#282a36",
	ButtonBG:      "#6e25d4",
	ButtonFG:      "#ffffff",
	DebugFG:       "#ffffff",
	DiffAddBG:     "#c4dcc8",
	DiffDeleteBG:  "#dcc4c6",
	DiffChangeBG:  "#c4c8dc",
	SearchBG:      "#dcd8b8",
	VisualBG:      "#aeb2cb",
	Black:         "#21222c",
	Red:           "#b10000",
	Green:         "#006519",
	Yellow:        "#6e5200",
	Blue:          "#6e25d4",
	Magenta:       "#aa0060",
	Cyan:          "#005e74",
	White:         "#f8f8f2",
	Orange:        "#864500",
	BrightBlack:   "#49557b",
	BrightRed:     "#be0000",
	BrightGreen:   "#006c1b",
	BrightYellow:  "#765900",
	BrightBlue:    "#762fda",
	BrightMagenta: "#b60067",
	BrightCyan:    "#00667d",
	BrightWhite:   "#ffffff",
}
