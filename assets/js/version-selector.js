// Fills the navbar version selector from versions.json in mike's format. It resolves URLs the way the
// MkDocs Material selector on legacy pages does, so both list and link versions identically:
// versions.json and every version live one level above this build's root, at ../versions.json
// and ../<version>/.
(function () {
  const selector = document.querySelector(".wt-version-selector");
  if (!selector) return;

  const label = selector.querySelector(".wt-version-label");
  const menu = selector.querySelector(".hextra-language-options");
  const check = selector.querySelector(".wt-version-check");
  const current = selector.dataset.version;
  const base = new URL(selector.dataset.base, window.location.href);

  // Longest prefix shared by every string, or "" for none.
  const commonPrefix = (strings) =>
    strings.reduce((prefix, s) => {
      let i = 0;
      while (i < prefix.length && prefix[i] === s[i]) i++;
      return prefix.slice(0, i);
    });

  // Returns the current page in the target version when its sitemap lists that page, and the target's
  // root otherwise. Sitemap entries can carry another host (legacy builds list
  // watchtower.nickfedor.com), so the page is looked up relative to the sitemap's own common prefix.
  const correspondingPage = async (target) => {
    const { pathname, search, hash } = window.location;
    if (!pathname.startsWith(base.pathname)) return target;
    const relative = pathname.slice(base.pathname.length);

    try {
      const response = await fetch(new URL("sitemap.xml", target));
      if (!response.ok) return target;
      const xml = new DOMParser().parseFromString(await response.text(), "application/xml");
      const locs = new Set(Array.from(xml.querySelectorAll("loc"), (loc) => loc.textContent.trim()));
      if (locs.size === 0) return target;

      const prefix = commonPrefix([...locs]);
      if (!locs.has(prefix) || !locs.has(new URL(relative, prefix).href)) return target;

      const page = new URL(relative, target);
      page.search = search;
      page.hash = hash;
      return page;
    } catch {
      return target;
    }
  };

  const render = (versions) => {
    for (const { version, title, aliases } of versions) {
      const target = new URL(`../${version}/`, base);

      const link = document.createElement("a");
      link.href = target.href;
      link.setAttribute("role", "menuitem");
      link.className =
        "hx:text-gray-700 hx:dark:text-gray-300 hx:hover:bg-gray-100 hx:hover:text-gray-900 hx:dark:hover:bg-neutral-800 hx:dark:hover:text-gray-100 hx:relative hx:cursor-pointer hx:whitespace-nowrap hx:rounded-sm hx:py-1.5 hx:transition-colors hx:ltr:pl-3 hx:ltr:pr-9 hx:rtl:pr-3 hx:rtl:pl-9";
      link.textContent = title;

      for (const alias of aliases) {
        const badge = document.createElement("span");
        badge.className = "wt-version-alias";
        badge.textContent = alias;
        link.append(badge);
      }

      if (version === current || aliases.includes(current)) {
        link.setAttribute("aria-current", "true");
        link.append(check.content.cloneNode(true));
        label.textContent = title;
      }

      // Modified clicks keep the browser's default, such as opening the version root in a new tab.
      link.addEventListener("click", async (event) => {
        if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
        event.preventDefault();
        window.location.assign(await correspondingPage(target));
      });

      const item = document.createElement("li");
      item.setAttribute("role", "none");
      item.className = "hx:flex hx:flex-col";
      item.append(link);
      menu.append(item);
    }
  };

  fetch(new URL("../versions.json", base))
    .then((response) => {
      if (!response.ok) throw new Error(`versions.json: HTTP ${response.status}`);
      return response.json();
    })
    .then(render)
    .catch((error) => {
      console.warn("Version selector disabled:", error);
      selector.classList.add("hx:hidden");
    });
})();
