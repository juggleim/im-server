package dbs

import (
	"reflect"
	"testing"
)

func TestSplitLegacyBanScopeValues(t *testing.T) {
	got := SplitLegacyBanScopeValues(string(UserBanScopeDevice), " device-2,device-1,device-2, ")
	want := []string{"device-1", "device-2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SplitLegacyBanScopeValues() = %v, want %v", got, want)
	}
	if got := SplitLegacyBanScopeValues("", "ignored"); !reflect.DeepEqual(got, []string{""}) {
		t.Fatalf("blank legacy scope = %v, want default empty value", got)
	}
}
