package tokens

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
)

func TestParseMalformedCiphertextNoPanic(t *testing.T) {
	key := bytes.Repeat([]byte{42}, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	// Construct a whole block with invalid padding, in addition to partial blocks.
	badPadding := make([]byte, aes.BlockSize)
	cipher.NewCBCEncrypter(block, key).CryptBlocks(badPadding, make([]byte, aes.BlockSize))
	for _, ciphertext := range [][]byte{nil, {1}, make([]byte, 15), make([]byte, 17), badPadding} {
		encoded, err := tools.PbMarshal(&pbobjs.TokenWrap{AppKey: "test-app", TokenValue: ciphertext})
		if err != nil {
			t.Fatal(err)
		}
		wrap, err := ParseTokenString(base64.URLEncoding.EncodeToString(encoded))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ParseToken(wrap, key); err == nil {
			t.Fatalf("accepted malformed token ciphertext: %x", ciphertext)
		}
	}
}
