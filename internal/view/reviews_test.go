package view

import (
	"fmt"
	"strings"
	"testing"
)

// visibleAt reports whether a column or card with the given classes renders
// at a breakpoint (0 = phones, 1 = md, 2 = lg), per the Tailwind classes
// reviewColumns emits.
func visibleAt(classes string, bp int) bool {
	has := func(c string) bool {
		for _, f := range strings.Fields(classes) {
			if f == c {
				return true
			}
		}
		return false
	}
	switch {
	case has("md:hidden") && bp >= 1:
		return false
	case has("lg:hidden") && bp >= 2:
		return false
	case has("hidden") && has("md:flex"):
		return bp >= 1
	case has("hidden") && has("lg:flex"):
		return bp >= 2
	}
	return true
}

func TestReviewColumnsShowEachReviewOncePerBreakpoint(t *testing.T) {
	for n := 1; n <= 20; n++ {
		reviews := make([]Review, n)
		for i := range reviews {
			reviews[i] = Review{Quote: fmt.Sprint(i)}
		}
		cols := reviewColumns(reviews)
		// Columns space their cards with flex gap; revealing one with
		// md:block / lg:block would switch it to block layout and stack the
		// cards flush against each other.
		for ci, col := range cols {
			if strings.Contains(col.Visibility, ":block") {
				t.Errorf("n=%d col %d: visibility %q reveals as block; use flex", n, ci+1, col.Visibility)
			}
		}
		for bp, name := range []string{"phone", "md", "lg"} {
			seen := map[string]int{}
			for ci, col := range cols {
				if !visibleAt(col.Visibility, bp) {
					continue
				}
				visible := 0
				for _, it := range col.Items {
					if !visibleAt(it.Hidden, bp) {
						continue
					}
					if !it.AriaHidden {
						seen[it.Quote]++
						visible++
					}
				}
				if visible == 0 {
					continue // an empty column has nothing to scroll
				}
				// The loop time must be set for this breakpoint and match
				// the cards actually showing in the column.
				want := fmt.Sprintf("%ds;", visible*[]int{reviewPaceCol1, reviewPaceCol2, reviewPaceCol3}[ci])
				if !strings.Contains(effectiveDuration(col.Style, bp), want) {
					t.Errorf("n=%d %s col %d: style %q, want duration %s", n, name, ci+1, col.Style, want)
				}
			}
			for _, r := range reviews {
				if seen[r.Quote] != 1 {
					t.Errorf("n=%d %s: review %s shown %d times, want 1", n, name, r.Quote, seen[r.Quote])
				}
			}
		}
	}
}

// effectiveDuration returns the --reviews-duration declaration the CSS in
// tailwind/input.css resolves to at a breakpoint: lg falls back to md, md to
// the base value.
func effectiveDuration(style string, bp int) string {
	decl := func(name string) string {
		for _, d := range strings.Split(style, ";") {
			if k, v, ok := strings.Cut(d, ":"); ok && strings.TrimSpace(k) == name {
				return strings.TrimSpace(v) + ";"
			}
		}
		return ""
	}
	names := []string{"--reviews-duration", "--reviews-duration-md", "--reviews-duration-lg"}
	for i := bp; i >= 0; i-- {
		if v := decl(names[i]); v != "" {
			return v
		}
	}
	return ""
}
