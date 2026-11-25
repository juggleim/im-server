package apis

import (
	"im-server/commons/tools"
	"testing"
)

func TestDualSecretSignature(t *testing.T) {
	for _, secret := range []string{"primary", "temporary"} {
		sig := tools.SHA1(secret + "nonce" + "123")
		if !isValidSignature(sig, "nonce", "123", "primary", "temporary") {
			t.Fatalf("rejected %s", secret)
		}
	}
	if isValidSignature(tools.SHA1("nonce123"), "nonce", "123", "", "") {
		t.Fatal("empty secrets accepted")
	}
	if isValidSignature(tools.SHA1("othernonce123"), "nonce", "123", "primary", "temporary") {
		t.Fatal("unknown secret accepted")
	}
}
