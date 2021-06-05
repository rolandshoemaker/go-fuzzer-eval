package main

import (
	"bytes"
	"io"
	"testing"

	"golang.org/x/net/html"
)

func FuzzParseEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		node, err := html.Parse(bytes.NewReader(b))
		if err != nil {
			return
		}
		// TODO: may want to walk the document
		if err := html.Render(io.Discard, node); err != nil {
			return
		}
		// TODO: may want to compare input and output documents match
	})
}

func FuzzParseWithOptionsEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		node, err := html.ParseWithOptions(bytes.NewReader(b), html.ParseOptionEnableScripting(false))
		if err != nil {
			return
		}
		if err := html.Render(io.Discard, node); err != nil {
			return
		}
	})
}

func FuzzParseFragmentEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		nodes, err := html.ParseFragment(bytes.NewReader(b), nil)
		if err != nil {
			return
		}
		for _, n := range nodes {
			if err := html.Render(io.Discard, n); err != nil {
				return
			}
		}
	})
}

func FuzzParseFragmentWithOptionsEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		nodes, err := html.ParseFragmentWithOptions(bytes.NewReader(b), nil, html.ParseOptionEnableScripting(false))
		if err != nil {
			return
		}
		for _, n := range nodes {
			if err := html.Render(io.Discard, n); err != nil {
				return
			}
		}
	})
}

func FuzzTokenizerEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		t := html.NewTokenizer(bytes.NewReader(b))
		for {
			tt := t.Next()
			if tt == html.ErrorToken {
				return
			}
			_ = t.Token()
			// TODO: may want to exercise lower level APIs by operating on the TokenType
		}
	})
}
