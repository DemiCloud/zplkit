package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "zplkit.toml")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestLoadMissingFileReturnsZeroValue(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "does-not-exist.toml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CompanyName != "" || cfg.AssetURL != "" || len(cfg.Aliases) != 0 {
		t.Fatalf("expected zero-value config, got %+v", cfg)
	}
}

func TestLoadValid(t *testing.T) {
	path := writeTemp(t, `
company_name = "Example Corp"
asset_url = "https://example.com/assets/"

[aliases]
default = "2x1/001"
qr = "2.25x1.25/001"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.CompanyName != "Example Corp" {
		t.Errorf("CompanyName = %q, want %q", cfg.CompanyName, "Example Corp")
	}
	if cfg.AssetURL != "https://example.com/assets/" {
		t.Errorf("AssetURL = %q, want %q", cfg.AssetURL, "https://example.com/assets/")
	}
	if cfg.Aliases["default"] != "2x1/001" {
		t.Errorf("Aliases[default] = %q, want %q", cfg.Aliases["default"], "2x1/001")
	}
	if cfg.Aliases["qr"] != "2.25x1.25/001" {
		t.Errorf("Aliases[qr] = %q, want %q", cfg.Aliases["qr"], "2.25x1.25/001")
	}
}

func TestLoadUnknownFieldRejected(t *testing.T) {
	path := writeTemp(t, `company_nmae = "typo"`)
	if _, err := Load(path); err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestLoadPrinterConfig(t *testing.T) {
	path := writeTemp(t, `
[printer]
backend         = "ipp"
endpoint        = "http://lprint-host:8000/ipp/print/zebra1"
document_format = "application/octet-stream"
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := PrinterConfig{
		Backend:        "ipp",
		Endpoint:       "http://lprint-host:8000/ipp/print/zebra1",
		DocumentFormat: "application/octet-stream",
	}
	if cfg.Printer != want {
		t.Errorf("Printer = %+v, want %+v", cfg.Printer, want)
	}
}

func TestLoadInvalidBackendRejected(t *testing.T) {
	path := writeTemp(t, `
[printer]
backend = "carrier-pigeon"
`)
	if _, err := Load(path); err == nil {
		t.Fatal("expected error for invalid printer.backend, got nil")
	}
}

func TestLoadRejectsZPLCommandCharacters(t *testing.T) {
	tests := []struct {
		name string
		toml string
	}{
		{"caret in company_name", `company_name = "Evil^XA^FO0,0^XZ"`},
		{"tilde in company_name", `company_name = "Evil~SD"`},
		{"caret in asset_url", `asset_url = "https://example.com/^XA"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTemp(t, tt.toml)
			if _, err := Load(path); err == nil {
				t.Fatal("expected error for ZPL command character, got nil")
			}
		})
	}
}
