package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	roleKey = "role"
	expKey  = "exp"
)

var (
	jwtAlg       = jwt.SigningMethodHS256
	jwtValidAlgs = []string{jwtAlg.Name}

	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

// Manager is a JWT token manager.
type Manager struct {
	secretKey []byte
	expire    time.Duration
}

// NewManager creates a new JWT token manager.
func NewManager(secretKey string, expire time.Duration) *Manager {
	return &Manager{
		secretKey: []byte(secretKey),
		expire:    expire,
	}
}

// GenerateToken generates a new JWT token with the given role.
func (m *Manager) GenerateToken(role string) (string, error) {
	token := jwt.NewWithClaims(jwtAlg,
		jwt.MapClaims{
			roleKey: role,
			expKey:  time.Now().Add(m.expire).Unix(),
		})

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies the given JWT token. It checks the token signature and expiration time.
func (m *Manager) VerifyToken(tokenString string) error {
	claims, err := m.getClaims(tokenString)
	if err != nil {
		return err
	}

	if err = m.checkExpire(claims); err != nil {
		return err
	}

	return nil
}

// GetRole extracts the role from the given JWT token. It returns a non-nil error if the token is invalid or expired.
func (m *Manager) GetRole(tokenString string) (string, error) {
	claims, err := m.getClaims(tokenString)
	if err != nil {
		return "", err
	}

	if err = m.checkExpire(claims); err != nil {
		return "", err
	}

	role, ok := claims[roleKey].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	return role, nil
}

// getClaims parses the given JWT token and returns the claims. It returns an error if the token is invalid.
func (m *Manager) getClaims(tokenString string) (jwt.MapClaims, error) {
	parserFunc := func(token *jwt.Token) (interface{}, error) { return m.secretKey, nil }
	token, err := jwt.Parse(tokenString, parserFunc, jwt.WithValidMethods(jwtValidAlgs))
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// checkExpire checks if the given JWT token is expired. It returns an error if the token is expired.
func (m *Manager) checkExpire(claims jwt.MapClaims) error {
	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return ErrInvalidToken
	}

	if time.Now().After(exp.UTC()) {
		return ErrTokenExpired
	}

	return nil
}
