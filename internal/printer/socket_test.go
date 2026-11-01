package printer

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestSendSocketWritesToListener(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer ln.Close()

	received := make(chan []byte, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			received <- nil
			return
		}
		defer conn.Close()
		b, _ := io.ReadAll(conn)
		received <- b
	}()

	zpl := "^XA^FDtest^FS^XZ"
	if err := sendSocket(ln.Addr().String(), zpl, 2*time.Second); err != nil {
		t.Fatalf("sendSocket: %v", err)
	}

	select {
	case got := <-received:
		if string(got) != zpl {
			t.Errorf("listener received %q, want %q", got, zpl)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for listener to receive data")
	}
}

func TestSendSocketConnectionRefused(t *testing.T) {
	// Bind and immediately close to get an address nothing is listening on.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	if err := sendSocket(addr, "^XA^XZ", 2*time.Second); err == nil {
		t.Fatal("expected error connecting to a closed port, got nil")
	}
}

func TestSendDispatchesToSocketBackend(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer ln.Close()
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	if err := Send(BackendSocket, ln.Addr().String(), "", "^XA^XZ", 2*time.Second); err != nil {
		t.Fatalf("Send: %v", err)
	}
}

func TestSendRejectsUnknownBackend(t *testing.T) {
	if err := Send("carrier-pigeon", "whatever", "", "^XA^XZ", time.Second); err == nil {
		t.Fatal("expected error for unknown backend, got nil")
	}
}
