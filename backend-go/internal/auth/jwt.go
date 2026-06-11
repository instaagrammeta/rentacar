package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// TokenType distinguishes access tokens from refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims is the JWT payload. Subject holds the user id (as in the Flask
// implementation: identity=str(user.id)); Role and Username are extra claims.
type Claims struct {
	Role     models.UserRole `json:"role"`
	Username string          `json:"username"`
	Type     TokenType       `json:"type"`
	jwt.RegisteredClaims
}

// Manager issues and validates JWT tokens.
type Manager struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewManager builds a token manager.
func NewManager(secret string, access, refresh time.Duration) *Manager {
	return &Manager{secret: []byte(secret), accessExpiry: access, refreshExpiry: refresh}
}

func (m *Manager) generate(user *models.User, tt TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Role:     user.Role,
		Username: user.Username,
		Type:     tt,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// GenerateTokens returns a fresh access + refresh token pair for the user.
func (m *Manager) GenerateTokens(user *models.User) (access string, refresh string, err error) {
	access, err = m.generate(user, AccessToken, m.accessExpiry)
	if err != nil {
		return "", "", err
	}
	refresh, err = m.generate(user, RefreshToken, m.refreshExpiry)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ErrInvalidToken is returned for any token validation failure.
var ErrInvalidToken = errors.New("invalid token")

// Parse validates a token string and returns its claims.
func (m *Manager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
