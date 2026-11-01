package printer

import (
	"fmt"
	"net"
	"time"
)

// sendSocket dials addr ("host:port") and writes zpl to the connection,
// bounded by timeout for both the connect and the write.
func sendSocket(addr, zpl string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", addr, err)
	}
	defer conn.Close()

	if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("setting write deadline: %w", err)
	}
	if _, err := conn.Write([]byte(zpl)); err != nil {
		return fmt.Errorf("sending to %s: %w", addr, err)
	}
	return nil
}
