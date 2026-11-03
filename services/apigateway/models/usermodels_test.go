package models

import (
	"encoding/json"
	"testing"
)

func TestUserOnlineStatusResponseIncludesDevices(t *testing.T) {
	resp := UserOnlineStatusResp{
		Items: []*UserOnlineStatusItem{
			{
				UserId:   "user-1",
				IsOnline: true,
				Devices: []*UserOnlineDevice{
					{DeviceId: "device-1", Platform: "Android"},
				},
			},
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal online status response: %v", err)
	}
	want := `{"items":[{"user_id":"user-1","is_online":true,"devices":[{"device_id":"device-1","platform":"Android"}]}]}`
	if string(data) != want {
		t.Fatalf("online status response = %s, want %s", data, want)
	}
}

func TestBanUserSupportsScopeValuesAndLegacyScopeValue(t *testing.T) {
	data := []byte(`{"user_id":"user-1","scope_key":"device","scope_value":"legacy-device","scope_values":["device-1","device-2"]}`)
	var user BanUser
	if err := json.Unmarshal(data, &user); err != nil {
		t.Fatalf("unmarshal ban user: %v", err)
	}
	if user.ScopeValue != "legacy-device" {
		t.Fatalf("scope_value = %q, want legacy-device", user.ScopeValue)
	}
	if len(user.ScopeValues) != 2 || user.ScopeValues[0] != "device-1" || user.ScopeValues[1] != "device-2" {
		t.Fatalf("scope_values = %v, want [device-1 device-2]", user.ScopeValues)
	}
}
