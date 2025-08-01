package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenInvalid         = errors.New("token inválido")
	ErrTokenExpired         = errors.New("token expirado")
	ErrInvalidSigningMethod = errors.New("algoritmo de assinatura inválido")
	ErrInvalidIssuer        = errors.New("issuer inválido")
	ErrInvalidAudience      = errors.New("audience inválido")
	ErrTokenMissing         = errors.New("token de autenticação ausente")
	ErrTokenInvalidFormat   = errors.New("formato do token inválido")
	ErrTokenRevoked         = errors.New("token revogado")
	ErrInvalidSignature     = errors.New("assinatura do token inválida")
	ErrInvalidExpClaim      = errors.New("claim 'exp' inválida ou ausente")
	ErrInternalAuth         = errors.New("erro interno na autenticação")
)

type JWTManager struct {
	SecretKey     string
	TokenDuration time.Duration
	Issuer        string
	Audience      string
}

type JWTService interface {
	Generate(uid int64, email string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

func NewJWTManager(secretKey string, duration time.Duration, issuer string, audience string) *JWTManager {
	return &JWTManager{
		SecretKey:     secretKey,
		TokenDuration: duration,
		Issuer:        issuer,
		Audience:      audience,
	}
}

func (j *JWTManager) Generate(uid int64, email string) (string, error) {
	now := time.Now()
	uidStr := strconv.FormatInt(uid, 10)

	claims := jwt.MapClaims{
		"sub":     uidStr,
		"user_id": uidStr,
		"email":   email,
		"iat":     now.Unix(),
		"exp":     now.Add(j.TokenDuration).Unix(),
		"iss":     j.Issuer,
		"aud":     j.Audience,
		"jti":     uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err // mantém original p/ outros tipos
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	// Validação manual de expiração, opcional se usar jwt.ErrTokenExpired
	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	if iss, ok := claims["iss"].(string); !ok || iss != j.Issuer {
		return nil, ErrInvalidIssuer
	}

	if aud, ok := claims["aud"].(string); !ok || aud != j.Audience {
		return nil, ErrInvalidAudience
	}

	return claims, nil
}

func (j *JWTManager) GetExpiration(token *jwt.Token) (time.Duration, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("claims inválidas")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return 0, errors.New("claim 'exp' ausente ou inválida")
	}

	expirationTime := time.Unix(int64(expFloat), 0)
	duration := time.Until(expirationTime)
	return duration, nil
}

func (j *JWTManager) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Garante que seja HMAC E que seja o algoritmo específico esperado (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(j.SecretKey), nil
	})
}
