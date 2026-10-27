# Watchtower Documentation

This is the source repository for [Watchtower's](https://github.com/nicholas-fedor/watchtower) documentation site.

The site is published via Cloudflare Pages and is available at [https://gowatchtower.org/](https://gowatchtower.org/).
It is built with [Hugo](https://gohugo.io/) and the [Hextra](https://github.com/imfing/hextra) theme.

## Prerequisites

- [Hugo Extended](https://gohugo.io/installation/)
- [Go](https://go.dev/) (for Hugo Modules)

## Local development

```bash
hugo mod get
hugo server
```

Open <http://localhost:1313/>

## Build

```bash
hugo mod get
hugo --gc --minify
```

Output is written to `public/`.

## Cloudflare Pages

| Setting                | Value                                                    |
|------------------------|----------------------------------------------------------|
| Production branch      | `main`                                                   |
| Build command          | `hugo mod get && hugo --gc --minify`                     |
| Build output directory | `public`                                                 |
| Environment variable   | `HUGO_VERSION` = `0.167.0` (or current Extended release) |

Hugo Modules require Go on the build image.
Cloudflare Pages provides Go by default for Hugo projects.
If module fetch fails, set `GO_VERSION` (e.g. `1.27.x`) in the project environment variables.

Custom domain: `gowatchtower.org` (configure in Cloudflare Pages → Custom domains).

## Versioned docs

Versioned docs follow mike's layout, so legacy MkDocs builds and Hugo builds share one site.
Each version lives in its own directory (`nightly/`, `v1.24.0/`, `latest/`).
A `versions.json` at the site root lists them, and the version selectors on both kinds of page read it.

`scripts/publish-version.sh` builds this checkout as one version and records it in `versions.json`:

```bash
./scripts/publish-version.sh --tree public --version v1.24.0 --base-url https://watchtower.nickfedor.com --alias latest
```

To preview the legacy versions together with this checkout as `nightly`, which replaces the legacy rolling `dev` build:

```bash
task preview:build
task preview:serve
```

Open <http://localhost:1314/>

The experimental `feat/hugh` branch is built by its own Cloudflare Pages project:

| Setting                | Value                        |
|------------------------|------------------------------|
| Production branch      | `feat/hugh`                  |
| Build command          | `./scripts/build-preview.sh` |
| Build output directory | `public`                     |
| Environment variable   | `HUGO_VERSION` = `0.167.0`   |

## Content layout

Documentation lives under `content/` with URL paths matching the legacy MkDocs site (no `/docs/` prefix):

- `/quickstart/`
- `/getting-started/`
- `/configuration/`
- `/notifications/`
- `/http-api/` (nav label: API)
- `/advanced-features/`

Static assets are in `static/` (logo, favicon, WASM template preview).

## Template preview

The notification template preview is a nested Go module at `tools/tplprev/` (`github.com/gowatchtower/tplprev`).
The `template-preview` shortcode inlines `static/template-preview/form.html`.

Rebuild the checked-in WASM after changing the module:

```bash
./scripts/build-tplprev.sh
```

Run module tests with `go -C tools/tplprev test ./...`.

CI runs tplprev tests, golangci-lint, gosec, govulncheck, a WASM compile, and a Hugo build. Deploy stays on Cloudflare Pages.

## Theme

Managed via Hugo Modules (`github.com/imfing/hextra`). Do not add files under `themes/`.
