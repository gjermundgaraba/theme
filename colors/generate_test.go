package colors

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

func TestContrastValidationAcrossHueRange(t *testing.T) {
	for _, dark := range []bool{true, false} {
		mode := "light"
		if dark {
			mode = "dark"
		}
		t.Run(mode, func(t *testing.T) {
			for surface := 0.0; surface <= 360; surface += 15 {
				for accent := 0.0; accent <= 360; accent += 15 {
					p := Generate(ThemeConfig{SurfaceHue: surface, AccentHue: accent}, dark)
					if issues := ValidateContrasts(p); len(issues) > 0 {
						t.Fatalf("surface=%.0f accent=%.0f issues=%v", surface, accent, issues)
					}
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

func TestDefaultConfigBackgroundIsLifted(t *testing.T) {
	c, err := colorful.Hex(Dark.BG)
	if err != nil {
		t.Fatal(err)
	}
	l, _, _ := c.OkLch()
	want := DefaultConfig.DarkSurfaceLightness
	if l < want-0.01 || l > want+0.01 {
		t.Fatalf("default dark BG OKLCH L = %.3f, want ~%.2f", l, want)
	}
	if want < MinDarkSurfaceLightness+0.05 {
		t.Fatalf("default dark BG L = %.2f sits at the floor; lift it above MinDarkSurfaceLightness+0.05", want)
	}
}

func TestSurfaceLightnessChangesBG(t *testing.T) {
	low := Generate(ThemeConfig{
		SurfaceHue:            143,
		AccentHue:             360,
		DarkSurfaceLightness:  0.15,
		LightSurfaceLightness: 0.95,
	}, true)
	high := Generate(ThemeConfig{
		SurfaceHue:            143,
		AccentHue:             360,
		DarkSurfaceLightness:  0.30,
		LightSurfaceLightness: 0.95,
	}, true)
	if low.BG == high.BG {
		t.Fatalf("dark BG did not change with DarkSurfaceLightness: %s vs %s", low.BG, high.BG)
	}
	if low.Surface == high.Surface || low.Overlay == high.Overlay {
		t.Fatalf("Surface/Overlay should follow BG: low=%+v high=%+v", low, high)
	}

	dim := Generate(ThemeConfig{
		SurfaceHue:            143,
		AccentHue:             360,
		DarkSurfaceLightness:  0.27,
		LightSurfaceLightness: 0.85,
	}, false)
	bright := Generate(ThemeConfig{
		SurfaceHue:            143,
		AccentHue:             360,
		DarkSurfaceLightness:  0.27,
		LightSurfaceLightness: 0.99,
	}, false)
	if dim.BG == bright.BG {
		t.Fatalf("light BG did not change with LightSurfaceLightness: %s vs %s", dim.BG, bright.BG)
	}
}

func TestZeroLightnessFallsBackToDefaults(t *testing.T) {
	partial := Generate(ThemeConfig{
		SurfaceHue: DefaultConfig.SurfaceHue,
		AccentHue:  DefaultConfig.AccentHue,
	}, true)
	full := Generate(DefaultConfig, true)
	if partial.BG != full.BG {
		t.Fatalf("partial config did not fall back to default lightness: %s vs %s", partial.BG, full.BG)
	}
}

func TestContrastValidationAcrossLightnessRange(t *testing.T) {
	for dl := MinDarkSurfaceLightness; dl <= MaxDarkSurfaceLightness+1e-9; dl += 0.01 {
		cfg := ThemeConfig{
			SurfaceHue:            143,
			AccentHue:             360,
			DarkSurfaceLightness:  dl,
			LightSurfaceLightness: DefaultConfig.LightSurfaceLightness,
		}
		if issues := ValidateContrasts(Generate(cfg, true)); len(issues) > 0 {
			t.Fatalf("dark L=%.2f issues=%v", dl, issues)
		}
	}
	for ll := MinLightSurfaceLightness; ll <= MaxLightSurfaceLightness+1e-9; ll += 0.01 {
		cfg := ThemeConfig{
			SurfaceHue:            143,
			AccentHue:             360,
			DarkSurfaceLightness:  DefaultConfig.DarkSurfaceLightness,
			LightSurfaceLightness: ll,
		}
		if issues := ValidateContrasts(Generate(cfg, false)); len(issues) > 0 {
			t.Fatalf("light L=%.2f issues=%v", ll, issues)
		}
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

		ratio := ContrastRatio(p.pal.FG, p.pal.BG)
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
			r := ContrastRatio(pair.fg, pair.bg)
			t.Logf("  %s: %.1f:1", pair.label, r)
		}
	}
}
