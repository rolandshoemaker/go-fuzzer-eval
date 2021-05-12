package main

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func FuzzZIPReaderEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
		if err != nil {
			return
		}
		for _, f := range r.File {
			if fr, err := f.Open(); err == nil {
				_, _ = io.ReadAll(fr)
			}
			if _, err := r.Open(f.Name); err != nil {
				continue
			}
		}
	})
}
