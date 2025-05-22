package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestIsAuthByBearerToken_SignatureInvalidAndContext(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("user")
		if claims == nil {
			http.Error(w, "claims ausentes", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := IsAuthByBearerToken(next)

	t.Run("Token ausente retorna 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Sem header Authorization
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Token ausente")
	})

	t.Run("Formato de token inválido retorna 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormatToken")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Formato de token inválido")
	})

	t.Run("Método assinatura inválido retorna erro 401", func(t *testing.T) {
		tokenString := "invalid.token.parts" // token inválido para forçar erro de parse/assinatura

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Token inválido")
	})

	t.Run("Token válido passa claims no contexto", func(t *testing.T) {
		claims := jwt.MapClaims{"foo": "bar"}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.LoadConfig().Jwt.SecretKey))
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Callback retorna jwt.ErrSignatureInvalid para método não HMAC", func(t *testing.T) {
		// JWT com alg RS256 no header, mas assinatura inválida (qualquer string qualquer)
		tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30.invalidsignature"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Token inválido")
	})

}
