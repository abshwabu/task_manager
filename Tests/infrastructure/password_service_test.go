package infrastructure_test

import (
	"example/task_manager/Infrastructure"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword, err := infrastructure.HashPassword(password)

	assert.NoError(t, err)
	assert.NotNil(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	isMatch := infrastructure.CheckPasswordHash(password, hashedPassword)
	assert.True(t, isMatch, "Hashed password should match original password")

	// Test with empty password
	emptyPasswordHash, err := infrastructure.HashPassword("")
	assert.Error(t, err)
	assert.EqualError(t, err, "password cannot be empty")
	assert.Empty(t, emptyPasswordHash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword, _ := infrastructure.HashPassword(password)

	// Test with correct password
	isMatch := infrastructure.CheckPasswordHash(password, hashedPassword)
	assert.True(t, isMatch, "Correct password should match hash")

	// Test with incorrect password
	isMatch = infrastructure.CheckPasswordHash("wrongpassword", hashedPassword)
	assert.False(t, isMatch, "Incorrect password should not match hash")

	// The HashPassword function now returns an error for empty password.
	// So, we cannot generate a valid hash for an empty password this way anymore.
	// We will test direct comparison with an empty password and an arbitrary valid hash.
	// Test with empty password and a non-empty hash (should be false)
	isMatch = infrastructure.CheckPasswordHash("", hashedPassword)
	assert.False(t, isMatch, "Empty password should not match a non-empty hash")

	// Test with empty hash
	isMatch = infrastructure.CheckPasswordHash(password, "")
	assert.False(t, isMatch, "Comparing with an empty hash should return false")

	// Test with malformed hash
	isMatch = infrastructure.CheckPasswordHash(password, "malformedhash")
	assert.False(t, isMatch, "Comparing with a malformed hash should return false")
}
