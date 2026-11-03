# name

A single centered `Name` field — a box or jar label, nothing else. See
the root README's Templates section for the full table of
paths/previews.

## Fields

- `Name` — via `-f Name=...`. Free text, no length bound.

## Sizing and line breaks

`Name` is rendered with `fitfont` (see the root README's Placeholders
section), shrinking from 48pt down to a 20pt floor to stay on one line.
Text too long even at 20pt still renders at 20pt, running past the
label edge rather than wrapping into overlapping text.

A real newline in the field's value renders as a second centered line,
sized together with the first so the longer of the two still fits:

    -f Name=$'Cables:\nUSB-C'

Use bash's `$'...'` quoting, not `"..."` — `"Cables:\nUSB-C"` leaves
the literal two characters `\` and `n` in the string instead of a line
break. This works via two more template functions, `zplbreak` and
`linecount` (also documented in the root README's Placeholders
section): `zplbreak` converts the newline into ZPL's own `\&`
explicit-break sequence (a raw newline byte has no defined meaning
inside `^FD...^FS` and breaks Labelary's parser), and `linecount`
picks the right `^FO` y-coordinate to re-center a 2-line render.

No border/square/rounded variants place `Name` identically and share
one (bordered-case) `fitfont` max-width — see
`templates/asset-tag/README.md`'s "Border variants" note; same
reasoning here.
