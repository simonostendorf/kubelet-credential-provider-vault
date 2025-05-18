package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/logger"
	"github.com/spf13/viper"
)

const DefaultConfigFile = "./kubelet-credential-provider-vault.yaml"

type Configuration struct {
	Log   LogConfiguration   `mapstructure:"log"`
	Vault VaultConfiguration `mapstructure:"vault"`
}

type LogConfiguration struct {
	File    string `mapstructure:"file"`
	Level   string `mapstructure:"level"`
	Enabled bool   `mapstructure:"enabled"`
}

type VaultConfiguration struct {
	Address            string                   `mapstructure:"address"`
	InsecureSkipVerify bool                     `mapstructure:"insecureSkipVerify"`
	Auth               VaultAuthConfiguration   `mapstructure:"auth"`
	Secret             VaultSecretConfiguration `mapstructure:"secret"`
}

type VaultAuthMethod string

const (
	VaultAuthMethodKubernetes VaultAuthMethod = "kubernetes"
)

func (v VaultAuthMethod) IsValid() bool {
	switch v {
	case VaultAuthMethodKubernetes:
		return true
	default:
		return false
	}
}

type VaultAuthConfiguration struct {
	Method VaultAuthMethod `mapstructure:"method"`
	Mount  string          `mapstructure:"mount"`
	Role   string          `mapstructure:"role"`
}

type VaultSecretConfiguration struct {
	Mount string `mapstructure:"mount"`
	Path  string `mapstructure:"path"`
}

func (c *Configuration) validate() error {
	var errs []error
	if c.Log.File == "" {
		errs = append(errs, fmt.Errorf("log file is required"))
	}
	if c.Log.Level == "" {
		errs = append(errs, fmt.Errorf("log level is required"))
	} else if _, err := logger.ParseLogLevel(c.Log.Level); err != nil {
		errs = append(errs, fmt.Errorf("log level is invalid. valid values are: debug, info, warn, error"))
	}
	if c.Vault.Address == "" {
		errs = append(errs, fmt.Errorf("vault address is required"))
	}
	if c.Vault.Auth.Method == "" {
		errs = append(errs, fmt.Errorf("vault auth method is required"))
	} else if !c.Vault.Auth.Method.IsValid() {
		errs = append(errs, fmt.Errorf("vault auth method is invalid. valid values are: %s", VaultAuthMethodKubernetes))
	}
	if c.Vault.Auth.Mount == "" {
		errs = append(errs, fmt.Errorf("vault auth mount is required"))
	}
	if c.Vault.Auth.Role == "" {
		errs = append(errs, fmt.Errorf("vault auth role is required"))
	}
	if c.Vault.Secret.Mount == "" {
		errs = append(errs, fmt.Errorf("vault secret mount is required"))
	}
	if c.Vault.Secret.Path == "" {
		errs = append(errs, fmt.Errorf("vault secret path is required"))
	}
	if len(errs) > 0 {
		err := ""
		for i, e := range errs {
			err += e.Error()
			if i < len(errs)-1 {
				err += "; "
			}
		}
		return errors.New(err)
	}
	return nil
}

func New(ctx context.Context, log logger.Logger, configFile string) (*Configuration, error) {
	// load .env file if present
	err := godotenv.Load()
	if err != nil {
		var pathErr *os.PathError
		// ignore not found error
		if errors.As(err, &pathErr) && !os.IsNotExist(pathErr) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
		log.Log(ctx, slog.LevelWarn, "No .env file found, using command line arguments, environment variables or config file")
	}

	// select config file
	if configFile != "" {
		viper.SetConfigFile(configFile)
		log.Log(ctx, slog.LevelInfo, "Using specified config file", "file", configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("kubelet-credential-provider-vault")
		viper.SetConfigType("yaml")
		log.Log(ctx, slog.LevelDebug, "No config file specified, using default config file", "file", DefaultConfigFile)
	}

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		// ignore not found error if config file is not explicitly set
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && configFile != "" {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Log(ctx, slog.LevelInfo, "Loaded config file", "file", viper.ConfigFileUsed())
	}

	// load config
	var cfg Configuration
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	// validate config
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
