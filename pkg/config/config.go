package config

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	AppConfigurationConnectionString string `yaml:"APPCONFIGURATION_CONNECTION_STRING"`
}

type Config struct {
	AppInsightsInstrumentationKey string
	Client                        *azappconfig.Client
}

// New configuration client from connection string
func NewConfigFromConnectionString(connectionString string) (*Config, error) {
	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}
	return &Config{
		Client: client,
	}, nil
}

// New configuration client from file
func NewConfigFromFile(configFile string) (*Config, error) {
	cfg := &Config{}

	if configFile == "" {
		return nil, errors.New("config file not provided")
	}

	// Create a new App Configuration client
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var yamlCfg yamlConfig
	err = yaml.Unmarshal(data, &yamlCfg)
	if err != nil {
		return nil, err
	}
	connectionString := yamlCfg.AppConfigurationConnectionString

	cfg.Client, err = azappconfig.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Get configuration value from key
func (cfg *Config) GetVar(key string) (string, error) {
	if cfg.Client == nil {
		return "", errors.New("app configuration client not initialized")
	}

	resp, err := cfg.Client.GetSetting(context.TODO(), key, nil)
	if err != nil {
		return "", err
	}

	return *resp.Value, nil
}

// Get configuration value from key, return empty string if not found
func (cfg *Config) GetVarOrDefault(key string, defaultValue string) string {
	value, err := cfg.GetVar(key)
	if err != nil {
		return defaultValue
	}
	return value
}
