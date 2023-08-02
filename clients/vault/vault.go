package vault

import (
	"context"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

var (
	vaultClient *VaultClient
)

func Client() *vault.Client {
	return vaultClient.client
}

type Config struct {
	MountPath string
	Token     string
	Host      string
}

type VaultClient struct {
	client *vault.Client
	cfg    *Config
}

func Init(cfg *Config) error {
	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = cfg.Host
	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return err
	}
	client.SetToken(cfg.Token)

	newClient := &VaultClient{
		client: client,
		cfg:    cfg,
	}
	vaultClient = newClient
	return nil
}

func GetCreds[T any](ctx context.Context, secretPath string) (*T, error) {
	kvSecret, err := vaultClient.client.KVv2(vaultClient.cfg.MountPath).Get(ctx, secretPath)
	if err != nil {
		return nil, err
	}
	var creds *T
	err = mapstructure.Decode(kvSecret.Data, &creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}
