package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ThemeData struct {
	Name string
	Palette
}

var funcMap = template.FuncMap{
	"hex": func(s string) string { return strings.TrimPrefix(s, "#") },
}

var themes = []ThemeData{
	{"gg-dark", Dark},
	{"gg-light", Light},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gg-theme <build|link>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "build":
		build()
	case "link":
		link()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func build() {
	type target struct {
		tmpl string
		out  func(string) string
	}

	perTheme := []target{
		{"templates/ghostty.tmpl", func(n string) string { return filepath.Join("build", "ghostty", n) }},
		{"templates/fish.tmpl", func(n string) string { return filepath.Join("build", "fish", n+".theme") }},
		{"templates/vscode/theme.json.tmpl", func(n string) string { return filepath.Join("build", "vscode", "gg-theme", "themes", n+".json") }},
		{"templates/typora.css.tmpl", func(n string) string { return filepath.Join("build", "typora", n+".css") }},
	}

	for _, t := range perTheme {
		tmpl := template.Must(template.New(filepath.Base(t.tmpl)).Funcs(funcMap).ParseFiles(t.tmpl))
		for _, theme := range themes {
			render(tmpl, t.out(theme.Name), theme)
		}
	}

	tmpl := template.Must(template.New("package.json.tmpl").Funcs(funcMap).ParseFiles("templates/vscode/package.json.tmpl"))
	render(tmpl, filepath.Join("build", "vscode", "gg-theme", "package.json"), nil)
}

func render(tmpl *template.Template, path string, data any) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  \033[32m✓\033[0m %s\n", path)
}

func link() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	links := [][2]string{
		{"build/ghostty/gg-dark", filepath.Join(home, ".config", "ghostty", "themes", "gg-dark")},
		{"build/ghostty/gg-light", filepath.Join(home, ".config", "ghostty", "themes", "gg-light")},
		{"build/fish/gg-dark.theme", filepath.Join(home, ".config", "fish", "themes", "gg-dark.theme")},
		{"build/fish/gg-light.theme", filepath.Join(home, ".config", "fish", "themes", "gg-light.theme")},
		{"build/vscode/gg-theme", filepath.Join(home, ".vscode", "extensions", "gg-theme")},
		{"build/vscode/gg-theme", filepath.Join(home, ".cursor", "extensions", "gg-theme")},
		{"build/typora/gg-dark.css", filepath.Join(home, "Library", "Application Support", "abnerworks.Typora", "themes", "gg-dark.css")},
		{"build/typora/gg-light.css", filepath.Join(home, "Library", "Application Support", "abnerworks.Typora", "themes", "gg-light.css")},
	}

	if err := linkAll(cwd, links); err != nil {
		log.Fatal(err)
	}
}

func linkAll(cwd string, links [][2]string) error {
	type target struct {
		src string
		dst string
	}

	targets := make([]target, 0, len(links))
	for _, link := range links {
		src := filepath.Join(cwd, link[0])
		if _, err := os.Stat(src); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("missing build output: %s", src)
			}
			return err
		}
		targets = append(targets, target{src: src, dst: link[1]})
	}

	for _, target := range targets {
		if err := os.MkdirAll(filepath.Dir(target.dst), 0o755); err != nil {
			return err
		}
		if err := os.RemoveAll(target.dst); err != nil {
			return err
		}
		if err := os.Symlink(target.src, target.dst); err != nil {
			return err
		}
		fmt.Printf("  \033[36m→\033[0m %s\n", target.dst)
	}

	return nil
}
