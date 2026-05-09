package colors

// PaletteToken describes one generated palette value for review artifacts.
type PaletteToken struct {
	Name  string `json:"name"`
	Hex   string `json:"hex"`
	OKLCH string `json:"oklch"`
	Role  string `json:"role"`
}

// PaletteTokens returns the palette in stable review order with role intent.
func PaletteTokens(p Palette) []PaletteToken {
	fields := []struct {
		name string
		hex  string
		role string
	}{
		{"BG", p.BG, "primary terminal, editor, Neovim, and document field"},
		{"Surface", p.Surface, "raised panels, sidebars, widgets, tabs, menus, and grouped surfaces"},
		{"Overlay", p.Overlay, "borders, dividers, guides, line numbers, ghost text, and inactive structure"},
		{"FG", p.FG, "primary foreground text"},
		{"FGDim", p.FGDim, "secondary text, comments, metadata, and subdued labels"},
		{"Cursor", p.Cursor, "cursor fill"},
		{"CursorText", p.CursorText, "text under cursor"},
		{"SelectionBG", p.SelectionBG, "selection and focused list background"},
		{"SelectionFG", p.SelectionFG, "text on selection background"},
		{"Accent", p.Accent, "accent text, links, and active emphasis"},
		{"ButtonBG", p.ButtonBG, "filled action, badge, progress, remote, and selected-control backgrounds"},
		{"ButtonFG", p.ButtonFG, "text on filled action backgrounds"},
		{"DebugFG", p.DebugFG, "readable foreground on saturated debug/progress/status fills"},
		{"DiffAddBG", p.DiffAddBG, "added diff background"},
		{"DiffDeleteBG", p.DiffDeleteBG, "deleted diff background"},
		{"DiffChangeBG", p.DiffChangeBG, "changed diff background"},
		{"SearchBG", p.SearchBG, "search match background"},
		{"VisualBG", p.VisualBG, "visual mode and secondary selection background"},
		{"Black", p.Black, "ANSI black slot, conventional terminal extreme"},
		{"Red", p.Red, "errors, destructive state, regexp, deleted git state, ANSI red"},
		{"Green", p.Green, "success, functions, added git state, ANSI green"},
		{"Yellow", p.Yellow, "strings, quotes, host markers, ANSI yellow"},
		{"Blue", p.Blue, "numbers, focus indicators, active structure, ANSI blue"},
		{"Magenta", p.Magenta, "keywords, operators, tags, ANSI magenta"},
		{"Cyan", p.Cyan, "types, links, commands, info state, ANSI cyan"},
		{"White", p.White, "ANSI white slot, conventional terminal extreme"},
		{"Orange", p.Orange, "warnings, parameters, search current match, changed git state"},
		{"BrightBlack", p.BrightBlack, "bright ANSI black and subdued terminal emphasis"},
		{"BrightRed", p.BrightRed, "bright ANSI red"},
		{"BrightGreen", p.BrightGreen, "bright ANSI green"},
		{"BrightYellow", p.BrightYellow, "bright ANSI yellow"},
		{"BrightBlue", p.BrightBlue, "bright ANSI blue"},
		{"BrightMagenta", p.BrightMagenta, "bright ANSI magenta"},
		{"BrightCyan", p.BrightCyan, "bright ANSI cyan"},
		{"BrightWhite", p.BrightWhite, "bright ANSI white slot, conventional terminal extreme"},
	}

	tokens := make([]PaletteToken, len(fields))
	for i, f := range fields {
		tokens[i] = PaletteToken{Name: f.name, Hex: f.hex, OKLCH: OKLCHString(f.hex), Role: f.role}
	}
	return tokens
}
