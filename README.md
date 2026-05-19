# gg-theme

A personal Go theme generator for keeping terminal, shell, editor, writing, and Neovim colors aligned from one OKLCH palette model.

The generator is the product center. It renders dark and light themes from shared color tokens into target-specific templates. The optional web generator is a local helper for hue exploration and output inspection.

## Supported targets

- Ghostty
- Fish shell
- VS Code / Cursor
- Zed
- Typora dark and light themes
- Neovim

## Usage

```sh
go run ./cli serve
go run ./cli link --dry-run
# then, when the plan looks right:
go run ./cli link
```

Or use Make:

```sh
make serve
make link
make test
```

The generator UI is the only build entry point. Start it with `make serve`, adjust hues and lightness, then click Apply to rewrite `colors/palette.go` and render every target into `build/` (including `build/palette.md` and `build/palette.json`).

### Commands

- `go run ./cli serve`: start the local web generator at `localhost:9090`. Apply rewrites `colors/palette.go` and renders all targets.
- `go run ./cli link --dry-run`: preview install symlink replacements without mutating local config.
- `go run ./cli link`: replace installed target paths with symlinks to generated output.
- `go test ./...`: run generator, renderer, linker, and web helper tests.

## Install targets

`link` removes each listed target path before creating a symlink. Use it only when you intend to replace the generated theme installs.

- Ghostty: `~/.config/ghostty/themes/`
- Fish: `~/.config/fish/themes/`
- VS Code: `~/.vscode/extensions/gg-theme`
- Cursor: `~/.cursor/extensions/gg-theme`
- Zed: `~/.config/zed/themes/gg-theme.json`
- Typora: `~/Library/Application Support/abnerworks.Typora/themes/`
- Neovim: `~/.config/nvim/colors/`

## Structure

```text
colors/                 palette contract, OKLCH generation, contrast checks
  palette.go            Palette, ThemeConfig, DefaultConfig
  generate.go           dark and light palette generation
  contrast.go           contrast validation and palette dump helpers

themes/                 embedded templates, build, link, render helpers
  build.go              render list and template execution
  link.go               install symlink list
  templates/            Ghostty, Fish, VS Code, Zed, Neovim, Typora templates

theme-generator/        optional local web generator
  server.go             HTTP server and palette API
  web/index.html        browser UI for hue exploration and output preview

cli/                    command entry point
build/                  generated output, do not hand-edit
```

## Changing the palette

1. Run `make serve` and open `localhost:9090`.
2. Adjust hues and lightness. Click Apply to rewrite `colors/palette.go` and render all targets into `build/`.
3. Inspect generated output, starting with `build/palette.md`.
4. Run `go test ./...`.
5. Run `go run ./cli link --dry-run` to preview install changes.
6. Run `go run ./cli link` only when you intentionally want to update installed local themes.

For deeper changes (new tokens, generation behavior), edit `colors/generate.go` directly, then re-apply through the UI.

## Adding a target

1. Add a template under `themes/templates/`.
2. Add render output in `themes.Build()`.
3. Add install symlinks in `themes.Link()`.
4. Add tests or update existing target coverage.
5. Run `go test ./...` and re-apply via the generator UI to refresh `build/`.

Keep `themes.Build()` and `themes.Link()` aligned when targets change.
