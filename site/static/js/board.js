/**
 * Personal check board: localStorage-backed kanban for catalog entries.
 * Shared by /board/ and entry-page add/remove controls.
 */
(function () {
  "use strict";

  var STORAGE_KEY = "xynova-model-kb.board.v1";
  var SCHEMA_VERSION = 1;
  var NOTE_MAX = 500;
  var STATUSES = ["to_check", "checking", "done"];
  var STATUS_LABELS = {
    to_check: "To check",
    checking: "Checking",
    done: "Done",
  };

  function emptyState() {
    return { version: SCHEMA_VERSION, items: {} };
  }

  function nowISO() {
    return new Date().toISOString();
  }

  function clampNote(note) {
    if (typeof note !== "string") return "";
    if (note.length <= NOTE_MAX) return note;
    return note.slice(0, NOTE_MAX);
  }

  function normalizeStatus(status) {
    if (STATUSES.indexOf(status) !== -1) return status;
    return "to_check";
  }

  function parseState(raw) {
    if (!raw) return emptyState();
    var data;
    try {
      data = JSON.parse(raw);
    } catch (e) {
      return emptyState();
    }
    if (!data || typeof data !== "object" || data.version !== SCHEMA_VERSION) {
      return emptyState();
    }
    if (!data.items || typeof data.items !== "object") {
      return emptyState();
    }
    var items = {};
    Object.keys(data.items).forEach(function (id) {
      if (!id || typeof id !== "string") return;
      var row = data.items[id];
      if (!row || typeof row !== "object") return;
      items[id] = {
        status: normalizeStatus(row.status),
        note: clampNote(typeof row.note === "string" ? row.note : ""),
        addedAt: typeof row.addedAt === "string" ? row.addedAt : nowISO(),
        updatedAt: typeof row.updatedAt === "string" ? row.updatedAt : nowISO(),
      };
    });
    return { version: SCHEMA_VERSION, items: items };
  }

  function loadState() {
    try {
      return parseState(localStorage.getItem(STORAGE_KEY));
    } catch (e) {
      return emptyState();
    }
  }

  function saveState(state, onError) {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
      return true;
    } catch (e) {
      if (typeof onError === "function") {
        onError(e);
      }
      return false;
    }
  }

  function text(el, value) {
    el.textContent = value == null ? "" : String(value);
  }

  function clear(el) {
    while (el.firstChild) el.removeChild(el.firstChild);
  }

  function el(tag, className, attrs) {
    var node = document.createElement(tag);
    if (className) node.className = className;
    if (attrs) {
      Object.keys(attrs).forEach(function (key) {
        if (attrs[key] == null) return;
        node.setAttribute(key, attrs[key]);
      });
    }
    return node;
  }

  function readCatalog() {
    var node = document.getElementById("board-catalog-data");
    if (!node) return [];
    try {
      var data = JSON.parse(node.textContent || "[]");
      // Defend against accidental double-encoding from HTML template escaping.
      if (typeof data === "string") {
        data = JSON.parse(data);
      }
      return Array.isArray(data) ? data : [];
    } catch (e) {
      return [];
    }
  }

  function catalogById(catalog) {
    var map = {};
    catalog.forEach(function (item) {
      if (item && item.id) map[item.id] = item;
    });
    return map;
  }

  /* ---------- Entry-page toggle ---------- */

  function initEntryControls() {
    var nodes = document.querySelectorAll("[data-board-entry]");
    if (!nodes.length) return;

    function refresh(control) {
      var id = control.getAttribute("data-board-id");
      var btn = control.querySelector("[data-board-toggle]");
      if (!btn || !id) return;
      var state = loadState();
      var onBoard = !!state.items[id];
      text(btn, onBoard ? "Remove from board" : "Add to board");
      btn.setAttribute("aria-pressed", onBoard ? "true" : "false");
    }

    nodes.forEach(function (control) {
      refresh(control);
      var btn = control.querySelector("[data-board-toggle]");
      if (!btn) return;
      btn.addEventListener("click", function () {
        var id = control.getAttribute("data-board-id");
        if (!id) return;
        var state = loadState();
        if (state.items[id]) {
          delete state.items[id];
        } else {
          var ts = nowISO();
          state.items[id] = {
            status: "to_check",
            note: "",
            addedAt: ts,
            updatedAt: ts,
          };
        }
        saveState(state);
        refresh(control);
      });
    });

    window.addEventListener("storage", function (ev) {
      if (ev.key !== STORAGE_KEY) return;
      nodes.forEach(refresh);
    });
  }

  /* ---------- Board page ---------- */

  function initBoardPage() {
    var root = document.querySelector("[data-board-root]");
    if (!root) return;

    root.hidden = false;
    var catalog = readCatalog();
    var byId = catalogById(catalog);
    var state = loadState();
    var messageEl = root.querySelector("[data-board-message]");
    var searchEl = root.querySelector("[data-board-search]");
    var sectionEl = root.querySelector("[data-board-section-filter]");
    var addSectionEl = root.querySelector("[data-board-add-section]");
    var sortEl = root.querySelector("[data-board-sort]");
    var catalogEl = root.querySelector("[data-board-catalog]");
    var dialogEl = root.querySelector("[data-board-add-dialog]");
    var openAddBtn = root.querySelector("[data-board-open-add]");
    var closeAddBtn = root.querySelector("[data-board-close-add]");
    var dragId = null;

    function setMessage(msg, isError) {
      if (!messageEl) return;
      text(messageEl, msg || "");
      messageEl.classList.toggle("board-message--error", !!isError);
    }

    function persist() {
      var ok = saveState(state, function () {
        setMessage("Could not save to local storage (quota or privacy mode).", true);
      });
      if (ok) setMessage("");
      return ok;
    }

    function matchesSection(entrySection, filter) {
      if (!filter || filter === "all") return true;
      return entrySection === filter;
    }

    function sortIds(ids) {
      var mode = sortEl ? sortEl.value : "updated";
      return ids.slice().sort(function (a, b) {
        var ia = state.items[a];
        var ib = state.items[b];
        var ca = byId[a];
        var cb = byId[b];
        if (mode === "title") {
          var ta = (ca && ca.title) || a;
          var tb = (cb && cb.title) || b;
          return ta.localeCompare(tb, undefined, { sensitivity: "base" });
        }
        var ka = mode === "added" ? ia.addedAt : ia.updatedAt;
        var kb = mode === "added" ? ib.addedAt : ib.updatedAt;
        if (ka === kb) return 0;
        return ka < kb ? 1 : -1;
      });
    }

    function renderCatalogList() {
      if (!catalogEl) return;
      clear(catalogEl);
      var q = (searchEl && searchEl.value ? searchEl.value : "").trim().toLowerCase();
      var sectionFilter = addSectionEl ? addSectionEl.value : "all";
      var shown = 0;
      catalog.forEach(function (item) {
        if (!item || !item.id) return;
        if (state.items[item.id]) return;
        if (!matchesSection(item.section, sectionFilter)) return;
        if (q) {
          var hay = [
            item.title,
            item.summary,
            item.section,
            item.catalogStatus,
            item.run,
            (item.tags || []).join(" "),
          ]
            .join(" ")
            .toLowerCase();
          if (hay.indexOf(q) === -1) return;
        }
        shown += 1;
        if (shown > 40) return;
        var li = el("li", "board-catalog__item");
        var meta = el("div", "board-catalog__meta");
        var title = el("span", "board-catalog__title");
        text(title, item.title);
        var side = el("span", "board-catalog__section");
        text(side, item.section);
        meta.appendChild(title);
        meta.appendChild(side);
        if (item.summary) {
          var sum = el("p", "board-catalog__summary");
          text(sum, item.summary);
          meta.appendChild(sum);
        }
        var addBtn = el("button", "board-btn board-btn--small", { type: "button" });
        text(addBtn, "Add");
        addBtn.addEventListener("click", function () {
          var ts = nowISO();
          state.items[item.id] = {
            status: "to_check",
            note: "",
            addedAt: ts,
            updatedAt: ts,
          };
          if (persist()) {
            render();
            setMessage('Added "' + (item.title || "entry") + '" to To check.');
          }
        });
        li.appendChild(meta);
        li.appendChild(addBtn);
        catalogEl.appendChild(li);
      });
      if (!catalogEl.firstChild) {
        var empty = el("li", "board-catalog__empty");
        text(empty, q ? "No matching catalog entries to add." : "All matching entries are already on the board.");
        catalogEl.appendChild(empty);
      }
    }

    function openAddDialog() {
      if (!dialogEl || typeof dialogEl.showModal !== "function") return;
      renderCatalogList();
      dialogEl.showModal();
      if (searchEl) {
        searchEl.focus();
        searchEl.select();
      }
    }

    function closeAddDialog() {
      if (!dialogEl || typeof dialogEl.close !== "function") return;
      if (dialogEl.open) dialogEl.close();
      if (openAddBtn) openAddBtn.focus();
    }

    function makeCard(id, row) {
      var meta = byId[id];
      var card = el("article", "board-card", {
        draggable: "true",
        "data-board-card": id,
      });
      if (!meta) card.classList.add("board-card--orphan");

      var head = el("div", "board-card__head");
      var title = el("a", "board-card__title", {
        href: meta ? meta.url : id,
      });
      text(title, meta ? meta.title : "Missing catalog entry");
      head.appendChild(title);

      var badges = el("div", "board-card__badges");
      var sec = el("span", "board-chip");
      text(sec, meta ? meta.section : "orphan");
      badges.appendChild(sec);
      if (meta && meta.catalogStatus) {
        var st = el("span", "board-chip board-chip--muted");
        text(st, "catalog: " + meta.catalogStatus);
        badges.appendChild(st);
      }
      if (meta && meta.run) {
        var run = el("span", "board-chip board-chip--muted");
        text(run, "run: " + meta.run);
        badges.appendChild(run);
      }
      if (!meta) {
        var warn = el("span", "board-chip board-chip--warn");
        text(warn, "not in current build");
        badges.appendChild(warn);
      }

      var statusLabel = el("label", "board-field board-field--compact");
      var statusText = el("span", "board-field__label");
      text(statusText, "Personal check status");
      var statusSelect = el("select", "board-select", { "data-board-status": id });
      STATUSES.forEach(function (s) {
        var opt = el("option", null, { value: s });
        text(opt, STATUS_LABELS[s]);
        if (s === row.status) opt.selected = true;
        statusSelect.appendChild(opt);
      });
      statusSelect.addEventListener("change", function () {
        row.status = normalizeStatus(statusSelect.value);
        row.updatedAt = nowISO();
        if (persist()) render();
      });
      statusLabel.appendChild(statusText);
      statusLabel.appendChild(statusSelect);

      var noteLabel = el("label", "board-field board-field--compact");
      var noteText = el("span", "board-field__label");
      text(noteText, "Personal note");
      var noteArea = el("textarea", "board-textarea", {
        rows: "2",
        maxlength: String(NOTE_MAX),
        "data-board-note": id,
        placeholder: "What to verify next",
      });
      noteArea.value = row.note || "";
      var noteTimer = null;
      noteArea.addEventListener("input", function () {
        clearTimeout(noteTimer);
        noteTimer = setTimeout(function () {
          row.note = clampNote(noteArea.value);
          row.updatedAt = nowISO();
          persist();
        }, 300);
      });
      noteLabel.appendChild(noteText);
      noteLabel.appendChild(noteArea);

      var actions = el("div", "board-card__actions");
      var removeBtn = el("button", "board-btn board-btn--small board-btn--danger", { type: "button" });
      text(removeBtn, "Remove");
      removeBtn.addEventListener("click", function () {
        delete state.items[id];
        if (persist()) render();
      });
      actions.appendChild(removeBtn);

      card.appendChild(head);
      card.appendChild(badges);
      card.appendChild(statusLabel);
      card.appendChild(noteLabel);
      card.appendChild(actions);

      card.addEventListener("dragstart", function (ev) {
        dragId = id;
        card.classList.add("board-card--dragging");
        if (ev.dataTransfer) {
          ev.dataTransfer.setData("text/plain", id);
          ev.dataTransfer.effectAllowed = "move";
        }
      });
      card.addEventListener("dragend", function () {
        dragId = null;
        card.classList.remove("board-card--dragging");
        root.querySelectorAll(".board-column__cards--over").forEach(function (n) {
          n.classList.remove("board-column__cards--over");
        });
      });

      return card;
    }

    function renderColumns() {
      var sectionFilter = sectionEl ? sectionEl.value : "all";
      STATUSES.forEach(function (status) {
        var drop = root.querySelector('[data-board-drop="' + status + '"]');
        var countEl = root.querySelector('[data-board-count="' + status + '"]');
        if (!drop) return;
        clear(drop);
        var ids = Object.keys(state.items).filter(function (id) {
          var row = state.items[id];
          if (!row || row.status !== status) return false;
          var meta = byId[id];
          var section = meta ? meta.section : null;
          if (section) return matchesSection(section, sectionFilter);
          return sectionFilter === "all";
        });
        ids = sortIds(ids);
        if (countEl) text(countEl, String(ids.length));
        if (!ids.length) {
          var empty = el("p", "board-column__empty");
          text(empty, "No cards");
          drop.appendChild(empty);
          return;
        }
        ids.forEach(function (id) {
          drop.appendChild(makeCard(id, state.items[id]));
        });
      });
    }

    function wireDrops() {
      root.querySelectorAll("[data-board-drop]").forEach(function (drop) {
        drop.addEventListener("dragover", function (ev) {
          ev.preventDefault();
          drop.classList.add("board-column__cards--over");
          if (ev.dataTransfer) ev.dataTransfer.dropEffect = "move";
        });
        drop.addEventListener("dragleave", function () {
          drop.classList.remove("board-column__cards--over");
        });
        drop.addEventListener("drop", function (ev) {
          ev.preventDefault();
          drop.classList.remove("board-column__cards--over");
          var id = dragId || (ev.dataTransfer && ev.dataTransfer.getData("text/plain"));
          var status = drop.getAttribute("data-board-drop");
          if (!id || !state.items[id] || !status) return;
          state.items[id].status = normalizeStatus(status);
          state.items[id].updatedAt = nowISO();
          if (persist()) render();
        });
      });
    }

    function render() {
      if (dialogEl && dialogEl.open) {
        renderCatalogList();
      }
      renderColumns();
    }

    if (searchEl) searchEl.addEventListener("input", renderCatalogList);
    if (addSectionEl) addSectionEl.addEventListener("change", renderCatalogList);
    if (sectionEl) sectionEl.addEventListener("change", render);
    if (sortEl) sortEl.addEventListener("change", render);

    if (openAddBtn) {
      openAddBtn.addEventListener("click", openAddDialog);
    }
    if (closeAddBtn) {
      closeAddBtn.addEventListener("click", closeAddDialog);
    }
    if (dialogEl) {
      dialogEl.addEventListener("click", function (ev) {
        if (ev.target === dialogEl) closeAddDialog();
      });
      dialogEl.addEventListener("cancel", function () {
        if (openAddBtn) {
          window.setTimeout(function () {
            openAddBtn.focus();
          }, 0);
        }
      });
    }

    wireDrops();
    render();

    window.addEventListener("storage", function (ev) {
      if (ev.key !== STORAGE_KEY) return;
      state = loadState();
      render();
      setMessage("Board updated from another tab.");
    });
  }

  function boot() {
    initEntryControls();
    initBoardPage();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot);
  } else {
    boot();
  }
})();
