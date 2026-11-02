// Package config loads zplkit's TOML config file: the company name and
// asset URL placeholders templates can reference, the default printer
// endpoint, and the alias table that resolves a short -t name to a
// template path.
package config

// Config holds the values templates can reference as {{.CompanyName}} and
// {{.AssetURL}}, the default printer settings, the alias -> template path
// table, and the declared template categories.
type Config struct {
	CompanyName string            `toml:"company_name"`
	AssetURL    string            `toml:"asset_url"`
	Printer     PrinterConfig     `toml:"printer"`
	Aliases     map[string]string `toml:"aliases"`
	// Categories is a purely informational list of the template category
	// directories in use (the first path segment under templates/, e.g.
	// "asset-tag"), shown by `print -l` for discoverability. It is never
	// enforced — templates/ is checked into the repo while zplkit.toml is
	// local and gitignored, so a category existing on disk without being
	// listed here (e.g. right after pulling an update) must still work.
	Categories []string `toml:"categories"`
}

// PrinterConfig is the lowest-priority source for how zpl gets delivered —
// a CLI flag or environment variable overrides any of these; see
// resolveBackend/resolveEndpoint/resolveDocumentFormat in cmd/zplkit/main.go.
type PrinterConfig struct {
	// Backend selects the delivery method: "socket" (raw TCP, the
	// traditional Zebra/CUPS-raw-queue wire format) or "ipp" (an IPP
	// Everywhere endpoint, e.g. one exposed by LPrint in front of a
	// printer that no longer accepts raw sockets directly).
	Backend string `toml:"backend"`
	// Endpoint's meaning depends on Backend: for "socket" it's "host:port"
	// dialed directly; for "ipp" it's the printer resource's full http(s)
	// URL (e.g. "http://lprint-host:8000/ipp/print/zebra1").
	Endpoint string `toml:"endpoint"`
	// DocumentFormat is the IPP document-format attribute sent with the
	// print job. Only meaningful for the "ipp" backend.
	DocumentFormat string `toml:"document_format"`
}
