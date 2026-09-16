# the-model-kb site

## Project overview

Hugo documentation site that publishes three repository catalogs:

| Section | Source of truth | Public path |
|---------|-----------------|-------------|
| Library | `library/` | `/library/` |
| Packaged models | `models/` Markdown docs + site overlays | `/models/` |
| Roster | `roster/` (via build-time remap) | `/roster/` |

Markdown catalogs stay outside `site/`. The site mounts them (or a generated roster tree) without treating Dockerfiles, Bake files, or weights as content.

## Core functionalities

- Mount `library/` at `/library/`
- Mount Markdown from `models/` at `/models/` (`includeFiles = ["**/*.md"]` only)
- Generate roster pages with `kind` remapped to `roster_kind` (Hugo reserves `kind`), then mount under `/roster/`
- Exclude catalog authoring files (`README.md`, `INDEX.md`, `taxonomy.md`, `_template.md`)
- Section overlays and package pages under `site/content/`
- Offline FlexSearch via Hextra
- Library section indexes: responsive skim rows with does / run / license chips (`layouts/_partials/library-section-index.html`, `static/css/library-catalog.css`)
- Library entry skim strip under the title (`layouts/_partials/library-entry-meta.html`)
- Library entry stills: GLightbox gallery grid injected under **What it is** from frontmatter `media` (`layouts/_partials/library-gallery.html`, vendored under `static/vendor/glightbox/`)
- Root `make serve` / `make build` / `make clean`
- GitHub Pages deploy from `.github/workflows/library-site.yml`

## Docs and libraries

- [Hugo documentation](https://gohugo.io/documentation/)
- [Hugo modules / mounts](https://gohugo.io/configuration/module/)
- [Hextra theme](https://imfing.github.io/hextra/docs/) (pinned module `github.com/imfing/hextra`)
- Library taxonomy: [`library/taxonomy.md`](../library/taxonomy.md)
- Roster taxonomy: [`roster/taxonomy.md`](../roster/taxonomy.md)
- Aggregate bake: [`models/docker-bake.hcl`](../models/docker-bake.hcl)

## Current file structure

```text
site/
  hugo.toml
  go.mod / go.sum
  Makefile                 # generate-roster, serve, build, clean
  scripts/generate-roster.py
  instructions.md
  content/
    _index.md              # Homepage
    library/<category>/    # Library section overlays
    models/<package>/      # Package pages
    roster/<category>/     # Roster section overlays
  generated/roster/        # Build output (gitignored)
  public/                  # Build output (gitignored)

library/                   # Mounted at content/library
roster/                    # Source for generated/roster
models/                    # Markdown docs mounted at content/models
  docker-bake.hcl          # Aggregate bake (invoke from repo root)
  **/INVOKE.md             # Published ops guides

.github/workflows/library-site.yml
```

## Local commands

From the repository root:

```bash
make serve   # http://localhost:1313/
make build   # production build into site/public/
make clean   # remove site/public/, resources/, generated/
```

Packaging (also from repo root):

```bash
make docker-print
make docker-push
make docker-push-target TARGET=gemma4-12B-serve
```

Hugo Extended, Go, and Python 3 are required for the site build (modules + roster remap).

## Why roster is generated

Hugo treats frontmatter `kind` as the page kind. Roster cards use `kind` for `sovereign` / `provider` / `service`. `scripts/generate-roster.py` copies `roster/` into `generated/roster/` and rewrites that key to `roster_kind` so the catalog files stay the authoring source of truth.

## GitHub Pages

After the first merge to `main`, set the repository Pages source to **GitHub Actions** (Settings → Pages → Build and deployment). The workflow builds from `site/` whenever `library/`, `roster/`, `models/`, or `site/` changes. The live URL is typically `https://xynova.github.io/the-model-kb/`.
