package view

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
	Body      []string     // intro paragraphs on the detail page
	Details   []DetailItem // the detail block
	Gallery   []GalleryImage
	CTAs      []CTA
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
			{Src: "nfork-dining.jpg", Alt: "The lodge dining room at North Fork Crossing"},
			{Src: "nfork-pond.jpg", Alt: "An artesian-fed fishing pond beside the lodge"},
		},
		CTAs: []CTA{
			{Label: "Check rates & availability", Href: "/contact", Primary: true},
			{Label: "Plan your trip", Href: "/contact"},
		},
	},
	{
		Slug:      "sharptail",
		Name:      "Sharptail Lodge",
		Eyebrow:   "Bird hunting · 150,000+ huntable acres · Big Sky Country",
		MetaTags:  []string{"Bird hunting", "Upland country"},
		Headline:  "Wild birds, wide open country",
		CardBlurb: "Glamping on a homesteaded ranch with 150,000+ huntable acres of wild bird country and Orvis-endorsed wingshooting.",
		CardImage: "lodge-sharptail.jpg",
		Hero:      "sharptail-hero.jpg",
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
			{Src: "sharptail-dogs.jpg", Alt: "Bird dogs and the hunting truck at Sharptail Lodge"},
			{Src: "sharptail-yurts.jpg", Alt: "Guest yurts under Big Sky Country"},
		},
		CTAs: []CTA{
			{Label: "Check rates & availability", Href: "/contact", Primary: true},
			{Label: "Plan your trip", Href: "/contact"},
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
			{Src: "yurt-wildlife.jpg", Alt: "Wildlife along the Missouri River"},
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

// TeamMember is one person on the About page.
type TeamMember struct {
	Name string
	Role string
	Bio  string
}

// Team is the Pro Outfitters roster, per the content spec.
var Team = []TeamMember{
	{Name: "Elli Ortloff", Role: "Office Manager", Bio: "Happiest camping with her son or out at the lake, Elli has spent years adventuring across Montana and can't think of a better way to spend her days than helping other people find adventures of their own."},
	{Name: "Joel Loran", Role: "Fishing Guide", Bio: "Guiding for 27 years. Lives in Missoula with his wife Andrea and their kids Jacob and Ava; guides spring and summer, teaches high school history in the Bitterroot the rest of the year."},
	{Name: "Pat Kane", Role: "Fishing & Bird Hunting Guide", Bio: "Guiding for 25 years. Lives in Missoula with his wife Karen and their kids Shannon and Seamus; rivers in spring and summer, birds in the fall."},
	{Name: "Michael Carlucci", Role: "Chef, Smith River & Lodges", Bio: "From a food-obsessed Italian family, Michael started out in the family business in NYC and has cooked for Pro Outfitters since 2000."},
	{Name: "Keith Kelly", Role: "Chef, North Fork Crossing Lodge", Bio: "Originally from Bourbon County, Keith has led culinary programs at premier lodges across the American West. His elevated take on rustic, wild, locally sourced cooking makes dinner its own reason to come back."},
}
