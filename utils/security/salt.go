package security

import (
	"crypto/rand"
	"os"
	"strconv"
)

func createSalt() ([]byte, error) {
	saltSize, err := strconv.Atoi(os.Getenv("SALT_SIZE"))
	if err != nil {
		return nil, err
	}
	b := make([]byte, saltSize)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
