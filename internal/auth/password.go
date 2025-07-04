package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("senha inválida")

// PasswordHasher define a abstração para hashing de senhas
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

// BcryptHasher é a implementação real usando bcrypt
type BcryptHasher struct{}

func (BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (BcryptHasher) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
