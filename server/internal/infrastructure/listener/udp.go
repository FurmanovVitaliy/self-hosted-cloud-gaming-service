package listener

import (
	"fmt"
	"net"
)

type UDPReader struct {
	Port       int
	bufferSize uint
	conn       *net.UDPConn
}

func NewUdpReader(port int, bufferSize uint, conn *net.UDPConn) *UDPReader {
	return &UDPReader{
		Port:       port,
		bufferSize: bufferSize,
		conn:       conn,
	}
}

func (r *UDPReader) Read() ([]byte, error) {
	buffer := make([]byte, r.bufferSize)
	n, _, err := r.conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading from UDP: %w", err)
	}
	if n == 0 {
		return nil, nil // No data received
	}
	// Return only the portion of the buffer that was actually read
	return buffer[:n], nil
}

func (r *UDPReader) GetPort() int {
	return r.Port
}

func (r *UDPReader) Close() error {
	return r.conn.Close()
}
