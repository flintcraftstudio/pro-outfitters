---
name: pro-outfitters-design
description: Use this skill to generate well-branded interfaces and assets for Pro Outfitters — an Orvis-endorsed Montana fly fishing & bird hunting outfitter (est. 1970) — either for production or throwaway prototypes/mocks. Contains editorial-minimal design guidelines, color + type tokens, logo assets, and a website UI kit for prototyping.
user-invocable: true
---

Read the `README.md` file within this skill, and explore the other available files.

If creating visual artifacts (slides, mocks, throwaway prototypes, etc), copy assets out and create
static HTML files for the user to view. If working on production code, you can copy assets and read
the rules here to become an expert in designing with this brand.

If the user invokes this skill without any other guidance, ask them what they want to build or
design, ask some questions, and act as an expert designer who outputs HTML artifacts _or_ production
code, depending on the need.

## Fast orientation
- **Tokens:** `colors_and_type.css` — all color, type, spacing, radius, shadow custom properties
  plus semantic element classes. Import it and build on top.
- **Voice:** confident, understated, place-rooted, a little literary. Never salesy. No emoji. See
  the CONTENT FUNDAMENTALS section of the README.
- **Look:** editorial-minimal. Paper-and-ink near-monochrome warmed for Montana light; one muted
  slate accent + heritage antique gold, both sparing. Photography carries the emotion; layout stays
  calm. No rugged-lodge clichés (no wood/leather/antler/rope).
- **Type:** Newsreader (display serif) + Libre Franklin (UI/body sans). Letter-spaced uppercase
  eyebrows.
- **Components & full page:** `ui_kits/website/` (reusable JSX) and `Homepage Mockup.html`
  (token-driven semantic HTML — the reference application of the whole system).
- **Logos:** `assets/logo-primary-boxed.png`, `assets/logo-secondary-gold.png`; brand gold `#A58F58`.

## Non-negotiables
- Keep it calm and spacious. Restraint is the brand.
- Use the gold rule (thin hairline + EST tick), not decorative icons, as the signature ornament.
- Icons are rare, thin (Lucide 1.5px). No emoji, no unicode-glyph icons.
- Token-driven and component-minded so output translates to Go + templ + htmx + Tailwind.
