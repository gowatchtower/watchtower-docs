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
| Environment variable   | `HUGO_VERSION` = `0.164.0` (or current Extended release) |

Hugo Modules require Go on the build image.
Cloudflare Pages provides Go by default for Hugo projects.
If module fetch fails, set `GO_VERSION` (e.g. `1.26.x`) in the project environment variables.

Custom domain: `gowatchtower.org` (configure in Cloudflare Pages → Custom domains).

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
