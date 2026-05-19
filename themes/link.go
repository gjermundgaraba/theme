package themes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Link symlinks built themes to their install locations.
// Existing target paths are removed before symlinks are created.
func Link() error {
	return LinkWithOptions(LinkOptions{})
}

// LinkOptions controls installation behavior.
type LinkOptions struct {
	DryRun bool
}

// LinkWithOptions symlinks built themes to their install locations.
func LinkWithOptions(opts LinkOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	links := installLinks(home)
	return LinkAllWithOptions(cwd, links, opts)
}

func installLinks(home string) [][2]string {
	return [][2]string{
		{"build/ghostty/gg-dark", filepath.Join(home, ".config", "ghostty", "themes", "gg-dark")},
		{"build/ghostty/gg-light", filepath.Join(home, ".config", "ghostty", "themes", "gg-light")},
		{"build/fish/gg-dark.theme", filepath.Join(home, ".config", "fish", "themes", "gg-dark.theme")},
		{"build/fish/gg-light.theme", filepath.Join(home, ".config", "fish", "themes", "gg-light.theme")},
		{"build/vscode/gg-theme", filepath.Join(home, ".vscode", "extensions", "gg-theme")},
		{"build/vscode/gg-theme", filepath.Join(home, ".cursor", "extensions", "gg-theme")},
		{"build/zed/gg-theme.json", filepath.Join(home, ".config", "zed", "themes", "gg-theme.json")},
		{"build/typora/gg-dark.css", filepath.Join(home, "Library", "Application Support", "abnerworks.Typora", "themes", "gg-dark.css")},
		{"build/typora/gg-dark", filepath.Join(home, "Library", "Application Support", "abnerworks.Typora", "themes", "gg-dark")},
		{"build/typora/gg-light.css", filepath.Join(home, "Library", "Application Support", "abnerworks.Typora", "themes", "gg-light.css")},
		{"build/typora/gg-light", filepath.Join(home, "Library", "Application Support", "abnerworks.Typora", "themes", "gg-light")},
		{"build/neovim/colors/gg-dark.lua", filepath.Join(home, ".config", "nvim", "colors", "gg-dark.lua")},
		{"build/neovim/colors/gg-light.lua", filepath.Join(home, ".config", "nvim", "colors", "gg-light.lua")},
	}
}

// LinkAll creates symlinks for each (src, dst) pair relative to cwd.
// It validates all sources before mutating destinations, then removes each destination before linking it.
func LinkAll(cwd string, links [][2]string) error {
	return LinkAllWithOptions(cwd, links, LinkOptions{})
}

// LinkAllWithOptions creates or previews symlinks for each (src, dst) pair relative to cwd.
func LinkAllWithOptions(cwd string, links [][2]string, opts LinkOptions) error {
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
		if opts.DryRun {
			fmt.Printf("  \033[36m?\033[0m %s -> %s\n", target.dst, target.src)
			continue
		}
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
