package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"
)

func TestAesDecryptRoundTrip(t *testing.T) {
	for _, keySize := range []int{16, 24, 32} {
		t.Run(fmt.Sprint(keySize), func(t *testing.T) {
			key := bytes.Repeat([]byte{42}, keySize)
			for _, plain := range [][]byte{nil, []byte("hello"), bytes.Repeat([]byte("a"), 16), bytes.Repeat([]byte("b"), 33)} {
				enc, err := AesEncrypt(plain, key)
				if err != nil {
					t.Fatal(err)
				}
				original := bytes.Clone(enc)
				dec, err := AesDecrypt(enc, key)
				if err != nil || !bytes.Equal(dec, plain) {
					t.Fatalf("plain=%q decrypted=%q err=%v", plain, dec, err)
				}
				if !bytes.Equal(enc, original) {
					t.Fatal("decrypt mutated ciphertext")
				}
			}
		})
	}
}

func TestAesDecryptInvalidLengthNoPanic(t *testing.T) {
	for _, keySize := range []int{0, 1, 15, 16, 24, 31, 32, 33} {
		for n := 0; n <= 65; n++ {
			if (keySize == 16 || keySize == 24 || keySize == 32) && n > 0 && n%aes.BlockSize == 0 {
				continue
			}
			if got, err := AesDecrypt(make([]byte, n), make([]byte, keySize)); err == nil || got != nil {
				t.Fatalf("key length=%d ciphertext length=%d: got=%x err=%v", keySize, n, got, err)
			}
		}
	}
}

func TestAesDecryptStrictPadding(t *testing.T) {
	key := bytes.Repeat([]byte{42}, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		data []byte
		ok   bool
	}{
		{"one byte", append(bytes.Repeat([]byte{'x'}, 15), 1), true},
		{"four bytes", append(bytes.Repeat([]byte{'x'}, 12), 4, 4, 4, 4), true},
		{"full block", bytes.Repeat([]byte{16}, 16), true},
		{"zero", append(bytes.Repeat([]byte{'x'}, 15), 0), false},
		{"over block", append(bytes.Repeat([]byte{'x'}, 15), 17), false},
		{"inconsistent", append(bytes.Repeat([]byte{'x'}, 13), 1, 3, 3), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			enc := make([]byte, len(tc.data))
			cipher.NewCBCEncrypter(block, key).CryptBlocks(enc, tc.data)
			got, err := AesDecrypt(enc, key)
			if (err == nil) != tc.ok || (!tc.ok && got != nil) {
				t.Fatalf("got=%x err=%v want valid=%v", got, err, tc.ok)
			}
		})
	}
	for _, data := range [][]byte{nil, {2}, {3, 3}} {
		if _, err := pkcs7UnPaddingStrict(data, aes.BlockSize); err == nil {
			t.Fatalf("accepted short padding: %x", data)
		}
	}
}

func FuzzAesDecryptNoPanic(f *testing.F) {
	f.Add([]byte{}, []byte("0123456789abcdef"))
	f.Add([]byte{1}, []byte{})
	f.Add(make([]byte, 16), []byte("0123456789abcdef"))
	f.Fuzz(func(t *testing.T, ciphertext, key []byte) {
		_, _ = AesDecrypt(ciphertext, key)
	})
}
