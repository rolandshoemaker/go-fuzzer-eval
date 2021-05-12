package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func FuzzGZIPReaderEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		r, err := gzip.NewReader(bytes.NewBuffer(b))
		if err != nil {
			return
		}
		if _, err := io.Copy(io.Discard, r); err != nil {
			return
		}

		if err := r.Close(); err != nil {
			return
		}
	})
}
