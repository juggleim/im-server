package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length <= 0 {
		return origData
	}
	unpadding := int(origData[length-1])
	if unpadding > length {
		return origData
	}
	return origData[:(length - unpadding)]
}

// AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// CBC panics on partial blocks. Tokens may contain untrusted ciphertext.
	if len(crypted) == 0 || len(crypted)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	return pkcs7UnPaddingStrict(origData, blockSize)
}

func pkcs7UnPaddingStrict(data []byte, blockSize int) ([]byte, error) {
	n := len(data)
	if n == 0 {
		return nil, errors.New("invalid padding")
	}
	p := int(data[n-1])
	if p == 0 || p > blockSize || p > n {
		return nil, errors.New("invalid padding")
	}
	for _, b := range data[n-p:] {
		if int(b) != p {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:n-p], nil
}
