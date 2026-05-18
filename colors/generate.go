package colors

import colorful "github.com/lucasb-eyer/go-colorful"

// Generate builds a Palette for the given config. Pass dark=true for dark mode.
// Partial configs are normalized so zero-valued lightness never produces a
// degenerate palette.
func Generate(cfg ThemeConfig, dark bool) Palette {
	cfg = cfg.Normalized()
	if dark {
		return generateDark(cfg)
	}
	return generateLight(cfg)
}

func oklchHex(l, c, h float64) string {
	return colorful.OkLch(l, c, h).Clamped().Hex()
}

func contrastText(bg, light, dark string) string {
	if ContrastRatio(light, bg) >= ContrastRatio(dark, bg) {
		return light
	}
	return dark
}

func generateDark(cfg ThemeConfig) Palette {
	h := cfg.SurfaceHue

	bgL, bgC := cfg.DarkSurfaceLightness, 0.026
	surfL, surfC := bgL+0.075, 0.030
	overL, overC := bgL+0.28, 0.040
	fgDimL, fgDimC := 0.64, 0.026
	fgL, fgC := 0.92, 0.008

	synL, synC := 0.76, 0.205
	brL, brC := 0.83, 0.220

	redH := 20.0
	orgH := 45.0
	yelH := 90.0
	grnH := 150.0
	cynH := 190.0
	bluH := 255.0
	magH := 325.0

	accent := oklchHex(0.76, 0.215, cfg.AccentHue)
	buttonBG := oklchHex(0.47, 0.235, cfg.AccentHue)
	buttonFG := contrastText(buttonBG, oklchHex(0.985, 0.006, h), oklchHex(0.12, 0.010, h))

	return Palette{
		BG:           oklchHex(bgL, bgC, h),
		Surface:      oklchHex(surfL, surfC, h),
		Overlay:      oklchHex(overL, overC, h),
		FG:           oklchHex(fgL, fgC, h),
		FGDim:        oklchHex(fgDimL, fgDimC, h),
		Cursor:       oklchHex(fgL, fgC, h),
		CursorText:   oklchHex(bgL, bgC, h),
		SelectionBG:  oklchHex(surfL+0.06, surfC+0.025, h),
		SelectionFG:  oklchHex(fgL, fgC, h),
		Accent:       accent,
		ButtonBG:     buttonBG,
		ButtonFG:     buttonFG,
		DebugFG:      oklchHex(0.17, 0.016, h),
		DiffAddBG:    oklchHex(bgL+0.045, 0.065, grnH),
		DiffDeleteBG: oklchHex(bgL+0.045, 0.065, redH),
		DiffChangeBG: oklchHex(bgL+0.045, 0.065, bluH),
		SearchBG:     oklchHex(bgL+0.075, 0.075, yelH),
		VisualBG:     oklchHex(surfL+0.06, surfC+0.025, h),

		Black:         oklchHex(0.145, 0.016, h),
		Red:           oklchHex(synL, synC, redH),
		Green:         oklchHex(synL, synC, grnH),
		Yellow:        oklchHex(synL, synC, yelH),
		Blue:          oklchHex(synL, synC, bluH),
		Magenta:       oklchHex(synL, synC, magH),
		Cyan:          oklchHex(synL, synC, cynH),
		White:         oklchHex(0.96, 0.007, h),
		Orange:        oklchHex(synL, synC, orgH),
		BrightBlack:   oklchHex(fgDimL, fgDimC, h),
		BrightRed:     oklchHex(brL, brC, redH),
		BrightGreen:   oklchHex(brL, brC, grnH),
		BrightYellow:  oklchHex(brL, brC, yelH),
		BrightBlue:    oklchHex(brL, brC, bluH),
		BrightMagenta: oklchHex(brL, brC, magH),
		BrightCyan:    oklchHex(brL, brC, cynH),
		BrightWhite:   oklchHex(0.985, 0.006, h),

		SyntaxKeyword:  oklchHex(synL, synC, magH),
		SyntaxString:   oklchHex(synL, synC, yelH),
		SyntaxNumber:   oklchHex(synL, synC, bluH),
		SyntaxComment:  oklchHex(fgDimL, fgDimC, h),
		SyntaxConstant: oklchHex(synL, synC, cynH),
		SyntaxFunction: oklchHex(brL, brC, cynH),
		SyntaxBuiltin:  oklchHex(synL, synC, grnH),
		SyntaxLink:     oklchHex(brL, brC, bluH),
		SyntaxError:    oklchHex(synL, synC, redH),
	}
}

func generateLight(cfg ThemeConfig) Palette {
	h := cfg.SurfaceHue

	bgL, bgC := cfg.LightSurfaceLightness, 0.014
	surfL, surfC := bgL-0.04, 0.018
	overL, overC := bgL-0.45, 0.026
	fgDimL, fgDimC := 0.38, 0.020
	fgL, fgC := 0.18, 0.012

	synL, synC := 0.44, 0.195
	brL, brC := 0.36, 0.215

	redH := 8.0
	orgH := 65.0
	yelH := 95.0
	grnH := 150.0
	cynH := 190.0
	bluH := 255.0
	magH := 325.0

	accent := oklchHex(0.42, 0.195, cfg.AccentHue)
	buttonBG := oklchHex(0.44, 0.25, cfg.AccentHue)
	buttonFG := contrastText(buttonBG, oklchHex(0.98, 0.006, h), oklchHex(fgL, fgC, h))

	return Palette{
		BG:           oklchHex(bgL, bgC, h),
		Surface:      oklchHex(surfL, surfC, h),
		Overlay:      oklchHex(overL, overC, h),
		FG:           oklchHex(fgL, fgC, h),
		FGDim:        oklchHex(fgDimL, fgDimC, h),
		Cursor:       oklchHex(fgL, fgC, h),
		CursorText:   oklchHex(bgL, bgC, h),
		SelectionBG:  oklchHex(surfL-0.09, surfC+0.026, h),
		SelectionFG:  oklchHex(fgL, fgC, h),
		Accent:       accent,
		ButtonBG:     buttonBG,
		ButtonFG:     buttonFG,
		DebugFG:      oklchHex(0.98, 0.006, h),
		DiffAddBG:    oklchHex(bgL-0.03, 0.04, grnH),
		DiffDeleteBG: oklchHex(bgL-0.03, 0.04, redH),
		DiffChangeBG: oklchHex(bgL-0.03, 0.04, bluH),
		SearchBG:     oklchHex(bgL-0.05, 0.055, yelH),
		VisualBG:     oklchHex(surfL-0.09, surfC+0.026, h),

		Black:         oklchHex(0.10, 0.01, h),
		Red:           oklchHex(synL, synC, redH),
		Green:         oklchHex(synL, synC, grnH),
		Yellow:        oklchHex(synL, synC, yelH),
		Blue:          oklchHex(synL, synC, bluH),
		Magenta:       oklchHex(synL, synC, magH),
		Cyan:          oklchHex(synL, synC, cynH),
		White:         oklchHex(0.95, 0.005, h),
		Orange:        oklchHex(synL, synC, orgH),
		BrightBlack:   oklchHex(fgDimL, fgDimC, h),
		BrightRed:     oklchHex(brL, brC, redH),
		BrightGreen:   oklchHex(brL, brC, grnH),
		BrightYellow:  oklchHex(brL, brC, yelH),
		BrightBlue:    oklchHex(brL, brC, bluH),
		BrightMagenta: oklchHex(brL, brC, magH),
		BrightCyan:    oklchHex(brL, brC, cynH),
		BrightWhite:   oklchHex(0.98, 0.005, h),

		SyntaxKeyword:  oklchHex(synL, synC, magH),
		SyntaxString:   oklchHex(synL, synC, yelH),
		SyntaxNumber:   oklchHex(synL, synC, bluH),
		SyntaxComment:  oklchHex(fgDimL, fgDimC, h),
		SyntaxConstant: oklchHex(synL, synC, cynH),
		SyntaxFunction: oklchHex(brL, brC, cynH),
		SyntaxBuiltin:  oklchHex(synL, synC, grnH),
		SyntaxLink:     oklchHex(brL, brC, bluH),
		SyntaxError:    oklchHex(synL, synC, redH),
	}
}
