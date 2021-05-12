package main

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

func FuzzImageDecodeEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		i, _, err := image.Decode(bytes.NewBuffer(b))
		if err != nil {
			return
		}
		_, _ = i.ColorModel(), i.Bounds()
	})
}
