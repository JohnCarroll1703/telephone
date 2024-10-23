package config

import (
	"context"
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

type (
	Config struct {
		GRPC     *GrpcConnections
		Service  *Service
		Database *Databases
		Nats     *Nats
		Vault    Vault
		Sentry   *Sentry
		Jeager   *Jeager
	}

	GrpcConnections struct {
		Auth     string `envconfig:"AUTH_URL" default:"auth:4040"`
		Acquirer string `envconfig:"ACQUIRER_URL" required:"true"`
	}

	Service struct {
		AppName      string         `envconfig:"APP_NAME" required:"true"`
		Environment  AppEnvironment `envconfig:"ENVIRONMENT" default:"local"`
		Port         string         `envconfig:"PORT" default:"8080"`
		GrpcPort     string         `envconfig:"GRPC_PORT" default:"4040"`
		Shared       string         `envconfig:"SHARED" required:"false"`
		Openapi      string         `envconfig:"OPENAPI_ENDPOINT" required:"false"`
		Domain       string         `envconfig:"DOMAIN" default:"localhost"`
		RunMigration bool           `envconfig:"RUN_MIGRATION"`
		Namespace    string         `envconfig:"NAMESPACE" required:"true"`
	}

	Databases struct {
		ReadDSNs     []string `envconfig:"DATABASE_READ_DSN"`
		WriteDSNs    []string `envconfig:"DATABASE_WRITE_DSN" required:"false"`
		Schema       string   `envconfig:"DATABASE_SCHEMA" required:"false"`
		PostgreDSN   string   `envconfig:"POSTGRE_DSN" required:"false"`
		PostgreTable string   `envconfig:"POSTGRE_TABLE" required:"false"`
		PostgreName  string   `envconfig:"POSTGRE_NAME" required:"false"`
	}

	Sentry struct {
		DSN string `envconfig:"SENTRY_DSN"`
	}

	Jeager struct {
		DSN string `envconfig:"JAEGER_DSN"`
	}

	Vault struct {
		URL           string `envconfig:"VAULT_URL" required:"true"`
		Namespace     string `envconfig:"NAMESPACE" required:"true"`
		VaultUser     string `envconfig:"VAULT_USER" default:""`
		VaultPassword string `envconfig:"VAULT_PASSWORD" default:""`
	}

	Nats struct {
		DSN string `envconfig:"NATS_DSN" required:"false"`
	}
)

type AppEnvironment string

const (
	PRODUCTION  AppEnvironment = "prod"
	SANDBOX     AppEnvironment = "sandbox"
	STAGE       AppEnvironment = "stage"
	DEVELOPMENT AppEnvironment = "dev"
	TEST        AppEnvironment = "test"
	LOCAL       AppEnvironment = "local"
)

func GetConfig(once *sync.Once, envfiles ...string) (*Config, error) {
	var (
		config *Config
		err    error
	)

	if once == nil {
		return nil, errors.New("received a nil value of sync.Once")
	}
	once.Do(
		func() {
			_ = godotenv.Load(envfiles...)
			var c Config
			err = envconfig.Process("", &c)
			if err != nil {
				err = fmt.Errorf("error parse config from env variables: %w\n", err)
				return
			}

			vaultErr := updateCfgFromVault(&c)
			if vaultErr != nil {
				// Есть возможность добить в env на случай если в vault едоступен
				fmt.Println(vaultErr.Error())
			}

			config = &c
		})
	return config, err
}

func updateCfgFromVault(cfg *Config) (err error) {
	vaultConfig := vault.DefaultConfig()

	vaultConfig.Address = cfg.Vault.URL

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return
	}

	if err = authorizeToVault(cfg, client); err != nil {
		return err
	}

	serviceSecret, err := client.KVv2(cfg.Vault.Namespace).Get(
		context.Background(),
		cfg.Service.AppName)
	if err != nil {
		return
	}

	shared, ok := serviceSecret.Data["shared"].(string)
	if !ok {
		return errors.New("could not read a valid string value from [shared]")
	}

	cfg.Service.Shared = shared

	cfg.Database.PostgreDSN, ok = serviceSecret.Data["POSTGRE_DSN"].(string)
	if !ok {
		return errors.New("could not read a valid string value from [POSTGRE_DSN]")
	}

	cfg.Database.PostgreTable, ok = serviceSecret.Data["POSTGRE_TABLE"].(string)
	if !ok {
		return errors.New("could not read a valid string value from [POSTGRE_TABLE]")
	}

	cfg.Database.PostgreName, ok = serviceSecret.Data["POSTGRE_NAME"].(string)
	if !ok {
		return errors.New("could not read a valid string value from [POSTGRE_NAME]")
	}

	return
}

func authorizeToVaultByToken(cfg *Config, client *vault.Client) error {
	k8sAuth, k8sErr := auth.NewKubernetesAuth(
		cfg.Vault.Namespace,
		auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
	)
	if k8sErr != nil {
		return fmt.Errorf("k8s auth error: %w", k8sErr)
	}

	_, k8sErr = client.Auth().Login(context.Background(), k8sAuth)
	if k8sErr != nil {
		return fmt.Errorf("unable to log in with Kubernetes auth: %w", k8sErr)
	}

	return nil
}

func authorizeToVault(cfg *Config, client *vault.Client) error {
	if cfg == nil || client == nil {
		return errors.New("received empty data to process the authorization to vault")
	}

	if cfg.Vault.VaultUser == "" {
		return authorizeToVaultByToken(cfg, client)
	}

	secret, err := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", cfg.Vault.VaultUser), map[string]any{
		"password": cfg.Vault.VaultPassword,
	})
	if err != nil {
		return err
	}

	if secret.Auth == nil {
		return errors.New("did not receive valid auth details from vault")
	}

	client.SetToken(secret.Auth.ClientToken)

	return nil
}

func (e AppEnvironment) IsProduction() bool {
	return e == PRODUCTION
}

func (e AppEnvironment) IsStage() bool {
	return e == STAGE
}

func (e AppEnvironment) IsDevelopment() bool {
	return e == DEVELOPMENT
}

func (e AppEnvironment) IsLocal() bool {
	return e == LOCAL
}

func (e AppEnvironment) IsTest() bool {
	return e == TEST
}

func (e AppEnvironment) String() string {
	return string(e)
}

func (e AppEnvironment) Validate() error {
	switch e {
	case LOCAL, DEVELOPMENT, STAGE, SANDBOX, PRODUCTION, TEST:
		return nil
	default:
		return fmt.Errorf("unexpected ENVIRONMENT in .env: %s", e)
	}
}
