package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

const (
	ENV_PRODUCTION = "PRODUCTION"
	ENV_PORT       = "PORT"

	ENV_VAULT_TOKEN               = "VAULT_TOKEN"
	ENV_VAULT_HOST                = "VAULT_HOST"
	ENV_VAULT_MOUNT_PATH          = "VAULT_MOUNT_PATH"
	ENV_VAULT_DB_CREDS_PATH       = "VAULT_DB_CREDS_PATH"
	ENV_VAULT_RABBITMQ_CREDS_PATH = "VAULT_RABBITMQ_CREDS_PATH"

	ENV_TRACER_URL  = "TRACER_URL"
	ENV_TRACER_NAME = "TRACER_NAME"
)

type VaultEnv struct {
	MountPath     string `mapstructure:"VAULT_MOUNT_PATH"`
	DbCredsPath   string `mapstructure:"VAULT_DB_CREDS_PATH"`
	RabbitMQCreds string `mapstructure:"VAULT_RABBITMQ_CREDS_PATH"`
	Token         string `mapstructure:"VAULT_TOKEN"`
	Host          string `mapstructure:"VAULT_HOST"`
}

type TracerConfig struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Vault      *VaultEnv
	Tracer     *TracerConfig
	Port       string
	Production bool
}

func Read() (*Config, error) {
	var envCfg Config
	viper.AutomaticEnv()
	isProduction := viper.GetBool(ENV_PRODUCTION)

	vaultToken := os.Getenv(ENV_VAULT_TOKEN)
	if vaultToken == "" {
		return nil, fmt.Errorf("env VAULT_TOKEN is empty")
	}

	vaultHost := os.Getenv(ENV_VAULT_HOST)
	if vaultHost == "" {
		return nil, fmt.Errorf("env VAULT_HOST is empty")
	}

	vaultHostUrl, err := url.Parse(vaultHost)
	if err != nil {
		return nil, fmt.Errorf("VAULT_HOST has invalid value: %w", err)
	}

	vaultMountPath := os.Getenv(ENV_VAULT_MOUNT_PATH)
	if vaultMountPath == "" {
		return nil, fmt.Errorf("env VAULT_MOUNT_PATH is empty")
	}

	vaultDbCredsPath := os.Getenv(ENV_VAULT_DB_CREDS_PATH)
	if vaultDbCredsPath == "" {
		return nil, fmt.Errorf("env VAULT_DB_CREDS_PATH is empty")
	}

	vaultRabbitMQCredsPath := os.Getenv(ENV_VAULT_RABBITMQ_CREDS_PATH)
	if vaultDbCredsPath == "" {
		return nil, fmt.Errorf("env VAULT_RABBITMQ_CREDS_PATH is empty")
	}

	envCfg = Config{
		Vault: &VaultEnv{
			Token:         vaultToken,
			MountPath:     vaultMountPath,
			DbCredsPath:   vaultDbCredsPath,
			Host:          vaultHostUrl.String(),
			RabbitMQCreds: vaultRabbitMQCredsPath,
		},
		Tracer: &TracerConfig{
			URL:  viper.GetString(ENV_TRACER_URL),
			Name: viper.GetString(ENV_TRACER_NAME),
		},
		Port:       viper.GetString(ENV_PORT),
		Production: isProduction,
	}

	return &envCfg, nil
}
