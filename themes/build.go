package themes

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gjermundgaraba/theme/colors"
	colorful "github.com/lucasb-eyer/go-colorful"
)

//go:embed templates
var templatesFS embed.FS

// ThemeData pairs a theme name with its palette for template rendering.
type ThemeData struct {
	Name string
	colors.Palette
}

var funcMap = template.FuncMap{
	"hex":  func(s string) string { return strings.TrimPrefix(s, "#") },
	"rgba": rgba,
}

func rgba(s string, alpha float64) string {
	c, err := colorful.Hex(s)
	if err != nil {
		return s
	}
	r, g, b := c.RGB255()
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, alpha)
}

func themeData() []ThemeData {
	return []ThemeData{
		{"gg-dark", colors.Dark},
		{"gg-light", colors.Light},
	}
}

func generatedThemeData(cfg colors.ThemeConfig) []ThemeData {
	return []ThemeData{
		{"gg-dark", colors.Generate(cfg, true)},
		{"gg-light", colors.Generate(cfg, false)},
	}
}

// Build renders all theme templates into build/.
func Build() error {
	return build(themeData())
}

// BuildWithConfig renders all theme templates with cfg without mutating colors.DefaultConfig.
func BuildWithConfig(cfg colors.ThemeConfig) error {
	return build(generatedThemeData(cfg))
}

func build(allThemes []ThemeData) error {
	type target struct {
		tmpl string
		out  func(string) string
	}

	perTheme := []target{
		{"templates/ghostty.tmpl", func(n string) string { return filepath.Join("build", "ghostty", n) }},
		{"templates/fish.tmpl", func(n string) string { return filepath.Join("build", "fish", n+".theme") }},
		{"templates/vscode/theme.json.tmpl", func(n string) string { return filepath.Join("build", "vscode", "gg-theme", "themes", n+".json") }},
		{"templates/neovim.lua.tmpl", func(n string) string { return filepath.Join("build", "neovim", "colors", n+".lua") }},
	}

	for _, t := range perTheme {
		tmpl, err := parseTemplate(t.tmpl)
		if err != nil {
			return err
		}
		for _, td := range allThemes {
			if err := render(tmpl, t.out(td.Name), td); err != nil {
				return err
			}
		}
	}

	tmpl, err := parseTemplate("templates/vscode/package.json.tmpl")
	if err != nil {
		return err
	}
	if err := render(tmpl, filepath.Join("build", "vscode", "gg-theme", "package.json"), nil); err != nil {
		return err
	}

	if err := renderPaletteReports(allThemes); err != nil {
		return err
	}

	// Typora folder-based themes.
	for _, td := range allThemes {
		typoraDir := filepath.Join("build", "typora", td.Name)
		if err := os.RemoveAll(typoraDir); err != nil {
			return err
		}
		typoraTemplates := []struct {
			tmpl string
			out  string
		}{
			{"templates/typora/theme.css.tmpl", filepath.Join("build", "typora", td.Name+".css")},
			{"templates/typora/codeblock.dark.css.tmpl", filepath.Join(typoraDir, "codeblock.css")},
			{"templates/typora/mermaid.dark.css.tmpl", filepath.Join(typoraDir, "mermaid.css")},
			{"templates/typora/sourcemode.dark.css.tmpl", filepath.Join(typoraDir, "sourcemode.css")},
		}
		for _, t := range typoraTemplates {
			tmpl, err := parseTemplate(t.tmpl)
			if err != nil {
				return err
			}
			if err := render(tmpl, t.out, td); err != nil {
				return err
			}
		}
		if err := copyEmbeddedDir("templates/typora/assets", typoraDir); err != nil {
			return err
		}
	}
	return nil
}

func parseTemplate(path string) (*template.Template, error) {
	data, err := templatesFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read embedded template %s: %w", path, err)
	}
	tmpl, err := template.New(filepath.Base(path)).Funcs(funcMap).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse embedded template %s: %w", path, err)
	}
	return tmpl, nil
}

func render(tmpl *template.Template, path string, data any) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	if err := atomicWriteFile(path, buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Printf("  \033[32m✓\033[0m %s\n", path)
	return nil
}

func copyEmbeddedDir(srcDir, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := fs.ReadDir(templatesFS, srcDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := templatesFS.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, e.Name())
		if err := atomicWriteFile(dstPath, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("  \033[32m✓\033[0m %s\n", dstPath)
	}
	return nil
}

func atomicWriteFile(path string, data []byte, perm fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
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

// RenderTemplateString renders a named embedded template file to a string.
func RenderTemplateString(tmplPath string, data ThemeData) (string, error) {
	t, err := parseTemplate(tmplPath)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderTemplate renders a named embedded template file to a string.
func RenderTemplate(tmplPath string, data ThemeData) string {
	s, err := RenderTemplateString(tmplPath, data)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return s
}
