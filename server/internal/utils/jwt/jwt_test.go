package jwt

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAndVerifyAccessToken(t *testing.T) {
	// Set secret env
	os.Setenv("JWT_ACCESS_SECRET", "super-secret-key")
	defer os.Unsetenv("JWT_ACCESS_SECRET")

	userID := uuid.New()
	role := "user"

	// Generate new access token
	accessKey, err := GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}
	if accessKey == "" {
		t.Fatalf("expected token string not to be empty")
	}

	// Verify the generated access token
	claims, err := VerifyAccessToken(accessKey)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected user ID %v, got %v", userID, claims.UserID)
	}
	if claims.UserRole != role {
		t.Errorf("expected role %v, got %v", role, claims.UserRole)
	}
	if claims.Type != TypeAccess {
		t.Errorf("expected token type %s, got %s", TypeAccess, claims.Type)
	}

}

func TestGenerateAndVerifyRefreshToken(t *testing.T) {
	// Set secret env
	os.Setenv("JWT_REFRESH_SECRET", "super-secret-key")
	defer os.Unsetenv("JWT_REFRESH_SECRET")

	userID := uuid.New()
	role := "user"

	// Generate new refresh token
	refreshKey, err := GenerateRefreshToken(userID, role)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}
	if refreshKey == "" {
		t.Fatalf("expected token string not to be empty")
	}

	// Verify the generated refresh token
	claims, err := VerifyRefreshToken(refreshKey)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected user ID %v, got %v", userID, claims.UserID)
	}
	if claims.UserRole != role {
		t.Errorf("expected role %v, got %v", role, claims.UserRole)
	}
	if claims.Type != TypeRefresh {
		t.Errorf("expected token type %s, got %s", TypeRefresh, claims.Type)
	}
}
