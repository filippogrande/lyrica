/* Interazioni minime del sito: toggle tema e dimensione del testo.
   Nessun framework, nessuno script inline (la CSP non lo consente).
   Riferimenti: DESIGN_SYSTEM.md */
(function () {
  "use strict";

  var THEME_KEY = "lyrica.theme";
  var TEXT_KEY = "lyrica.textScale";
  var SCALES = [0.9, 1, 1.15, 1.3];

  document.addEventListener("click", function (event) {
    var target = event.target.closest("[data-theme-toggle], [data-text-size]");
    if (!target) {
      return;
    }
    if (target.hasAttribute("data-theme-toggle")) {
      toggleTheme();
      return;
    }
    changeTextSize(target.getAttribute("data-text-size"));
  });

  function currentTheme() {
    return document.documentElement.getAttribute("data-bs-theme") === "dark" ? "dark" : "light";
  }

  function toggleTheme() {
    var next = currentTheme() === "dark" ? "light" : "dark";
    document.documentElement.setAttribute("data-bs-theme", next);
    store(THEME_KEY, next);
  }

  function changeTextSize(direction) {
    var index = SCALES.indexOf(readScale());
    if (index === -1) {
      index = SCALES.indexOf(1);
    }
    var next = direction === "bigger" ? index + 1 : index - 1;
    if (next < 0 || next >= SCALES.length) {
      return;
    }
    applyScale(SCALES[next]);
  }

  function readScale() {
    var stored = null;
    try {
      stored = parseFloat(window.localStorage.getItem(TEXT_KEY));
    } catch (err) {
      stored = null;
    }
    return isNaN(stored) ? 1 : stored;
  }

  function applyScale(scale) {
    document.documentElement.style.setProperty("--text-scale", String(scale));
    store(TEXT_KEY, String(scale));
  }

  function store(key, value) {
    try {
      window.localStorage.setItem(key, value);
    } catch (err) {
      /* Preferenza non salvabile: la pagina resta usabile per questa sessione. */
    }
  }

  applyScale(readScale());
})();
