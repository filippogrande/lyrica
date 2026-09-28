/* Interazioni minime del sito: toggle tema, dimensione del testo, filtro
   alfabetico della pagina Bands e scelta della lingua del brano.
   Nessun framework, nessuno script inline (la CSP non lo consente).
   Riferimenti: docs/FRONTEND.md */
(function () {
  "use strict";

  var THEME_KEY = "lyrica.theme";
  var TEXT_KEY = "lyrica.textScale";
  var SCALES = [0.9, 1, 1.15, 1.3];
  /* Parametro dell'URL per ogni lato del testo: il link condiviso riapre la
     pagina con la lingua scelta. */
  var LANG_PARAMS = { original: "orig", translation: "lang" };

  document.addEventListener("click", function (event) {
    var letter = event.target.closest("[data-band-letter]");
    if (letter) {
      filterBands(letter);
      return;
    }
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

  document.addEventListener("change", function (event) {
    var select = event.target.closest("[data-lang-select]");
    if (select) {
      showTextLang(select, true);
    }
  });

  /* Filtro alfabetico: tutti i gruppi stanno nel DOM, si mostra solo quello
     scelto con l'attributo hidden. Senza JS restano visibili tutti, quindi la
     pagina non perde contenuto. */
  function filterBands(button) {
    var index = button.closest("[data-band-index]");
    if (!index) {
      return;
    }
    var chosen = button.getAttribute("data-band-letter");
    var groups = index.querySelectorAll("[data-band-group]");
    for (var i = 0; i < groups.length; i++) {
      var shown = chosen === "*" || groups[i].getAttribute("data-band-group") === chosen;
      groups[i].hidden = !shown;
    }
    var buttons = index.querySelectorAll("[data-band-letter]");
    for (var j = 0; j < buttons.length; j++) {
      var active = buttons[j] === button;
      buttons[j].classList.toggle("is-active", active);
      buttons[j].setAttribute("aria-pressed", active ? "true" : "false");
    }
  }

  /* Scelta della lingua del brano: i blocchi di tutte le lingue sono nel DOM e
     il menu mostra quello scelto, cambiando l'attributo data-text-selected. Il
     CSS nasconde le altre lingue solo quando JS c'è (regola .js in
     assets/css/lyrica.css): senza JS restano tutte visibili. */
  function showTextLang(select, writeURL) {
    var side = select.getAttribute("data-lang-select");
    var code = select.value;
    var blocks = document.querySelectorAll('[data-text-block="' + side + '"]');
    for (var i = 0; i < blocks.length; i++) {
      var shown = blocks[i].getAttribute("data-text-lang") === code;
      blocks[i].setAttribute("data-text-selected", shown ? "1" : "0");
    }
    if (writeURL) {
      writeLangParam(LANG_PARAMS[side], code);
    }
  }

  /* Un lato senza menu ha una lingua sola e non c'è niente da commutare (D95).
     Il blocco è già quello mostrato; se però il contenuto ne dichiarasse due, il
     CSS nasconderebbe il secondo senza che nessuno possa riaprirlo: meglio
     mostrarlo che perderlo. */
  function showSidesWithoutSelect() {
    var blocks = document.querySelectorAll("[data-text-block]");
    for (var i = 0; i < blocks.length; i++) {
      var side = blocks[i].getAttribute("data-text-block");
      if (!document.querySelector('[data-lang-select="' + side + '"]')) {
        blocks[i].setAttribute("data-text-selected", "1");
      }
    }
  }

  function writeLangParam(name, code) {
    if (!name) {
      return;
    }
    try {
      var url = new URL(window.location.href);
      url.searchParams.set(name, code);
      window.history.replaceState(null, "", url);
    } catch (err) {
      /* URL non aggiornabile: la lingua scelta resta quella mostrata. */
    }
  }

  /* All'apertura si applica la lingua chiesta nell'URL (?orig=, ?lang=): è la
     scelta scritta nel link condiviso. Se il parametro non c'è o non corrisponde
     a nessuna lingua del brano, resta quella decisa alla generazione. */
  function applyLangParams() {
    var selects = document.querySelectorAll("[data-lang-select]");
    for (var i = 0; i < selects.length; i++) {
      applyLangParam(selects[i], LANG_PARAMS[selects[i].getAttribute("data-lang-select")]);
      showTextLang(selects[i], false);
    }
  }

  function applyLangParam(select, name) {
    var wanted = name ? param(name) : "";
    if (!wanted) {
      return;
    }
    for (var i = 0; i < select.options.length; i++) {
      if (select.options[i].value === wanted) {
        select.value = wanted;
        return;
      }
    }
  }

  function param(name) {
    try {
      return new URLSearchParams(window.location.search).get(name) || "";
    } catch (err) {
      return "";
    }
  }

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
  applyLangParams();
  showSidesWithoutSelect();
})();
