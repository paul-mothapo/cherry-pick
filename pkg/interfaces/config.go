package interfaces

import "github.com/cherry-pick/pkg/types"

// ConfigManager defines the interface for handling system configuration.
type ConfigManager interface {
	LoadConfig(configPath string) error
	GetConfig() *types.Config
	UpdateConfig(config *types.Config) error
	ValidateConfig(config *types.Config) error
}
