package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/demicloud/zplkit/internal/config"
	"github.com/demicloud/zplkit/internal/label"
	"github.com/demicloud/zplkit/internal/printer"
)

func printCmd(args []string) {
	fs := flag.NewFlagSet("print", flag.ContinueOnError)
	fs.Usage = func() { usageFor("print") }

	var (
		template           string
		dryRun             bool
		list               bool
		backendFlag        string
		endpointFlag       string
		documentFormatFlag string
		templateDir        string
		configPath         string
		timeoutSec         int
	)
	fs.StringVarP(&template, "template", "t", "default", "template name, alias, or path")
	fs.BoolVarP(&dryRun, "dry-run", "d", false, "write rendered ZPL to stdout instead of printing")
	fs.BoolVarP(&list, "list", "l", false, "list available templates and aliases")
	fs.StringVar(&backendFlag, "backend", "", `delivery backend: "socket" or "ipp" (env ZPLKIT_BACKEND; config: printer.backend; default "socket")`)
	fs.StringVar(&endpointFlag, "endpoint", "", `printer endpoint: "host:port" for socket, a printer resource URL for ipp (env ZPLKIT_ENDPOINT; config: printer.endpoint)`)
	fs.StringVar(&documentFormatFlag, "document-format", "", `IPP document-format attribute, ipp backend only (env ZPLKIT_DOCUMENT_FORMAT; config: printer.document_format)`)
	fs.StringVar(&templateDir, "templates", envOr("ZPLKIT_TEMPLATES", "templates"), "template directory (env ZPLKIT_TEMPLATES)")
	fs.StringVar(&configPath, "config", envOr("ZPLKIT_CONFIG", "zplkit.toml"), "config file path (env ZPLKIT_CONFIG)")
	fs.IntVar(&timeoutSec, "timeout", envOrInt("ZPLKIT_TIMEOUT", 5), "connect/send timeout in seconds (env ZPLKIT_TIMEOUT)")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "zplkit:", err)
		os.Exit(1)
	}

	if list {
		if err := listTemplates(templateDir, cfg); err != nil {
			fmt.Fprintln(os.Stderr, "zplkit:", err)
			os.Exit(1)
		}
		return
	}

	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "print: a TAG argument is required\n\n")
		usageFor("print")
		os.Exit(1)
	}
	tag := fs.Arg(0)
	if !tagRE.MatchString(tag) {
		fmt.Fprintln(os.Stderr, "zplkit: invalid tag:", tag)
		os.Exit(1)
	}

	path, err := label.Resolve(templateDir, template, cfg.Aliases)
	if err != nil {
		fmt.Fprintln(os.Stderr, "zplkit:", err)
		os.Exit(1)
	}

	zpl, err := label.Render(path, label.Data{
		Tag:         tag,
		CompanyName: cfg.CompanyName,
		AssetURL:    cfg.AssetURL,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "zplkit:", err)
		os.Exit(1)
	}

	if dryRun {
		fmt.Println(zpl)
		return
	}

	backend := resolveBackend(backendFlag, cfg)
	endpoint := resolveEndpoint(endpointFlag, cfg)
	documentFormat := resolveDocumentFormat(documentFormatFlag, cfg)
	if err := printer.Send(backend, endpoint, documentFormat, zpl, time.Duration(timeoutSec)*time.Second); err != nil {
		fmt.Fprintln(os.Stderr, "zplkit:", err)
		os.Exit(1)
	}
}

func listTemplates(templateDir string, cfg *config.Config) error {
	names, err := label.List(templateDir)
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Println(name)
	}

	if len(cfg.Aliases) == 0 {
		return nil
	}

	fmt.Println("\nAliases:")
	aliasNames := make([]string, 0, len(cfg.Aliases))
	for a := range cfg.Aliases {
		aliasNames = append(aliasNames, a)
	}
	sort.Strings(aliasNames)
	for _, a := range aliasNames {
		fmt.Printf("  %s -> %s\n", a, cfg.Aliases[a])
	}
	return nil
}
