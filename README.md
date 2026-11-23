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

To preview the published site with this checkout as `nightly`:

```bash
task preview:build
task preview:serve
```

Open <http://localhost:1314/>

## Build

```bash
hugo mod get
hugo --gc --minify
```

Output is written to `public/`.

## Versioned docs

Versions follow mike's layout, so the legacy MkDocs builds and the Hugo builds share one site.
Each version has its own directory (`nightly/`, `v1.24.0/`, `latest/`), and both version selectors read the root `versions.json`.

The `cf-pages` branch holds the published site.
The Build and Publish Site workflow (`.github/workflows/hugo.yaml`) updates it.
A push to `main` rebuilds `nightly`.

Each tag, such as `v1.24.0`, marks the commit that documents that Watchtower release.
Pushing a tag builds it as its own version, and as `latest` when it is the newest release.

The workflow uses two composite actions in `.github/actions/`.
`build-versioned-docs` builds one version into a worktree of the site branch, and `publish-versioned-docs` commits and pushes it.

## Cloudflare Pages

| Setting                | Value      |
|------------------------|------------|
| Production branch      | `cf-pages` |
| Build command          | (none)     |
| Build output directory | `/`        |

Custom domain: `gowatchtower.org` (configure in Cloudflare Pages → Custom domains).

## Content layout

Documentation lives under `content/` with URL paths matching the legacy MkDocs site (no `/docs/` prefix), within each version's directory:

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

CI runs tplprev tests, golangci-lint, gosec, govulncheck, a WASM compile, and a Hugo build. The same workflow publishes the site to `cf-pages`.

## Theme

Managed via Hugo Modules (`github.com/imfing/hextra`). Do not add files under `themes/`.
