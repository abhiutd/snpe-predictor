package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	KeyLength    = 32
	CryptoHeader = fmt.Sprintf("==AES%d==", KeyLength)
	KeyBuffer    = []byte(strings.Repeat("=", 32))
)

func IsEncryptedString(s string) bool {
	return strings.HasPrefix(s, CryptoHeader)
}

func IsEncrypted(b []byte) bool {
	return IsEncryptedString(string(b))
}

func EncryptStringBase64(key, text string) (string, error) {
	e, err := EncryptString(key, text)
	if err != nil {
		return "", err
	}
	return addCryptoHeaderS(base64.StdEncoding.EncodeToString([]byte(e))), nil
}

func EncryptString(key, text string) (string, error) {
	b, err := Encrypt([]byte(key), []byte(text))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Encrypt(key, text []byte) ([]byte, error) {
	if len(key) != KeyLength {
		key = append(key, KeyBuffer[:KeyLength-len(key)]...)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return addCryptoHeader(ciphertext), nil
}

func removeCryptoHeader(b []byte) []byte {
	return []byte(removeCryptoHeaderS(string(b)))
}

func removeCryptoHeaderS(b string) string {
	return strings.TrimPrefix(string(b), CryptoHeader)
}

func addCryptoHeader(b []byte) []byte {
	return append([]byte(CryptoHeader), b...)
}

func addCryptoHeaderS(b string) string {
	return CryptoHeader + b
}

func DecryptStringBase64(key, text string) (string, error) {
	q := removeCryptoHeader([]byte(text))
	s, err := base64.StdEncoding.DecodeString(string(q))
	if err != nil {
		s = q
	}

	b, err := Decrypt([]byte(key), s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func DecryptString(key, text string) (string, error) {
	b, err := Decrypt([]byte(key), []byte(text))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Decrypt(key, text []byte) ([]byte, error) {
	if len(key) != KeyLength {
		key = append(key, KeyBuffer[:KeyLength-len(key)]...)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	text = removeCryptoHeader(text)
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
