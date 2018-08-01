package codeviewer

import (
	"os"
	"path/filepath"
)

var (
	// ConfigDir is the used configuration dir (just caching stuff)
	ConfigDir string
	// StyleDir is the subdirectory for styles
	StyleDir string = "styles"
	// LangDir is the subdirectory for languages
	LangDir string = "languages"
)

// GetConfigDir returns the config directory which is default in the cli tool
func GetConfigDir() string {
	// $XDG_CONFIG_HOME
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, "codeviewer")
	}
	// $HOME/.config
	if d := os.Getenv("HOME"); d != "" {
		return filepath.Join(d, ".config", "codeviewer")
	}
	// nothing suitable found
	return ""
}
