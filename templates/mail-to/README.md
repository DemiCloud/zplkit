# mail-to

A recipient address label: name and street address, nothing automatic.
See the root README's Templates section for the full table of
paths/previews.

## Fields

- `Name`, `Address1`, `City`, `State`, `ZIP` — all via `-f`, e.g.
  `-f Name="Jane Doe" -f Address1="123 Main St" -f City=Springfield
  -f State=IL -f ZIP=62701`. None are automatic — this is the
  recipient, not you.

All four are free text with no length bound, sized with `fitfont` (see
the root README's Placeholders section) so a long value shrinks instead
of running past the label edge or a bordered variant's border. No
border/square/rounded variants place fields identically and share one
(bordered-case) `fitfont` max-width — see `templates/asset-tag/README.md`'s
"Border variants" note; same reasoning here.
