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
		type file struct {
			header  *tar.Header
			content []byte
		}
		files := []file{}
		for {
			hdr, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, r); err != nil {
				return // or continue?
			}
			files = append(files, file{header: hdr, content: buf.Bytes()})
		}

		w := tar.NewWriter(io.Discard)
		for _, f := range files {
			if err := w.WriteHeader(f.header); err != nil {
				continue
			}
			if _, err := w.Write(f.content); err != nil {
				continue
			}
		}
		w.Close()
	})
}
