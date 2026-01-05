package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 12)
	return string(bytes), err
}

func VerifyHashed(hashed, str string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str))
	return err
}
