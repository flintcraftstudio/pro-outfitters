// Pro Outfitters — minimal, restrained interactions.
// 1. Sticky nav: transparent over the hero, paper-blur once scrolled past it.
// 2. Gentle reveal on scroll-in. Motion respects prefers-reduced-motion.
// 3. Re-render Turnstile in forms htmx swaps in (htmx and Turnstile load
//    only on the pages with forms; this listener is inert elsewhere).

(function () {
  var nav = document.getElementById("nav");
  if (!nav) return;
  var hero = document.querySelector("[data-hero]");
  // Pages without a hero start in the scrolled (ink-on-paper) state.
  if (!hero) {
    nav.dataset.state = "scrolled";
    return;
  }
  // Switch once the hero's bottom edge rises within 90px of the top of the
  // viewport. An IntersectionObserver with a -90px top margin reports exactly
  // that, without reading layout: measuring hero.offsetHeight at startup
  // forced a full synchronous layout of the page (~300 ms on a mid-range
  // phone), and it also tracks the hero resizing as fonts and images settle.
  if (!("IntersectionObserver" in window)) {
    nav.dataset.state = "scrolled";
    return;
  }
  new IntersectionObserver(
    function (entries) {
      nav.dataset.state = entries[0].isIntersecting ? "top" : "scrolled";
    },
    { rootMargin: "-90px 0px 0px 0px" }
  ).observe(hero);
})();

(function () {
  var els = document.querySelectorAll(".reveal");
  if (!els.length) return;
  if (!("IntersectionObserver" in window)) {
    els.forEach(function (e) {
      e.classList.add("in");
    });
    return;
  }
  var io = new IntersectionObserver(
    function (entries) {
      entries.forEach(function (en) {
        if (en.isIntersecting) {
          en.target.classList.add("in");
          io.unobserve(en.target);
        }
      });
    },
    { threshold: 0.16, rootMargin: "0px 0px -8% 0px" }
  );
  els.forEach(function (e) {
    io.observe(e);
  });
})();

// 3. Turnstile after htmx swaps. Turnstile renders its widgets once, on
// page load. When the contact form comes back with an error, htmx swaps
// in a fresh form whose widget container is empty, so the next submit
// would carry no token and fail verification. Render any container that
// has no widget yet after each swap. (Tokens are single-use, so a fresh
// widget is what a retry needs anyway.)
document.addEventListener("htmx:afterSwap", function () {
  if (!window.turnstile) return;
  document.querySelectorAll(".cf-turnstile").forEach(function (el) {
    // A rendered widget holds its response input; a fresh one doesn't.
    if (!el.querySelector('input[name="cf-turnstile-response"]')) {
      window.turnstile.render(el);
    }
  });
});

// Reviews marquee: start the scroll and card fade once the grid is in view,
// and wire the Pause/Play toggle (CSS shows it only under .js, and keeps it
// hidden under reduced motion, where nothing moves).
(function () {
  var grid = document.querySelector("[data-reviews-grid]");
  if (!grid) return;
  var toggle = document.querySelector("[data-reviews-toggle]");
  if (toggle) {
    var label = toggle.querySelector("[data-reviews-toggle-label]");
    toggle.addEventListener("click", function () {
      var paused = grid.classList.toggle("is-paused");
      if (label) label.textContent = paused ? "Play reviews" : "Pause reviews";
    });
  }
  if (!("IntersectionObserver" in window)) {
    grid.classList.add("in-view");
    return;
  }
  var io = new IntersectionObserver(
    function (entries) {
      entries.forEach(function (en) {
        if (en.isIntersecting) {
          grid.classList.add("in-view");
          io.unobserve(grid);
        }
      });
    },
    { threshold: 0.25 }
  );
  io.observe(grid);
})();
