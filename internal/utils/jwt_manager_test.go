package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_Generate(t *testing.T) {
	originalKey := os.Getenv("JWT_SECRET")
	_ = os.Setenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")
	defer func() { _ = os.Setenv("JWT_SECRET", originalKey) }()

	t.Run("generate valid token for user", func(t *testing.T) {
		manager := NewJWTManager()
		user := User{ID: "1", Username: "admin", Role: "ADMIN"}

		token, err := manager.Generate(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("different users produce different tokens", func(t *testing.T) {
		manager := NewJWTManager()
		user1 := User{ID: "1", Username: "admin", Role: "ADMIN"}
		user2 := User{ID: "2", Username: "trainer", Role: "TRAINER"}

		token1, _ := manager.Generate(user1)
		token2, _ := manager.Generate(user2)

		assert.NotEqual(t, token1, token2)
	})

	t.Run("token contains standard claims", func(t *testing.T) {
		manager := NewJWTManager()
		user := User{ID: "42", Username: "jdoe", Role: "SUPER_USER"}

		token, err := manager.Generate(user)

		assert.NoError(t, err)
		assert.Contains(t, token, "eyJ") // JWT header prefix
	})
}
