package colors

import colorful "github.com/lucasb-eyer/go-colorful"

// hueStop defines one position on the syntax wheel: an absolute hue
// plus per-hue corrections from the syntax baseline (synL, synC).
// Uniform L,C across hues makes yellow olive, blue dusty, and green
// neon; the per-hue deltas keep each color visually itself while
// preserving a single global vibrance knob (synL/synC).
type hueStop struct {
	H, DL, DC float64
}

// syntaxStops groups the seven canonical syntax hues. Hue identity is
// stable across modes (red is always red); only the L/C corrections
// differ between dark and light to track white-on-paper vs
// light-on-ink readability.
type syntaxStops struct {
	Red, Orange, Yellow, Green, Cyan, Blue, Magenta hueStop
}

// Hue positions are stable across modes (red is always red); only L/C
// corrections differ between dark and light.
// Red sits at 12° and Orange at 58° to keep them clearly distinct after
// sRGB gamut clamping in light mode — at L<0.5 they share an sRGB corner
// and collapse together if the gap is narrower.
var darkSyntaxStops = syntaxStops{
	Red:     hueStop{H: 12, DL: +0.00, DC: -0.025},
	Orange:  hueStop{H: 58, DL: -0.02, DC: -0.005},
	Yellow:  hueStop{H: 92, DL: +0.10, DC: -0.045},
	Green:   hueStop{H: 145, DL: -0.04, DC: -0.025},
	Cyan:    hueStop{H: 195, DL: -0.03, DC: -0.020},
	Blue:    hueStop{H: 250, DL: -0.08, DC: +0.015},
	Magenta: hueStop{H: 325, DL: +0.00, DC: -0.030},
}

var lightSyntaxStops = syntaxStops{
	Red:     hueStop{H: 12, DL: -0.02, DC: -0.020},
	Orange:  hueStop{H: 58, DL: +0.00, DC: +0.005},
	Yellow:  hueStop{H: 92, DL: -0.06, DC: -0.020},
	Green:   hueStop{H: 145, DL: -0.02, DC: -0.005},
	Cyan:    hueStop{H: 195, DL: -0.04, DC: +0.000},
	Blue:    hueStop{H: 250, DL: +0.02, DC: +0.020},
	Magenta: hueStop{H: 325, DL: -0.01, DC: -0.005},
}

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

	// Surface chroma sits at whisper level. The brand hue is present
	// across every layer but never announces itself — beloved dark
	// themes (Tokyo Night, Catppuccin, Rosé Pine, Solarized) all keep
	// their BG near-neutral and let identity come from syntax and
	// accent. A clearly chromatic field competes with the syntax
	// wheel and produces "themed soup" instead of a stage.
	bgL, bgC := cfg.DarkSurfaceLightness, 0.006
	surfL, surfC := bgL+0.07, 0.012
	overL, overC := bgL+0.22, 0.018
	fgDimL, fgDimC := 0.64, 0.018
	fgL, fgC := 0.92, 0.008

	// Syntax baseline. The per-hue stop table corrects from here.
	// Bright lift is additive on top of the corrected per-hue values
	// so bright yellow ends up the brightest and bright blue the
	// deepest — matching how each hue actually peaks in OKLCH.
	synL, synC := 0.76, 0.205
	const brightLiftL, brightLiftC = 0.07, 0.015

	syn := darkSyntaxStops
	synColor := func(s hueStop) string {
		return oklchHex(synL+s.DL, synC+s.DC, s.H)
	}
	brColor := func(s hueStop) string {
		return oklchHex(synL+s.DL+brightLiftL, synC+s.DC+brightLiftC, s.H)
	}

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
		SelectionBG:  oklchHex(surfL+0.06, 0.060, cfg.AccentHue),
		SelectionFG:  oklchHex(fgL, fgC, h),
		Accent:       accent,
		ButtonBG:     buttonBG,
		ButtonFG:     buttonFG,
		DebugFG:      oklchHex(0.17, 0.016, h),
		DiffAddBG:    oklchHex(bgL+0.045, 0.065, syn.Green.H),
		DiffDeleteBG: oklchHex(bgL+0.045, 0.065, syn.Red.H),
		DiffChangeBG: oklchHex(bgL+0.045, 0.065, syn.Blue.H),
		SearchBG:     oklchHex(bgL+0.075, 0.075, syn.Yellow.H),
		VisualBG:     oklchHex(surfL+0.06, 0.055, h),

		Black:         oklchHex(0.145, 0.016, h),
		Red:           synColor(syn.Red),
		Green:         synColor(syn.Green),
		Yellow:        synColor(syn.Yellow),
		Blue:          synColor(syn.Blue),
		Magenta:       synColor(syn.Magenta),
		Cyan:          synColor(syn.Cyan),
		White:         oklchHex(0.96, 0.007, h),
		Orange:        synColor(syn.Orange),
		BrightBlack:   oklchHex(fgDimL, fgDimC, h),
		BrightRed:     brColor(syn.Red),
		BrightGreen:   brColor(syn.Green),
		BrightYellow:  brColor(syn.Yellow),
		BrightBlue:    brColor(syn.Blue),
		BrightMagenta: brColor(syn.Magenta),
		BrightCyan:    brColor(syn.Cyan),
		BrightWhite:   oklchHex(0.985, 0.006, h),

		SyntaxKeyword:  synColor(syn.Magenta),
		SyntaxString:   synColor(syn.Yellow),
		SyntaxNumber:   synColor(syn.Blue),
		SyntaxComment:  oklchHex(fgDimL, fgDimC, h),
		SyntaxConstant: synColor(syn.Cyan),
		SyntaxFunction: brColor(syn.Cyan),
		SyntaxBuiltin:  synColor(syn.Green),
		SyntaxLink:     brColor(syn.Blue),
		SyntaxError:    synColor(syn.Red),
	}
}

func generateLight(cfg ThemeConfig) Palette {
	h := cfg.SurfaceHue

	bgL, bgC := cfg.LightSurfaceLightness, 0.005
	surfL, surfC := bgL-0.04, 0.010
	overL, overC := bgL-0.45, 0.015
	fgDimL, fgDimC := 0.38, 0.015
	fgL, fgC := 0.18, 0.012

	// Light-mode bright lift goes DARKER, not lighter — on paper,
	// "bright" means more saturated and deeper-set, not whiter.
	synL, synC := 0.44, 0.195
	const brightLiftL, brightLiftC = -0.08, 0.020

	syn := lightSyntaxStops
	synColor := func(s hueStop) string {
		return oklchHex(synL+s.DL, synC+s.DC, s.H)
	}
	brColor := func(s hueStop) string {
		return oklchHex(synL+s.DL+brightLiftL, synC+s.DC+brightLiftC, s.H)
	}

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
		SelectionBG:  oklchHex(surfL-0.09, 0.060, cfg.AccentHue),
		SelectionFG:  oklchHex(fgL, fgC, h),
		Accent:       accent,
		ButtonBG:     buttonBG,
		ButtonFG:     buttonFG,
		DebugFG:      oklchHex(0.98, 0.006, h),
		DiffAddBG:    oklchHex(bgL-0.03, 0.04, syn.Green.H),
		DiffDeleteBG: oklchHex(bgL-0.03, 0.04, syn.Red.H),
		DiffChangeBG: oklchHex(bgL-0.03, 0.04, syn.Blue.H),
		SearchBG:     oklchHex(bgL-0.05, 0.055, syn.Yellow.H),
		VisualBG:     oklchHex(surfL-0.09, 0.055, h),

		Black:         oklchHex(0.10, 0.01, h),
		Red:           synColor(syn.Red),
		Green:         synColor(syn.Green),
		Yellow:        synColor(syn.Yellow),
		Blue:          synColor(syn.Blue),
		Magenta:       synColor(syn.Magenta),
		Cyan:          synColor(syn.Cyan),
		White:         oklchHex(0.95, 0.005, h),
		Orange:        synColor(syn.Orange),
		BrightBlack:   oklchHex(fgDimL, fgDimC, h),
		BrightRed:     brColor(syn.Red),
		BrightGreen:   brColor(syn.Green),
		BrightYellow:  brColor(syn.Yellow),
		BrightBlue:    brColor(syn.Blue),
		BrightMagenta: brColor(syn.Magenta),
		BrightCyan:    brColor(syn.Cyan),
		BrightWhite:   oklchHex(0.98, 0.005, h),

		SyntaxKeyword:  synColor(syn.Magenta),
		SyntaxString:   synColor(syn.Yellow),
		SyntaxNumber:   synColor(syn.Blue),
		SyntaxComment:  oklchHex(fgDimL, fgDimC, h),
		SyntaxConstant: synColor(syn.Cyan),
		SyntaxFunction: brColor(syn.Cyan),
		SyntaxBuiltin:  synColor(syn.Green),
		SyntaxLink:     brColor(syn.Blue),
		SyntaxError:    synColor(syn.Red),
	}
}
