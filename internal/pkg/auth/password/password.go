package auth

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("senha inválida")

// PasswordHasher define a abstração para hashing de senhas
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

// BcryptHasher é a implementação real usando bcrypt
type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

func (BcryptHasher) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// password.go
func (h BcryptHasher) Hash(password string) (string, error) {
	// Usa o custo configurado, com fallback para DefaultCost
	cost := h.cost
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func (h *BcryptHasher) IsHash(s string) bool {
	// Bcrypt hashes começam com $2a$, $2b$, etc
	// E devem ter pelo menos 60 caracteres (hash bcrypt típico)
	return (strings.HasPrefix(s, "$2a$") ||
		strings.HasPrefix(s, "$2b$") ||
		strings.HasPrefix(s, "$2y$")) && len(s) >= 60
}
