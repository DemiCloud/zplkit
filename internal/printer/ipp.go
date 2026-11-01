package printer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	ipp "github.com/phin1x/go-ipp"
)

// sendIPP submits zpl as an IPP Print-Job to endpoint (a printer
// resource's full http(s) URL, e.g.
// "http://lprint-host:8000/ipp/print/zebra1"), with documentFormat as the
// IPP document-format attribute (e.g. "application/octet-stream").
//
// go-ipp's high-level PrintJob/PrintDocuments helpers hardcode a
// CUPS-style "/printers/<name>" resource path, which doesn't match
// LPrint's PAPPL-style routes (e.g. "/ipp/print/<name>") — so this calls
// the lower-level Request/SendRequestContext API directly against the
// caller-supplied endpoint instead.
func sendIPP(endpoint, documentFormat, zpl string, timeout time.Duration) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid IPP endpoint %q: %w", endpoint, err)
	}

	useTLS := u.Scheme == "https"
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		port = 631
		if useTLS {
			port = 443
		}
	}

	adapter := ipp.NewHttpAdapter(u.Hostname(), port, "", "", useTLS, ipp.WithHttpClient(&http.Client{Timeout: timeout}))
	client := ipp.NewIPPClientWithAdapter("zplkit", adapter)

	req := ipp.NewRequest(ipp.OperationPrintJob, 1)
	req.OperationAttributes[ipp.AttributePrinterURI] = printerURIAttribute(u)
	req.OperationAttributes[ipp.AttributeRequestingUserName] = "zplkit"
	req.OperationAttributes[ipp.AttributeJobName] = "zplkit"
	req.OperationAttributes[ipp.AttributeDocumentFormat] = documentFormat
	req.File = strings.NewReader(zpl)
	req.FileSize = len(zpl)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := client.SendRequestContext(ctx, endpoint, req, nil); err != nil {
		return fmt.Errorf("IPP print job to %s: %w", endpoint, err)
	}
	return nil
}

// printerURIAttribute returns endpoint's URL with its scheme swapped to
// the ipp/ipps convention IPP servers expect in the printer-uri
// attribute, even though the actual HTTP request (made with the original
// http(s) endpoint) is what's really dialed.
func printerURIAttribute(u *url.URL) string {
	uri := *u
	switch u.Scheme {
	case "http":
		uri.Scheme = "ipp"
	case "https":
		uri.Scheme = "ipps"
	}
	return uri.String()
}
