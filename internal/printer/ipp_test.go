package printer

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ipp "github.com/phin1x/go-ipp"
)

// newFakeIPPServer records the decoded request and the document body sent
// to it, and responds with a minimal successful IPP response.
func newFakeIPPServer(t *testing.T, gotReq **ipp.Request, gotBody *bytes.Buffer) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := ipp.NewRequestDecoder(r.Body).Decode(gotBody)
		if err != nil {
			t.Fatalf("server: decoding request: %v", err)
		}
		*gotReq = req

		resp := ipp.NewResponse(ipp.StatusOk, req.RequestId)
		payload, err := resp.Encode()
		if err != nil {
			t.Fatalf("server: encoding response: %v", err)
		}
		w.Header().Set("Content-Type", ipp.ContentTypeIPP)
		w.Write(payload)
	}))
}

func TestSendIPPSubmitsDocumentFormatAndBody(t *testing.T) {
	var gotReq *ipp.Request
	var gotBody bytes.Buffer

	srv := newFakeIPPServer(t, &gotReq, &gotBody)
	defer srv.Close()

	zpl := "^XA^FDtest^FS^XZ"
	endpoint := srv.URL + "/ipp/print/zebra1"
	if err := sendIPP(endpoint, "application/octet-stream", zpl, 2*time.Second); err != nil {
		t.Fatalf("sendIPP: %v", err)
	}

	if gotReq == nil {
		t.Fatal("server never received a request")
	}
	if got := gotReq.OperationAttributes[ipp.AttributeDocumentFormat]; got != "application/octet-stream" {
		t.Errorf("document-format = %v, want %q", got, "application/octet-stream")
	}
	wantURI := "ipp" + endpoint[len("http"):] // http://... -> ipp://...
	if got := gotReq.OperationAttributes[ipp.AttributePrinterURI]; got != wantURI {
		t.Errorf("printer-uri = %v, want %q", got, wantURI)
	}
	if gotBody.String() != zpl {
		t.Errorf("document body = %q, want %q", gotBody.String(), zpl)
	}
}

func TestSendIPPServerErrorIsReturned(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := ipp.NewRequestDecoder(r.Body).Decode(nil)
		if err != nil {
			t.Fatalf("server: decoding request: %v", err)
		}
		resp := ipp.NewResponse(ipp.StatusErrorDocumentFormatNotSupported, req.RequestId)
		payload, err := resp.Encode()
		if err != nil {
			t.Fatalf("server: encoding response: %v", err)
		}
		w.Header().Set("Content-Type", ipp.ContentTypeIPP)
		w.Write(payload)
	}))
	defer srv.Close()

	err := sendIPP(srv.URL+"/ipp/print/zebra1", "application/octet-stream", "^XA^XZ", 2*time.Second)
	if err == nil {
		t.Fatal("expected error for a rejected document-format, got nil")
	}
}

func TestSendIPPInvalidEndpoint(t *testing.T) {
	if err := sendIPP("://not a url", "application/octet-stream", "^XA^XZ", time.Second); err == nil {
		t.Fatal("expected error for an invalid endpoint, got nil")
	}
}
