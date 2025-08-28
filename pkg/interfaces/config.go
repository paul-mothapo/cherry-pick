// Package interfaces defines configuration management interfaces.
package interfaces

import "github.com/cherry-pick/pkg/types"

// ConfigManager defines the interface for handling system configuration.
type ConfigManager interface {
	// LoadConfig loads configuration from various sources.
	LoadConfig(configPath string) error

	// GetConfig returns the current configuration.
	GetConfig() *types.Config

	// UpdateConfig updates the configuration.
	UpdateConfig(config *types.Config) error

	// ValidateConfig validates the configuration.
	ValidateConfig(config *types.Config) error
}
