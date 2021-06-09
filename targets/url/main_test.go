package url

import (
	"net/url"
	"reflect"
	"testing"
)

func fuzzParseQuery(t *testing.T, b []byte) {
	query, _ := url.ParseQuery(string(b))
	if len(query) == 0 {
		t.Skip()
	}
	queryStr2 := query.Encode()
	query2, _ := url.ParseQuery(queryStr2)
	if len(query2) == 0 {
		t.Skip()
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
