package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Print(`zplkit — render a ZPL label template for an asset tag and send it to a network label printer

Usage:
  zplkit <subcommand> [flags]

Subcommands:
  print    Render a template for a tag and send it to the printer
  admin    Send a raw ZPL file to the printer (provisioning/configuration)
  version  Show version information
  help     Show help for a subcommand
`)
}

func usageFor(sub string) {
	switch sub {
	case "print":
		fmt.Print(`zplkit print — render a ZPL label template for TAG and send it to the printer

Usage:
  zplkit print [flags] TAG
  zplkit print -l

TAG must match ^[A-Z]{2,4}-[0-9]{6}$ (e.g. ABC-000001).

Flags:
  -t, --template string         Template name, alias, or path (default "default")
  -d, --dry-run                 Write rendered ZPL to stdout instead of printing
  -l, --list                    List available templates and aliases
      --backend string          Delivery backend: "socket" or "ipp" (env ZPLKIT_BACKEND; config: printer.backend; default "socket")
      --endpoint string         Printer endpoint (env ZPLKIT_ENDPOINT; config: printer.endpoint)
      --document-format string  IPP document-format attribute, ipp backend only (env ZPLKIT_DOCUMENT_FORMAT; config: printer.document_format)
      --templates string        Template directory (env ZPLKIT_TEMPLATES; default "templates")
      --config string           Config file path (env ZPLKIT_CONFIG; default "zplkit.toml")
      --timeout int             Connect/send timeout in seconds (env ZPLKIT_TIMEOUT; default 5)
  -h, --help                    Show this help
`)
	case "admin":
		fmt.Print(`zplkit admin — send a raw ZPL file straight to the printer (no template rendering)

Use this for printer-side provisioning/configuration: storing a logo
(~DG), setting the real-time clock (^SC), calibration, and similar
one-off setup that isn't tied to a specific asset tag.

Usage:
  zplkit admin [flags] FILE

Flags:
  -d, --dry-run                 Write the file's contents to stdout instead of sending it
      --backend string          Delivery backend: "socket" or "ipp" (env ZPLKIT_BACKEND; config: printer.backend; default "socket")
      --endpoint string         Printer endpoint (env ZPLKIT_ENDPOINT; config: printer.endpoint)
      --document-format string  IPP document-format attribute, ipp backend only (env ZPLKIT_DOCUMENT_FORMAT; config: printer.document_format)
      --config string           Config file path (env ZPLKIT_CONFIG; default "zplkit.toml")
      --timeout int             Connect/send timeout in seconds (env ZPLKIT_TIMEOUT; default 5)
  -h, --help                    Show this help
`)
	case "version":
		fmt.Print(`zplkit version — show version information

Usage:
  zplkit version
  zplkit --version
  zplkit -V
`)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", sub)
		usage()
		os.Exit(1)
	}
}
