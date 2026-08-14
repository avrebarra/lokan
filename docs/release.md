# Releases

Prebuilt binaries are published to GitHub releases by tagging. The pipeline
mirrors `git-workpool`'s release routine: tag → GitHub Actions → GoReleaser →
`tar.gz` assets + checksums + release notes.

## How to cut a release

1. Bump nothing by hand — the tag *is* the version.
2. From `main`, cut the release with the bumper:

   ```sh
   ./runtask bump [patch|minor|major]   # tags the next version and pushes it
   ./runtask bump dry                   # preview the next tag, no changes
   ```

   `bump` reads the highest existing `v*` tag, bumps the requested level
   (bare `bump` = patch), refuses on a dirty tree, then tags and pushes.
   Tagging manually works the same: `git tag v0.1.0 && git push origin v0.1.0`.

3. The `release` workflow (`.github/workflows/release.yml`) runs GoReleaser,
   which builds `lokan` for darwin/linux × amd64/arm64, archives each as a
   `tar.gz` (binary only), writes `checksums.txt`, and creates a GitHub
   release with the assets attached.

## What you get

| Asset | Content |
|---|---|
| `lokan_darwin_amd64.tar.gz` / `lokan_darwin_arm64.tar.gz` | macOS binary |
| `lokan_linux_amd64.tar.gz` / `lokan_linux_arm64.tar.gz` | Linux binary |
| `checksums.txt` | SHA-256 sums for every asset |

Users install from `latest`:

```sh
curl -fsSL https://raw.githubusercontent.com/avressatelier/lokan/main/install.sh | sh
```

`install.sh` maps the machine's OS/arch to the matching asset, extracts the
binary into `~/.local/bin` (override with `LOKAN_INSTALL_DIR`), and warns when
that dir isn't on `PATH`.

## Config files

| File | Role |
|---|---|
| `.goreleaser.yaml` | Build matrix, archive naming, release target (`dir: ./engine`, `main: ./cmd/lokan` — matches `runtask build`) |
| `.github/workflows/release.yml` | Tag-triggered CI: checkout → Go → GoReleaser `release --clean` |
| `install.sh` | One-liner installer used in the README |

## Gotchas

- **Versioning:** there is no version flag in the binary; the tag is the
  single source of truth. Do not invent a separate `--version` number.
- **The release needs `engine/go.mod`** for `setup-go` (go-version-file) — the
  module root is `engine/`, not the repo root.
- **Dry-run locally:** `goreleaser release --clean --snapshot` builds the
  archives without publishing.
