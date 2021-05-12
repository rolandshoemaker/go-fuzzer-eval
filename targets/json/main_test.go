package json

import (
	"encoding/json"
	"testing"
)

func FuzzUnmarshalJSONEmptySeed(f *testing.F) {
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, b []byte) {
		var i interface{}
		json.Unmarshal(b, &i)
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
		var i interface{}
		json.Unmarshal(b, &i)
	})
}

// TODO: add a target with a saturated corpus
