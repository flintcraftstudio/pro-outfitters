package view

// A gallery reads as a paced editorial sequence rather than a grid: on md+
// it cycles through an asymmetric pair (a large 3:2 frame beside a smaller
// portrait set lower, with an empty column between them), then a full-width
// panorama, and the next pair mirrors the first. On phones every frame
// stacks at 3:2. The class strings are literal so Tailwind's scan of this
// package keeps them; per-image object-position comes from imgFocus.

// gallerySlot is one image placed in the sequence.
type gallerySlot struct {
	GalleryImage
	Wrap  string // grid placement on the 12-column md+ grid
	Img   string // aspect ratio and crop
	Sizes string // srcset sizes hint for the placement's rendered width
}

// Placements. Widths in Sizes assume the 75rem container less its 3rem
// side padding (1104px of content at full width).
var (
	slotWideLeft = gallerySlot{
		Wrap:  "md:col-span-7 md:col-start-1",
		Img:   "aspect-[3/2]",
		Sizes: "(min-width: 1200px) 640px, (min-width: 768px) 56vw, 100vw",
	}
	slotNarrowRight = gallerySlot{
		Wrap:  "md:col-span-4 md:col-start-9 md:mt-[clamp(4rem,10vw,9rem)]",
		Img:   "aspect-[3/2] md:aspect-[4/5]",
		Sizes: "(min-width: 1200px) 360px, (min-width: 768px) 32vw, 100vw",
	}
	slotNarrowLeft = gallerySlot{
		Wrap:  "md:col-span-4 md:col-start-1 md:mt-[clamp(4rem,10vw,9rem)]",
		Img:   "aspect-[3/2] md:aspect-[4/5]",
		Sizes: "(min-width: 1200px) 360px, (min-width: 768px) 32vw, 100vw",
	}
	slotWideRight = gallerySlot{
		Wrap:  "md:col-span-7 md:col-start-6",
		Img:   "aspect-[3/2]",
		Sizes: "(min-width: 1200px) 640px, (min-width: 768px) 56vw, 100vw",
	}
	slotPanorama = gallerySlot{
		Wrap:  "md:col-span-12",
		Img:   "aspect-[3/2] md:aspect-[2/1]",
		Sizes: "(min-width: 1200px) 1104px, 100vw",
	}
	// slotSolo takes an image left alone at the start of a cycle (the last
	// of four, say), centred so the sequence doesn't end on a lopsided half.
	slotSolo = gallerySlot{
		Wrap:  "md:col-span-8 md:col-start-3",
		Img:   "aspect-[3/2]",
		Sizes: "(min-width: 1200px) 736px, (min-width: 768px) 64vw, 100vw",
	}
)

// galleryLayout assigns each image its placement in the sequence.
func galleryLayout(images []GalleryImage) []gallerySlot {
	slots := make([]gallerySlot, len(images))
	for i, im := range images {
		mirrored := (i/3)%2 == 1
		var s gallerySlot
		switch {
		case i%3 == 2:
			s = slotPanorama
		case i%3 == 0 && i == len(images)-1:
			s = slotSolo
		case i%3 == 0 && mirrored:
			s = slotNarrowLeft
		case i%3 == 0:
			s = slotWideLeft
		case mirrored:
			s = slotWideRight
		default:
			s = slotNarrowRight
		}
		s.GalleryImage = im
		slots[i] = s
	}
	return slots
}
