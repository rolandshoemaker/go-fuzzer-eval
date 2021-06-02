package json

import (
	"encoding/json"
	"testing"
)

func fuzzUnmarshal(b []byte) {
	for _, t := range []func() interface{}{
		func() interface{} { return new(interface{}) },
		func() interface{} { return new(map[string]interface{}) },
		func() interface{} { return new([]interface{}) },
	} {
		if err := json.Unmarshal(b, t()); err != nil {
			return
		}

		encoded, err := json.Marshal(t())
		if err != nil {
			return
		}

		if err := json.Unmarshal(encoded, t()); err != nil {
			return
		}
	}
}

func FuzzUnmarshalJSONEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		fuzzUnmarshal(b)
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

	f.Fuzz(func(_ *testing.T, b []byte) {
		fuzzUnmarshal(b)
	})
}

// TODO: add a target with a saturated corpus

// TODO: add a json.Decoder target
