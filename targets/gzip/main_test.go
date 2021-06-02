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
		for _, multistream := range []bool{true, false} {
			r, err := gzip.NewReader(bytes.NewBuffer(b))
			if err != nil {
				continue
			}

			r.Multistream(multistream)

			decompressed := bytes.NewBuffer(nil)
			if _, err := io.Copy(decompressed, r); err != nil {
				continue
			}

			if err := r.Close(); err != nil {
				continue
			}

			for _, level := range []int{gzip.NoCompression, gzip.BestSpeed, gzip.BestCompression, gzip.DefaultCompression, gzip.HuffmanOnly} {
				w, err := gzip.NewWriterLevel(io.Discard, level)
				if err != nil {
					continue
				}
				w.Write(decompressed.Bytes())
				if err := w.Flush(); err != nil {
					continue
				}
				w.Close()
			}
		}
	})
}
