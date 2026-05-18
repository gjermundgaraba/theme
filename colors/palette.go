package colors

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
	SurfaceHue:            150,
	AccentHue:             150,
	DarkSurfaceLightness:  0.30,
	LightSurfaceLightness: 0.88,
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
// DefaultConfig so partial configs never produce degenerate palettes.
func (cfg ThemeConfig) Normalized() ThemeConfig {
	if cfg.DarkSurfaceLightness == 0 {
		cfg.DarkSurfaceLightness = DefaultConfig.DarkSurfaceLightness
	}
	if cfg.LightSurfaceLightness == 0 {
		cfg.LightSurfaceLightness = DefaultConfig.LightSurfaceLightness
	}
	return cfg
}
