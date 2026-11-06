(function () {
  "use strict";

  function parseJobs(el) {
    var raw = el.getAttribute("data-jobs") || "";
    return raw.split(/\s+/).filter(Boolean);
  }

  function aaFiled(el) {
    return el.getAttribute("data-aa") === "filed";
  }

  function applyFilters(root) {
    var active = [];
    root.querySelectorAll("[data-roster-filter][aria-pressed='true']").forEach(function (btn) {
      active.push(btn.getAttribute("data-roster-filter"));
    });

    root.querySelectorAll("[data-roster-row]").forEach(function (row) {
      var show = true;
      if (active.indexOf("aa-filed") !== -1 && !aaFiled(row)) {
        show = false;
      }
      if (show && active.indexOf("job-coding") !== -1) {
        if (parseJobs(row).indexOf("coding") === -1) {
          show = false;
        }
      }
      row.hidden = !show;
      row.classList.toggle("roster-filter-hidden", !show);
    });
  }

  function initBar(bar) {
    bar.querySelectorAll("[data-roster-filter]").forEach(function (btn) {
      btn.addEventListener("click", function () {
        if (btn.disabled) {
          return;
        }
        var on = btn.getAttribute("aria-pressed") === "true";
        btn.setAttribute("aria-pressed", on ? "false" : "true");
        btn.classList.toggle("roster-filters__chip--active", !on);
        var root = bar.closest("[data-roster-filter-root]") || document;
        applyFilters(root);
      });
    });
  }

  document.querySelectorAll("[data-roster-filters]").forEach(initBar);
})();
