# AGENTS.md

This repo is a small Go CLI that renders shared theme data into Ghostty, Fish, and VS Code/Cursor outputs.

- Start with `main.go`, `colors.go`, and `templates/`.
- Validate normal changes with `go test ./...` and `go run . build`.
- `make build` and `make link` are thin wrappers around `go run . build` and `go run . link`.
- Treat `build/` as generated output. Regenerate it after changing `colors.go` or anything under `templates/`, and do not hand-edit files in `build/`.
- `go run . link` replaces the installed theme targets under the user home directory. Run it only when you intentionally want to update local Ghostty, Fish, VS Code, or Cursor theme installs.
- When adding or removing a target, keep the render list in `build()` and the symlink list in `link()` aligned.
