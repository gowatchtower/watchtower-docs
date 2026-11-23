# AGENTS.md

Hugo + [Hextra](https://github.com/imfing/hextra) docs site for Watchtower. Published at <https://gowatchtower.org/>.

Root `go.mod` / `go.sum` exist only for Hugo Modules (theme pin). Do not add a root `main.go` or treat this repo as a Go application.

`tools/tplprev/` is a nested `github.com/gowatchtower/tplprev` module (CLI + WASM). Use `go -C tools/tplprev` for that module only.

## Commands

Hugo Extended **v0.146+** required (developed on **0.170**). Fetch the theme and build the template preview before any build or server:

```bash
hugo mod get
./scripts/build-tplprev.sh  # writes static/tplprev.wasm and static/wasm_exec.js
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

GitHub Actions (`.github/workflows/`, with the site's composite actions in `.github/actions/`):

- `hugo.yaml` builds the template preview and the site on pull requests, and publishes `nightly` and tagged versions to the `cf-pages` branch
- `test.yaml` runs the tplprev tests and vet
- `lint-go.yaml` — golangci-lint on `tools/tplprev`
- `lint-gh.yaml` — actionlint
- `security.yaml` — gosec + govulncheck on tplprev

Renovate config is `.github/renovate.json`.

Do not commit `public/`, `resources/`, `_vendor/`, or `.hugo_build.lock`.

## Deploy

The `cf-pages` branch holds the published site, with one directory per version beside `versions.json`. Cloudflare Pages deploys that branch as it is, with no build step. Do not merge it into `main` or commit source files to it.

`hugo.yaml` updates `cf-pages` through the `build-versioned-docs` and `publish-versioned-docs` actions. A push to `main` rebuilds `nightly`. A push of a `v*` tag builds that tag's checkout as its own version. Each version is built once, from its own tag, so later layout changes do not reach published versions.

## Versioning

Versioned docs use mike's layout. Each version has its own directory, such as `nightly/` or `v1.24.0/`, beside a `versions.json` of `{version, title, aliases}` entries. Legacy MkDocs builds, copied from the `gh-pages` branch of `nicholas-fedor/watchtower`, sit in the same tree. Their Material selector and this site's selector both read `../versions.json` and link to `../<version>/`, so each lists every version.

- `params.versioning.version` names the version a build is published as. It is empty by default, which hides the selector, and is set through `HUGO_PARAMS_VERSIONING_VERSION`.
- Each `v*` tag marks the commit that documents that Watchtower release, such as `v1.23.2`. Tag the commit that finishes a release's sync. Pushing the tag publishes it.
- `.github/actions/build-versioned-docs` builds the checkout as one version into a worktree of the site branch. It also rewrites `versions.json`. The newest release listed holds the `latest` alias, and is built a second time into `latest/`.
- `.github/actions/publish-versioned-docs` commits that worktree and pushes it to the site branch. The workflow runs it only for a push or manual run on `main` or a `v*` tag.
- Keep both versioned-docs actions free of Watchtower specifics, so other Hugo sites can use them. Pass those as inputs.
- `scripts/build-preview.sh` (`task preview:build`, then `task preview:serve` on <http://localhost:1314/>) copies `cf-pages` from origin into `public/` and builds this checkout over its `nightly` version. `scripts/` holds local tools only. CI logic belongs in the actions.

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

`{{< template-preview >}}` inlines `static/template-preview/form.html`. `static/tplprev.wasm` and `static/wasm_exec.js` are build outputs and are not committed. The `build-tplprev` action compiles them from `tools/tplprev/` in `hugo.yaml` before each site build, so every published version carries a preview built from its own source.

For local work, build them once, and again after changing `tools/tplprev/`. The Taskfile's `serve`, `build`, and `preview:build` tasks do this first.

```bash
./scripts/build-tplprev.sh
```

## Outputs

`hugo.yaml` emits extra formats used by the UI — do not drop them:

- home: `html`, `llms` → `llms.txt`
- pages/sections: `markdown` (copy-as-markdown / context menu)
