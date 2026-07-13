package dbs

import (
	"testing"

	"im-server/commons/pbdefines/pbobjs"
)

func TestParseConverExtsWithOnlyGlobalConverTags(t *testing.T) {
	want := &pbobjs.ConverExts{
		GlobalConverTags: map[string]bool{
			"global_tag": true,
		},
	}

	got := parseConverExts(converExts2Bs(want))
	if got == nil {
		t.Fatal("parseConverExts() returned nil for global conversation tags")
	}
	if !got.GlobalConverTags["global_tag"] {
		t.Fatalf("parseConverExts() GlobalConverTags = %v, want global_tag=true", got.GlobalConverTags)
	}
}
