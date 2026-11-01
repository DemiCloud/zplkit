package main

import (
	"testing"

	"github.com/demicloud/zplkit/internal/config"
	"github.com/demicloud/zplkit/internal/printer"
)

func TestResolveBackendPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		flagVal string
		env     string
		cfg     string
		want    printer.Backend
	}{
		{"flag wins over everything", "ipp", "socket", "ipp", "ipp"},
		{"env wins over config", "", "ipp", "socket", "ipp"},
		{"config wins over default", "", "", "ipp", "ipp"},
		{"built-in default when nothing set", "", "", "", printer.BackendSocket},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				t.Setenv("ZPLKIT_BACKEND", tt.env)
			}
			cfg := &config.Config{Printer: config.PrinterConfig{Backend: tt.cfg}}
			if got := resolveBackend(tt.flagVal, cfg); got != tt.want {
				t.Errorf("resolveBackend(%q, %+v) = %q, want %q", tt.flagVal, cfg, got, tt.want)
			}
		})
	}
}

func TestResolveEndpointPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		flagVal string
		env     string
		cfg     string
		want    string
	}{
		{"flag wins over everything", "flag:1", "env:2", "cfg:3", "flag:1"},
		{"env wins over config", "", "env:2", "cfg:3", "env:2"},
		{"config wins over default", "", "", "cfg:3", "cfg:3"},
		{"built-in default when nothing set", "", "", "", "printer.example.net:9100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				t.Setenv("ZPLKIT_ENDPOINT", tt.env)
			}
			cfg := &config.Config{Printer: config.PrinterConfig{Endpoint: tt.cfg}}
			if got := resolveEndpoint(tt.flagVal, cfg); got != tt.want {
				t.Errorf("resolveEndpoint(%q, %+v) = %q, want %q", tt.flagVal, cfg, got, tt.want)
			}
		})
	}
}

func TestResolveDocumentFormatPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		flagVal string
		env     string
		cfg     string
		want    string
	}{
		{"flag wins over everything", "flag/fmt", "env/fmt", "cfg/fmt", "flag/fmt"},
		{"env wins over config", "", "env/fmt", "cfg/fmt", "env/fmt"},
		{"config wins over default", "", "", "cfg/fmt", "cfg/fmt"},
		{"built-in default when nothing set", "", "", "", "application/octet-stream"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				t.Setenv("ZPLKIT_DOCUMENT_FORMAT", tt.env)
			}
			cfg := &config.Config{Printer: config.PrinterConfig{DocumentFormat: tt.cfg}}
			if got := resolveDocumentFormat(tt.flagVal, cfg); got != tt.want {
				t.Errorf("resolveDocumentFormat(%q, %+v) = %q, want %q", tt.flagVal, cfg, got, tt.want)
			}
		})
	}
}
