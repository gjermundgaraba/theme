# Product

## Register

product

## Users

This is primarily for the owner of the theme system: a developer maintaining a personal color environment across terminal, shell, editor, writing, and Neovim targets.

The normal workflow is local and code-first: adjust the shared palette model, validate contrast, render all target themes, inspect generated output, and link the results into application config directories. The optional web generator exists to help explore hues and preview output, but it is not the primary product surface.

## Product Purpose

Generate coherent dark and light themes from one shared OKLCH-driven palette contract. The Go generator should be the source of truth for color decisions, and every supported target should render from that same token set without manual edits to generated files.

Success means the user can change a small number of palette inputs, regenerate all targets, verify contrast, and keep Ghostty, Fish, VS Code, Cursor, Typora, and Neovim visually aligned with low friction.

## Primary Surfaces

- **Color generator:** `colors/` defines `ThemeConfig`, `Palette`, `Generate()`, contrast helpers, and canonical dark and light palettes.
- **Theme renderer:** `themes/` embeds target templates and renders build output.
- **CLI:** `cli/` dispatches `build`, `link`, and `serve` commands.
- **Generated targets:** Ghostty, Fish, VS Code/Cursor, Typora, and Neovim consume the shared palette.
- **Optional web generator:** `theme-generator/` provides a local browser interface for hue exploration and output inspection.

## Brand Personality

Expressive, color-forward, playful, precise, and developer-native. The system should feel like a personal color instrument in code, not a generic admin product. Visual choices should support palette quality, target coherence, and trust in generated output.

## Anti-references

- Generic SaaS dashboard chrome.
- Drab monochrome developer utility.
- Decorative UI that obscures color judgment or generated output.
- Control-panel sprawl that turns quick theme iteration into design-system maintenance.
- Documentation that makes the optional web tool appear more important than the generated themes.

## Design Principles

- Treat the palette contract as the product center: every visual decision should be traceable to generated tokens.
- Keep generated outputs trustworthy: templates and build output must match the shared palette without hand edits.
- Treat contrast as core product quality: accessibility checks protect both UI readability and theme correctness.
- Prefer target-native mappings over decoration: each app should feel coherent with its platform while sharing the same color DNA.
- Preserve a personal, opinionated feel: the themes can be colorful and expressive without becoming arbitrary.
- Keep the optional web generator useful but secondary: it should inspect and apply the palette, not redefine the product.

## Accessibility & Inclusion

Target WCAG AA for relevant text and controls. Prioritize strong contrast validation, keyboard usability in the optional web generator, visible focus states, and reduced-motion-safe interactions. Because this product generates themes, contrast validation is both accessibility work and generator correctness.
