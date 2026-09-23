package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TypeAccess  TokenType = "access"
	TypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserRole string    `json:"user_role"`
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a short lived access token valid for 15 minutes
func GenerateAccessToken(uid uuid.UUID, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_ACCESS_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("jwt access secret string not set")
	}

	claims := Claims{
		UserID:   uid,
		UserRole: role,
		Type:     TypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // Expires after 15 minutes after creation
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "tickr-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return signedString, nil
}

// GenerateRefreshToken generates a long lived refresh token valid for 7 days
func GenerateRefreshToken(uid uuid.UUID, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("jwt refresh secret string not set")
	}

	claims := Claims{
		UserID:   uid,
		UserRole: role,
		Type:     TypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Expires after 7 days after creation
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "tickr-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return signedString, nil
}

// VerifyAccessToken validates the signature of the access token
func VerifyAccessToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("jwt access secret is not set")
	}

	claims, err := parseAndValidateToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	// Return error if access token is not passed
	if claims.Type != TypeAccess {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// VerifyRefreshToken validates the signature of the refresh token
func VerifyRefreshToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("jwt refresh secret is not set")
	}

	claims, err := parseAndValidateToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	// Return error if refresh token is not passed
	if claims.Type != TypeRefresh {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Helper function to handle parsing and signature validation
func parseAndValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
