package clients

import (
	"context"
	"fmt"

	"git.zonatelecom.ru/fsin/censor/config"
	vault "github.com/hashicorp/vault/api"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
)

type VaultClient struct {
	client *vault.Client
	cfg    *config.VaultEnv
}

type DBCreds struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (d *DBCreds) SourceString(production bool) string {
	if !production {
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=disable",
			d.Username, d.Password, d.DBName, d.Host,
		)
	}
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s",
		d.Username, d.Password, d.DBName, d.Host,
	)
}

type RabbitMQCreds struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (r *RabbitMQCreds) SourceString() string {
	return fmt.Sprintf("amqp://%s:%s@%s/", r.Username, r.Password, r.Host)
}

// New Vault client
func Vault(cfg *config.VaultEnv) (*VaultClient, error) {
	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = cfg.Host
	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, err
	}
	client.SetToken(cfg.Token)

	newClient := &VaultClient{
		client: client,
		cfg:    cfg,
	}

	return newClient, nil
}

func (v *VaultClient) Client() *vault.Client {
	return v.client
}

func (v *VaultClient) DBCreds(ctx context.Context) (*DBCreds, error) {
	kvSecret, err := v.client.KVv2(v.cfg.MountPath).Get(ctx, v.cfg.DbCredsPath)
	if err != nil {
		return nil, err
	}
	var creds *DBCreds
	err = mapstructure.Decode(kvSecret.Data, &creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func (v *VaultClient) RabbitMQCreds(ctx context.Context) (*RabbitMQCreds, error) {
	kvSecret, err := v.client.KVv2(v.cfg.MountPath).Get(ctx, v.cfg.RabbitMQCreds)
	if err != nil {
		return nil, err
	}
	var creds *RabbitMQCreds
	err = mapstructure.Decode(kvSecret.Data, &creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func DB(ctx context.Context, production bool, creds *DBCreds) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", creds.SourceString(production))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
