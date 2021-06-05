package json

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func fuzzUnmarshal(t *testing.T, b []byte) {
	for _, typ := range []func() interface{}{
		func() interface{} { return new(interface{}) },
		func() interface{} { return new(map[string]interface{}) },
		func() interface{} { return new([]interface{}) },
	} {
		i := typ()
		if err := json.Unmarshal(b, i); err != nil {
			return
		}

		encoded, err := json.Marshal(i)
		if err != nil {
			t.Fatalf("failed to marshal: %s", err)
		}

		if err := json.Unmarshal(encoded, i); err != nil {
			t.Fatalf("failed to roundtrip: %s", err)
		}
	}
}

func FuzzUnmarshalJSONEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, b []byte) {
		fuzzUnmarshal(t, b)
	})
}

func FuzzUnmarshalJSONBasicSeed(f *testing.F) {
	f.Add([]byte(`{
    "object": {
        "slice": [
            1,
            2.0,
            "3",
            [4],
            {5: {}}
        ]
    },
    "slice": [[]],
    "string": ":)",
    "int": 1e5,
    "float": 3e-9"
}`))

	f.Fuzz(func(t *testing.T, b []byte) {
		fuzzUnmarshal(t, b)
	})
}

// TODO: add a target with a saturated corpus

func fuzzDecoderToken(b []byte) {
	r := bytes.NewReader(b)
	d := json.NewDecoder(r)
	for {
		_, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}
}

func FuzzDecoderTokenEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		fuzzDecoderToken(b)
	})
}

func FuzzDecoderTokenBasicSeed(f *testing.F) {
	f.Add([]byte(`{
    "object": {
        "slice": [
            1,
            2.0,
            "3",
            [4],
            {5: {}}
        ]
    },
    "slice": [[]],
    "string": ":)",
    "int": 1e5,
    "float": 3e-9"
}`))

	f.Fuzz(func(_ *testing.T, b []byte) {
		fuzzDecoderToken(b)
	})
}
