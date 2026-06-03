/* Sections.jsx — composed page sections built from primitives. */
const { useState: useS2, useEffect: useE2 } = React;

const Nav = ({ over, onPlan }) => {
  const [scrolled, setScrolled] = useS2(false);
  useE2(() => {
    const onScroll = () => setScrolled(window.scrollY > window.innerHeight * 0.7);
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);
  const cls = "k-nav " + (scrolled ? "scrolled" : (over ? "over" : "scrolled"));
  return (
    <nav className={cls}>
      <div className="inner">
        <a href="#top"><Wordmark /></a>
        <div className="links">
          <a href="#fishing">Fly Fishing</a>
          <a href="#hunting">Bird Hunting</a>
          <a href="#lodges">The Lodges</a>
          <a href="#stewardship">Stewardship</a>
        </div>
        <Button variant="primary" onClick={onPlan}>Plan Your Trip</Button>
      </div>
    </nav>
  );
};

const Hero = ({ onPlan }) => (
  <header id="top" style={{ position: "relative", minHeight: "94vh", display: "flex", alignItems: "flex-end", color: "var(--on-dark)", isolation: "isolate" }}>
    <Photo tone="dusk" style={{ position: "absolute", inset: 0, zIndex: 0 }} />
    <div style={{ position: "absolute", inset: 0, zIndex: 1, background: "linear-gradient(180deg,rgba(16,14,11,.42),rgba(16,14,11,.08) 30%,rgba(16,14,11,.2) 62%,rgba(16,14,11,.72))" }} />
    <div className="container" style={{ position: "relative", zIndex: 2, width: "100%", paddingBottom: "clamp(3rem,7vh,5rem)", paddingTop: "8rem" }}>
      <Eyebrow tone="light">Est. 1970&nbsp;&nbsp;·&nbsp;&nbsp;Orvis-Endorsed&nbsp;&nbsp;·&nbsp;&nbsp;Montana</Eyebrow>
      <h1 style={{ fontFamily: "var(--font-display)", fontWeight: 400, fontSize: "clamp(2.9rem,6.4vw,5.4rem)", lineHeight: 1.02, letterSpacing: "-.015em", color: "var(--on-dark)", maxWidth: "16ch", margin: "1.4rem 0 0" }}>
        Montana, the way it was meant to be fished.
      </h1>
      <p style={{ fontFamily: "var(--font-display)", fontWeight: 300, fontSize: "clamp(1.25rem,2.1vw,1.5rem)", lineHeight: 1.5, color: "var(--on-dark-soft)", maxWidth: "46ch", marginTop: "1.5rem" }}>
        Owner-operated since 1970. Guided float trips, upland wingshooting, and three lodges on water worth protecting.
      </p>
      <div style={{ display: "flex", gap: "1rem", marginTop: "2.3rem", flexWrap: "wrap" }}>
        <Button variant="primary" onClick={onPlan}>Plan Your Trip</Button>
        <Button variant="ghost">Our Story</Button>
      </div>
    </div>
  </header>
);

const ExperienceCard = ({ id, tone, eyebrow, title, body, link }) => (
  <article id={id} style={{ display: "flex", flexDirection: "column" }}>
    <Photo tone={tone} caption={eyebrow} style={{ aspectRatio: "5 / 4" }} />
    <div style={{ paddingTop: "1.4rem" }}>
      <Eyebrow tone="gold">{eyebrow}</Eyebrow>
      <h3 style={{ fontFamily: "var(--font-display)", fontWeight: 400, fontSize: "clamp(1.6rem,2.6vw,2.1rem)", letterSpacing: "-.01em", margin: ".6rem 0" }}>{title}</h3>
      <p style={{ color: "var(--gray)", maxWidth: "40ch", marginBottom: "1.1rem" }}>{body}</p>
      <LinkArrow>{link}</LinkArrow>
    </div>
  </article>
);

const LodgeCard = ({ tone, meta, title, body }) => (
  <article style={{ background: "var(--surface)", border: "1px solid var(--border)", borderRadius: 2, overflow: "hidden", display: "flex", flexDirection: "column" }}>
    <Photo tone={tone} style={{ aspectRatio: "4 / 3" }} />
    <div style={{ padding: "1.5rem 1.4rem 1.6rem" }}>
      <div style={{ display: "flex", gap: "1rem", fontSize: ".72rem", fontWeight: 600, letterSpacing: ".1em", textTransform: "uppercase", color: "var(--gold-deep)", marginBottom: ".8rem" }}>
        {meta.map((m, i) => <span key={i}>{m}</span>)}
      </div>
      <h3 style={{ fontFamily: "var(--font-display)", fontWeight: 400, fontSize: "1.4rem", marginBottom: ".5rem" }}>{title}</h3>
      <p style={{ fontSize: ".9375rem", color: "var(--gray)", marginBottom: "1rem" }}>{body}</p>
      <LinkArrow small>Explore</LinkArrow>
    </div>
  </article>
);

const StatBand = () => (
  <section className="k-band dark">
    <div className="container" style={{ display: "grid", gridTemplateColumns: "repeat(4,1fr)", gap: 0 }}>
      {[["1970", "Established"], ["50+", "Seasons guided"], ["3", "Owner-run lodges"], ["Orvis", "Endorsed outfitter"]].map(([n, c], i) => (
        <div key={i} style={{ padding: i ? "0 clamp(1rem,2.2vw,2rem)" : "0 clamp(1rem,2.2vw,2rem) 0 0", borderLeft: i ? "1px solid rgba(243,238,229,.16)" : "0" }}>
          <div style={{ fontFamily: "var(--font-display)", fontWeight: 300, fontSize: "clamp(2.2rem,4vw,3.2rem)", lineHeight: 1, color: "var(--on-dark)" }}>{n}</div>
          <div style={{ fontSize: ".72rem", fontWeight: 600, letterSpacing: ".13em", textTransform: "uppercase", color: "var(--on-dark-soft)", marginTop: ".8rem" }}>{c}</div>
        </div>
      ))}
    </div>
  </section>
);

const PullQuote = ({ children, cite }) => (
  <section className="k-band alt" style={{ textAlign: "center" }}>
    <div className="container" style={{ maxWidth: "52rem" }}>
      <span style={{ fontFamily: "var(--font-display)", fontSize: "5rem", lineHeight: 0, color: "var(--gold)", display: "block", height: "2.2rem", marginBottom: "1rem" }}>&ldquo;</span>
      <blockquote style={{ fontFamily: "var(--font-display)", fontWeight: 300, fontStyle: "italic", fontSize: "clamp(1.6rem,3.4vw,2.6rem)", lineHeight: 1.32, letterSpacing: "-.01em", color: "var(--ink)" }}>{children}</blockquote>
      <cite style={{ display: "block", marginTop: "1.8rem", fontStyle: "normal", fontSize: ".8125rem", fontWeight: 600, letterSpacing: ".13em", textTransform: "uppercase", color: "var(--gray-light)" }}>{cite}</cite>
    </div>
  </section>
);

Object.assign(window, { Nav, Hero, ExperienceCard, LodgeCard, StatBand, PullQuote });
