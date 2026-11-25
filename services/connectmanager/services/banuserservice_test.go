package services

import (
	"im-server/commons/tools"
	"im-server/services/connectmanager/dbs"
	"reflect"
	"testing"
)

func TestNormalizeBanScopeValuesSupportsLegacyAndArrayInputs(t *testing.T) {
	got := normalizeBanScopeValues(string(dbs.UserBanScopeDevice), " device-2,device-1,device-2 ", []string{"device-3", " device-1 ", ""})
	want := []string{"device-1", "device-2", "device-3"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeBanScopeValues() = %v, want %v", got, want)
	}
}

func TestNormalizeBanScopeValuesUsesEmptyValueForDefaultScope(t *testing.T) {
	got := normalizeBanScopeValues(string(dbs.UserBanScopeDefault), "ignored", []string{"ignored-too"})
	if !reflect.DeepEqual(got, []string{""}) {
		t.Fatalf("normalizeBanScopeValues(default) = %v, want empty value", got)
	}
}

func TestBanUserItemActiveBanMatchesExactScopeValue(t *testing.T) {
	now := int64(1000)
	banUser := &BanUserItem{UserId: "user-1", Items: map[string]map[string]*BanItem{
		string(dbs.UserBanScopeDevice): {
			"device-12":  {EndTime: 0, Ext: "blocked-device-12"},
			"device-old": {EndTime: now, Ext: "expired"},
		},
	}}
	if _, blocked := banUser.ActiveBan(string(dbs.UserBanScopeDevice), "device-123", now); blocked {
		t.Fatal("device-123 must not match the device-12 rule")
	}
	if rule, blocked := banUser.ActiveBan(string(dbs.UserBanScopeDevice), "device-12", now); !blocked || rule.Ext != "blocked-device-12" {
		t.Fatalf("device-12 rule = (%+v, %v), want active exact match", rule, blocked)
	}
	if _, blocked := banUser.ActiveBan(string(dbs.UserBanScopeDevice), "device-old", now); blocked {
		t.Fatal("a rule ending at the current time must be expired")
	}
}

func TestNewBanRuleOverridesMatchingLegacyCacheRule(t *testing.T) {
	banUser := &BanUserItem{UserId: "user-1", Items: make(map[string]map[string]*BanItem)}
	addBanItem(banUser, string(dbs.UserBanScopeDevice), "device-1", 0, "legacy")
	addBanItem(banUser, string(dbs.UserBanScopeDevice), "device-1", 0, "new")
	rule, blocked := banUser.ActiveBan(string(dbs.UserBanScopeDevice), "device-1", 1000)
	if !blocked || rule.Ext != "new" {
		t.Fatalf("active rule = (%+v, %v), want the new rule to take precedence", rule, blocked)
	}
}

func TestBanQueryOffsetSupportsLegacyAndRulePhases(t *testing.T) {
	encodedId, err := tools.EncodeInt(42)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, offset, wantPhase string
		wantId                  int64
	}{
		{name: "old numeric offset", offset: encodedId, wantPhase: banQueryPhaseLegacy, wantId: 42},
		{name: "legacy phase", offset: encodeBanQueryOffset(banQueryPhaseLegacy, 42), wantPhase: banQueryPhaseLegacy, wantId: 42},
		{name: "rules phase", offset: encodeBanQueryOffset(banQueryPhaseRules, 42), wantPhase: banQueryPhaseRules, wantId: 42},
	} {
		t.Run(test.name, func(t *testing.T) {
			phase, id := decodeBanQueryOffset(test.offset)
			if phase != test.wantPhase || id != test.wantId {
				t.Fatalf("decodeBanQueryOffset(%q) = (%q, %d), want (%q, %d)", test.offset, phase, id, test.wantPhase, test.wantId)
			}
		})
	}
}
