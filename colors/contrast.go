package colors

import (
	"fmt"
	"math"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func sRGBToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func luminance(hex string) float64 {
	c, _ := colorful.Hex(hex)
	r, g, b := c.RGB255()
	return 0.2126729*sRGBToLinear(float64(r)/255) +
		0.7151522*sRGBToLinear(float64(g)/255) +
		0.0721750*sRGBToLinear(float64(b)/255)
}

// ContrastRatio computes the WCAG contrast ratio between two hex colors.
func ContrastRatio(hex1, hex2 string) float64 {
	l1 := luminance(hex1)
	l2 := luminance(hex2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// ContrastCheck describes a foreground/background contrast pair.
type ContrastCheck struct {
	Label   string
	Text    string
	BG      string
	Minimum float64
	Ratio   float64
	Pass    bool
}

// ContrastChecks returns the canonical palette contrast checks.
func ContrastChecks(p Palette) []ContrastCheck {
	pairs := []struct {
		label, text, bg string
		minRatio        float64
	}{
		{"FG on BG", p.FG, p.BG, 7.0},
		{"FGDim on BG", p.FGDim, p.BG, 3.0},
		{"Cursor on BG", p.Cursor, p.BG, 4.5},
		{"SelectionFG on SelectionBG", p.SelectionFG, p.SelectionBG, 4.5},
		{"Accent on BG", p.Accent, p.BG, 4.5},
		{"ButtonFG on ButtonBG", p.ButtonFG, p.ButtonBG, 4.5},
		{"FG on Surface", p.FG, p.Surface, 7.0},
		{"FGDim on Surface", p.FGDim, p.Surface, 3.0},
		{"Overlay on BG", p.Overlay, p.BG, 1.5},
		{"FG on SearchBG", p.FG, p.SearchBG, 4.5},
		{"BG on Orange search/current match", p.BG, p.Orange, 4.5},
		{"FG on DiffAddBG", p.FG, p.DiffAddBG, 4.5},
		{"FG on DiffDeleteBG", p.FG, p.DiffDeleteBG, 4.5},
		{"FG on DiffChangeBG", p.FG, p.DiffChangeBG, 4.5},
		{"DebugFG on Orange debug/status", p.DebugFG, p.Orange, 4.5},
		{"DebugFG on Cyan pager progress", p.DebugFG, p.Cyan, 4.5},
		{"Red on BG", p.Red, p.BG, 3.0},
		{"Green on BG", p.Green, p.BG, 3.0},
		{"Yellow on BG", p.Yellow, p.BG, 3.0},
		{"Blue on BG", p.Blue, p.BG, 3.0},
		{"Magenta on BG", p.Magenta, p.BG, 3.0},
		{"Cyan on BG", p.Cyan, p.BG, 3.0},
		{"Orange on BG", p.Orange, p.BG, 3.0},
		{"BrightBlack on BG", p.BrightBlack, p.BG, 3.0},
		{"BrightRed on BG", p.BrightRed, p.BG, 3.0},
		{"BrightGreen on BG", p.BrightGreen, p.BG, 3.0},
		{"BrightYellow on BG", p.BrightYellow, p.BG, 3.0},
		{"BrightBlue on BG", p.BrightBlue, p.BG, 3.0},
		{"BrightMagenta on BG", p.BrightMagenta, p.BG, 3.0},
		{"BrightCyan on BG", p.BrightCyan, p.BG, 3.0},
	}

	checks := make([]ContrastCheck, len(pairs))
	for i, c := range pairs {
		r := ContrastRatio(c.text, c.bg)
		checks[i] = ContrastCheck{
			Label:   c.label,
			Text:    c.text,
			BG:      c.bg,
			Minimum: c.minRatio,
			Ratio:   r,
			Pass:    r >= c.minRatio,
		}
	}
	return checks
}

// ValidateContrasts checks key palette pairs and returns issues.
func ValidateContrasts(p Palette) []string {
	var issues []string
	for _, c := range ContrastChecks(p) {
		if !c.Pass {
			issues = append(issues, fmt.Sprintf("%s: %.1f:1 (need %.1f:1)", c.Label, c.Ratio, c.Minimum))
		}
	}
	return issues
}

// OKLCHString returns a compact OKLCH annotation for a hex color.
func OKLCHString(hex string) string {
	c, _ := colorful.Hex(hex)
	l, ch, h := c.OkLch()
	return fmt.Sprintf("oklch(%.2f, %.3f, %.0f°)", l, ch, h)
}

// DumpPalette prints all palette fields with OKLCH annotations.
func DumpPalette(p Palette) {
	fields := []string{
		"BG", "Surface", "Overlay", "FG", "FGDim",
		"Cursor", "CursorText", "SelectionBG", "SelectionFG",
		"Accent", "ButtonBG", "ButtonFG", "DebugFG",
		"DiffAddBG", "DiffDeleteBG", "DiffChangeBG", "SearchBG", "VisualBG",
		"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White", "Orange",
		"BrightBlack", "BrightRed", "BrightGreen", "BrightYellow", "BrightBlue",
		"BrightMagenta", "BrightCyan", "BrightWhite",
	}
	vals := map[string]string{
		"BG": p.BG, "Surface": p.Surface, "Overlay": p.Overlay,
		"FG": p.FG, "FGDim": p.FGDim,
		"Cursor": p.Cursor, "CursorText": p.CursorText,
		"SelectionBG": p.SelectionBG, "SelectionFG": p.SelectionFG,
		"Accent": p.Accent, "ButtonBG": p.ButtonBG, "ButtonFG": p.ButtonFG, "DebugFG": p.DebugFG,
		"DiffAddBG": p.DiffAddBG, "DiffDeleteBG": p.DiffDeleteBG,
		"DiffChangeBG": p.DiffChangeBG, "SearchBG": p.SearchBG, "VisualBG": p.VisualBG,
		"Black": p.Black, "Red": p.Red, "Green": p.Green, "Yellow": p.Yellow,
		"Blue": p.Blue, "Magenta": p.Magenta, "Cyan": p.Cyan, "White": p.White,
		"Orange":      p.Orange,
		"BrightBlack": p.BrightBlack, "BrightRed": p.BrightRed,
		"BrightGreen": p.BrightGreen, "BrightYellow": p.BrightYellow,
		"BrightBlue": p.BrightBlue, "BrightMagenta": p.BrightMagenta,
		"BrightCyan": p.BrightCyan, "BrightWhite": p.BrightWhite,
	}

	var sb strings.Builder
	for _, f := range fields {
		v := vals[f]
		fmt.Fprintf(&sb, "  %-15s %s  %s\n", f, v, OKLCHString(v))
	}
	fmt.Print(sb.String())
}
