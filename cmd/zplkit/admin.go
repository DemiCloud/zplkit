package main

import (
	"fmt"
	"os"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/demicloud/zplkit/internal/config"
	"github.com/demicloud/zplkit/internal/printer"
)

// adminCmd sends printer-side provisioning/configuration that isn't tied
// to a specific asset tag. Its one positional argument is either a known
// action name (see adminActions in admin_actions.go — clock, darkness,
// calibrate, speed), whose ZPL/SGD is built here, or — for anything not
// in that map, e.g. storing a logo (~DG) — a raw .zpl file path, read
// and sent exactly as-is: no template resolution, no TAG, no
// placeholder rendering.
func adminCmd(args []string) {
	fs := flag.NewFlagSet("admin", flag.ContinueOnError)
	fs.Usage = func() { usageFor("admin") }

	var (
		dryRun             bool
		backendFlag        string
		endpointFlag       string
		documentFormatFlag string
		configPath         string
		timeoutSec         int
	)
	fs.BoolVarP(&dryRun, "dry-run", "d", false, "write the ZPL/file contents to stdout instead of sending it")
	fs.StringVar(&backendFlag, "backend", "", `delivery backend: "socket" or "ipp" (env ZPLKIT_BACKEND; config: printer.backend; default "socket")`)
	fs.StringVar(&endpointFlag, "endpoint", "", `printer endpoint: "host:port" for socket, a printer resource URL for ipp (env ZPLKIT_ENDPOINT; config: printer.endpoint)`)
	fs.StringVar(&documentFormatFlag, "document-format", "", `IPP document-format attribute, ipp backend only (env ZPLKIT_DOCUMENT_FORMAT; config: printer.document_format)`)
	fs.StringVar(&configPath, "config", envOr("ZPLKIT_CONFIG", "zplkit.toml"), "config file path (env ZPLKIT_CONFIG)")
	fs.IntVar(&timeoutSec, "timeout", envOrInt("ZPLKIT_TIMEOUT", 5), "connect/send timeout in seconds (env ZPLKIT_TIMEOUT)")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "admin: an action or a .zpl file path is required\n\n")
		usageFor("admin")
		os.Exit(1)
	}

	var raw string
	if build, ok := adminActions[fs.Arg(0)]; ok {
		zpl, err := build(fs.Args()[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "zplkit:", err)
			os.Exit(1)
		}
		raw = zpl
	} else {
		if fs.NArg() != 1 {
			fmt.Fprintf(os.Stderr, "admin: unknown action %q, and a file path takes no extra arguments\n\n", fs.Arg(0))
			usageFor("admin")
			os.Exit(1)
		}
		b, err := os.ReadFile(fs.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "zplkit:", err)
			os.Exit(1)
		}
		raw = string(b)
	}

	if dryRun {
		fmt.Println(raw)
		return
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "zplkit:", err)
		os.Exit(1)
	}

	backend := resolveBackend(backendFlag, cfg)
	endpoint := resolveEndpoint(endpointFlag, cfg)
	documentFormat := resolveDocumentFormat(documentFormatFlag, cfg)
	if err := printer.Send(backend, endpoint, documentFormat, raw, time.Duration(timeoutSec)*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "zplkit:", err)
		os.Exit(1)
	}
}
