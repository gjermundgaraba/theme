package generator

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gjermundgaraba/theme/colors"
)

func TestParseHueFallsBackAndClamps(t *testing.T) {
	for _, tt := range []struct {
		name     string
		raw      string
		fallback float64
		want     float64
	}{
		{"empty", "", 157, 157},
		{"invalid", "nope", 157, 157},
		{"negative", "-10", 157, 0},
		{"too high", "480", 157, 360},
		{"valid", "210", 157, 210},
		{"decimal", "210.7", 157, 211},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseHue(tt.raw, tt.fallback); got != tt.want {
				t.Fatalf("parseHue(%q, %.0f) = %.0f, want %.0f", tt.raw, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestHandleIndexRendersCurrentConfig(t *testing.T) {
	old := colors.DefaultConfig
	defer colors.SetConfig(old)
	colors.SetConfig(colors.ThemeConfig{
		SurfaceHue:            210,
		AccentHue:             45,
		DarkSurfaceLightness:  0.22,
		LightSurfaceLightness: 0.93,
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handleIndex(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		`id="surface-hue-val">210°`,
		`value="210" class="hue-slider" id="surface-hue"`,
		`id="accent-hue-val">45°`,
		`value="45" class="hue-slider" id="accent-hue"`,
		`id="dark-lightness-val">22%`,
		`value="22" class="lightness-slider" id="dark-lightness"`,
		`id="light-lightness-val">93%`,
		`value="93" class="lightness-slider" id="light-lightness"`,
		`surfaceHue: 210,`,
		`accentHue: 45,`,
		`darkLightness: 22 / 100,`,
		`lightLightness: 93 / 100,`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("rendered index missing %q", want)
		}
	}
}

func TestRewritePaletteGoIsAtomicAndCleansTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "palette.go")
	original := []byte("package colors\n\nvar DefaultConfig = ThemeConfig{\n\tSurfaceHue:            157,\n\tAccentHue:             360,\n\tDarkSurfaceLightness:  0.27,\n\tLightSurfaceLightness: 0.95,\n}\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := rewritePaletteGo(path, colors.ThemeConfig{
		SurfaceHue:            210,
		AccentHue:             45,
		DarkSurfaceLightness:  0.22,
		LightSurfaceLightness: 0.93,
	}); err != nil {
		t.Fatal(err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(updated)
	for _, want := range []string{
		"SurfaceHue:            210",
		"AccentHue:             45",
		"DarkSurfaceLightness:  0.22",
		"LightSurfaceLightness: 0.93",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("updated palette missing %q: %s", want, content)
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %v, want 0600", got)
	}
	matches, err := filepath.Glob(filepath.Join(dir, ".palette.go.*.tmp"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files were not cleaned: %v", matches)
	}
}

func TestHandleGenerateRequiresGet(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/generate", nil)
	rec := httptest.NewRecorder()

	handleGenerate(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandleGenerateReturnsNamedPalettes(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/generate?surfaceHue=210&accentHue=45", nil)
	rec := httptest.NewRecorder()

	handleGenerate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got map[string]paletteJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["dark"].Name != "gg-dark" || got["light"].Name != "gg-light" {
		t.Fatalf("unexpected names: dark=%q light=%q", got["dark"].Name, got["light"].Name)
	}
	for _, name := range []string{"vscode/package.json", "zed"} {
		if !sliceContains(got["dark"].TemplateNames, name) {
			t.Fatalf("dark template names missing %q: %v", name, got["dark"].TemplateNames)
		}
	}
}

func sliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func TestRewritePaletteGoRejectsMissingFields(t *testing.T) {
	for _, tt := range []struct {
		name    string
		content string
		wantSub string
	}{
		{
			"missing all",
			"package colors\n\nvar DefaultConfig = ThemeConfig{}\n",
			"SurfaceHue",
		},
		{
			"missing dark lightness",
			"package colors\n\nvar DefaultConfig = ThemeConfig{\n\tSurfaceHue: 157,\n\tAccentHue:  360,\n\tLightSurfaceLightness: 0.95,\n}\n",
			"DarkSurfaceLightness",
		},
		{
			"missing light lightness",
			"package colors\n\nvar DefaultConfig = ThemeConfig{\n\tSurfaceHue: 157,\n\tAccentHue:  360,\n\tDarkSurfaceLightness: 0.27,\n}\n",
			"LightSurfaceLightness",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "palette.go")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
			err := rewritePaletteGo(path, colors.ThemeConfig{
				SurfaceHue:            210,
				AccentHue:             45,
				DarkSurfaceLightness:  0.22,
				LightSurfaceLightness: 0.93,
			})
			if err == nil {
				t.Fatal("expected missing field error")
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("error %q does not mention %q", err.Error(), tt.wantSub)
			}
		})
	}
}

func TestParseLightnessClampsAndRounds(t *testing.T) {
	for _, tt := range []struct {
		name             string
		raw              string
		fallback, lo, hi float64
		want             float64
	}{
		{"empty falls back", "", 0.27, 0.12, 0.30, 0.27},
		{"invalid falls back", "nope", 0.27, 0.12, 0.30, 0.27},
		{"below clamp", "0.05", 0.27, 0.12, 0.30, 0.12},
		{"above clamp", "0.99", 0.27, 0.12, 0.30, 0.30},
		{"valid rounds", "0.275", 0.27, 0.12, 0.30, 0.28},
		{"valid passes through", "0.22", 0.27, 0.12, 0.30, 0.22},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLightness(tt.raw, tt.fallback, tt.lo, tt.hi)
			if got != tt.want {
				t.Fatalf("parseLightness(%q) = %.3f, want %.3f", tt.raw, got, tt.want)
			}
		})
	}
}

func TestValidateSameOrigin(t *testing.T) {
	missing := httptest.NewRequest("POST", "/api/apply", nil)
	if err := validateSameOrigin(missing); err == nil {
		t.Fatal("expected missing origin to be rejected")
	}

	allowed := httptest.NewRequest("POST", "/api/apply", nil)
	allowed.Header.Set("Origin", "http://localhost:9090")
	if err := validateSameOrigin(allowed); err != nil {
		t.Fatalf("allowed origin rejected: %v", err)
	}
}

func TestHandleApplySuccessMutatesConfigAndBuilds(t *testing.T) {
	old := colors.DefaultConfig
	defer colors.SetConfig(old)

	t.Chdir(t.TempDir())
	palettePath := filepath.Join("colors", "palette.go")
	if err := os.MkdirAll(filepath.Dir(palettePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(palettePath, []byte("package colors\n\nvar DefaultConfig = ThemeConfig{\n\tSurfaceHue:            157,\n\tAccentHue:             360,\n\tDarkSurfaceLightness:  0.27,\n\tLightSurfaceLightness: 0.95,\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/apply?surfaceHue=210&accentHue=45&darkSurfaceLightness=0.22&lightSurfaceLightness=0.93", nil)
	req.Header.Set("Origin", "http://localhost:9090")
	rec := httptest.NewRecorder()

	handleApply(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status = %d, body = %q", rec.Code, rec.Body.String())
	}
	got := colors.DefaultConfig
	if got.SurfaceHue != 210 || got.AccentHue != 45 || got.DarkSurfaceLightness != 0.22 || got.LightSurfaceLightness != 0.93 {
		t.Fatalf("config = %+v, want surface 210 accent 45 dark 0.22 light 0.93", got)
	}
	rewritten, err := os.ReadFile(palettePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"DarkSurfaceLightness:  0.22", "LightSurfaceLightness: 0.93"} {
		if !strings.Contains(string(rewritten), want) {
			t.Fatalf("rewritten palette.go missing %q:\n%s", want, rewritten)
		}
	}
	if _, err := os.Stat(filepath.Join("build", "ghostty", "gg-dark")); err != nil {
		t.Fatalf("expected build output: %v", err)
	}
}

func TestHandleApplyRollsBackWhenBuildFails(t *testing.T) {
	old := colors.DefaultConfig
	defer colors.SetConfig(old)
	previousBuilder := buildWithConfig
	defer func() { buildWithConfig = previousBuilder }()

	t.Chdir(t.TempDir())
	palettePath := filepath.Join("colors", "palette.go")
	if err := os.MkdirAll(filepath.Dir(palettePath), 0o755); err != nil {
		t.Fatal(err)
	}
	original := "package colors\n\nvar DefaultConfig = ThemeConfig{\n\tSurfaceHue:            157,\n\tAccentHue:             270,\n\tDarkSurfaceLightness:  0.27,\n\tLightSurfaceLightness: 0.95,\n}\n"
	if err := os.WriteFile(palettePath, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	colors.SetConfig(colors.ThemeConfig{
		SurfaceHue:            157,
		AccentHue:             270,
		DarkSurfaceLightness:  0.27,
		LightSurfaceLightness: 0.95,
	})
	buildWithConfig = func(cfg colors.ThemeConfig) error {
		if cfg.SurfaceHue == 210 && cfg.AccentHue == 45 {
			return errors.New("forced build failure")
		}
		return nil
	}

	req := httptest.NewRequest("POST", "/api/apply?surfaceHue=210&accentHue=45", nil)
	req.Header.Set("Origin", "http://localhost:9090")
	rec := httptest.NewRecorder()

	handleApply(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if colors.DefaultConfig.SurfaceHue != 157 || colors.DefaultConfig.AccentHue != 270 {
		t.Fatalf("config was not rolled back: %+v", colors.DefaultConfig)
	}
	content, err := os.ReadFile(palettePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != original {
		t.Fatalf("palette.go was not rolled back:\n%s", string(content))
	}
}

func TestHandleApplyRejectsCrossOrigin(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/apply?surfaceHue=210&accentHue=45", nil)
	req.Header.Set("Origin", "http://evil.test")
	rec := httptest.NewRecorder()

	handleApply(rec, req)

	if rec.Code != 403 {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "origin") {
		t.Fatalf("expected origin error, got %q", rec.Body.String())
	}
}

func TestPaletteJSONExposesEveryGeneratedColor(t *testing.T) {
	cfg := colors.ThemeConfig{
		SurfaceHue:            157,
		AccentHue:             360,
		DarkSurfaceLightness:  0.27,
		LightSurfaceLightness: 0.95,
	}
	got := toPaletteJSON("gg-dark", cfg, colors.Generate(cfg, true))
	wantKeys := []string{
		"BG", "Surface", "Overlay", "FG", "FGDim",
		"Cursor", "CursorText", "SelectionBG", "SelectionFG",
		"Accent", "ButtonBG", "ButtonFG", "DebugFG",
		"DiffAddBG", "DiffDeleteBG", "DiffChangeBG", "SearchBG", "VisualBG",
		"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White", "Orange",
		"BrightBlack", "BrightRed", "BrightGreen", "BrightYellow", "BrightBlue", "BrightMagenta", "BrightCyan", "BrightWhite",
	}

	if got.Config.SurfaceHue != cfg.SurfaceHue || got.Config.AccentHue != cfg.AccentHue {
		t.Fatalf("config = %+v, want %+v", got.Config, cfg)
	}
	if got.Config.DarkSurfaceLightness != cfg.DarkSurfaceLightness || got.Config.LightSurfaceLightness != cfg.LightSurfaceLightness {
		t.Fatalf("lightness = (%.2f, %.2f), want (%.2f, %.2f)",
			got.Config.DarkSurfaceLightness, got.Config.LightSurfaceLightness,
			cfg.DarkSurfaceLightness, cfg.LightSurfaceLightness)
	}
	if len(got.Colors) != len(wantKeys) {
		t.Fatalf("color count = %d, want %d", len(got.Colors), len(wantKeys))
	}
	for _, key := range wantKeys {
		if got.Colors[key] == "" {
			t.Fatalf("missing generated color %q", key)
		}
	}
}

func TestPaletteJSONIncludesTyporaForDarkAndLight(t *testing.T) {
	cfg := colors.ThemeConfig{SurfaceHue: 157, AccentHue: 360}
	dark := toPaletteJSON("gg-dark", cfg, colors.Generate(cfg, true))
	for _, key := range []string{"typora.css", "typora/codeblock", "typora/mermaid", "typora/sourcemode"} {
		if !sliceContains(dark.TemplateNames, key) {
			t.Fatalf("dark template names missing %q: %v", key, dark.TemplateNames)
		}
	}

	light := toPaletteJSON("gg-light", cfg, colors.Generate(cfg, false))
	for _, key := range []string{"typora.css", "typora/codeblock", "typora/mermaid", "typora/sourcemode"} {
		if !sliceContains(light.TemplateNames, key) {
			t.Fatalf("light template names missing %q: %v", key, light.TemplateNames)
		}
	}
}

func TestHandleTemplateRendersNamedTarget(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/template?name=ghostty&mode=dark&surfaceHue=210&accentHue=45", nil)
	rec := httptest.NewRecorder()

	handleTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", rec.Code, rec.Body.String())
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["content"] == "" {
		t.Fatal("template content is empty")
	}
}

func TestHandleTemplateRejectsUnknownName(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/template?name=bogus", nil)
	rec := httptest.NewRecorder()

	handleTemplate(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandleTemplateRequiresGet(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/template?name=ghostty", nil)
	rec := httptest.NewRecorder()

	handleTemplate(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
