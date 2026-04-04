# gg-theme

A personal theme generator that produces dark and light themes for multiple applications from a single set of color definitions.

## Themes

- **gg-dark** — based on Dracula
- **gg-light** — light grey variant with adjusted accent colors

## Supported targets

- Ghostty
- Fish shell
- VSCode / Cursor

## Usage

```sh
go run . build    # render all templates to build/
go run . install  # symlink build/ outputs to config dirs
```

### Install targets

| Target | Location |
|--------|----------|
| Ghostty | `~/.config/ghostty/themes/` |
| Fish | `~/.config/fish/themes/` |
| VSCode | `~/.vscode/extensions/gg-theme` |
| Cursor | `~/.cursor/extensions/gg-theme` |

## Structure

```
colors.go          palette definitions
main.go            build + install commands
templates/         go text/template files
  ghostty.tmpl
  fish.tmpl
  vscode/
    package.json.tmpl
    theme.json.tmpl
build/             generated output (gitignored)
```

## Adding a new target

1. Create a template in `templates/`
2. Add the template and output path to the `perTheme` list in `main.go`
3. Add symlink entries to the `install` function
