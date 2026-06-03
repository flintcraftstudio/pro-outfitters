// Pro Outfitters — minimal, restrained interactions.
// 1. Sticky nav: transparent over the hero, paper-blur once scrolled past it.
// 2. Gentle reveal on scroll-in. Motion respects prefers-reduced-motion.

(function () {
  var nav = document.getElementById("nav");
  if (!nav) return;
  var hero = document.querySelector("[data-hero]");
  // Pages without a hero start in the scrolled (ink-on-paper) state.
  if (!hero) {
    nav.dataset.state = "scrolled";
    return;
  }
  // Cache the threshold so the scroll handler never reads layout geometry
  // (offsetHeight) — reading it per scroll event forces a reflow. Recompute
  // only when the hero's height can actually change: on resize and once
  // fonts/images have settled after load.
  var threshold = hero.offsetHeight - 90;
  function recompute() {
    threshold = hero.offsetHeight - 90;
    onScroll();
  }
  function onScroll() {
    nav.dataset.state = window.scrollY > threshold ? "scrolled" : "top";
  }
  onScroll();
  window.addEventListener("scroll", onScroll, { passive: true });
  window.addEventListener("resize", recompute, { passive: true });
  window.addEventListener("load", recompute);
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
