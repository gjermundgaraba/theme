package main

import colorful "github.com/lucasb-eyer/go-colorful"

type ThemeConfig struct {
	SurfaceHue float64
	AccentHue  float64
}

func Generate(cfg ThemeConfig, dark bool) Palette {
	if dark {
		return generateDark(cfg)
	}
	return generateLight(cfg)
}

func oklchHex(l, c, h float64) string {
	return colorful.OkLch(l, c, h).Clamped().Hex()
}

func generateDark(cfg ThemeConfig) Palette {
	h := cfg.SurfaceHue

	bgL, bgC := 0.13, 0.015
	surfL, surfC := 0.19, 0.018
	overL, overC := 0.38, 0.025
	fgDimL, fgDimC := 0.58, 0.015
	fgL, fgC := 0.90, 0.005

	synL, synC := 0.72, 0.15
	brL, brC := 0.80, 0.17

	redH := 20.0
	orgH := 45.0
	yelH := 90.0
	grnH := 150.0
	cynH := 190.0
	bluH := 255.0
	magH := 325.0

	return Palette{
		BG:            oklchHex(bgL, bgC, h),
		Surface:       oklchHex(surfL, surfC, h),
		Overlay:       oklchHex(overL, overC, h),
		FG:            oklchHex(fgL, fgC, h),
		FGDim:         oklchHex(fgDimL, fgDimC, h),
		Cursor:        oklchHex(fgL, fgC, h),
		CursorText:    oklchHex(bgL, bgC, h),
		SelectionBG:   oklchHex(surfL+0.05, surfC+0.01, h),
		SelectionFG:   oklchHex(fgL, fgC, h),
		ButtonBG:      oklchHex(0.50, 0.18, cfg.AccentHue),
		ButtonFG:      oklchHex(0.95, 0.005, h),
		DebugFG:       oklchHex(0.13, 0.01, h),
		DiffAddBG:     oklchHex(bgL+0.03, 0.03, grnH),
		DiffDeleteBG:  oklchHex(bgL+0.03, 0.03, redH),
		DiffChangeBG:  oklchHex(bgL+0.03, 0.03, bluH),
		SearchBG:      oklchHex(bgL+0.05, 0.04, yelH),
		VisualBG:      oklchHex(surfL+0.05, surfC+0.01, h),

		Black:       oklchHex(0.10, 0.01, h),
		Red:         oklchHex(synL, synC, redH),
		Green:       oklchHex(synL, synC, grnH),
		Yellow:      oklchHex(synL, synC, yelH),
		Blue:        oklchHex(synL, synC, bluH),
		Magenta:     oklchHex(synL, synC, magH),
		Cyan:        oklchHex(synL, synC, cynH),
		White:       oklchHex(0.95, 0.005, h),
		Orange:      oklchHex(synL, synC, orgH),
		BrightBlack: oklchHex(fgDimL, fgDimC, h),
		BrightRed:   oklchHex(brL, brC, redH),
		BrightGreen: oklchHex(brL, brC, grnH),
		BrightYellow:  oklchHex(brL, brC, yelH),
		BrightBlue:    oklchHex(brL, brC, bluH),
		BrightMagenta: oklchHex(brL, brC, magH),
		BrightCyan:    oklchHex(brL, brC, cynH),
		BrightWhite:   oklchHex(0.98, 0.005, h),
	}
}

func generateLight(cfg ThemeConfig) Palette {
	h := cfg.SurfaceHue

	bgL, bgC := 0.95, 0.012
	surfL, surfC := 0.91, 0.015
	overL, overC := 0.50, 0.020
	fgDimL, fgDimC := 0.38, 0.015
	fgL, fgC := 0.18, 0.010

	synL, synC := 0.45, 0.16
	brL, brC := 0.38, 0.18

	redH := 15.0
	orgH := 55.0
	yelH := 90.0
	grnH := 150.0
	cynH := 190.0
	bluH := 255.0
	magH := 325.0

	return Palette{
		BG:            oklchHex(bgL, bgC, h),
		Surface:       oklchHex(surfL, surfC, h),
		Overlay:       oklchHex(overL, overC, h),
		FG:            oklchHex(fgL, fgC, h),
		FGDim:         oklchHex(fgDimL, fgDimC, h),
		Cursor:        oklchHex(fgL, fgC, h),
		CursorText:    oklchHex(bgL, bgC, h),
		SelectionBG:   oklchHex(surfL-0.09, surfC+0.02, h),
		SelectionFG:   oklchHex(fgL, fgC, h),
		ButtonBG:      oklchHex(0.45, 0.22, cfg.AccentHue),
		ButtonFG:      oklchHex(0.98, 0.005, h),
		DebugFG:       oklchHex(0.98, 0.005, h),
		DiffAddBG:     oklchHex(bgL-0.03, 0.03, grnH),
		DiffDeleteBG:  oklchHex(bgL-0.03, 0.03, redH),
		DiffChangeBG:  oklchHex(bgL-0.03, 0.03, bluH),
		SearchBG:      oklchHex(bgL-0.05, 0.04, yelH),
		VisualBG:      oklchHex(surfL-0.09, surfC+0.02, h),

		Black:       oklchHex(0.10, 0.01, h),
		Red:         oklchHex(synL, synC, redH),
		Green:       oklchHex(synL, synC, grnH),
		Yellow:      oklchHex(synL, synC, yelH),
		Blue:        oklchHex(synL, synC, bluH),
		Magenta:     oklchHex(synL, synC, magH),
		Cyan:        oklchHex(synL, synC, cynH),
		White:       oklchHex(0.95, 0.005, h),
		Orange:      oklchHex(synL, synC, orgH),
		BrightBlack: oklchHex(fgDimL, fgDimC, h),
		BrightRed:   oklchHex(brL, brC, redH),
		BrightGreen: oklchHex(brL, brC, grnH),
		BrightYellow:  oklchHex(brL, brC, yelH),
		BrightBlue:    oklchHex(brL, brC, bluH),
		BrightMagenta: oklchHex(brL, brC, magH),
		BrightCyan:    oklchHex(brL, brC, cynH),
		BrightWhite:   oklchHex(0.98, 0.005, h),
	}
}
