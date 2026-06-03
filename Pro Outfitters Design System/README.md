# Pro Outfitters — Design System

An editorial-minimal design system for a premium website rebuild of **Pro Outfitters**, an
Orvis-endorsed fly fishing and bird hunting outfitter in Montana, established 1970.
Owner-operated lodges, guided river float trips, and upland wingshooting.

This repository is a **pitch deliverable**: a token-driven design system plus a fully realized
homepage mockup that sells the vision. It is structured to translate cleanly into a future
Go + templ + htmx + Tailwind rebuild — token-first, component-minded, near-zero JS.

---

## The company

- **Who:** Pro Outfitters — Orvis-endorsed fly fishing + bird hunting outfitter, Montana, est. 1970.
- **What they sell:** Guided river float trips, upland wingshooting, and owner-operated lodge stays.
  Three properties: **North Fork Crossing Lodge**, **Sharptail Lodge**, and **The Yurt at Craig**.
- **Who buys:** Affluent sporting travelers booking $3k–$10k+ guided lodge experiences. They value
  authenticity, conservation ethic, and earned access over flash. Discerning, not loud.
- **Positioning:** "An authentic taste of Montana before." Fifty-plus years of earned access,
  owner-operated, conservation-minded (Trout Unlimited partnership).

## Design direction — editorial-minimal

The photography carries the emotion (Montana rivers, golden-hour drift boats, dogs on point, big
sky). The layout stays calm and gets out of the way. **Restraint is the whole point.**

Deliberately **avoided**: rugged-lodge clichés — faux-wood textures, distressed leather, antler
motifs, rope borders. Sophistication comes from space and type, not props.

- Generous whitespace; asymmetric editorial layouts.
- Large, confident serif headlines; highly readable sans body at a comfortable measure.
- Letter-spaced uppercase eyebrows/labels.
- Near-monochrome paper-and-ink palette, warmed for Montana light.
- One restrained functional accent (muted slate, drawn from river/sky) + the heritage **antique
  gold** of the logo, both used sparingly.

---

## Sources provided

- `uploads/logo-whiteprimaryboxed (1).png` — primary boxed wordmark (black on white, gold rule
  frame). Copied to `assets/logo-primary-boxed.png`.
- `uploads/PRO-secondary-gold-and-white-e1610751463192 (1).png` — secondary gold "PRO" mark.
  Copied to `assets/logo-secondary-gold.png`.

No codebase, Figma, or existing site was provided — the visual system is derived from the two logo
assets and the written brief. The exact brand gold was sampled from the secondary mark: **#A58F58**.

---

## Index — what's in this folder

| Path | What it is |
|---|---|
| `README.md` | This file — context, content + visual foundations, iconography, manifest |
| `colors_and_type.css` | All design tokens: color, type, spacing, radii, shadow as CSS custom properties + semantic element styles |
| `SKILL.md` | Agent Skill manifest for reuse in Claude Code |
| `assets/` | Logos and imagery used across the system |
| `preview/` | Design System tab cards (colors, type, spacing, components) |
| `ui_kits/website/` | Website UI kit — reusable JSX components + interactive `index.html` |
| `Homepage Mockup.html` | The full homepage mockup applying the system end-to-end |

### `ui_kits/website/` contents
- `index.html` — interactive demo: full page with a working 3-step "Plan Your Trip" inquiry modal.
- `ui.css` — kit-level component styles layered on the tokens.
- `Primitives.jsx` — Eyebrow, GoldRule, Button, LinkArrow, Photo, Wordmark, SectionHead.
- `Sections.jsx` — Nav (sticky/over-hero state), Hero, ExperienceCard, LodgeCard, StatBand, PullQuote.
- `InquiryModal.jsx` — the multi-step inquiry flow (the kit's one real interaction).

### `preview/` — Design System tab cards
Colors (neutrals, ink, accent, gold, semantic), Type (display, headings, body, eyebrow, quote),
Spacing (scale, radii, elevation), Components (buttons, section header, card, stats, nav), and
Brand (primary logo, gold mark, wordmark lockup).

---

## CONTENT FUNDAMENTALS

How copy is written for Pro Outfitters. The voice is **confident, understated, and place-rooted** —
a little literary, never salesy.

**Tone & vibe.** Calm authority earned over fifty years. The brand never shouts. It assumes the
reader is discerning and treats them that way. Think a well-read river guide, not a brochure.

**Person.** Mostly third-person and plural-first ("the waters we fish," "our guides"). Addresses the
reader as **"you"** at moments of invitation ("Plan your trip"). Avoids breathless first-person hype.

**Casing.** Sentence case for body and headlines. **UPPERCASE, letter-spaced** reserved for eyebrows
and labels only ("EST. 1970", "FLY FISHING", "THE LODGES"). Never all-caps a full sentence.

**Sentence rhythm.** Short, declarative. Fragments are allowed for cadence. One idea per line in
headlines. Body sentences run long enough to breathe but never purple.

**What to avoid:** superlative-stuffing ("the most incredible," "world-class," "ultimate"),
adventure-brochure language ("epic," "adrenaline," "bucket-list"), exclamation marks, and emoji.
No emoji, ever.

**Examples (house voice):**
- Headline: *"Montana, the way it was meant to be fished."*
- Eyebrow: `EST. 1970 · ORVIS-ENDORSED`
- Supporting line: *"Owner-operated since 1970. Guided float trips, upland wingshooting, and three
  lodges on water worth protecting."*
- Story: *"We've spent fifty years earning access to water most people never see — and learning
  when to leave it alone."*
- Stewardship: *"You can't sell a river twice. We fish like we intend to keep it."*
- CTA: *"Plan your trip"* / *"Begin an inquiry"* (never "Book now!!!").

---

## VISUAL FOUNDATIONS

**Palette.** Near-monochrome paper-and-ink, warmed for Montana light. A warm bone off-white
(`--bg` #FAF7F1) ground; a deep warm near-black ink (`--ink` #1A1815) for text; warm grays for
secondary text and hairline borders. **One functional accent** — a muted desaturated slate
(`--accent` #3F6175), drawn from river water under big sky — used sparingly for links, focus, and
quiet emphasis. The **heritage antique gold** (`--gold` #A58F58) from the logo appears only as fine
rules, the EST mark, and occasional eyebrow accents. No bright or playful color anywhere.

**Typography.** A literary editorial serif — **Newsreader** — carries display and headlines; it has
real character at large sizes and a calm bookish texture in pull-quotes. Body and UI run in
**Libre Franklin**, a neutral, highly readable American grotesque. Eyebrows and labels are Libre
Franklin uppercase, letter-spaced ~0.18em. Modular scale on a major-third feel; body 18px at a
~66-character measure with generous leading (1.65).

**Spacing.** 4px base unit. Generous vertical rhythm — section padding is large (96–160px) to let
content breathe. Asymmetric grids: 12-column thinking with intentional empty columns rather than
edge-to-edge fills.

**Backgrounds.** Full-bleed photography is the emotional engine — Montana rivers, golden-hour drift
boats, dogs on point, big sky. Imagery skews **warm and natural**, golden-hour, never oversaturated
or filtered; a faint warm grade ties it to the paper. No textures, gradients-as-decoration,
patterns, or props. Backgrounds are paper, ink, or a photograph — nothing in between. A subtle
dark-to-transparent **protection gradient** sits over hero/footer imagery so type stays legible.

**Borders & cards.** Cards are gallery-like: defined by a **hairline border** (`--border` #E4DDD0)
and whitespace, not by shadow or heavy rounding. Corner radius is near-zero — **2px** at most;
much of the system is square-cornered for editorial crispness.

**Shadow & elevation.** Almost none. Elevation is communicated through space and hairlines, not
drop-shadows. A single soft shadow (`--shadow-soft`) exists for the rare floating element (sticky
nav on scroll); it is barely perceptible.

**Animation.** Restrained. Gentle opacity/transform fades on scroll-in (≤500ms,
`cubic-bezier(0.22,1,0.36,1)` ease-out). No bounces, no parallax theatrics, no infinite loops.
Motion respects `prefers-reduced-motion`.

**Hover & press.** Links and text buttons shift to the accent or gold and reveal a thin underline.
Primary buttons darken slightly and the gold/ink inverts on hover; no scaling. Press states use a
subtle 1px downward nudge or a slightly deeper tone — never a cartoon shrink. Image cards reveal a
caption and a faint zoom (scale 1.03) under the protection gradient.

**Transparency & blur.** Used only for the sticky nav, which goes from transparent over the hero to
a translucent paper backdrop (`backdrop-filter: blur(10px)`) once scrolled. Otherwise surfaces are
opaque.

**Layout rules.** Persistent light nav with an always-visible primary CTA ("Plan Your Trip").
Max content width ~1280px; text columns capped near 66ch. Eyebrow → headline → supporting line is
the recurring section-header pattern.

---

## ICONOGRAPHY

Icons are **rare and quiet** in this system — editorial-minimal means type and space do the work, so
icons appear only where they genuinely aid navigation or contact. There is no brand icon font and no
proprietary icon set in the source material (only the two logos).

- **Set:** [Lucide](https://lucide.dev) — thin **1.5px** open-stroke line icons, rounded joins.
  Their restrained hairline weight matches the editorial feel far better than filled or heavy icons.
  Loaded from CDN (`lucide@latest`); no icons were committed to `assets/`.
  *Substitution flag:* the brand provided no icons, so Lucide is a chosen default — swap if a
  house set emerges.
- **Where used:** directional affordances (`arrow-up-right`, `chevron-down`), contact/footer
  (`phone`, `mail`, `map-pin`), social (`instagram`, `facebook`), and the inquiry form. Sized
  16–20px, stroked in `--ink` or `--gray`, never filled, never colored except on hover (→ accent).
- **No emoji. No unicode glyph icons.** The only decorative mark is the fly-hook detail living
  inside the logo's "R" — it is part of the wordmark and is never extracted or reused as an icon.
- **The gold rule** (a thin hairline, sometimes a short centered tick like the logo's `EST 1970`
  frame) functions as the brand's signature non-icon ornament — use it instead of decorative icons.

### Brand assets in `assets/`
- `logo-primary-boxed.png` — primary boxed wordmark, black on white with the gold rule frame.
- `logo-secondary-gold.png` — secondary gold "PRO" mark (for dark grounds / small lockups).
- **Wordmark lockup (CSS):** for inline use in nav/footer the wordmark is reproduced as a small
  HTML/CSS lockup — "PRO OUTFITTERS" in the display face with a hairline gold rule and the
  `EST. 1970` tick — so it inherits `currentColor` and works on paper or dark grounds. See the
  `Wordmark` component in `ui_kits/website/`.
