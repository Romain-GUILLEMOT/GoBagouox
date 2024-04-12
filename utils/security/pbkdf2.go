package security

import (
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strconv"
)

func pbkdf2Encode(salt []byte) ([]byte, error) {
	iteration := os.Getenv("PBKDF2_ITERATIONS")
	key := os.Getenv("WEBSERVER_SHARED_KEY")

	iterationInt, err := strconv.Atoi(iteration)
	if err != nil {
		return nil, err
	}
	hashedKey := pbkdf2.Key([]byte(key), salt, iterationInt, 32, sha256.New)
	return hashedKey, nil
}
