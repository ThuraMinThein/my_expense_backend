package helper

import (
	"os"
	"testing"

	"github.com/ThuraMinThein/my_expense_backend/config"
)

func TestEncryption(t *testing.T) {
	// Setup config
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012") // 32 bytes key
	config.LoadConfig()

	originalText := "secret data"

	// Test Encrypt
	encryptedText, err := Encrypt(originalText)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if encryptedText == "" {
		t.Fatal("Encrypted text is empty")
	}
	if originalText == encryptedText {
		t.Fatal("Encrypted text is same as original")
	}

	// Test Decrypt
	decryptedText, err := Decrypt(encryptedText)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if originalText != decryptedText {
		t.Fatalf("Decrypted text '%s' does not match original '%s'", decryptedText, originalText)
	}
}

func TestEncryption_InvalidKey(t *testing.T) {
	// Setup invalid config
	os.Setenv("ENCRYPTION_KEY", "shortkey")
	config.LoadConfig()

	_, err := Encrypt("data")
	if err == nil {
		t.Fatal("Expected error with invalid key, got nil")
	}

	_, err = Decrypt("someencodeddata")
	if err == nil {
		t.Fatal("Expected error with invalid key, got nil")
	}
}
