package themes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gjermundgaraba/theme/colors"
)

func TestBuildWithConfigDoesNotMutateGlobalConfigAndUsesTyporaPalettes(t *testing.T) {
	t.Chdir(t.TempDir())
	old := colors.DefaultConfig
	defer colors.SetConfig(old)

	cfg := colors.ThemeConfig{SurfaceHue: 210, AccentHue: 45}
	if err := BuildWithConfig(cfg); err != nil {
		t.Fatal(err)
	}
	if colors.DefaultConfig != old {
		t.Fatalf("DefaultConfig mutated to %+v, want %+v", colors.DefaultConfig, old)
	}

	darkCSS, err := os.ReadFile(filepath.Join("build", "typora", "gg-dark.css"))
	if err != nil {
		t.Fatal(err)
	}
	lightCSS, err := os.ReadFile(filepath.Join("build", "typora", "gg-light.css"))
	if err != nil {
		t.Fatal(err)
	}
	darkBG := colors.Generate(cfg, true).BG
	lightBG := colors.Generate(cfg, false).BG
	if !strings.Contains(string(darkCSS), darkBG) {
		t.Fatalf("Typora dark CSS missing dark background %s", darkBG)
	}
	if strings.Contains(string(darkCSS), lightBG) {
		t.Fatalf("Typora dark CSS contains light background %s", lightBG)
	}
	if !strings.Contains(string(lightCSS), lightBG) {
		t.Fatalf("Typora light CSS missing light background %s", lightBG)
	}
	if strings.Contains(string(lightCSS), darkBG) {
		t.Fatalf("Typora light CSS contains dark background %s", darkBG)
	}
}

func TestBuildOutputsValidVSCodeJSON(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := Build(); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join("build", "vscode", "gg-theme", "package.json"),
		filepath.Join("build", "vscode", "gg-theme", "themes", "gg-dark.json"),
		filepath.Join("build", "vscode", "gg-theme", "themes", "gg-light.json"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var decoded map[string]any
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("%s is not valid JSON: %v", path, err)
		}
	}
}

func TestBuildProducesEveryInstallLinkSource(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := Build(); err != nil {
		t.Fatal(err)
	}
	for _, link := range installLinks(t.TempDir()) {
		if _, err := os.Stat(link[0]); err != nil {
			t.Fatalf("install source %s missing after build: %v", link[0], err)
		}
	}
}

func TestBuildRemovesStaleTyporaFolderAssets(t *testing.T) {
	t.Chdir(t.TempDir())
	stale := filepath.Join("build", "typora", "gg-dark", "stale-asset.txt")
	if err := os.MkdirAll(filepath.Dir(stale), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Build(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("expected stale Typora asset to be removed, got %v", err)
	}
	if _, err := os.Stat(filepath.Join("build", "typora", "gg-dark", "codeblock.css")); err != nil {
		t.Fatalf("expected rendered Typora dark CSS after cleanup: %v", err)
	}
	if _, err := os.Stat(filepath.Join("build", "typora", "gg-light", "codeblock.css")); err != nil {
		t.Fatalf("expected rendered Typora light CSS after cleanup: %v", err)
	}
	if _, err := os.Stat(filepath.Join("build", "palette.md")); err != nil {
		t.Fatalf("expected generated palette markdown report: %v", err)
	}
	if _, err := os.Stat(filepath.Join("build", "palette.json")); err != nil {
		t.Fatalf("expected generated palette JSON report: %v", err)
	}
}
