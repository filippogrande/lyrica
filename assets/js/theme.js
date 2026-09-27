/* Applica il tema PRIMA del primo paint, per non mostrare un lampo di tema
   sbagliato (DESIGN_SYSTEM.md). Va caricato in modo sincrono nell'head.

   Nota: se localStorage è bloccato dal browser la preferenza salvata non è
   leggibile e si applica il tema di sistema. È il comportamento previsto, non
   un errore nascosto. */
(function () {
  "use strict";

  var KEY = "lyrica.theme";
  var saved = null;

  try {
    saved = window.localStorage.getItem(KEY);
  } catch (err) {
    saved = null;
  }

  var theme = saved === "light" || saved === "dark" ? saved : systemTheme();
  document.documentElement.setAttribute("data-bs-theme", theme);

  function systemTheme() {
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
})();
