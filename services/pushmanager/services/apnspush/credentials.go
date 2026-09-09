package apnspush

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const maxPrivateKeySize = 16 * 1024

func validateCredentials(privateKey []byte, keyID, teamID string) error {
	if !validAppleID(keyID) || !validAppleID(teamID) {
		return errors.New("P8 Key ID and Team ID must be 10 uppercase letters or digits")
	}
	invalid := errors.New("invalid P8 private key: expected PKCS#8 PEM ECDSA P-256, at most 16 KiB")
	if len(privateKey) == 0 || len(privateKey) > maxPrivateKeySize {
		return invalid
	}
	input := bytes.TrimSpace(privateKey)
	if !bytes.HasPrefix(input, []byte("-----BEGIN PRIVATE KEY-----")) {
		return invalid
	}
	block, rest := pem.Decode(input)
	if block == nil || block.Type != "PRIVATE KEY" || len(block.Headers) != 0 || len(bytes.TrimSpace(rest)) != 0 {
		return invalid
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return invalid
	}
	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok || key.Curve != elliptic.P256() || key.D == nil || key.D.Sign() <= 0 || key.D.Cmp(key.Params().N) >= 0 || !key.Curve.IsOnCurve(key.X, key.Y) {
		return invalid
	}
	x, y := key.Curve.ScalarBaseMult(key.D.Bytes())
	if x.Cmp(key.X) != 0 || y.Cmp(key.Y) != 0 {
		return invalid
	}
	return nil
}

func validAppleID(value string) bool {
	if len(value) != 10 {
		return false
	}
	for i := range value {
		c := value[i]
		if !(c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

// validTopic validates a bundle before the optional .voip suffix is appended.
func validTopic(value string) bool {
	if len(value) == 0 || len(value) > 100 {
		return false
	}
	for i := range value {
		c := value[i]
		if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '.' || c == '-') {
			return false
		}
	}
	return true
}
