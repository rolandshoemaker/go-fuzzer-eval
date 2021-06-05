package dns

import (
	"bytes"
	"context"
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

func FuzzLookupHostEmptySeed(f *testing.F) {
	f.Fuzz(func(_ *testing.T, b []byte) {
		r := &net.Resolver{
			Dial: func(_ context.Context, network, _ string) (net.Conn, error) {
				return &mockConn{buf: bytes.NewReader(b)}, nil
			},
		}
		r.LookupHost(context.Background(), "")
	})
}
