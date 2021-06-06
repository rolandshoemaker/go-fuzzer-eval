package url

import (
	"net/url"
	"reflect"
	"testing"
)

func fuzzParseQuery(t *testing.T, b []byte) {
	query, err := url.ParseQuery(string(b))
	if err != nil {
		t.Skip()
	}
	queryStr2 := query.Encode()
	query2, err := url.ParseQuery(queryStr2)
	if err != nil {
		t.Fatalf("ParseQuery failed to decode a valid encoded query %s: %v", queryStr2, err)
	}
	if !reflect.DeepEqual(query, query2) {
		t.Errorf("ParseQuery gave different query after being encoded\nbefore: %v\nafter: %v", query, query2)
	}
}

func FuzzParseQueryEmptySeed(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		fuzzParseQuery(t, b)
	})
}

func FuzzParseQueryBasicSeed(f *testing.F) {
	f.Add([]byte("x=1&y=2&y=3;z"))
	f.Fuzz(func(t *testing.T, b []byte) {
		fuzzParseQuery(t, b)
	})
}
