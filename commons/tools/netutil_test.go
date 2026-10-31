package tools

import "testing"

func TestParseRemoteIP(t *testing.T) {
	tests := []struct {
		remoteAddress string
		want          string
		ok            bool
	}{
		{remoteAddress: "203.0.113.9:1234", want: "203.0.113.9", ok: true},
		{remoteAddress: "203.0.113.9", want: "203.0.113.9", ok: true},
		{remoteAddress: "[2001:db8::1]:1234", want: "2001:db8::1", ok: true},
		{remoteAddress: "unknown", ok: false},
	}
	for _, test := range tests {
		address, ok := ParseRemoteIP(test.remoteAddress)
		if ok != test.ok {
			t.Fatalf("ParseRemoteIP(%q) ok=%v, want=%v", test.remoteAddress, ok, test.ok)
		}
		if ok && address.String() != test.want {
			t.Fatalf("ParseRemoteIP(%q)=%q, want=%q", test.remoteAddress, address, test.want)
		}
	}
}
