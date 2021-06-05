package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"
)

func FuzzImageDecodeEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, b []byte) {
		cfg, _, err := image.DecodeConfig(bytes.NewReader(b))
		if err != nil {
			return
		}
		if cfg.Width*cfg.Height > 1e6 {
			return
		}
		img, typ, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return
		}
		switch typ {
		case "png":
			levels := []png.CompressionLevel{
				png.DefaultCompression,
				png.NoCompression,
				png.BestSpeed,
				png.BestCompression,
			}
			for _, l := range levels {
				var w bytes.Buffer
				e := &png.Encoder{CompressionLevel: l}
				err = e.Encode(&w, img)
				if err != nil {
					t.Fatalf("failed to encode valid image: %s", err)
				}
				img1, err := png.Decode(&w)
				if err != nil {
					t.Fatalf("failed to decode roundtripped image: %s", err)
				}
				got := img1.Bounds()
				want := img.Bounds()
				if !got.Eq(want) {
					t.Fatalf("roundtripped image bounds have changed, got: %s, want: %s", got, want)
				}
			}
		case "gif":
			for q := 1; q <= 256; q++ {
				var w bytes.Buffer
				err := gif.Encode(&w, img, &gif.Options{NumColors: q})
				if err != nil {
					t.Fatalf("failed to encode valid image: %s", err)
				}
				img1, err := gif.Decode(&w)
				if err != nil {
					t.Fatalf("failed to decode roundtripped image: %s", err)
				}
				got := img1.Bounds()
				want := img.Bounds()
				if !got.Eq(want) {
					t.Fatalf("roundtripped image bounds have changed, got: %s, want: %s", got, want)
				}
			}
		case "jpeg":
			for q := 1; q <= 100; q++ {
				var w bytes.Buffer
				err := jpeg.Encode(&w, img, &jpeg.Options{Quality: q})
				if err != nil {
					t.Fatalf("failed to encode valid image: %s", err)
				}
				img1, err := jpeg.Decode(&w)
				if err != nil {
					t.Fatalf("failed to decode roundtripped image: %s", err)
				}
				got := img1.Bounds()
				want := img.Bounds()
				if !got.Eq(want) {
					t.Fatalf("roundtripped image bounds have changed, got: %s, want: %s", got, want)
				}
			}
		}
		return
	})
}
