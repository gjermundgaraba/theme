package generator

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/gjermundgaraba/theme/colors"
	"github.com/gjermundgaraba/theme/themes"
)

//go:embed web
var webEmbed embed.FS

var indexTemplate = template.Must(template.New("index.html").ParseFS(webEmbed, "web/index.html"))

var (
	applyMu                 sync.RWMutex
	buildWithConfig         = themes.BuildWithConfig
	surfaceHueRE            = regexp.MustCompile(`(SurfaceHue:\s*)\d+`)
	accentHueRE             = regexp.MustCompile(`(AccentHue:\s*)\d+`)
	darkSurfaceLightnessRE  = regexp.MustCompile(`(DarkSurfaceLightness:\s*)[\d.]+`)
	lightSurfaceLightnessRE = regexp.MustCompile(`(LightSurfaceLightness:\s*)[\d.]+`)
)

type paletteJSON struct {
	Name      string            `json:"name"`
	Config    configJSON        `json:"config"`
	Colors    map[string]string `json:"colors"`
	Contrasts []contrastJSON    `json:"contrasts"`
	Templates map[string]string `json:"templates"`
}

type configJSON struct {
	SurfaceHue            float64 `json:"surfaceHue"`
	AccentHue             float64 `json:"accentHue"`
	DarkSurfaceLightness  float64 `json:"darkSurfaceLightness"`
	LightSurfaceLightness float64 `json:"lightSurfaceLightness"`
}

type contrastJSON struct {
	Label   string  `json:"label"`
	Ratio   float64 `json:"ratio"`
	Minimum float64 `json:"minimum"`
	Pass    bool    `json:"pass"`
}

// Serve starts the theme editor web server.
func Serve() {
	webFS, err := fs.Sub(webEmbed, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	staticFiles := http.FileServer(http.FS(webFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			handleIndex(w, r)
			return
		}
		staticFiles.ServeHTTP(w, r)
	})
	mux.HandleFunc("/api/generate", handleGenerate)
	mux.HandleFunc("/api/apply", handleApply)

	addr := "localhost:9090"
	fmt.Printf("Theme editor: http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	applyMu.RLock()
	cfg := colors.DefaultConfig
	palette := colors.Dark
	applyMu.RUnlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTemplate.Execute(w, struct {
		SurfaceHue                  int
		AccentHue                   int
		DarkSurfaceLightnessPct     int
		LightSurfaceLightnessPct    int
		MinDarkSurfaceLightnessPct  int
		MaxDarkSurfaceLightnessPct  int
		MinLightSurfaceLightnessPct int
		MaxLightSurfaceLightnessPct int
		Palette                     colors.Palette
	}{
		SurfaceHue:                  int(math.Round(cfg.SurfaceHue)),
		AccentHue:                   int(math.Round(cfg.AccentHue)),
		DarkSurfaceLightnessPct:     int(math.Round(cfg.DarkSurfaceLightness * 100)),
		LightSurfaceLightnessPct:    int(math.Round(cfg.LightSurfaceLightness * 100)),
		MinDarkSurfaceLightnessPct:  int(math.Round(colors.MinDarkSurfaceLightness * 100)),
		MaxDarkSurfaceLightnessPct:  int(math.Round(colors.MaxDarkSurfaceLightness * 100)),
		MinLightSurfaceLightnessPct: int(math.Round(colors.MinLightSurfaceLightness * 100)),
		MaxLightSurfaceLightnessPct: int(math.Round(colors.MaxLightSurfaceLightness * 100)),
		Palette:                     palette,
	}); err != nil {
		log.Printf("render index: %v", err)
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	applyMu.RLock()
	defaults := colors.DefaultConfig
	applyMu.RUnlock()

	q := r.URL.Query()
	cfg := configFromQuery(q, defaults)

	dark, err := toPaletteJSON("gg-dark", cfg, colors.Generate(cfg, true))
	if err != nil {
		jsonError(w, "render templates: "+err.Error())
		return
	}
	light, err := toPaletteJSON("gg-light", cfg, colors.Generate(cfg, false))
	if err != nil {
		jsonError(w, "render templates: "+err.Error())
		return
	}
	resp := map[string]*paletteJSON{
		"dark":  dark,
		"light": light,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("encode generate response: %v", err)
	}
}

func toPaletteJSON(name string, cfg colors.ThemeConfig, p colors.Palette) (*paletteJSON, error) {
	colorMap := map[string]string{
		"BG": p.BG, "Surface": p.Surface, "Overlay": p.Overlay,
		"FG": p.FG, "FGDim": p.FGDim,
		"Cursor": p.Cursor, "CursorText": p.CursorText,
		"SelectionBG": p.SelectionBG, "SelectionFG": p.SelectionFG,
		"Accent": p.Accent, "ButtonBG": p.ButtonBG, "ButtonFG": p.ButtonFG,
		"DebugFG":   p.DebugFG,
		"DiffAddBG": p.DiffAddBG, "DiffDeleteBG": p.DiffDeleteBG,
		"DiffChangeBG": p.DiffChangeBG, "SearchBG": p.SearchBG,
		"VisualBG": p.VisualBG,
		"Black":    p.Black, "Red": p.Red, "Green": p.Green,
		"Yellow": p.Yellow, "Blue": p.Blue, "Magenta": p.Magenta,
		"Cyan": p.Cyan, "White": p.White, "Orange": p.Orange,
		"BrightBlack": p.BrightBlack, "BrightRed": p.BrightRed,
		"BrightGreen": p.BrightGreen, "BrightYellow": p.BrightYellow,
		"BrightBlue": p.BrightBlue, "BrightMagenta": p.BrightMagenta,
		"BrightCyan": p.BrightCyan, "BrightWhite": p.BrightWhite,
	}

	checks := colors.ContrastChecks(p)
	contrasts := make([]contrastJSON, len(checks))
	for i, c := range checks {
		contrasts[i] = contrastJSON{
			Label:   c.Label,
			Ratio:   math.Round(c.Ratio*10) / 10,
			Minimum: c.Minimum,
			Pass:    c.Pass,
		}
	}

	data := themes.ThemeData{Name: name, Palette: p}
	templates := make(map[string]string)
	for key, tmplPath := range map[string]string{
		"ghostty":             "templates/ghostty.tmpl",
		"fish":                "templates/fish.tmpl",
		"vscode":              "templates/vscode/theme.json.tmpl",
		"vscode/package.json": "templates/vscode/package.json.tmpl",
		"neovim":              "templates/neovim.lua.tmpl",
	} {
		if err := addTemplate(templates, key, tmplPath, data); err != nil {
			return nil, err
		}
	}
	for key, tmplPath := range map[string]string{
		"typora.css":        "templates/typora/theme.css.tmpl",
		"typora/codeblock":  "templates/typora/codeblock.dark.css.tmpl",
		"typora/mermaid":    "templates/typora/mermaid.dark.css.tmpl",
		"typora/sourcemode": "templates/typora/sourcemode.dark.css.tmpl",
	} {
		if err := addTemplate(templates, key, tmplPath, data); err != nil {
			return nil, err
		}
	}

	return &paletteJSON{
		Name: name,
		Config: configJSON{
			SurfaceHue:            cfg.SurfaceHue,
			AccentHue:             cfg.AccentHue,
			DarkSurfaceLightness:  cfg.DarkSurfaceLightness,
			LightSurfaceLightness: cfg.LightSurfaceLightness,
		},
		Colors:    colorMap,
		Contrasts: contrasts,
		Templates: templates,
	}, nil
}

func configFromQuery(q url.Values, defaults colors.ThemeConfig) colors.ThemeConfig {
	return colors.ThemeConfig{
		SurfaceHue: parseHue(q.Get("surfaceHue"), defaults.SurfaceHue),
		AccentHue:  parseHue(q.Get("accentHue"), defaults.AccentHue),
		DarkSurfaceLightness: parseLightness(
			q.Get("darkSurfaceLightness"),
			defaults.DarkSurfaceLightness,
			colors.MinDarkSurfaceLightness,
			colors.MaxDarkSurfaceLightness,
		),
		LightSurfaceLightness: parseLightness(
			q.Get("lightSurfaceLightness"),
			defaults.LightSurfaceLightness,
			colors.MinLightSurfaceLightness,
			colors.MaxLightSurfaceLightness,
		),
	}
}

func addTemplate(dst map[string]string, key, tmplPath string, data themes.ThemeData) error {
	s, err := themes.RenderTemplateString(tmplPath, data)
	if err != nil {
		return err
	}
	dst[key] = s
	return nil
}

func handleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := validateSameOrigin(r); err != nil {
		jsonErrorStatus(w, http.StatusForbidden, err.Error())
		return
	}

	applyMu.Lock()
	defer applyMu.Unlock()

	previous := colors.DefaultConfig
	q := r.URL.Query()
	cfg := configFromQuery(q, previous)
	paletteGo := filepath.Join("colors", "palette.go")
	if err := rewritePaletteGo(paletteGo, cfg); err != nil {
		jsonError(w, "rewrite palette.go: "+err.Error())
		return
	}

	colors.SetConfig(cfg)
	if err := buildWithConfig(cfg); err != nil {
		rollbackErr := rollbackApply(paletteGo, previous)
		msg := "build themes: " + err.Error()
		if rollbackErr != nil {
			msg += "; rollback failed: " + rollbackErr.Error()
		}
		jsonError(w, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("encode apply response: %v", err)
	}
}

func rollbackApply(paletteGo string, previous colors.ThemeConfig) error {
	var errs []string
	if err := rewritePaletteGo(paletteGo, previous); err != nil {
		errs = append(errs, "rewrite palette.go: "+err.Error())
	}
	colors.SetConfig(previous)
	if err := buildWithConfig(previous); err != nil {
		errs = append(errs, "build themes: "+err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

func parseHue(raw string, fallback float64) float64 {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return fallback
	}
	if v < 0 {
		return 0
	}
	if v > 360 {
		return 360
	}
	return math.Round(v)
}

// parseLightness reads an OKLCH lightness from a query string, clamps it to
// [min,max], and rounds to two decimals so palette.go stores stable values.
func parseLightness(raw string, fallback, min, max float64) float64 {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return fallback
	}
	if v < min {
		v = min
	}
	if v > max {
		v = max
	}
	return math.Round(v*100) / 100
}

func rewritePaletteGo(path string, cfg colors.ThemeConfig) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	s := string(content)
	if !surfaceHueRE.MatchString(s) {
		return fmt.Errorf("palette.go does not contain SurfaceHue")
	}
	if !accentHueRE.MatchString(s) {
		return fmt.Errorf("palette.go does not contain AccentHue")
	}
	if !darkSurfaceLightnessRE.MatchString(s) {
		return fmt.Errorf("palette.go does not contain DarkSurfaceLightness")
	}
	if !lightSurfaceLightnessRE.MatchString(s) {
		return fmt.Errorf("palette.go does not contain LightSurfaceLightness")
	}
	s = surfaceHueRE.ReplaceAllString(s, fmt.Sprintf("${1}%d", int(cfg.SurfaceHue)))
	s = accentHueRE.ReplaceAllString(s, fmt.Sprintf("${1}%d", int(cfg.AccentHue)))
	s = darkSurfaceLightnessRE.ReplaceAllString(s, fmt.Sprintf("${1}%.2f", cfg.DarkSurfaceLightness))
	s = lightSurfaceLightnessRE.ReplaceAllString(s, fmt.Sprintf("${1}%.2f", cfg.LightSurfaceLightness))

	if !strings.Contains(s, "DefaultConfig") {
		return fmt.Errorf("palette.go does not contain DefaultConfig")
	}

	return atomicWriteFile(path, []byte(s), info.Mode().Perm())
}

func atomicWriteFile(path string, data []byte, perm fs.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	tmp, err := os.CreateTemp(dir, "."+base+".*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func validateSameOrigin(r *http.Request) error {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return fmt.Errorf("missing request origin")
	}
	u, err := url.Parse(origin)
	if err != nil {
		return fmt.Errorf("invalid request origin")
	}
	allowed := map[string]bool{
		"http://localhost:9090": true,
		"http://127.0.0.1:9090": true,
	}
	if !allowed[u.Scheme+"://"+u.Host] {
		return fmt.Errorf("request origin is not allowed")
	}
	return nil
}

func jsonError(w http.ResponseWriter, msg string) {
	jsonErrorStatus(w, http.StatusInternalServerError, msg)
}

func jsonErrorStatus(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		log.Printf("encode error response: %v", err)
	}
}
