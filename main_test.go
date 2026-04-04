package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLinkAllFailsBeforeCreatingLinksWhenAnySourceIsMissing(t *testing.T) {
	cwd := t.TempDir()
	existingSrc := filepath.Join(cwd, "build", "ghostty", "gg-dark")
	if err := os.MkdirAll(filepath.Dir(existingSrc), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(existingSrc, []byte("theme"), 0o644); err != nil {
		t.Fatal(err)
	}

	dstRoot := t.TempDir()
	firstDst := filepath.Join(dstRoot, "ghostty", "gg-dark")
	secondDst := filepath.Join(dstRoot, "ghostty", "gg-light")

	err := linkAll(cwd, [][2]string{
		{"build/ghostty/gg-dark", firstDst},
		{"build/ghostty/gg-light", secondDst},
	})
	if err == nil {
		t.Fatal("expected missing build output error")
	}
	if !strings.Contains(err.Error(), "missing build output") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Lstat(firstDst); !os.IsNotExist(err) {
		t.Fatalf("expected no link to be created, got %v", err)
	}
}

func TestLinkAllReplacesExistingDirectoryWithSymlink(t *testing.T) {
	cwd := t.TempDir()
	src := filepath.Join(cwd, "build", "vscode", "gg-theme")
	if err := os.MkdirAll(filepath.Join(src, "themes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	dstRoot := t.TempDir()
	dst := filepath.Join(dstRoot, ".vscode", "extensions", "gg-theme")
	if err := os.MkdirAll(filepath.Join(dst, "stale"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dst, "stale", "old.txt"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := linkAll(cwd, [][2]string{{"build/vscode/gg-theme", dst}}); err != nil {
		t.Fatal(err)
	}

	info, err := os.Lstat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected symlink, got mode %v", info.Mode())
	}

	target, err := os.Readlink(dst)
	if err != nil {
		t.Fatal(err)
	}
	if target != src {
		t.Fatalf("expected symlink to %q, got %q", src, target)
	}

	content, err := os.ReadFile(filepath.Join(dst, "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "{}" {
		t.Fatalf("unexpected linked file contents: %q", content)
	}
}
