/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: [
    "./internal/view/**/*.templ",
    // Some shared class strings live in Go const/vars (e.g. form control
    // classes in content.go). Scan the view package's .go files too, or
    // those utilities get purged.
    "./internal/view/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        // Pro Outfitters — editorial-minimal, paper-and-ink, warmed for
        // Montana light. Drives the public site. See colors_and_type.css.
        po: {
          paper:          "#FAF7F1",
          "paper-deep":   "#F3EEE4",
          surface:        "#FEFCF8",
          ink:            "#1A1815",
          "ink-soft":     "#423D36",
          gray:           "#6E665B",
          "gray-light":   "#9A9285",
          border:         "#E4DDD0",
          "border-strong":"#CFC6B6",
          accent:         "#3F6175",
          "accent-deep":  "#2E4B5C",
          "accent-wash":  "#DCE7ED",
          gold:           "#A58F58",
          "gold-deep":    "#8A7745",
          "on-dark":      "#F3EEE5",
          "on-dark-soft": "#C9C2B5",
          "ground-dark":  "#14120F",
        },
        // Legacy admin/auth palette (dark theme). Kept so the CMS admin
        // and login views render unchanged while the public site is
        // reskinned to the Pro Outfitters brand.
        ff: {
          dark:    "#0f1117",
          dark2:   "#161a24",
          panel:   "#13161f",
          panel2:  "#1a1e2b",
          dusk:    "#5b7ec4",
          "dusk-hover": "#6d8fd4",
          "dusk-mid": "#2a3d6e",
          "dusk-lo": "#1e2d52",
          ice:     "#a8c0e8",
          moon:    "#c8d5ed",
          paper:   "#edf0f7",
          cream:   "#dde2ed",
          ash:     "#7a8099",
          stone:   "#3d4459",
          border:  "rgba(255,255,255,0.07)",
          border2: "rgba(255,255,255,0.12)",
        },
      },
      fontFamily: {
        display: ['"Newsreader"', "Georgia", ...defaultTheme.fontFamily.serif],
        body:    ['"Libre Franklin"', ...defaultTheme.fontFamily.sans],
      },
      letterSpacing: {
        eyebrow: "0.18em",
      },
      maxWidth: {
        measure: "38rem",   // ~66ch comfortable text column
        container: "75rem", // 1200px content max
      },
      transitionTimingFunction: {
        // A gentle, symmetric ease (easeInOutSine). The old curve started at
        // y=1, snapping each property most of the way instantly — which read
        // as abrupt. This eases in and out evenly for a slow, "lazy river"
        // glide. Pair it with the longer durations below.
        editorial: "cubic-bezier(0.37, 0, 0.63, 1)",
      },
    },
  },
  plugins: [],
}
