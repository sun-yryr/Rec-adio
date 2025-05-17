package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sun-yryr/recoto/internal/fileutil"
)

// TestLoadConfig_FileNotExists ensures that LoadConfig creates a missing file
// and applies the correct default file permissions.
func TestLoadConfig_FileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Attempt to load config, expecting a new file to be created
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check that config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if perms := info.Mode().Perm(); perms != fileutil.DefaultFilePerm {
		t.Errorf("Expected config file perm %#o, got %#o", fileutil.DefaultFilePerm, perms)
	}

	// Ensure we got a non-nil config back
	if cfg == nil {
		t.Fatal("Expected non-nil config, got nil")
	}
}

// TestSaveConfig ensures that SaveConfig writes a file with the correct
// permissions.
func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Prepare a dummy config to save
	cfg := &Config{}

	// Save the config, expecting the file to be created
	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Check that config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if perms := info.Mode().Perm(); perms != fileutil.DefaultFilePerm {
		t.Errorf("Expected config file perm %#o, got %#o", fileutil.DefaultFilePerm, perms)
	}
}