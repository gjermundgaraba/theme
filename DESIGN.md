---
name: gg-theme
description: A local OKLCH theme generator for terminal, shell, editor, writing, and Neovim targets.
palette: build/palette.md
typography:
  display:
    fontFamily: "Space Grotesk, system-ui, sans-serif"
    fontSize: "2.5rem"
    fontWeight: 400
    lineHeight: "2.75rem"
    letterSpacing: "-1.5px"
  headline:
    fontFamily: "Space Grotesk, system-ui, sans-serif"
    fontSize: "1.63rem"
    fontWeight: 700
    lineHeight: "1.875rem"
    letterSpacing: "-1px"
  body:
    fontFamily: "Satoshi, system-ui, sans-serif"
    fontSize: "16px"
    fontWeight: 400
    lineHeight: "1.625rem"
  mono:
    fontFamily: "Berkeley Mono, ui-monospace, monospace"
    fontSize: "12.5px"
    fontWeight: 400
    lineHeight: 1.65
rounded:
  xs: "3px"
  sm: "4px"
  md: "6px"
  lg: "8px"
spacing:
  sm: "8px"
  md: "12px"
  lg: "16px"
  xl: "24px"
components:
  editor-surface:
    backgroundColor: "{Palette.BG}"
    textColor: "{Palette.FG}"
    rounded: "{rounded.sm}"
    padding: "0"
  editor-selection:
    backgroundColor: "{Palette.SelectionBG}"
    textColor: "{Palette.SelectionFG}"
    rounded: "{rounded.xs}"
    padding: "0"
  action-fill:
    backgroundColor: "{Palette.ButtonBG}"
    textColor: "{Palette.ButtonFG}"
    rounded: "{rounded.sm}"
    padding: "0"
  typora-prose:
    backgroundColor: "{Palette.BG}"
    textColor: "{Palette.FG}"
    typography: "{typography.body}"
    padding: "20px"
---

# Design System: gg-theme

## 1. Overview

**Creative North Star: "The Local Color Lab"**

gg-theme is a generated theme system. The product is the theme itself: the shared OKLCH palette contract, the target-native mappings, and the installed results in Ghostty, Fish, VS Code, Cursor, Neovim, and Typora. The optional browser tool is only a local inspector for the contract.

The physical scene is a developer living inside their own tools, tuning the same color atmosphere across terminal, shell, editor, and writing surfaces. The theme should feel precise, colorful, developer-native, and personal. It should make a terminal prompt, a VS Code diff, a Neovim buffer, and a Typora document feel like one environment without flattening every target into the same skin.

The color strategy is Full palette. Tinted neutral surfaces hold the working field. Syntax, ANSI, diagnostics, search, selection, and action tokens provide deliberate color range. The system rejects generic SaaS dashboard chrome, drab monochrome developer utility, decorative UI that obscures color judgment or generated output, control-panel sprawl, and documentation that makes the optional web tool appear more important than the generated themes.

**Key Characteristics:**

- Theme-first: generated targets are the primary product surface.
- Token-governed: every color decision traces back to `colors.Palette`.
- Target-native: each app receives mappings that fit its own rendering model.
- Color-forward: syntax and ANSI hues are vivid enough to be personal, never arbitrary.
- Contrast-safe: readability checks are generator correctness.

## 2. Colors

The default palette is an Ink Green Field with a Raspberry Accent. The surface hue gives both dark and light themes a green-tinted neutral base; the accent hue gives action and active state a raspberry signal. Current values live in `colors.DefaultConfig` and the generated `build/palette.md` report.

### Primary

- **Raspberry Accent** (`Accent`): readable accent text for active states, links, Typora primary color, and generated theme emphasis. It is not the same token as filled action color.
- **Raspberry Action Fill** (`ButtonBG`, `ButtonFG`): filled action, badge, remote, progress, and selected-control roles. Always pair `ButtonBG` with generated `ButtonFG`.

### Secondary

- **Terminal Spectrum** (`Red`, `Orange`, `Yellow`, `Green`, `Cyan`, `Blue`, `Magenta`): syntax, ANSI, git state, diagnostics, search, shell roles, and code preview roles.
- **Bright Terminal Spectrum** (`BrightRed`, `BrightGreen`, `BrightYellow`, `BrightBlue`, `BrightMagenta`, `BrightCyan`, `BrightWhite`, `BrightBlack`): elevated terminal emphasis for bright ANSI slots. These should feel more energetic than base syntax colors, not like separate brand colors.

### Neutral

- **Ink Green Field** (`BG` dark): the dark terminal, editor, Neovim, and Typora base. It must remain tinted and must never become pure black.
- **Raised Ink Surface** (`Surface` dark): editor sidebars, widgets, title bars, fold regions, Typora side panels, and grouped shell/editor surfaces.
- **Moss Overlay Line** (`Overlay` dark): borders, dividers, line numbers, ghost text, indent guides, and inactive structure.
- **Mist Text** (`FG`, `FGDim` dark): primary and secondary foregrounds for dark targets.
- **Pale Green Paper** (`BG` light): the light generated base, matched to the same surface hue.
- **Soft Paper Surface** (`Surface` light): sidebars, panels, widgets, and grouped editor surfaces in light targets.
- **Leaf Grey Structure** (`Overlay` light): visible but subordinate structure in light targets.
- **Ink Text** (`FG`, `FGDim` light): primary and secondary foregrounds for light targets.

### Named Rules

**The Palette Contract Rule.** If a generated target needs a color, map it to a named `colors.Palette` token or add a token to the generator. Raw target-local color drift is prohibited.

**The Tinted Neutral Rule.** Pure black, pure white, and untinted grey are forbidden. Every neutral must keep a relationship to `SurfaceHue`.

**The Contrast Is Correctness Rule.** `colors.ContrastChecks()` is part of the theme contract. A beautiful token that fails its required pair is not shippable.

## 3. Typography

**Display Font:** Space Grotesk with system sans fallback.  
**Body Font:** Satoshi with system sans fallback.  
**Label/Mono Font:** Berkeley Mono for code and generated output.

**Character:** Typography matters most in Typora and code surfaces. Typora gets the editorial voice: Space Grotesk headings, Satoshi prose, Berkeley Mono code. Terminals and editors inherit platform-native UI typography but receive the same palette and syntax semantics.

### Hierarchy

- **Display** (400, `2.5rem`, `2.75rem`): Typora h1 and document-level identity.
- **Headline** (700, `1.63rem`, `1.875rem`): Typora h2 and strong prose sections.
- **Title** (700, roughly `1.17rem` to `1.63rem`): Typora h3 and h4 hierarchy.
- **Body** (400, `16px`, `1.625rem`): Typora prose, lists, tables, and long reading surfaces. Keep prose close to 65 to 75 characters where the target permits it.
- **Label** (target-native): editor, terminal, and shell labels should use the host app vocabulary rather than inventing a custom display layer.
- **Mono** (400 to 600, `12.5px` equivalent where controlled): code blocks, inline code, generated output, and source mode.

### Named Rules

**The Target Owns Type Rule.** The theme may style Typora typography directly. Editors, terminals, and shell targets receive color mappings without pretending to own their UI font stack.

**The Code Must Stay Exact Rule.** Code surfaces prioritize spacing, contrast, and token distinction over decorative type treatments.

## 4. Elevation

Elevation is layered but low drama. The theme uses tonal layering first: `BG`, `Surface`, `Overlay`, selection fills, search fills, visual mode fills, and diff fills. Shadows exist only where a target expects them, primarily Typora menus and small document affordances.

### Shadow Vocabulary

- **Typora Menu Shadow** (`4px 4px 20px rgba(BG, 0.79)`): dark target-native menu depth. It is not a global elevation style.
- **Typora Inset Mark** (`inset 0 -1px 0 rgba(BG, 0.25)`): inline document structure, subtle enough to read as typography support.
- **Editor Shadow Proxy** (`Black` with alpha suffixes in VS Code): widget and scrollbar shadow roles should use palette-derived black, not raw black.

### Named Rules

**The Tonal Layer Rule.** Reach for `Surface`, `Overlay`, and state fills before shadows. Theme depth should read as palette structure.

**The Target-Native Depth Rule.** A shadow is allowed only when the target has a native role for it. Do not invent decorative elevation in generated themes.

## 5. Components

In this system, components are generated theme surfaces, not web widgets. Document the rendered roles that users actually live in: terminal palettes, editor chrome, syntax, selections, diffs, diagnostics, and Typora prose.

### Terminal ANSI

- **Shape:** no shape. Terminal color is a 16-slot transport.
- **Primary:** `BG`, `FG`, `Cursor`, `CursorText`, `SelectionBG`, and `SelectionFG` define the terminal field.
- **Base ANSI:** `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White` map directly to Ghostty palette slots `0` through `7`.
- **Bright ANSI:** bright tokens map directly to slots `8` through `15` and must remain distinct from base slots.
- **Rule:** terminal colors should preserve semantic separation at small sizes and low font weights.

### Shell Syntax

- **Style:** Fish maps command language onto the same syntax spectrum: commands to `Cyan`, keywords to `Magenta`, quotes to `Yellow`, errors to `Red`, autosuggestions to `FGDim`, and selections to `SelectionBG`.
- **State:** shell feedback should feel like the terminal theme, not a separate shell skin.

### Editor Chrome

- **Background:** editor base uses `BG`; sidebars, title bars, widgets, suggestions, status bar, tabs, and dropdowns use `Surface`.
- **Structure:** borders, indent guides, whitespace marks, rulers, inactive tabs, and ghost text use `Overlay` with target-native alpha where needed.
- **Action:** badges, progress, remote status, and buttons use `ButtonBG` and `ButtonFG`.
- **Focus:** focus borders and active indicators use syntax tokens such as `Blue` or `Accent` according to target conventions.

### Syntax

- **Keywords and operators:** `Magenta`.
- **Strings and search alternates:** `Yellow`.
- **Numbers and structural editor focus:** `Blue`.
- **Functions and success state:** `Green`.
- **Types and links:** `Cyan`.
- **Parameters and warnings:** `Orange`.
- **Errors and destructive state:** `Red`.
- **Comments and inactive metadata:** `FGDim` or `Overlay`, depending on density.

### Selection, Search, and Visual State

- **Selection:** `SelectionBG` with `SelectionFG`, validated for text contrast.
- **Search:** `SearchBG` carries search state without overpowering syntax.
- **Visual mode:** `VisualBG` uses the surface hue so visual selections feel connected to the neutral field.
- **Rule:** state backgrounds must be visible without washing out foreground text.

### Diffs and Diagnostics

- **Diff add:** `DiffAddBG` plus `Green` gutters or markers.
- **Diff delete:** `DiffDeleteBG` plus `Red` gutters or markers.
- **Diff change:** `DiffChangeBG` plus `Blue` or `Orange` according to target convention.
- **Diagnostics:** errors use `Red`, warnings use `Orange`, info uses `Cyan`, hints use `Green`.
- **Debug:** `DebugFG` exists because debug status surfaces need a dedicated readable foreground.

### Typora Document

- **Background:** `BG` with `Surface` sidebars and panels.
- **Text:** prose uses `FG` and secondary text uses `FGDim`.
- **Headings:** Space Grotesk creates the document voice without changing the shared color system.
- **Code:** Berkeley Mono uses the same syntax family as editor targets.
- **Mermaid and source mode:** must consume palette tokens and preserve contrast. Raw greys and unmanaged rgba values are prohibited.

## 6. Do's and Don'ts

### Do:

- **Do** treat the generated theme as the product. The web tool is secondary.
- **Do** keep `colors/` as the visual source of truth: `ThemeConfig`, `Generate()`, `Palette`, and contrast checks define the system.
- **Do** regenerate `build/` after changing `colors/` or `themes/templates/`.
- **Do** map every target role back to named tokens, including Typora Mermaid and source mode CSS.
- **Do** preserve target-native behavior: Ghostty wants ANSI slots, Fish wants shell roles, VS Code wants workbench roles, Neovim wants highlight groups, Typora wants document typography.
- **Do** validate contrast before increasing chroma or lowering lightness.
- **Do** keep `themes.Build()` and `themes.Link()` aligned when targets change.
- **Do** treat `themes.Link()` as destructive replacement of configured install paths.

### Don't:

- **Don't** make the optional web tool appear more important than the generated themes.
- **Don't** hand-edit files in `build/`.
- **Don't** introduce pure black, pure white, raw greys, or unmanaged raw rgba literals into templates.
- **Don't** create generic SaaS dashboard chrome.
- **Don't** make a drab monochrome developer utility.
- **Don't** add decorative UI that obscures color judgment or generated output.
- **Don't** let control-panel sprawl turn quick theme iteration into design-system maintenance.
- **Don't** use color as decoration when it does not explain syntax, state, contrast, or target behavior.
- **Don't** use side-stripe borders greater than `1px`, gradient text, decorative glass, hero metrics, repeated generic card grids, or modals as the first solution.
