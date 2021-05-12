package main

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"
)

func FuzzTarReaderEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		r := tar.NewReader(bytes.NewReader(b))
		for {
			_, err := r.Next()
			if err == io.EOF {
				break // End of archive
			}
			if err != nil {
				return
			}
			if _, err := io.Copy(io.Discard, r); err != nil {
				return
			}
		}
	})
}
