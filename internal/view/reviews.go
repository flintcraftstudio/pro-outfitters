package view

import "fmt"

// Review is one guest quote in the home page reviews marquee. Quote is the
// pull line set in the display face; Body, Trip and Source are optional.
// Quotes are guests' own words: trim with an ellipsis, never reword.
type Review struct {
	Quote  string
	Body   string
	Author string // first name and last initial
	Trip   string // the lodge or river, when the review names it
	Source string // e.g. "Google", for reviews quoted from a public listing
}

// Reviews are the quotes the marquee cycles through: Google reviews the
// client supplied (2023-2026), plus four testimonials from the old WordPress
// home page (Lamar W., Tom S., Eric D., David R.). Order matters: reviews are
// dealt round-robin into columns, so neighbours here land in different
// columns; keep trips mixed down each column.
var Reviews = []Review{
	{
		Quote:  "The perfect place to disconnect.",
		Body:   "A huge thank you to our guides, Joel and Pat. Despite challenging fishing conditions, they never stopped working to put us on fish. … Chef Keith absolutely knocked it out of the park every single night.",
		Author: "Lisa F.",
		Trip:   "North Fork Crossing Lodge",
		Source: "Google",
	},
	{
		Quote:  "Worth every penny.",
		Body:   "The service was phenomenal and they worked very hard to ensure all of us enjoyed our trip as much as possible. … We especially looked forward to all of the food that came from a Dutch oven or anything served in the woven basket.",
		Author: "Ethan P.",
		Trip:   "Smith River",
		Source: "Google",
	},
	{
		Quote:  "Beautiful, clean, and comfortable.",
		Body:   "I've been to PRO Outfitters' bird hunting lodge and stayed in the yurts several times. … The folks associated with PRO are the best, from the owner to the chef and all the guides and their dogs.",
		Author: "Kay K.",
		Trip:   "Sharptail Lodge",
		Source: "Google",
	},
	{
		Quote:  "He had the time of his life.",
		Body:   "We hope that he will be able to attend in some capacity next year; to camp again or perhaps a trip on the Smith River. He loved the counselor, Elijah.",
		Author: "Kira R.",
		Trip:   "Kids Camp",
		Source: "Google",
	},
	{
		Quote:  "The fishing guides were knowledgeable, kind, and engaging.",
		Body:   "Incredible experience fly fishing on the Blackfoot River. Lodging was top notch. … Going to book again for next year.",
		Author: "Lauren S.",
		Trip:   "Blackfoot River",
		Source: "Google",
	},
	{
		Quote:  "Our ladies' group looks forward to this trip every year.",
		Body:   "I've been a guest at Pro Outfitters/North Fork Crossing Lodge for over 20 years.",
		Author: "Renee M.",
		Trip:   "North Fork Crossing Lodge",
		Source: "Google",
	},
	{
		Quote:  "This was my 13th year with Pro Outfitters and they all seem to get better every year.",
		Author: "Lamar W.",
	},
	{
		Quote:  "A day that all fishermen dream about.",
		Body:   "Myself (82 years old), my son (54), my brother (80) and longtime family friend (80) just finished our five day/four night excursion on the Smith River. … On our third day, when the weather warmed up and the bugs came out, we had a day that all fishermen dream about. We were having so much luck that dinner had to be held up because we didn't get back to camp until 6:30.",
		Author: "Francis G.",
		Trip:   "Smith River",
		Source: "Google",
	},
	{
		Quote:  "Clean, quiet, and peaceful.",
		Body:   "My husband and I have been coming to PRO for years, every time has been an amazing experience. … The food for all meals, including packed lunches, is delicious. On top of that, the fishing is great, and the guides are so welcoming, helpful, and top notch.",
		Author: "Elliot D.",
		Source: "Google",
	},
	{
		Quote:  "A fantastic memory for my son and me.",
		Body:   "This was our second father and son trip out West. The first was in 2002. … The Guides and support team were very experienced, professional, polite and patient with us.",
		Author: "Clayton J.",
		Trip:   "Smith River",
		Source: "Google",
	},
	{
		Quote:  "Best bird hunting in America.",
		Author: "Tom S.",
		Trip:   "Sharptail Lodge",
	},
	{
		Quote:  "A memorable father-son trip.",
		Body:   "The accommodations were comfortable, the food was excellent, and every fishing guide we worked with was knowledgeable, patient, and genuinely enjoyable to spend time with.",
		Author: "Frank T.",
		Trip:   "North Fork Crossing Lodge",
		Source: "Google",
	},
	{
		Quote:  "The best trip of our lives.",
		Body:   "PRO Outfitters was responsible for giving us an experience that we will remember as the best trip of our lives. Five stars is not nearly enough.",
		Author: "Eric D.",
	},
	{
		Quote:  "The food and hospitality were fantastic.",
		Body:   "Our fishing guide, Pride, did everything he could to get us on fish after several storms blew most of the rivers out.",
		Author: "Fred G.",
		Source: "Google",
	},
	{
		Quote:  "Truly a great little gem.",
		Body:   "Above expectations. Truly a great little gem and a wonderful Montana experience.",
		Author: "David R.",
	},
}

// The marquee is ported from the schluters-metal-art reviews component:
// reviews are dealt across three columns that scroll upward at different
// speeds. Each column renders its cards twice and animates by -50% for a
// seamless loop; the second copy is aria-hidden so screen readers hear each
// review once. Phones show one column holding every review, md two, lg three
// once there are enough reviews to fill them (see reviewThreeColumnMin);
// cards that belong to a wider layout's extra column are hidden from the
// narrower ones with md:hidden / lg:hidden.

type reviewItem struct {
	Review
	Hidden     string // responsive hide class
	Delay      string // fade-in animation-delay
	AriaHidden bool   // true for the duplicated half of the loop
}

type reviewColumn struct {
	Items      []reviewItem
	Visibility string // wrapper class, e.g. "hidden md:block"
	Style      string // per-breakpoint loop durations, read by input.css
}

var reviewFadeDelays = []string{"0s", "0.15s", "0.3s", "0.45s", "0.6s", "0.75s"}

// Seconds of loop per visible card. A column shows a different number of
// cards at each breakpoint, so its loop time is set per breakpoint from that
// count, keeping the scroll at a steady, readable pace everywhere. The
// columns differ slightly so they drift past each other.
const (
	reviewPaceCol1 = 24
	reviewPaceCol2 = 30
	reviewPaceCol3 = 27
)

// reviewLoopStyle sets --reviews-duration (phones) and the md / lg overrides
// for a column showing base, md and lg cards; 0 leaves a breakpoint unset
// for a column that is hidden there.
func reviewLoopStyle(pace, base, md, lg int) string {
	style := ""
	if base > 0 {
		style += fmt.Sprintf("--reviews-duration: %ds;", base*pace)
	}
	if md > 0 {
		style += fmt.Sprintf("--reviews-duration-md: %ds;", md*pace)
	}
	if lg > 0 {
		style += fmt.Sprintf("--reviews-duration-lg: %ds;", lg*pace)
	}
	return style
}

func dealReviews(reviews []Review, n int) [][]Review {
	out := make([][]Review, n)
	for i, r := range reviews {
		out[i%n] = append(out[i%n], r)
	}
	return out
}

// reviewThreeColumnMin is the review count below which lg keeps two columns:
// with fewer, a three-column loop is short enough that the same card shows
// twice on screen at once.
const reviewThreeColumnMin = 10

// reviewGridCols is the grid's column classes for the layout reviewColumns
// produces.
func reviewGridCols(reviews []Review) string {
	if len(reviews) < reviewThreeColumnMin {
		return "md:grid-cols-2"
	}
	return "md:grid-cols-2 lg:grid-cols-3"
}

// reviewColumns distributes reviews for the breakpoints. With three columns,
// column 3's share is split in two so that at md each half folds into one of
// the two visible columns; on phones everything folds into column 1.
func reviewColumns(reviews []Review) []reviewColumn {
	build := func(rs []Review, hide func(int) string) []reviewItem {
		items := make([]reviewItem, 0, len(rs)*2)
		for copy := 0; copy < 2; copy++ {
			for j, r := range rs {
				items = append(items, reviewItem{
					Review:     r,
					Hidden:     hide(j),
					Delay:      reviewFadeDelays[(j*7)%len(reviewFadeDelays)],
					AriaHidden: copy == 1,
				})
			}
		}
		return items
	}

	if len(reviews) < reviewThreeColumnMin {
		cols := dealReviews(reviews, 2)
		col1All := append(append([]Review{}, cols[0]...), cols[1]...)
		return []reviewColumn{
			{
				Items: build(col1All, func(j int) string {
					if j >= len(cols[0]) {
						return "md:hidden"
					}
					return ""
				}),
				Style: reviewLoopStyle(reviewPaceCol1, len(col1All), len(cols[0]), 0),
			},
			{
				Items:      build(cols[1], func(int) string { return "" }),
				Visibility: "hidden md:block",
				Style:      reviewLoopStyle(reviewPaceCol2, 0, len(cols[1]), 0),
			},
		}
	}

	cols := dealReviews(reviews, 3)
	col1, col2 := cols[0], cols[1]
	sub := dealReviews(cols[2], 2)
	col3a, col3b := sub[0], sub[1]

	col1All := append(append(append(append([]Review{}, col1...), col3a...), col3b...), col2...)
	col2All := append(append([]Review{}, col2...), col3b...)
	col3All := append(append([]Review{}, col3a...), col3b...)

	return []reviewColumn{
		{
			Items: build(col1All, func(j int) string {
				switch {
				case j >= len(col1)+len(col3a):
					return "md:hidden"
				case j >= len(col1):
					return "lg:hidden"
				}
				return ""
			}),
			Style: reviewLoopStyle(reviewPaceCol1, len(col1All), len(col1)+len(col3a), len(col1)),
		},
		{
			Items: build(col2All, func(j int) string {
				if j >= len(col2) {
					return "lg:hidden"
				}
				return ""
			}),
			Visibility: "hidden md:block",
			Style:      reviewLoopStyle(reviewPaceCol2, 0, len(col2All), len(col2)),
		},
		{
			Items:      build(col3All, func(int) string { return "" }),
			Visibility: "hidden lg:block",
			Style:      reviewLoopStyle(reviewPaceCol3, 0, 0, len(col3All)),
		},
	}
}
