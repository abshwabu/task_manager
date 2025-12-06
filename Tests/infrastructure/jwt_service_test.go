package infrastructure_test

import (
	"example/task_manager/Infrastructure"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	username := "testuser"
	isAdmin := false
	tokenString, err := infrastructure.GenerateJWT(username, isAdmin)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := infrastructure.ValidateJWT(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, isAdmin, claims.IsAdmin)

	// Test admin user
	adminUsername := "adminuser"
	adminIsAdmin := true
	adminTokenString, err := infrastructure.GenerateJWT(adminUsername, adminIsAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, adminTokenString)

	adminClaims, err := infrastructure.ValidateJWT(adminTokenString)
	assert.NoError(t, err)
	assert.NotNil(t, adminClaims)
	assert.Equal(t, adminUsername, adminClaims.Username)
	assert.Equal(t, adminIsAdmin, adminClaims.IsAdmin)
}

func TestValidateJWT(t *testing.T) {
	username := "testuser"
	isAdmin := false
	tokenString, _ := infrastructure.GenerateJWT(username, isAdmin)

	// Test with a valid token
	claims, err := infrastructure.ValidateJWT(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, isAdmin, claims.IsAdmin)

	// Test with an invalid token string
	invalidTokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiaXNfYWRtaW4iOmZhbHNlLCJleHAiOjE2NzgwNTI0MjN9.invalid_signature"
	claims, err = infrastructure.ValidateJWT(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Test with an expired token
	expiredClaims := &infrastructure.Claims{
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired an hour ago
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte("secret_key"))

	claims, err = infrastructure.ValidateJWT(expiredTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Test with empty token string
	claims, err = infrastructure.ValidateJWT("")
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Test with a token signed with a different key
	differentKey := []byte("different_secret")
	differentToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &infrastructure.Claims{
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	differentSignedToken, _ := differentToken.SignedString(differentKey)
	claims, err = infrastructure.ValidateJWT(differentSignedToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
