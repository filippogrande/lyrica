// Ricerca live (D110): indice JSON generato a build time, filtro lato client.
// Nessun framework, nessun endpoint: il dropdown cerca nell'indice scaricato
// una volta. La normalizzazione (minuscole, accenti, ß) è la stessa su query e
// testo, e vale anche per chi scrive "Fur Elise" cercando "Für Elise".
//
// Il campo e il dropdown stanno nel DOM ma si mostrano solo con la classe .js
// (messa da theme.js prima del primo paint): senza JS non c'è niente da
// nascondere.
(function () {
  "use strict";

  var box = document.querySelector(".search-box");
  if (!box) {
    return;
  }
  var input = document.getElementById("search-input");
  var list = document.getElementById("search-results");
  var lang = box.getAttribute("data-lang") || "it";
  var emptyText = list.getAttribute("data-empty") || "";
  var MAX = 8;

  // Normalizza per il confronto: minuscole, accenti separati e rimossi, ß
  // sciolto come nella traslitterazione del repo.
  function strip(s) {
    return s
      .toLowerCase()
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/ß/g, "ss");
  }

  var index = null;
  var items = [];
  var active = -1;
  var timer = null;

  function loadIndex() {
    if (index) {
      return Promise.resolve(index);
    }
    return fetch("/" + lang + "/search-index.json", { cache: "force-cache" })
      .then(function (res) {
        if (!res.ok) {
          return { tracks: [], bands: [], albums: [] };
        }
        return res.json();
      })
      .then(function (data) {
        // Pre-normalizza le righe una volta sola: il confronto a ogni tasto
        // non deve rifare il lavoro per ogni verso.
        index = {
          bands: data.bands.map(function (b) {
            return { kind: "band", name: b.name, path: b.path, norm: strip(b.name) };
          }),
          albums: data.albums.map(function (a) {
            return { kind: "album", title: a.title, band: a.band, path: a.path, norm: strip(a.title + " " + a.band) };
          }),
          tracks: data.tracks.map(function (t) {
            return {
              kind: "track",
              title: t.title,
              band: t.band,
              album: t.album,
              path: t.path,
              lines: t.lines,
              norm: t.lines.map(strip),
            };
          }),
        };
        return index;
      });
  }

  function search(q) {
    var out = [];
    index.bands.forEach(function (b) {
      if (b.norm.indexOf(q) !== -1) {
        out.push(b);
      }
    });
    index.albums.forEach(function (a) {
      if (a.norm.indexOf(q) !== -1) {
        out.push(a);
      }
    });
    index.tracks.forEach(function (t) {
      if (out.length >= MAX) {
        return;
      }
      if (strip(t.title).indexOf(q) !== -1 || strip(t.band).indexOf(q) !== -1 || strip(t.album).indexOf(q) !== -1) {
        out.push({ kind: t.kind, title: t.title, band: t.band, album: t.album, path: t.path, verse: "" });
        return;
      }
      for (var i = 0; i < t.norm.length; i++) {
        if (t.norm[i].indexOf(q) !== -1) {
          out.push({ kind: t.kind, title: t.title, band: t.band, album: t.album, path: t.path, verse: t.lines[i] });
          break;
        }
      }
    });
    return out.slice(0, MAX);
  }

  function render() {
    list.textContent = "";
    active = -1;
    if (items.length === 0) {
      var empty = document.createElement("li");
      empty.className = "search-empty";
      empty.textContent = emptyText;
      list.appendChild(empty);
      return;
    }
    items.forEach(function (item, i) {
      var li = document.createElement("li");
      li.className = "search-item";
      var a = document.createElement("a");
      a.href = item.path;
      var title = document.createElement("span");
      title.className = "search-item-title";
      title.textContent = item.title || item.name;
      a.appendChild(title);
      if (item.album) {
        var meta = document.createElement("span");
        meta.className = "search-item-meta";
        meta.textContent = " " + item.band + " · " + item.album;
        a.appendChild(meta);
      } else if (item.band) {
        var meta2 = document.createElement("span");
        meta2.className = "search-item-meta";
        meta2.textContent = " " + item.band;
        a.appendChild(meta2);
      }
      if (item.verse) {
        var verse = document.createElement("span");
        verse.className = "search-item-verse";
        verse.textContent = "\u201c" + item.verse + "\u201d";
        a.appendChild(verse);
      }
      li.appendChild(a);
      li.setAttribute("data-index", String(i));
      list.appendChild(li);
    });
  }

  function setActive(i) {
    var els = list.querySelectorAll(".search-item");
    if (els.length === 0) {
      return;
    }
    if (i < 0) {
      i = 0;
    }
    if (i >= els.length) {
      i = els.length - 1;
    }
    active = i;
    els.forEach(function (el, j) {
      el.classList.toggle("is-active", j === active);
    });
  }

  function openDropdown() {
    list.removeAttribute("hidden");
    input.setAttribute("aria-expanded", "true");
  }

  function closeDropdown() {
    list.setAttribute("hidden", "");
    input.setAttribute("aria-expanded", "false");
    active = -1;
  }

  input.addEventListener("input", function () {
    var q = strip(input.value.trim());
    if (q.length < 2) {
      closeDropdown();
      return;
    }
    clearTimeout(timer);
    timer = setTimeout(function () {
      loadIndex().then(function () {
        items = search(q);
        render();
        openDropdown();
        setActive(0);
      });
    }, 180);
  });

  input.addEventListener("focus", function () {
    if (items.length > 0) {
      openDropdown();
    }
  });

  input.addEventListener("keydown", function (e) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActive(active + 1);
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActive(active - 1);
    } else if (e.key === "Escape") {
      closeDropdown();
    } else if (e.key === "Enter") {
      var el = list.querySelector(".search-item.is-active a");
      if (el) {
        e.preventDefault();
        window.location.href = el.getAttribute("href");
      }
    }
  });

  document.addEventListener("click", function (e) {
    if (!box.contains(e.target)) {
      closeDropdown();
    }
  });
})();
