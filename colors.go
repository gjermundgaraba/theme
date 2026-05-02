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

var DefaultConfig = ThemeConfig{
	SurfaceHue: 157,
	AccentHue:  360,
}

var Dark = Generate(DefaultConfig, true)
var Light = Generate(DefaultConfig, false)

func SetConfig(cfg ThemeConfig) {
	DefaultConfig = cfg
	Dark = Generate(cfg, true)
	Light = Generate(cfg, false)
}
