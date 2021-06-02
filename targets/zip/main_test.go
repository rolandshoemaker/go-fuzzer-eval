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

		type file struct {
			header  *zip.FileHeader
			content []byte
		}
		files := []file{}

		for _, f := range r.File {
			fr, err := f.Open()
			if err != nil {
				continue
			}
			content, err := io.ReadAll(fr)
			if err != nil {
				continue
			}
			files = append(files, file{header: &f.FileHeader, content: content})
			if _, err := r.Open(f.Name); err != nil {
				continue
			}
		}

		w := zip.NewWriter(io.Discard)
		for _, f := range files {
			ww, err := w.CreateHeader(f.header)
			if err != nil {
				continue
			}
			if _, err := ww.Write(f.content); err != nil {
				continue
			}
		}
		w.Close()

		// TODO: check roundtrip?
	})
}
