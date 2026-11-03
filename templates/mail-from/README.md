# mail-from

A return-address label: your company name plus the same name/address
fields as `mail-to`. See the root README's Templates section for the
full table of paths/previews.

## Fields

- `CompanyName` — automatic, from `zplkit.toml`'s `company_name`.
- `Name`, `Address1`, `City`, `State`, `ZIP` — via `-f`, same as
  `mail-to` (this is your own name/address, not the recipient's).

Same `fitfont` sizing and border-variant field-parity notes as `mail-to`
— see `templates/mail-to/README.md` / `templates/asset-tag/README.md`.
