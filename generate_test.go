package main

import (
	"testing"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func TestContrastValidation(t *testing.T) {
	for _, p := range []struct {
		name string
		pal  Palette
	}{
		{"Dark", Dark},
		{"Light", Light},
	} {
		t.Run(p.name, func(t *testing.T) {
			issues := ValidateContrasts(p.pal)
			if len(issues) > 0 {
				for _, issue := range issues {
					t.Error(issue)
				}
			}
		})
	}
}

func TestSyntaxColorSeparation(t *testing.T) {
	for _, p := range []struct {
		name string
		pal  Palette
	}{
		{"Dark", Dark},
		{"Light", Light},
	} {
		t.Run(p.name, func(t *testing.T) {
			colors := map[string]string{
				"Red": p.pal.Red, "Orange": p.pal.Orange, "Yellow": p.pal.Yellow,
				"Green": p.pal.Green, "Cyan": p.pal.Cyan, "Blue": p.pal.Blue, "Magenta": p.pal.Magenta,
			}
			names := []string{"Red", "Orange", "Yellow", "Green", "Cyan", "Blue", "Magenta"}

			for i := 0; i < len(names); i++ {
				for j := i + 1; j < len(names); j++ {
					c1, _ := colorful.Hex(colors[names[i]])
					c2, _ := colorful.Hex(colors[names[j]])
					dist := c1.DistanceCIEDE2000(c2)
				if dist < 0.15 {
					t.Errorf("%s vs %s: CIEDE2000=%.3f (want >= 0.15)", names[i], names[j], dist)
					}
				}
			}
		})
	}
}

func TestDumpPalette(t *testing.T) {
	for _, p := range []struct {
		name string
		pal  Palette
	}{
		{"Dark", Dark},
		{"Light", Light},
	} {
		t.Logf("=== %s ===", p.name)
		DumpPalette(p.pal)

		ratio := contrastRatio(p.pal.FG, p.pal.BG)
		t.Logf("  FG↔BG contrast: %.1f:1", ratio)

		for _, pair := range []struct {
			label, fg, bg string
		}{
			{"FGDim on BG", p.pal.FGDim, p.pal.BG},
			{"Overlay on BG", p.pal.Overlay, p.pal.BG},
			{"Red on BG", p.pal.Red, p.pal.BG},
			{"Blue on BG", p.pal.Blue, p.pal.BG},
			{"ButtonFG on ButtonBG", p.pal.ButtonFG, p.pal.ButtonBG},
		} {
			r := contrastRatio(pair.fg, pair.bg)
			t.Logf("  %s: %.1f:1", pair.label, r)
		}
	}
}
