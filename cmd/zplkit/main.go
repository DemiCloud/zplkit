// Command zplkit renders a ZPL label template for an asset tag and sends
// it to a network label printer, or sends a raw ZPL file straight through
// for printer-side provisioning (see "zplkit help admin").
package main

import (
	"fmt"
	"os"
	"regexp"
)

var version = "dev"

var tagRE = regexp.MustCompile(`^[A-Z]{2,4}-[0-9]{6}$`)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "print":
		printCmd(os.Args[2:])
	case "admin":
		adminCmd(os.Args[2:])
	case "version", "--version", "-V":
		fmt.Println("zplkit", version)
	case "help", "--help", "-h":
		if len(os.Args) > 2 {
			usageFor(os.Args[2])
		} else {
			usage()
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
