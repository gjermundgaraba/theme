package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//go:embed web
var webEmbed embed.FS

type paletteJSON struct {
	Colors    map[string]string  `json:"colors"`
	Contrasts []contrastJSON     `json:"contrasts"`
	Templates map[string]string  `json:"templates"`
}

type contrastJSON struct {
	Label   string  `json:"label"`
	Ratio   float64 `json:"ratio"`
	Minimum float64 `json:"minimum"`
	Pass    bool    `json:"pass"`
}

func serve() {
	webFS, err := fs.Sub(webEmbed, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(webFS)))
	mux.HandleFunc("/api/generate", handleGenerate)
	mux.HandleFunc("/api/apply", handleApply)

	addr := "localhost:9090"
	fmt.Printf("Theme editor: http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	surfaceHue, _ := strconv.ParseFloat(q.Get("surfaceHue"), 64)
	accentHue, _ := strconv.ParseFloat(q.Get("accentHue"), 64)

	cfg := ThemeConfig{SurfaceHue: surfaceHue, AccentHue: accentHue}

	resp := map[string]*paletteJSON{
		"dark":  toPaletteJSON("gg-dark", Generate(cfg, true)),
		"light": toPaletteJSON("gg-light", Generate(cfg, false)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func toPaletteJSON(name string, p Palette) *paletteJSON {
	colors := map[string]string{
		"BG": p.BG, "Surface": p.Surface, "Overlay": p.Overlay,
		"FG": p.FG, "FGDim": p.FGDim,
		"Cursor": p.Cursor, "CursorText": p.CursorText,
		"SelectionBG": p.SelectionBG, "SelectionFG": p.SelectionFG,
		"ButtonBG": p.ButtonBG, "ButtonFG": p.ButtonFG,
		"Black": p.Black, "Red": p.Red, "Green": p.Green,
		"Yellow": p.Yellow, "Blue": p.Blue, "Magenta": p.Magenta,
		"Cyan": p.Cyan, "White": p.White, "Orange": p.Orange,
		"BrightBlack": p.BrightBlack, "BrightRed": p.BrightRed,
		"BrightGreen": p.BrightGreen, "BrightYellow": p.BrightYellow,
		"BrightBlue": p.BrightBlue, "BrightMagenta": p.BrightMagenta,
		"BrightCyan": p.BrightCyan, "BrightWhite": p.BrightWhite,
	}

	checks := []struct {
		label, text, bg string
		min             float64
	}{
		{"FG on BG", p.FG, p.BG, 7.0},
		{"FGDim on BG", p.FGDim, p.BG, 3.0},
		{"Red on BG", p.Red, p.BG, 3.0},
		{"Green on BG", p.Green, p.BG, 3.0},
		{"Yellow on BG", p.Yellow, p.BG, 3.0},
		{"Blue on BG", p.Blue, p.BG, 3.0},
		{"Magenta on BG", p.Magenta, p.BG, 3.0},
		{"Cyan on BG", p.Cyan, p.BG, 3.0},
		{"Orange on BG", p.Orange, p.BG, 3.0},
		{"ButtonFG on ButtonBG", p.ButtonFG, p.ButtonBG, 4.5},
	}

	contrasts := make([]contrastJSON, len(checks))
	for i, c := range checks {
		ratio := contrastRatio(c.text, c.bg)
		contrasts[i] = contrastJSON{
			Label:   c.label,
			Ratio:   math.Round(ratio*10) / 10,
			Minimum: c.min,
			Pass:    ratio >= c.min,
		}
	}

	data := ThemeData{Name: name, Palette: p}
	templates := map[string]string{
		"ghostty": renderTemplateToStr("templates/ghostty.tmpl", data),
		"fish":    renderTemplateToStr("templates/fish.tmpl", data),
		"vscode":  renderTemplateToStr("templates/vscode/theme.json.tmpl", data),
		"neovim":  renderTemplateToStr("templates/neovim.lua.tmpl", data),
	}

	return &paletteJSON{
		Colors:    colors,
		Contrasts: contrasts,
		Templates: templates,
	}
}

func renderTemplateToStr(tmplPath string, data any) string {
	t, err := template.New(filepath.Base(tmplPath)).Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return buf.String()
}

func handleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	surfaceHue, _ := strconv.ParseFloat(q.Get("surfaceHue"), 64)
	accentHue, _ := strconv.ParseFloat(q.Get("accentHue"), 64)

	cfg := ThemeConfig{SurfaceHue: surfaceHue, AccentHue: accentHue}

	if err := rewriteColorsGo(cfg); err != nil {
		jsonError(w, "rewrite colors.go: "+err.Error())
		return
	}

	SetConfig(cfg)
	build()

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func rewriteColorsGo(cfg ThemeConfig) error {
	data, err := os.ReadFile("colors.go")
	if err != nil {
		return err
	}

	content := string(data)

	content = regexp.MustCompile(`(SurfaceHue:\s*)\d+`).ReplaceAllString(content, fmt.Sprintf("${1}%d", int(cfg.SurfaceHue)))
	content = regexp.MustCompile(`(AccentHue:\s*)\d+`).ReplaceAllString(content, fmt.Sprintf("${1}%d", int(cfg.AccentHue)))

	themeVar := fmt.Sprintf("\t\"gg-dark\",\n\tGenerate(DefaultConfig, true),\n\t\"gg-light\",\n\tGenerate(DefaultConfig, false),\n")
	if !strings.Contains(content, "DefaultConfig") {
		return fmt.Errorf("colors.go does not contain DefaultConfig")
	}

	_ = themeVar
	return os.WriteFile("colors.go", []byte(content), 0o644)
}

func jsonError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
