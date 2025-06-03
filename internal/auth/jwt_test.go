package auth

import (
	"strings"
	"testing"
	"time"
)

func TestJWTManager_Generate(t *testing.T) {
	jwtManager := NewJWTManager("minha-chave-secreta", time.Minute*15)

	token, err := jwtManager.Generate(123, "teste@email.com")
	if err != nil {
		t.Fatalf("esperado token válido, mas erro ocorreu: %v", err)
	}

	if token == "" {
		t.Fatal("token gerado está vazio")
	}

	if !strings.Contains(token, ".") {
		t.Error("token gerado não parece ser um JWT válido (sem '.' delimitadores)")
	}
}
