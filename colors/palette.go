package colors

import "math"

// Palette contains every generated color token.
type Palette struct {
	BG, Surface, Overlay     string
	FG, FGDim                string
	Cursor, CursorText       string
	SelectionBG, SelectionFG string
	Accent                   string
	ButtonBG, ButtonFG       string
	DebugFG                  string
	DiffAddBG, DiffDeleteBG  string
	DiffChangeBG, SearchBG   string
	VisualBG                 string

	Black, Red, Green, Yellow           string
	Blue, Magenta, Cyan, White          string
	Orange                              string
	BrightBlack, BrightRed, BrightGreen string
	BrightYellow, BrightBlue            string
	BrightMagenta, BrightCyan           string
	BrightWhite                         string

	// Semantic syntax tokens. Decouple editor-theme roles from the raw
	// ANSI hues so templates name the role, not the color. Guaranteed
	// >=3:1 contrast against Surface in both light and dark modes
	// (enforced by ContrastChecks).
	SyntaxKeyword, SyntaxString, SyntaxNumber  string
	SyntaxComment, SyntaxConstant              string
	SyntaxFunction, SyntaxBuiltin              string
	SyntaxLink, SyntaxError                    string
}

// ThemeConfig controls palette generation.
type ThemeConfig struct {
	SurfaceHue            float64
	AccentHue             float64
	DarkSurfaceLightness  float64
	LightSurfaceLightness float64
}

// Lightness clamps for the dark and light surface fields. The dark upper
// bound is the empirical limit where FGDim-on-Surface still meets WCAG AA;
// going higher compresses Surface against FGDim.
const (
	MinDarkSurfaceLightness  = 0.12
	MaxDarkSurfaceLightness  = 0.30
	MinLightSurfaceLightness = 0.82
	MaxLightSurfaceLightness = 1.00
)

// DefaultConfig is the project's canonical hue and lightness configuration.
var DefaultConfig = ThemeConfig{
	SurfaceHue:            135,
	AccentHue:             16,
	DarkSurfaceLightness:  0.28,
	LightSurfaceLightness: 0.91,
}

// Dark and Light are pre-built palettes from DefaultConfig.
var (
	Dark  = Generate(DefaultConfig, true)
	Light = Generate(DefaultConfig, false)
)

// SetConfig updates DefaultConfig and rebuilds Dark/Light. Zero-valued
// lightness fields are filled from the existing DefaultConfig before storing.
func SetConfig(cfg ThemeConfig) {
	cfg = cfg.Normalized()
	DefaultConfig = cfg
	Dark = Generate(cfg, true)
	Light = Generate(cfg, false)
}

// Normalized returns cfg with zero-valued lightness fields filled from
// DefaultConfig and AccentHue nudged out of the forbidden zones so partial
// or colliding inputs never produce a degenerate palette.
func (cfg ThemeConfig) Normalized() ThemeConfig {
	if cfg.DarkSurfaceLightness == 0 {
		cfg.DarkSurfaceLightness = DefaultConfig.DarkSurfaceLightness
	}
	if cfg.LightSurfaceLightness == 0 {
		cfg.LightSurfaceLightness = DefaultConfig.LightSurfaceLightness
	}
	cfg.SurfaceHue = wrapHue(cfg.SurfaceHue)
	cfg.AccentHue = nudgeAccentHue(wrapHue(cfg.AccentHue), cfg.SurfaceHue)
	return cfg
}

// minAccentSeparation is the minimum hue distance (deg) between Accent and
// SurfaceHue / Red / Green / Magenta. Below this, the accent dissolves into
// the field or aliases as error/success/keyword state.
const minAccentSeparation = 25.0

// accentForbiddenHues are the syntax hues the accent must stay clear of.
// Red, Green, Magenta carry strong semantic signal in generated targets and
// must remain distinguishable from the one accent token.
var accentForbiddenHues = [...]float64{20, 150, 325}

func wrapHue(h float64) float64 {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return h
}

func hueDistance(a, b float64) float64 {
	d := math.Abs(wrapHue(a) - wrapHue(b))
	if d > 180 {
		d = 360 - d
	}
	return d
}

func accentHueValid(accent, surface float64) bool {
	if hueDistance(accent, surface) < minAccentSeparation {
		return false
	}
	for _, f := range accentForbiddenHues {
		if hueDistance(accent, f) < minAccentSeparation {
			return false
		}
	}
	return true
}

// nudgeAccentHue returns accent if it is already a valid distance from
// SurfaceHue and the forbidden syntax hues; otherwise it walks outward in 1°
// increments and returns the closest valid hue. Forbidden zones never fill
// 360°, so the loop always terminates.
func nudgeAccentHue(accent, surface float64) float64 {
	if accentHueValid(accent, surface) {
		return accent
	}
	for delta := 1.0; delta <= 180.0; delta++ {
		if c := wrapHue(accent + delta); accentHueValid(c, surface) {
			return c
		}
		if c := wrapHue(accent - delta); accentHueValid(c, surface) {
			return c
		}
	}
	return accent
}
