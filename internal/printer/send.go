// Package printer delivers rendered ZPL to a label printer, either by
// dialing it directly over raw TCP or by submitting an IPP print job (e.g.
// to an LPrint instance fronting a printer that no longer accepts raw
// sockets — see CUPS's raw-queue deprecation).
package printer

import (
	"fmt"
	"time"
)

// Backend selects how Send delivers zpl.
type Backend string

const (
	// BackendSocket dials endpoint ("host:port") directly and writes zpl
	// as-is — the traditional Zebra/CUPS-raw-queue wire format.
	BackendSocket Backend = "socket"
	// BackendIPP submits zpl as an IPP Print-Job to endpoint, a printer
	// resource's full http(s) URL.
	BackendIPP Backend = "ipp"
)

// Send delivers zpl using backend. endpoint and documentFormat are
// interpreted as described on BackendSocket/BackendIPP; documentFormat is
// ignored by BackendSocket.
func Send(backend Backend, endpoint, documentFormat, zpl string, timeout time.Duration) error {
	switch backend {
	case BackendSocket:
		return sendSocket(endpoint, zpl, timeout)
	case BackendIPP:
		return sendIPP(endpoint, documentFormat, zpl, timeout)
	default:
		return fmt.Errorf("unknown printer backend: %q (want %q or %q)", backend, BackendSocket, BackendIPP)
	}
}
