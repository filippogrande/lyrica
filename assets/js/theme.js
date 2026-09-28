/* Applica il tema PRIMA del primo paint, per non mostrare un lampo di tema
   sbagliato (docs/FRONTEND.md). Va caricato in modo sincrono nell'head.

   Mette anche la classe .js su <html>: da lì il CSS sa che lo script c'è, e
   regole come quella delle lingue del brano (che nasconde le lingue non scelte)
   valgono solo con JS. Senza JS il testo resta tutto visibile.

   Nota: se localStorage è bloccato dal browser la preferenza salvata non è
   leggibile e si applica il tema di sistema. È il comportamento previsto, non
   un errore nascosto. */
(function () {
  "use strict";

  var KEY = "lyrica.theme";
  var saved = null;

  document.documentElement.classList.add("js");

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
