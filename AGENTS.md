# AGENTS.md

Hugo + [Hextra](https://github.com/imfing/hextra) docs site for Watchtower. Published at <https://gowatchtower.org/>.

Root `go.mod` / `go.sum` exist only for Hugo Modules (theme pin). Do not add a root `main.go` or treat this repo as a Go application.

`tools/tplprev/` is a nested `github.com/gowatchtower/tplprev` module (CLI + WASM). Use `go -C tools/tplprev` for that module only.

## Commands

Hugo Extended **v0.146+** required (developed on **0.164**). Fetch the theme before any build or server:

```bash
hugo mod get
hugo server                 # http://localhost:1313/
hugo --gc --minify          # writes public/
```

Verify the site with `hugo --gc --minify`. For tplprev:

```bash
go -C tools/tplprev test -timeout 60s ./...
go -C tools/tplprev vet ./...
golangci-lint run --timeout 60s ./...   # from tools/tplprev
./scripts/build-tplprev.sh
```

GitHub Actions (`.github/workflows/`):

- `hugo.yaml` — Hugo Extended build
- `test.yaml` — tplprev tests, vet, and WASM build
- `lint-go.yaml` — golangci-lint on `tools/tplprev`
- `lint-gh.yaml` — actionlint
- `security.yaml` — gosec + govulncheck on tplprev

Renovate config is `.github/renovate.json`. Deploy stays on Cloudflare Pages.

Do not commit `public/`, `resources/`, `_vendor/`, or `.hugo_build.lock`.

## Deploy

Cloudflare Pages (`main`):

- Build: `hugo mod get && hugo --gc --minify`
- Output: `public`
- `HUGO_VERSION=0.167.0` (Extended). Go is needed for module fetch; set `GO_VERSION` if fetch fails.

## Versioning

Versioned docs use mike's layout: each version in its own directory (`nightly/`, `v1.24.0/`, `latest/`) beside a `versions.json` of `{version, title, aliases}` entries. Legacy MkDocs builds from the `gh-pages` branch of `nicholas-fedor/watchtower` sit in the same tree unchanged. Their Material selector and this site's selector both read `../versions.json` and link to `../<version>/`, so each lists every version.

- `params.versioning.version` names the version a build is published as. It is empty by default, which hides the selector, and is set through `HUGO_PARAMS_VERSIONING_VERSION`.
- `scripts/publish-version.sh --tree DIR --version V --base-url URL [--alias latest] [--replaces OLD]` builds this checkout into `DIR/V/` and updates `DIR/versions.json`. With `--alias`, it moves the alias to V and builds a second full copy into `DIR/<alias>/`. With `--replaces`, it removes OLD's directory and entry and moves OLD's aliases to V. It replaces `mike deploy`.
- `scripts/build-preview.sh` (`task preview:build`, then `task preview:serve` on <http://localhost:1314/>) assembles the legacy tree plus this checkout as `nightly` into `public/`. `nightly` replaces the legacy rolling `dev` build.

The experimental `feat/hugh` branch is built by its own Cloudflare Pages project. That project uses production branch `feat/hugh`, build command `./scripts/build-preview.sh`, output `public`, and `HUGO_VERSION=0.167.0`.

## Theme and layouts

Theme is `github.com/imfing/hextra` v0.13.0 via Hugo Modules. `themes/` is unused — do not add files there.

Hugo 0.146+ layout paths use underscore prefixes. Overrides belong here, not under `themes/`:

- `layouts/_shortcodes/` — site shortcodes
- `layouts/_partials/` — Hextra overrides (`custom/head-end.html`, `components/page-context-menu.html`, `navbar-title.html`)
- `layouts/hextra-home.html` — homepage layout (`content/_index.md` sets `layout: hextra-home`)
- `assets/css/custom.css` — Hextra auto-loads this; site is **dark-only** (`params.theme.displayToggle: false`). Do not add a light theme or toggle.

`navbar-title.html` copies Hextra's partial and adds the version selector (`assets/js/version-selector.js`). The selector reuses Hextra's `hextra-language-switcher` and `hextra-language-options` classes, so the theme's `lang.js` drives the menu. Recheck both after a Hextra upgrade.

`markup.goldmark.renderer.unsafe: true` is required for raw HTML on the homepage.

## Content

Legacy MkDocs URL paths: **no `/docs/` prefix**. Pages are bundles:

| Kind | Path | Front matter |
| --- | --- | --- |
| Section | `content/<section>/_index.md` | `type: docs` + `cascade.type: docs`, `weight` |
| Page | `content/<section>/<page>/index.md` | `title` + `weight` only |

New top-level sections need all three: `content/<section>/_index.md`, a `menu.main` entry in `hugo.yaml`, and a card on `content/_index.md`.

Prefer Hextra shortcodes with `{{< >}}`: `tabs`/`tab`, `callout`, `cards`/`card`. Nested tabs are used (see `content/quickstart/`).

Versioned builds are served from a subpath such as `/nightly/`, so links and asset URLs must not start with `/`. Use relative links in content and shortcode parameters (`link="quickstart/"`), and `relURL` without a leading slash in templates. Hugo's `relURL` leaves a path that starts with `/` relative to the host root.

## Product references

Keep these distinct:

- Docker image: `nickfedor/watchtower`
- Source repo: `nicholas-fedor/watchtower` (also Shoutrrr: `nicholas-fedor/shoutrrr`, docs at <https://shoutrrr.nickfedor.com>)
- This docs repo: `gowatchtower/watchtower-docs`

Do not point readers at `containrrr/watchtower` (archived).

## Template preview

`{{< template-preview >}}` inlines `static/template-preview/form.html`. Rebuild WASM after changing `tools/tplprev/`:

```bash
./scripts/build-tplprev.sh
```

Commit the updated `static/tplprev.wasm` and `static/wasm_exec.js`. Cloudflare Pages does not compile WASM.

## Outputs

`hugo.yaml` emits extra formats used by the UI — do not drop them:

- home: `html`, `llms` → `llms.txt`
- pages/sections: `markdown` (copy-as-markdown / context menu)
