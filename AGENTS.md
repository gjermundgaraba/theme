# AGENTS.md

This repo is a Go CLI that renders shared theme data into Ghostty, Fish, VS Code/Cursor, Neovim, and Typora outputs.

- `colors/`: Palette type, ThemeConfig, Generate(), contrast functions, DefaultConfig.
- `themes/`: embedded templates, Build(), Link(), RenderTemplate().
- `theme-generator/`: optional interactive web editor (Serve()); do not treat it as the primary product surface.
- `cli/`: entry point dispatching link/serve.
- Validate with `go test ./...` or `make test`.
- Build outputs from the generator UI: `make serve`, open `localhost:9090`, click Apply. That rewrites `colors/palette.go` and renders all targets into `build/`.
- `make link` replaces installed theme targets under the user home directory.
- Treat `build/` as generated output. Re-apply via the UI after changing `colors/` or `themes/templates/`.
- When adding or removing a target, keep the render list in `Build()` and the symlink list in `Link()` aligned.
