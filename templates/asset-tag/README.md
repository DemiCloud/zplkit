# asset-tag

Asset tags: a company name, a Code 128 barcode, and a human-readable tag —
some variants add a QR code or a border. See the root README's Templates
section for the full table of paths/previews.

## Fields

- `CompanyName` — automatic, from `zplkit.toml`'s `company_name`.
- `AssetURL` — automatic, from `zplkit.toml`'s `asset_url` (`2.25x1.25`'s
  QR variants only).
- `Tag` — from `print`'s positional `TAG` argument or `-f Tag=...`; must
  match `^[A-Z]{2,4}-[0-9]{6}$` (enforced in `cmd/zplkit/print.go`, not
  here).

`2x1/003`/`004`/`006`'s "If found: ..." line is fixed text baked into the
template, not a field — there's no `Contact` placeholder for it.

## Border variants

Every layout ships three ways: no border, square border, rounded border.
All three place fields identically — only the `^GB` border graphic
differs — so switching between them never reflows the label.

## `CompanyName` sizing

`CompanyName` is free text with no length bound, so every variant sizes
it with `fitfont` (see the root README's Placeholders section) instead
of a fixed font size — a long configured company name shrinks instead of
running past the label edge, or on `2.25x1.25`, into the QR code.

## `2.25x1.25`'s QR/barcode gotcha

`^BY2,2,72` (set for the Code 128 field) leaves `72` as persistent
printer state. The `^BQ` (QR) field that follows picks that up as an
extra Y-offset, even though `^BC` already specified its own explicit
height and QR codes have no height parameter of their own to need one —
confirmed on both Labelary and a real printer, so it's firmware
behavior, not a preview artifact. Each `2.25x1.25` template resets it
first: `^BY2,2,0` right before the QR's `^FO`/`^BQ`. If you add a QR
field after a barcode field in a new template, you need the same reset.
