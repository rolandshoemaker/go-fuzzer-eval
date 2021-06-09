package tls

import (
	"bytes"
	"crypto/tls"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	buf *bytes.Reader
}

func (mc *mockConn) Read(b []byte) (n int, err error) {
	return mc.buf.Read(b)
}
func (mc *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (mc *mockConn) Close() error                       { return nil }
func (mc *mockConn) LocalAddr() net.Addr                { return &net.IPAddr{} }
func (mc *mockConn) RemoteAddr() net.Addr               { return &net.IPAddr{} }
func (mc *mockConn) SetDeadline(t time.Time) error      { return nil }
func (mc *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (mc *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func FuzzTLSClientHandshakeEmptySeed(f *testing.F) {
	f.Fuzz(func(_ *testing.T, b []byte) {
		c := tls.Client(&mockConn{buf: bytes.NewReader(b)}, &tls.Config{ServerName: "non-existent.golang.org"})
		c.Handshake()
		c = tls.Client(&mockConn{buf: bytes.NewReader(b)}, &tls.Config{InsecureSkipVerify: true})
		c.Handshake()
	})
}

func FuzzTLSServerHandshakeEmptySeed(f *testing.F) {
	f.Fuzz(func(_ *testing.T, b []byte) {
		s := tls.Server(&mockConn{buf: bytes.NewReader(b)}, &tls.Config{})
		s.Handshake()
	})
}
