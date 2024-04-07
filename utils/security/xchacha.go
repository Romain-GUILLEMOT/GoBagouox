package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
)

func EncryptXChaCha(plainText string, key []byte) (string, error) {
	if len(key) != 32 { // XChaCha20-Poly1305 nécessite une clé de 32 octets
		return "", errors.New("Invalid key length")
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptXChaCha(encodedText string, key []byte) (string, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}

	if len(cipherText) < chacha20poly1305.NonceSizeX {
		return "", errors.New("Malformed ciphertext")
	}
	nonce := cipherText[:chacha20poly1305.NonceSizeX]
	cipherText = cipherText[chacha20poly1305.NonceSizeX:]

	plaintext, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
