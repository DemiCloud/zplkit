# zplkit — Agent / Contributor Instructions

## Collaboration Philosophy

Development is a collaborative process. Agents are expected to exercise independent
technical judgment, not passive agreement. When a different solution better aligns with
industry standards or formal specifications (e.g. ZPL semantics, Go idioms), agents
should raise concerns and advocate for the more appropriate approach.

## Project Overview

`zplkit` renders a ZPL label template and sends it to a network label printer (raw TCP
or IPP — see `internal/printer`), and separately can send a raw `.zpl` file straight
through for printer-side provisioning (`admin`). Templates aren't asset-tag-specific —
`templates/<category>/<size>/<NNN>.zpl` groups them by what they're for (`asset-tag`,
`mail-to`, `mail-from`, `name` today; nothing stops adding more) and by physical label
size within that.

- Written in **Go 1.25**. Always **`CGO_ENABLED=0`**. Cross-platform (linux/darwin/
  windows) — this is a one-shot CLI, not a daemon, so there's no systemd/install story
  like `config-watch`/`anycast-sentinel` have.
- Produces a **single static binary** from `cmd/zplkit/`.
- **The repo's real contribution surface is `templates/`, not Go code.** Adding a label
  layout — even for a brand-new category — is "drop a `.zpl` file + a preview PNG,
  maybe add a TOML alias" — see the README's "Adding a template" section. Don't make
  template changes require touching Go; this is also why `label.Data` is a generic
  `map[string]string` rather than a fixed struct (see Key Conventions below) — a
  per-category struct would mean a Go change every time a template needed a new field.
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
    print.go            ← printCmd: render a template and send/dry-run/list it
    admin.go             ← adminCmd: read a raw .zpl file and send/dry-run it, no rendering
    resolve.go            ← envOr/envOrInt + resolveBackend/resolveEndpoint/resolveDocumentFormat,
                             shared by printCmd and adminCmd
    usage.go               ← usage() (top-level) and usageFor(sub) (per-subcommand help text)
    main_test.go             ← tests for the resolve* precedence functions
internal/
  config/
    config.go          ← Config struct (company_name/asset_url/printer/aliases/categories), toml tags
    load.go             ← Load(path): decode + DisallowUnknownFields + ZPL-injection check
    load_test.go
  label/
    label.go            ← Resolve (alias/path lookup), Render (text/template + field validation), List
    label_test.go
  printer/
    send.go              ← Backend type + Send(backend, endpoint, documentFormat, zpl, timeout) dispatcher
    socket.go             ← sendSocket: raw TCP via net.DialTimeout
    socket_test.go
    ipp.go                 ← sendIPP: IPP Print-Job via github.com/phin1x/go-ipp's low-level API
    ipp_test.go
templates/
  <category>/            ← one directory per label category, e.g. "asset-tag" — what the label is for
    <size>/                ← one directory per physical label size, e.g. "2x1", "2.25x1.25"
      <NNN>.zpl              ← catalog number, not a descriptive name — see README's "Templates" section
      <NNN>.png               ← Labelary-rendered preview, same name as the .zpl file
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
- `config.Load` rejects `^`/`~` in `company_name`/`asset_url` at load time, for an early,
  specific error pointing at `zplkit.toml` — those are ZPL command prefixes. This is
  belt-and-suspenders with `label.Render`'s own check (see below), which is what's
  actually relied on now that field values can also come from `print -f`.
  `label.Resolve` separately validates the `-t` name/alias itself (character allowlist,
  no `..`) so a crafted `-t` can't escape `TEMPLATE_DIR`.
- `label.Data` is `map[string]string`, not a fixed struct — see Project Overview above
  for why. `label.Render` parses templates with `.Option("missingkey=error")`, so a
  placeholder that doesn't resolve to a key in `Data` (typo'd, or simply a field nobody
  supplied via `-f`) is still a hard template-execution error, the same guarantee the
  old struct gave for free. Don't drop that `Option` call when touching `Render`.
- `label.Render` also runs every `Data` value through `validateFieldValues` (same
  `^`/`~` ZPL-breaker check as `config.Load`'s, kept as a separate copy rather than a
  shared helper — see the "do not import each other" note above) before executing the
  template. This is now the one enforced choke point for "does this value end up
  verbatim in ZPL," regardless of whether it came from config or a CLI `-f` flag — see
  Security below.

### Testing

- `internal/config`'s tests are table-driven over raw TOML strings written to a temp
  file (see `writeTemp` in `load_test.go`) — valid case, unknown-field rejection, and
  each ZPL-breaker-character case.
- `internal/label`'s tests exercise `Resolve`/`Render`/`List` against real temp-dir
  template files, not mocks — covers alias lookup, explicit-path lookup, flat-name
  fallback, path-traversal rejection, an unknown-placeholder render failure, a field
  nobody supplied (`missingkey=error`), and each ZPL-breaker-character case in a
  field value.
- `internal/printer`'s socket tests (`socket_test.go`) run a real `net.Listen`
  loopback listener and assert on the bytes it received, plus a connection-refused
  case. Its IPP tests (`ipp_test.go`) run a real `httptest.Server` that decodes the
  request with `ipp.NewRequestDecoder` and asserts on the `document-format`/
  `printer-uri` attributes and the document body, then replies with an encoded
  `ipp.Response` — not a mock of `go-ipp`'s types, an actual wire-format round trip.
  Neither package mocks the network/HTTP layer, since the whole point of `printer` is
  getting that layer right.

### Security

Every value in a `label.Data` ends up verbatim inside ZPL sent to a printer, where
`^`/`~` start new commands — `CompanyName`/`AssetURL` (from `zplkit.toml`), `Tag` (from
`print`'s positional `TAG` argument or `-f Tag=...`), and any other field a template
references, set via `print -f Field=value`. `label.Render`'s `validateFieldValues`
enforces the `^`/`~` rejection on every key in `Data` in one place, so this is now
structural rather than something each value's source needs to remember — see Key
Conventions above. `Tag` specifically also gets the stronger `tagRE` shape check
(`^[A-Z]{2,4}-[0-9]{6}$`) in `print.go`, applied whenever `Data["Tag"]` is set, from
either source. Apply extra scrutiny to any change that:

- Adds a path where a value reaches `label.Render` without going through `Data` (and
  therefore without going through `validateFieldValues`) — e.g. any future templating
  that doesn't route through `label.Data`.
- Loosens or removes `label.Render`'s `.Option("missingkey=error")` or
  `validateFieldValues` call.
- `admin.go` is the one deliberate exception to all of this: it reads an arbitrary local
  file and sends its bytes to the printer completely unvalidated — no `tagRE`, no
  `label.Resolve`'s path-traversal check, no field validation. That's intentional (it's
  an explicit operator action on a file they chose, not something driven by a scanned
  tag or other external input), not an oversight — don't add `print`-style validation
  to it, and don't let `print.go` grow an "send raw file" mode that skips validation the
  way `admin.go` does.
- Changes `label.Resolve`'s name/path validation — this is what stands between a
  `-t` argument and reading an arbitrary file via `TEMPLATE_DIR/../../...`.
