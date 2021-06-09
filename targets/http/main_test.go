package html

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
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

func FuzzClientResponseEmptySeed(f *testing.F) {
	f.Fuzz(func(_ *testing.T, b []byte) {
		c := &http.Client{
			Transport: &http.Transport{
				Dial: func(_, _ string) (net.Conn, error) {
					return &mockConn{buf: bytes.NewReader(b)}, nil
				},
			},
		}

		// Supress output from default logger, since the net/http
		// package is very loud
		log.Default().SetOutput(io.Discard)

		for _, method := range []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		} {
			req, err := http.NewRequest(method, "http://non-existent.golang.org", nil)
			if err != nil {
				continue
			}
			c.Do(req)
		}
	})
}
