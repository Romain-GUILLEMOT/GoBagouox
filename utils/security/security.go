package security

import (
	"GoBagouox/utils"
	"encoding/base64"
)

func Encrypt(content string) (string, string, error) {
	salt, err := createSalt()
	if err != nil {
		utils.Error("Error during salt creation.", err, 0)
		return "", "", err
	}
	encryptedKey, err := pbkdf2Encode(salt)
	if err != nil {
		utils.Error("Error during PBKDF2 creation.", err, 0)
		return "", "", err
	}
	encryptedText, err := encryptXChaCha(content, encryptedKey)
	if err != nil {
		utils.Error("Error during XCha encoding.", err, 0)
		return "", "", err
	}
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	return encryptedText, saltBase64, nil

}
func Decrypt(content string, saltBase64 string) (string, error) {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		utils.Error("Error during Salt decoding.", err, 0)
		return "", err
	}
	encryptedKey, err := pbkdf2Encode(salt)
	if err != nil {
		utils.Error("Error during PBKDF2 creation.", err, 0)
		return "", err
	}
	encryptedText, err := decryptXChaCha(content, encryptedKey)
	if err != nil {
		utils.Error("Error during XCha encoding.", err, 0)
		return "", err
	}
	return encryptedText, nil
}
