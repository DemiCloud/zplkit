package main

import (
	"os"
	"strconv"

	"github.com/demicloud/zplkit/internal/config"
	"github.com/demicloud/zplkit/internal/printer"
)

// envOr returns the environment variable key's value, or def if unset.
// Flag defaults for ZPLKIT_TEMPLATES/ZPLKIT_CONFIG/ZPLKIT_TIMEOUT are
// computed this way so their precedence is CLI flag > env var > built-in
// default. The printer settings have one more layer below the env var —
// the config file — so they're resolved separately via resolveBackend/
// resolveEndpoint/resolveDocumentFormat below, shared by printCmd and
// adminCmd.
func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func envOrInt(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// resolveBackend returns the first non-empty value, in priority order: an
// explicit --backend flag, the ZPLKIT_BACKEND environment variable, the
// config file's printer.backend, then the built-in default.
func resolveBackend(flagVal string, cfg *config.Config) printer.Backend {
	for _, v := range []string{flagVal, os.Getenv("ZPLKIT_BACKEND"), cfg.Printer.Backend} {
		if v != "" {
			return printer.Backend(v)
		}
	}
	return printer.BackendSocket
}

// resolveEndpoint is resolveBackend's counterpart for
// --endpoint/ZPLKIT_ENDPOINT/config. The built-in default only makes sense
// for the socket backend; an "ipp" backend with no endpoint configured
// fails when printer.Send rejects the empty host:port.
func resolveEndpoint(flagVal string, cfg *config.Config) string {
	for _, v := range []string{flagVal, os.Getenv("ZPLKIT_ENDPOINT"), cfg.Printer.Endpoint} {
		if v != "" {
			return v
		}
	}
	return "printer.example.net:9100"
}

// resolveDocumentFormat is resolveBackend's counterpart for
// --document-format/ZPLKIT_DOCUMENT_FORMAT/config. Only meaningful for the
// "ipp" backend.
func resolveDocumentFormat(flagVal string, cfg *config.Config) string {
	for _, v := range []string{flagVal, os.Getenv("ZPLKIT_DOCUMENT_FORMAT"), cfg.Printer.DocumentFormat} {
		if v != "" {
			return v
		}
	}
	return "application/octet-stream"
}
