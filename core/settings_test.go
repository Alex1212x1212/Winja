package core

import (
	"os"
	"testing"
)

func TestJSONSettings(t *testing.T) {
	// Test setting and getting
	err := SetPrivateMode(true)
	if err != nil {
		t.Fatalf("Failed to set PrivateMode: %v", err)
	}
	if !IsPrivateMode() {
		t.Fatal("Expected PrivateMode to be true")
	}

	SetPrivateMode(false)
	if IsPrivateMode() {
		t.Fatal("Expected PrivateMode to be false")
	}

	SetAPIKey("QA_TEST_KEY_123")
	if GetAPIKey() != "QA_TEST_KEY_123" {
		t.Fatal("API Key did not match")
	}
	
	SetNotifLevel(2)
	if GetNotifLevel() != 2 {
		t.Fatal("Notification level did not match")
	}
}

func TestRegistrySettings(t *testing.T) {
	// Startup
	err := SetRunAtStartup(false)
	if err != nil {
		t.Fatalf("Failed to write to registry (Startup): %v", err)
	}

	// Context menu
	err = SetExplorerExtension(false)
	if err != nil {
		t.Fatalf("Failed to write to registry (Explorer): %v", err)
	}
}

func TestHashing(t *testing.T) {
	// Create a dummy file
	tmpFile := "qa_test_file.txt"
	os.WriteFile(tmpFile, []byte("Winja QA Test"), 0644)
	defer os.Remove(tmpFile)

	hash, err := HashFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to hash file: %v", err)
	}
	if len(hash) != 64 {
		t.Fatalf("Expected SHA-256 hash to be 64 chars, got %d", len(hash))
	}
}
