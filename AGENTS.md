# zplkit — Agent / Contributor Instructions

## Collaboration Philosophy

Development is a collaborative process. Agents are expected to exercise independent
technical judgment, not passive agreement. When a different solution better aligns with
industry standards or formal specifications (e.g. ZPL semantics, Go idioms), agents
should raise concerns and advocate for the more appropriate approach.

## Project Overview

`zplkit` renders a ZPL label template for an asset tag and sends it to a network label
printer (raw TCP or IPP — see `internal/printer`), and separately can send a raw `.zpl`
file straight through for printer-side provisioning (`admin`).

- Written in **Go 1.25**. Always **`CGO_ENABLED=0`**. Cross-platform (linux/darwin/
  windows) — this is a one-shot CLI, not a daemon, so there's no systemd/install story
  like `config-watch`/`anycast-sentinel` have.
- Produces a **single static binary** from `cmd/zplkit/`.
- **The repo's real contribution surface is `templates/`, not Go code.** Adding a label
  layout is "drop a `.zpl` file + a preview PNG, maybe add a TOML alias" — see the
  README's "Adding a template" section. Don't make template changes require touching Go.
- CLI is subcommand-based (`print`/`admin`/`version`/`help`), following the same
  per-subcommand-`FlagSet` + `usage()`/`usageFor(sub)` pattern as `config-watch` — see
  that project's `cmd/config-watch/main.go` if extending this further. Unlike
  `config-watch`, there's no `install`/`uninstall` (no daemon to install).

---

## Repository Layout

```text
cmd/
  zplkit/
    main.go            ← dispatch on os.Args[1] only; var version; tagRE
    print.go            ← printCmd: render a template for TAG and send/dry-run/list it
    admin.go             ← adminCmd: read a raw .zpl file and send/dry-run it, no rendering
    resolve.go            ← envOr/envOrInt + resolveBackend/resolveEndpoint/resolveDocumentFormat,
                             shared by printCmd and adminCmd
    usage.go               ← usage() (top-level) and usageFor(sub) (per-subcommand help text)
    main_test.go             ← tests for the resolve* precedence functions
internal/
  config/
    config.go          ← Config struct (company_name/asset_url/printer/aliases), toml tags
    load.go             ← Load(path): decode + DisallowUnknownFields + ZPL-injection check
    load_test.go
  label/
    label.go            ← Resolve (alias/path lookup), Render (text/template), List
    label_test.go
  printer/
    send.go              ← Backend type + Send(backend, endpoint, documentFormat, zpl, timeout) dispatcher
    socket.go             ← sendSocket: raw TCP via net.DialTimeout
    socket_test.go
    ipp.go                 ← sendIPP: IPP Print-Job via github.com/phin1x/go-ipp's low-level API
    ipp_test.go
templates/
  <size>/                ← one directory per physical label size, e.g. "2x1", "2.25x1.25"
    <name>.zpl            ← see README's "Template style" for the comment/formatting convention
    <name>.png             ← Labelary-rendered preview, same name as the .zpl file
```

`internal/config`, `internal/label`, and `internal/printer` do not import each other.
`cmd/zplkit`'s `print.go`/`admin.go` are the only places that wire `config.Load`'s output
into `label.Resolve`/`label.Render` (print.go only — admin.go never touches `label`) and
into `printer.Send`.

Templates are **not** `go:embed`-ed into the binary — they're read from disk at runtime
(`--templates`/`ZPLKIT_TEMPLATES`, default `./templates`) specifically so a new template
can be added without a rebuild or release. Don't change this to embedding.

---

## Build Commands

```bash
# Development build → build/zplkit
make build

# Run tests / static analysis
make test
make vet

# Stripped, trimpath, multi-arch release archives → dist/
make release

# Remove build/ and dist/
make clean
```

**Never build into the repo root.** Always target `build/` explicitly or use `make`.
Compiled binaries belong only in `build/` or `dist/`, both of which are gitignored.

---

## Git Workflow

**After every meaningful change, make a commit.** This applies to:

- Bug fixes (including single-line corrections)
- New features or behaviour changes
- Refactors or code reorganisation
- Documentation / AGENTS.md updates
- Build / Makefile changes
- New or edited templates

**Before committing, run the test suite and confirm it passes:**

```bash
make test && make vet
```

Do not commit if any test fails. If a change intentionally removes functionality,
update or delete the affected tests in the same commit. If you add new logic, add
corresponding tests before committing.

Do **not** push automatically; only commit locally unless the user explicitly asks
to push.

Use conventional-commits style:

- `feat:` — new user-visible feature
- `fix:` — bug fix
- `refactor:` — internal restructure, no behaviour change
- `docs:` — documentation only
- `build:` — Makefile, CI, or toolchain changes
- `chore:` — everything else (dependency updates, file renames, etc.)
- `template:` — adding or editing a `.zpl` template/preview, no Go change

---

## Key Conventions

### Go

- **Always `CGO_ENABLED=0`.**
- Version is injected at link time via `-ldflags "-X main.version=$(VERSION)"` (see
  `Makefile`) into the package-level `var version` in `main.go` — mirrors `xp-calc`'s
  pattern, not `config-watch`'s fuller `internal/version` package, since there's no
  install step whose metadata needs Commit/BuildDate/BuiltBy.
- CLI flags use `github.com/spf13/pflag`, one `pflag.NewFlagSet` per subcommand
  (`print.go`/`admin.go`) — there's no top-level flag set, since `print`/`admin` take
  different arguments. `usage()`/`usageFor(sub)` in `usage.go` hold the actual help
  text; each subcommand's `FlagSet.Usage` is wired to `usageFor` so `-h`/`--help`
  prints the right subcommand help. When adding a flag, update both the `FlagSet` and
  its `usageFor` case — they aren't generated from each other.
- Flag defaults that have an env var equivalent but no config-file equivalent
  (`ZPLKIT_TEMPLATES`, `ZPLKIT_CONFIG`, `ZPLKIT_TIMEOUT`) are computed via `envOr`/
  `envOrInt` in `resolve.go` *before* being passed as the flag's default, so precedence
  is always CLI flag > env var > built-in default. Don't add a second env-reading pass
  after `fs.Parse()` for these.
- The printer settings (`--backend`/`--endpoint`/`--document-format`, flags on *both*
  `print` and `admin`) have a fourth source below the env var — `zplkit.toml`'s
  `[printer]` table — so they can't use the `envOr` default trick (the config has to be
  loaded first, and loading it needs `--config`, which is itself one of the `envOr`-style
  flags). Their flags default to `""` and are resolved *after* `config.Load` via
  `resolveBackend`/`resolveEndpoint`/`resolveDocumentFormat` in `resolve.go`, which walk
  flag → env → config → built-in default explicitly. If you add another setting that
  should be config-overridable, follow this pattern, not `envOr`'s.
- `internal/printer` is a `Backend`-dispatched package, not just a socket sender:
  `Send(backend, endpoint, documentFormat, zpl, timeout)` routes to `sendSocket`
  (`socket.go`) or `sendIPP` (`ipp.go`). `sendIPP` deliberately bypasses `go-ipp`'s
  high-level `PrintJob`/`PrintDocuments` — those hardcode a CUPS-style
  `/printers/<name>` resource path via `adapter.GetHttpUri("printers", printer)`,
  which doesn't match LPrint's PAPPL-style routes (e.g. `/ipp/print/<name>`) — and
  instead builds a `*ipp.Request` and calls `client.SendRequestContext(ctx, endpoint,
  req, nil)` directly against the caller-supplied endpoint. Don't "simplify" this back
  to `PrintJob`; it would silently stop working against LPrint.
- Config is TOML via `github.com/pelletier/go-toml/v2`, decoded with
  `DisallowUnknownFields()` — a typo'd key in `zplkit.toml` is a load error, not a
  silently-ignored field.
- `config.Load` rejects `^`/`~` in `company_name`/`asset_url` — those are ZPL command
  prefixes, and these values get substituted into ZPL the same way an operator-supplied
  template field does. `label.Resolve` separately validates the `-t` name/alias itself
  (character allowlist, no `..`) so a crafted `-t` can't escape `TEMPLATE_DIR`. If you
  add a new config field that ends up inside rendered ZPL, give it the same check.
- Template rendering uses `text/template` against a `label.Data{Tag, CompanyName,
  AssetURL}` struct, not manual string replacement — an unknown field reference
  (typo'd placeholder) is then a hard template-execution error for free, rather than
  something `label.Render` has to detect itself.

### Testing

- `internal/config`'s tests are table-driven over raw TOML strings written to a temp
  file (see `writeTemp` in `load_test.go`) — valid case, unknown-field rejection, and
  each ZPL-breaker-character case.
- `internal/label`'s tests exercise `Resolve`/`Render`/`List` against real temp-dir
  template files, not mocks — covers alias lookup, explicit-path lookup, flat-name
  fallback, path-traversal rejection, and a missing-placeholder render failure.
- `internal/printer`'s socket tests (`socket_test.go`) run a real `net.Listen`
  loopback listener and assert on the bytes it received, plus a connection-refused
  case. Its IPP tests (`ipp_test.go`) run a real `httptest.Server` that decodes the
  request with `ipp.NewRequestDecoder` and asserts on the `document-format`/
  `printer-uri` attributes and the document body, then replies with an encoded
  `ipp.Response` — not a mock of `go-ipp`'s types, an actual wire-format round trip.
  Neither package mocks the network/HTTP layer, since the whole point of `printer` is
  getting that layer right.

### Security

`company_name`/`asset_url` (from `zplkit.toml`) and the rendered `{{.Tag}}` all end up
verbatim inside ZPL sent to a printer, where `^`/`~` start new commands. `TAG` is
constrained by `tagRE` in `print.go` before it ever reaches `label.Render`; config values
are constrained by `config.Load`'s validation. Apply extra scrutiny to any change that:

- Adds a new value that gets substituted into a template (a new `label.Data` field,
  a new config field) without threading it through an equivalent character check.
- `admin.go` is the one deliberate exception to all of this: it reads an arbitrary local
  file and sends its bytes to the printer completely unvalidated — no `tagRE`, no
  `label.Resolve`'s path-traversal check, no `config.Load`'s `^`/`~` rejection. That's
  intentional (it's an explicit operator action on a file they chose, not something
  driven by a scanned tag or other external input), not an oversight — don't add
  `print`-style validation to it, and don't let `print.go` grow an "send raw file" mode
  that skips validation the way `admin.go` does.
- Changes `label.Resolve`'s name/path validation — this is what stands between a
  `-t` argument and reading an arbitrary file via `TEMPLATE_DIR/../../...`.
