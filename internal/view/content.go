package view

import "net/url"

// Brand contact + identity constants. These come from the content spec
// (pro-outfitters-content-spec.md). Phone, mailing address, license, and
// social handle are real; keep them in sync with the spec.
const (
	BrandPhone           = "(406) 442-5489"
	BrandPhoneHref       = "tel:+14064425489"
	BrandMailAddress     = "PO Box 621, Helena, MT 59624"
	BrandInstagram       = "https://www.instagram.com/prooutfittersmontana"
	BrandInstagramHandle = "@prooutfittersmontana"
	BrandFacebook        = "https://www.facebook.com/prooutfitters"
	BrandLicense         = "Brandon Boedecker, Montana Outfitter #6481"
	BrandPermitNote      = "Smith River trips operated under special-use permit with the Helena–Lewis and Clark National Forest."
)

// Inquiry form control classes. Editorial-minimal but clearly fillable: a
// subtly recessed paper-deep field, a defined border, and a visible slate
// focus ring (WCAG 2.4.7). selectClass swaps px for pl/pr so the custom
// chevron has room.
const fieldClass = "w-full rounded-sm border border-po-border-strong bg-po-paper-deep px-4 py-3 text-[1rem] text-po-ink placeholder-po-gray-light outline-none transition focus:border-po-accent focus:ring-2 focus:ring-po-accent/30"
const selectClass = "w-full appearance-none rounded-sm border border-po-border-strong bg-po-paper-deep py-3 pl-4 pr-11 text-[1rem] text-po-ink outline-none transition focus:border-po-accent focus:ring-2 focus:ring-po-accent/30"

// NavLink is one item in the primary navigation / footer columns.
type NavLink struct {
	Label string
	Href  string
}

// PrimaryNav is the public site's main navigation.
var PrimaryNav = []NavLink{
	{Label: "The Lodges", Href: "/lodges"},
	{Label: "Smith River", Href: "/smith-river"},
	{Label: "About", Href: "/about"},
}

// InquiryInterests are the options for the contact form's "What interests
// you?" select, per the content spec.
var InquiryInterests = []string{
	"Fly fishing",
	"Bird hunting",
	"Smith River float",
	"Kids Camp",
	"Not sure yet",
}

// ContactMaxLen caps each contact form field, in characters. The form
// sets maxlength from it and handler.validateContact enforces it, so a
// bot posting directly can't fill the inbox or the notification email.
var ContactMaxLen = map[string]int{
	"name":            120,
	"email":           254,
	"phone":           40,
	"interest":        60,
	"preferred_dates": 120,
	"group_size":      60,
	"message":         5000,
}

// CTA is a call-to-action link. Primary renders as a filled button.
type CTA struct {
	Label   string
	Href    string
	Primary bool
}

// DetailItem is one row in a lodge/trip detail block (label → value).
type DetailItem struct {
	Label string
	Value string
}

// GalleryImage is one photo in a subpage gallery.
type GalleryImage struct {
	Src string
	Alt string
}

// Lodge carries everything the overview card and the detail page need.
type Lodge struct {
	Slug      string
	Name      string
	Eyebrow   string       // letter-spaced label above the headline
	MetaTags  []string     // small uppercase chips on the overview card
	Headline  string       // serif display headline
	CardBlurb string       // short text on the overview card
	CardImage string       // card thumbnail filename (under /static/img)
	Hero      string       // hero image filename (under /static/img)
	HeroAlt   string       // what the hero shows; empty falls back to Name
	Body      []string     // intro paragraphs on the detail page
	Details   []DetailItem // the detail block
	Gallery   []GalleryImage
	CTAs      []CTA
	// Address and Directions fill the "Getting here" section; both are
	// optional. Directions are transcribed from the old site's PDF sheets.
	Address    string
	Directions []DirectionRoute
}

// DirectionRoute is driving directions to a lodge from one town.
type DirectionRoute struct {
	From   string
	Steps  []string
	Arrive string // the closing landmark line
}

// mapsURL links the lodge's address to Google Maps.
func (l Lodge) mapsURL() string {
	return "https://www.google.com/maps/search/?api=1&query=" + url.QueryEscape(l.Address)
}

// heroAlt is the hero image's alt text.
func (l Lodge) heroAlt() string {
	if l.HeroAlt != "" {
		return l.HeroAlt
	}
	return l.Name
}

var lodges = []Lodge{
	{
		Slug:      "north-fork-crossing",
		Name:      "North Fork Crossing Lodge",
		Eyebrow:   "Fly fishing · Blackfoot River · Ovando, MT",
		MetaTags:  []string{"Fly fishing", "Blackfoot River"},
		Headline:  "The river that ran through it",
		CardBlurb: "Tent cabins on the bank of the North Fork of the Blackfoot, with five artesian-fed ponds and an Orvis-endorsed kids camp.",
		CardImage: "lodge-northfork.jpg",
		Hero:      "nfork-hero.jpg",
		HeroAlt:   "Tent cabins lit beneath the Milky Way at North Fork Crossing",
		Body: []string{
			"This is the water Norman Maclean wrote into A River Runs Through It — the North Fork of the Blackfoot, just outside Ovando. We've guided anglers down it for close to forty years, and our guides know every seam and riffle of it.",
			"The lodge sits on the riverbank in the Blackfoot Valley, remote and lush, holding strong populations of native West Slope cutthroat. Our tent cabins were built to keep the intimacy of camping without giving up comfort: wood floors, custom feather beds, heat, electricity, and a private bathhouse. The main lodge opens onto a lounge and dining room overlooking five artesian-fed ponds — stocked, and perfect for an evening cast before dinner.",
			"Crisp, bright mornings that dissolve into long Montana afternoons. An Orvis-endorsed fly fishing experience, and an Orvis-endorsed kids camp besides.",
		},
		Details: []DetailItem{
			{Label: "Species", Value: "West Slope cutthroat, rainbow, brown trout"},
			{Label: "Season", Value: "June through mid-October"},
			{Label: "Endorsements", Value: "Orvis Endorsed Fly Fishing Lodge · Orvis Endorsed Kids Camp"},
		},
		Gallery: []GalleryImage{
			{Src: "nfork-grounds.jpg", Alt: "Tent cabins and the lodge above an artesian-fed pond"},
			{Src: "nfork-tents.jpg", Alt: "Two tent cabins joined by a gravel path through the grass"},
			{Src: "nfork-pond-aerial.jpg", Alt: "An artesian-fed pond in the pines, the tent cabins along its far bank"},
			{Src: "nfork-dining.jpg", Alt: "The lodge dining room at North Fork Crossing"},
		},
		CTAs: []CTA{
			{Label: "Check rates & availability", Href: "/contact", Primary: true},
			{Label: "Kids Camp", Href: "/kids-camp"},
		},
		Address: "328 North Fork Hill Road, Ovando, MT 59854",
		Directions: []DirectionRoute{
			{
				From: "Helena",
				Steps: []string{
					"Take Highway 12 west out of town, about 30 miles to the junction with Highway 141 in Avon.",
					"Turn right on Highway 141 and travel 33 miles to where it ends at a T with Highway 200.",
					"Turn left on Highway 200 and continue west about 12 miles to Ovando, which lies on the south side of the highway. Trixie's Restaurant and Bar, on the left, is a good landmark.",
					"Turn left at Trixie's and stay on the Ovando–Helmville Road. (If you miss it, turn left at the sign for Ovando and bear right through town onto the same road.)",
					"It's the only paved road running south from the highway; you may see signs for Brown's Lake, Helmville, or Harry Morgan Fishing Access.",
					"After about 2 miles the road turns from paved to dirt. Continue about another mile, until you cross the North Fork of the Blackfoot.",
				},
				Arrive: "North Fork Crossing is immediately on the left, just after the bridge.",
			},
			{
				From: "Missoula",
				Steps: []string{
					"Take Interstate 90 east to Highway 200, exiting at Bonner/Milltown (about 10 miles).",
					"Continue east on Highway 200 about 60 miles to Ovando, which lies on the south side (your right) of the highway.",
					"Turn right at the sign for Ovando and bear right through town onto the Ovando–Helmville Road. (If you miss it, turn right at Trixie's Restaurant and Bar and keep left onto the same road.)",
					"It's the only paved road running south from the highway; you may see signs for Brown's Lake, Helmville, or Harry Morgan Fishing Access.",
					"After about 2 miles the road turns from paved to dirt. Continue about another mile, until you cross the North Fork of the Blackfoot.",
				},
				Arrive: "North Fork Crossing is immediately on the left, just after the bridge.",
			},
			{
				From: "Bozeman",
				Steps: []string{
					"Take Interstate 90 west to Garrison Junction (about 130 miles) and exit onto Highway 12 east.",
					"Continue about 15 miles to Avon and the junction with Highway 141.",
					"In Avon, turn left onto Highway 141 and travel 33 miles to where it ends at a T with Highway 200.",
					"Turn left on Highway 200 and continue west about 12 miles to Ovando. Trixie's Restaurant and Bar, on the left, is a good landmark.",
					"Turn left at Trixie's and stay on the Ovando–Helmville Road. (If you miss it, turn left at the sign for Ovando and bear right through town onto the same road.)",
					"After about 2 miles the road turns from paved to dirt. Continue about another mile, until you cross the North Fork of the Blackfoot.",
				},
				Arrive: "North Fork Crossing is immediately on the left, just after the bridge.",
			},
		},
	},
	{
		Slug:      "sharptail",
		Name:      "Sharptail Lodge",
		Eyebrow:   "Bird hunting · 150,000+ huntable acres · Big Sky Country",
		MetaTags:  []string{"Bird hunting", "Upland country"},
		Headline:  "Wild birds, wide open country",
		CardBlurb: "Glamping on a homesteaded ranch with 150,000+ huntable acres of wild bird country and Orvis-endorsed wingshooting.",
		CardImage: "sharptail-yurts.jpg",
		Hero:      "sharptail-point.jpg",
		HeroAlt:   "A bird dog on point in tall autumn grass",
		Body: []string{
			"Watch wild birds rise into a brilliant sky over more than 150,000 huntable acres of some of the most prolific wild bird habitat in America. Our guests hunt wild birds only — sharptail grouse and Hungarian partridge to open the season, ring-neck pheasant after the first week of October.",
			"The ranch was homesteaded in the late 1800s. The original red barn now houses our dogs; guests stay in yurts with two queen beds, private baths, and gas fireplaces, gathering in two larger yurts joined by a deck and an outdoor fire pit for meals and cocktails. The dining is excellent and the bar is stocked with beer, wine, and mixers (guests bring their own hard liquor if they'd like).",
			"The views are what define this place. Panoramic country runs north as far as you can see — coulees, grasslands, farm ground, and timbered slopes where the sky meets the horizon and the term \"Big Sky\" finally makes sense. Bring your own dog or hunt behind ours. Either way, the hunting starts right out the front door.",
			"Glamping, at its best. An Orvis-endorsed wingshooting experience, and we've guided bird hunters in this country for close to forty years.",
		},
		Details: []DetailItem{
			{Label: "Species by season", Value: "Sharptail grouse & Hungarian partridge early; ring-neck pheasant after ~Oct 7"},
			{Label: "Season", Value: "September 1 through mid-November"},
			{Label: "Endorsement", Value: "Orvis Endorsed Wingshooting Lodge"},
		},
		Gallery: []GalleryImage{
			{Src: "sharptail-yurts.jpg", Alt: "Guest yurts under Big Sky Country"},
			{Src: "sharptail-hunters.jpg", Alt: "Two hunters with their birds, a flat-topped butte behind them"},
			{Src: "sharptail-dining.jpg", Alt: "The dining yurt, with its long table and antler chandelier"},
			{Src: "sharptail-lab-birds.jpg", Alt: "A black Lab sitting beside pheasants and a shotgun"},
			{Src: "sharptail-bedroom.jpg", Alt: "A guest yurt at sunset, its two beds lit by bedside lamps"},
			{Src: "sharptail-dogs.jpg", Alt: "Three setters looking out of the hunting truck"},
		},
		CTAs: []CTA{
			{Label: "Check rates & availability", Href: "/contact", Primary: true},
			{Label: "Gear list", Href: "/packing-lists/bird-hunting"},
		},
		Address: "445 Sharptail Lane, Geyser, MT 59447",
		Directions: []DirectionRoute{
			{
				From: "Great Falls",
				Steps: []string{
					"Head east on Highway 87/200 about 45 miles to Geyser.",
					"From the Geyser schoolhouse (on your left), continue 1.5 miles toward Stanford on the same highway.",
					"Turn right (south) on the gravel road immediately past mile marker 23, marked with a brown Forest Service sign for Lonetree Road.",
					"Go straight 4.1 miles to a T, turn right, and continue 2.2 miles.",
					"Turn left at the mailbox marked PRO.",
				},
				Arrive: "Sharptail Lodge is on your right.",
			},
			{
				From: "Lewistown",
				Steps: []string{
					"Travel west on Highway 87/200 to Stanford (about 45 miles), then continue about 13 miles past Stanford.",
					"Turn left on the gravel road just before mile marker 23.",
					"Go straight 4.1 miles to a T, turn right, and continue 2.2 miles.",
					"Turn left at the mailbox marked PRO.",
				},
				Arrive: "Sharptail Lodge is on your right.",
			},
		},
	},
	{
		Slug:      "yurt-at-craig",
		Name:      "The Yurt at Craig",
		Eyebrow:   "Fly fishing · Missouri River · Craig, MT",
		MetaTags:  []string{"Fly fishing", "Missouri River"},
		Headline:  "Coffee on the deck, the Missouri out front",
		CardBlurb: "A self-catered base on the Missouri with the fly shops, restaurants, and bars of Craig a short walk away.",
		CardImage: "lodge-yurt.jpg",
		Hero:      "yurt-hero.jpg",
		HeroAlt:   "The yurt's kitchen, with a breakfast bar under the lattice walls",
		Body: []string{
			"Some of the finest trout water in the world runs past the front door. The Yurt at Craig is a self-catered base on the Missouri — coffee on the deck over the river in the morning, the fly shops, restaurants, and bars of Craig a short walk down the street.",
			"Eight minutes from the Wolf Creek take-out, nineteen from Holter Lake. The main yurt has a full kitchen, living room, and two bedrooms (one queen each) plus a covered deck. Two detached guest bedrooms sleep two apiece, and the bathhouse holds two private, fully furnished bath and shower rooms.",
			"A fly fishing getaway with room for the whole party.",
		},
		Details: []DetailItem{
			{Label: "Water", Value: "Missouri River"},
			{Label: "Best for", Value: "Self-guided stays, groups, anglers who want the town close"},
			{Label: "Sleeps", Value: "Up to 6 across the main yurt and guest house"},
		},
		Gallery: []GalleryImage{
			{Src: "yurt-bedroom.jpg", Alt: "A guest bedroom at The Yurt at Craig"},
			{Src: "yurt-wildlife.jpg", Alt: "A deer at the water's edge on the Missouri"},
		},
		CTAs: []CTA{
			{Label: "Check rates & availability", Href: "/contact", Primary: true},
			{Label: "Plan your trip", Href: "/contact"},
		},
	},
}

// Lodges returns all lodges for the overview page.
func Lodges() []Lodge { return lodges }

// LodgeBySlug returns the lodge with the given slug, or false if none.
func LodgeBySlug(slug string) (Lodge, bool) {
	for _, l := range lodges {
		if l.Slug == slug {
			return l, true
		}
	}
	return Lodge{}, false
}

// TeamMember is one person on the About page. Photo is a filename under
// /static/img (a square headshot); the .avif/.webp siblings are derived.
type TeamMember struct {
	Name  string
	Role  string
	Bio   string
	Photo string
}

// Team is the Pro Outfitters roster, per the content spec.
var Team = []TeamMember{
	{Name: "Elli Ortloff", Role: "Office Manager", Photo: "team-elli.jpg", Bio: "Happiest camping with her son or out at the lake, Elli has spent years adventuring across Montana and can't think of a better way to spend her days than helping other people find adventures of their own."},
	{Name: "Joel Loran", Role: "Fishing Guide", Photo: "team-joel.jpg", Bio: "Guiding for 27 years. Lives in Missoula with his wife Andrea and their kids Jacob and Ava; guides spring and summer, teaches high school history in the Bitterroot the rest of the year."},
	{Name: "Pat Kane", Role: "Fishing & Bird Hunting Guide", Photo: "team-pat.jpg", Bio: "Guiding for 25 years. Lives in Missoula with his wife Karen and their kids Shannon and Seamus; rivers in spring and summer, birds in the fall."},
	{Name: "Michael Carlucci", Role: "Chef, Smith River & Lodges", Photo: "team-michael.jpg", Bio: "From a food-obsessed Italian family, Michael started out in the family business in NYC and has cooked for Pro Outfitters since 2000."},
	{Name: "Keith Kelly", Role: "Chef, North Fork Crossing Lodge", Photo: "team-keith.jpg", Bio: "Originally from Bourbon County, Keith has led culinary programs at premier lodges across the American West. His elevated take on rustic, wild, locally sourced cooking makes dinner its own reason to come back."},
}

// PackingSection is one titled group of items on a packing list. Note is an
// optional line under the title (what the outfitter provides, say).
type PackingSection struct {
	Title string
	Note  string
	Items []string
}

// PackingList is a "what to bring" page. The lists are carried over from the
// old site, lightly cleaned up; a list with no Sections renders as a draft.
type PackingList struct {
	Slug        string
	Name        string
	Eyebrow     string
	Lead        string
	Description string // meta description
	LegacyURL   string // the list on the old WordPress site, for reference
	Sections    []PackingSection
	BackLabel   string
	BackHref    string
}

var packingLists = []PackingList{
	{
		Slug:        "smith-river",
		Name:        "Smith River packing list",
		Eyebrow:     "Before you go",
		Lead:        "What to bring for a multi-day float on the Smith.",
		Description: "What to pack for a guided Smith River float trip with Pro Outfitters.",
		LegacyURL:   "https://prooutfitters.com/smith-river-packing-list/",
		Sections: []PackingSection{
			{Title: "Personal", Items: []string{
				"Biodegradable soap",
				"Personal towel or washcloth",
				"Hand wipes",
				"Personal toiletries",
			}},
			{Title: "Clothing", Items: []string{
				"Socks (polypro or wool)",
				"Swimsuit",
				"Lightweight shirts, long and short sleeved",
				"Quick-dry shorts",
				"Quick-dry pants, plus a casual pair (jeans or chinos)",
				"Fleece for under waders or around camp",
				"Polypro or Thermax long underwear",
				"Fleece pullover or sweater",
				"Warm jacket (no down)",
				"Gloves",
				"A good rain jacket",
				"Vest or chest pack",
				"Bootfoot waders and wading belt (recommended in cooler weather)",
				"Wide-brimmed hat",
			}},
			{Title: "Footwear", Items: []string{
				"Camp shoes, for trips to the bathroom",
				"Wading boots or river sandals",
				"Hiking shoes",
			}},
			{Title: "Sleeping bag", Note: "We'll provide one if you don't have your own. Let us know in advance."},
			{Title: "Fishing gear", Note: "All gear can be provided at no additional cost.", Items: []string{
				"8- or 9-foot rod, 4, 5 or 6 weight",
				"Floating line",
				"Hemostats and clippers",
				"Leaders and tippet",
				"Flies are provided, but bring your own if you like",
				"A Montana fishing license, bought in advance",
			}},
			{Title: "Miscellaneous", Items: []string{
				"Sunscreen",
				"Polarized sunglasses",
				"Lip balm with sunscreen",
				"Insect repellent",
				"Camera",
				"Book or journal, in a zip-lock bag",
				"Flashlight or headlamp, with spare batteries",
			}},
			{Title: "Beverages", Note: "We provide beverages, including beer and wine. We can't provide hard liquor, so bring your own if you'd like cocktails; we can stop for it locally on request."},
		},
		BackLabel: "Back to the Smith River",
		BackHref:  "/smith-river",
	},
	{
		Slug:        "kids-camp",
		Name:        "Kids Camp: what to bring",
		Eyebrow:     "Before you go",
		Lead:        "What to pack for Kids Camp at North Fork Crossing Lodge.",
		Description: "What to pack for the Pro Outfitters Kids Camp at North Fork Crossing Lodge.",
		LegacyURL:   "https://prooutfitters.com/kids-camp-what-to-bring/",
		Sections: []PackingSection{
			{Title: "Personal", Note: "North Fork provides towels, soap, shampoo, lotion, hairdryers and bathrobes. Bring everything else for the week, like a toothbrush and toothpaste."},
			{Title: "Clothing", Items: []string{
				"Socks (polypro, wool or cotton)",
				"Lightweight shirts (quick-dry fishing shirts are ideal; T-shirts are fine)",
				"Fleece pullover or sweater",
				"Warm jacket (mornings can be below 40 degrees)",
				"Quick-dry pants or shorts for fishing and the day's activities",
				"One or two pairs of casual pants, like jeans",
				"Athletic shorts for games",
				"Gloves and a warm hat for early mornings",
				"Hat with a brim",
				"Polarized sunglasses",
				"Rain jacket",
				"Vest or chest pack",
			}},
			{Title: "Footwear", Items: []string{
				"Slip-ons for trips to the bathroom",
				"River sandals or shoes",
				"Sneakers or hiking boots",
			}},
			{Title: "Fishing gear", Note: "We'll outfit campers with gear at no additional cost if they need it.", Items: []string{
				"8- or 9-foot fly rod, 4 to 6 weight",
				"Reel with floating line",
				"Fly floatant",
				"Hemostat and clippers",
				"Leaders and tippet",
			}},
			{Title: "Miscellaneous", Items: []string{
				"Camera",
				"Book or journal",
				"Sunscreen and lip protection",
			}},
		},
		BackLabel: "Back to Kids Camp",
		BackHref:  "/kids-camp",
	},
	{
		Slug:        "bird-hunting",
		Name:        "Bird hunting gear list",
		Eyebrow:     "Before you go",
		Lead:        "What to pack for a hunt at Sharptail Lodge.",
		Description: "What to pack for a guided wild-bird hunt at Sharptail Lodge with Pro Outfitters.",
		LegacyURL:   "https://prooutfitters.com/bird-hunting-equipment-list/",
		Sections: []PackingSection{
			{Title: "Boots", Note: "The most important item: there's a lot of walking. Bring boots that are already broken in. A comfortable pair of lightweight hiking boots beats a new pair of leather ones."},
			{Title: "Clothing", Items: []string{
				"Comfortable hunting pants (heavy briar pants aren't necessary)",
				"Loose-fitting hunting shirts, light and medium weight",
				"A hunting vest to carry shells and a water bottle (orange preferred)",
				"Hunting cap (orange preferred)",
				"A weather- and wind-resistant hunting coat",
				"A good rain coat",
				"Good walking socks",
			}},
			{Title: "Gear", Items: []string{
				"Shotgun, 20 or 12 gauge. We shoot 7s for sharptail and Huns, 6s for pheasant",
				"A soft or hard cooler to take your birds home",
			}},
			{Title: "Other items", Items: []string{
				"Fleece vest, jacket or pullover",
				"Capilene long underwear",
				"Shooting gloves",
				"Shooting glasses",
				"Stocking cap",
			}},
		},
		BackLabel: "Back to Sharptail Lodge",
		BackHref:  "/lodges/sharptail",
	},
}

// PackingLists returns every packing list.
func PackingLists() []PackingList { return packingLists }

// PackingListBySlug returns the packing list with the given slug, or false.
func PackingListBySlug(slug string) (PackingList, bool) {
	for _, p := range packingLists {
		if p.Slug == slug {
			return p, true
		}
	}
	return PackingList{}, false
}
