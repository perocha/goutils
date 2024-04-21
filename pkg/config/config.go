package config

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
)

type Config struct {
	AppInsightsInstrumentationKey string
	Client                        *azappconfig.Client
}

// New configuration client
func NewConfig(connectionString string) (*Config, error) {
	client, err := azappconfig.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}
	return &Config{
		Client: client,
	}, nil
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
