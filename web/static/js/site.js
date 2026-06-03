// Pro Outfitters — minimal, restrained interactions.
// 1. Sticky nav: transparent over the hero, paper-blur once scrolled past it.
// 2. Gentle reveal on scroll-in. Motion respects prefers-reduced-motion.

(function () {
  var nav = document.getElementById("nav");
  if (nav) {
    var hero = document.querySelector("[data-hero]");
    function onScroll() {
      var threshold = hero ? hero.offsetHeight - 90 : 120;
      nav.dataset.state = window.scrollY > threshold ? "scrolled" : "top";
    }
    // Pages without a hero start in the scrolled (ink-on-paper) state.
    if (!hero) {
      nav.dataset.state = "scrolled";
    } else {
      onScroll();
      window.addEventListener("scroll", onScroll, { passive: true });
    }
  }
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
