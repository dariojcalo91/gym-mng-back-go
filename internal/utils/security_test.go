package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("hash and verify valid password", func(t *testing.T) {
		password := "mySecretPassword123"
		hash, err := HashPassword(password)

		assert.NoError(t, err)
		assert.True(t, CheckPasswordHash(password, hash))
		assert.False(t, CheckPasswordHash("wrongPassword", hash))
	})

	t.Run("empty password", func(t *testing.T) {
		hash, err := HashPassword("")
		assert.NoError(t, err)
		assert.True(t, CheckPasswordHash("", hash))
	})

	t.Run("different passwords produce different hashes", func(t *testing.T) {
		hash1, _ := HashPassword("password1")
		hash2, _ := HashPassword("password2")
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("same password produces different hashes (bcrypt salt)", func(t *testing.T) {
		hash1, _ := HashPassword("samepassword")
		hash2, _ := HashPassword("samepassword")
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestEncryptEmail(t *testing.T) {
	originalKey := os.Getenv("AES_MASTER_KEY")
	_ = os.Setenv("AES_MASTER_KEY", "abcdefghijklmnopqrstuvwxyz123456")
	defer func() { _ = os.Setenv("AES_MASTER_KEY", originalKey) }()

	t.Run("encrypt and decrypt email", func(t *testing.T) {
		email := "user@example.com"
		encrypted, err := EncryptEmail(email)

		assert.NoError(t, err)
		assert.NotEqual(t, email, encrypted)
		assert.NotEmpty(t, encrypted)
	})

	t.Run("different encryptions produce different ciphertexts", func(t *testing.T) {
		email := "same@example.com"
		enc1, _ := EncryptEmail(email)
		enc2, _ := EncryptEmail(email)

		assert.NotEqual(t, enc1, enc2)
	})

	t.Run("empty email", func(t *testing.T) {
		encrypted, err := EncryptEmail("")

		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)
	})

	t.Run("fails without master key", func(t *testing.T) {
		_ = os.Unsetenv("AES_MASTER_KEY")
		_, err := EncryptEmail("test@mail.com")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid key size")
	})
}
