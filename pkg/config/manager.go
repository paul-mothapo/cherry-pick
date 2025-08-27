package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
)

type ConfigManagerImpl struct {
	config *types.Config
}

func NewConfigManager() interfaces.ConfigManager {
	return &ConfigManagerImpl{}
}

func (cm *ConfigManagerImpl) LoadConfig(configPath string) error {
	cm.config = &types.Config{
		DatabaseConnections: make(map[string]string),
		AnalysisSettings: types.AnalysisSettings{
			SampleSize:           1000,
			LargeTableThreshold:  1000000,
			QualityScoreMinimum:  0.7,
			AutoAnalysisInterval: "24h",
		},
		AlertSettings: types.AlertSettings{
			EnableAlerts: true,
		},
		SecuritySettings: types.SecuritySettings{
			EnablePIIDetection: true,
			PIIPatterns:        []string{"email", "phone", "ssn", "address"},
		},
	}

	if configPath == "" {
		return nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cm.saveConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var fileConfig types.Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	cm.mergeWithDefaults(&fileConfig)
	cm.config = &fileConfig

	return nil
}

func (cm *ConfigManagerImpl) GetConfig() *types.Config {
	if cm.config == nil {
		_ = cm.LoadConfig("")
	}
	return cm.config
}

func (cm *ConfigManagerImpl) UpdateConfig(config *types.Config) error {
	if err := cm.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	cm.config = config
	return nil
}

func (cm *ConfigManagerImpl) ValidateConfig(config *types.Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if config.AnalysisSettings.SampleSize < 1 {
		return fmt.Errorf("sample size must be greater than 0")
	}

	if config.AnalysisSettings.LargeTableThreshold < 0 {
		return fmt.Errorf("large table threshold cannot be negative")
	}

	if config.AnalysisSettings.QualityScoreMinimum < 0 || config.AnalysisSettings.QualityScoreMinimum > 1 {
		return fmt.Errorf("quality score minimum must be between 0 and 1")
	}

	for name, connectionString := range config.DatabaseConnections {
		if name == "" {
			return fmt.Errorf("database connection name cannot be empty")
		}
		if connectionString == "" {
			return fmt.Errorf("database connection string cannot be empty for %s", name)
		}
	}

	return nil
}

func (cm *ConfigManagerImpl) saveConfig(configPath string) error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (cm *ConfigManagerImpl) mergeWithDefaults(fileConfig *types.Config) {
	if fileConfig.DatabaseConnections == nil {
		fileConfig.DatabaseConnections = make(map[string]string)
	}

	if fileConfig.AnalysisSettings.SampleSize == 0 {
		fileConfig.AnalysisSettings.SampleSize = 1000
	}
	if fileConfig.AnalysisSettings.LargeTableThreshold == 0 {
		fileConfig.AnalysisSettings.LargeTableThreshold = 1000000
	}
	if fileConfig.AnalysisSettings.QualityScoreMinimum == 0 {
		fileConfig.AnalysisSettings.QualityScoreMinimum = 0.7
	}
	if fileConfig.AnalysisSettings.AutoAnalysisInterval == "" {
		fileConfig.AnalysisSettings.AutoAnalysisInterval = "24h"
	}

	if len(fileConfig.SecuritySettings.PIIPatterns) == 0 {
		fileConfig.SecuritySettings.PIIPatterns = []string{"email", "phone", "ssn", "address"}
	}
}
