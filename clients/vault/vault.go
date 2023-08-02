package vault

import (
	"git.zonatelecom.ru/fsin/censor/config"
	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *vault.Client
	cfg    *config.VaultEnv
}
