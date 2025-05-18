package config

import "testing"

func TestVaultAuthMethodIsValid(t *testing.T) {
	tests := []struct {
		name   string
		method VaultAuthMethod
		want   bool
	}{
		{
			name:   "kubernetes",
			method: "kubernetes",
			want:   true,
		},
		{
			name:   "invalid",
			method: "invalid",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method.IsValid(); got != tt.want {
				t.Errorf("unexpected result: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	defaultConfig := Configuration{
		Log: LogConfiguration{
			Enabled: true,
			File:    "test.log",
			Level:   "info",
		},
		Vault: VaultConfiguration{
			Address:            "http://localhost:8200",
			InsecureSkipVerify: false,
			Auth: VaultAuthConfiguration{
				Method: VaultAuthMethodKubernetes,
				Mount:  "kubernetes",
				Role:   "example",
			},
			Secret: VaultSecretConfiguration{
				Mount: "secret",
				Path:  "example",
			},
		},
	}

	tests := []struct {
		name       string
		config     Configuration
		wantErrMsg string
	}{
		{
			name:       "valid config",
			config:     defaultConfig,
			wantErrMsg: "",
		},
		{
			name: "missing log file",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Log.File = ""
				return cfg
			}(),
			wantErrMsg: "log file is required",
		},
		{
			name: "missing log level",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Log.Level = ""
				return cfg
			}(),
			wantErrMsg: "log level is required",
		},
		{
			name: "invalid log level",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Log.Level = "invalid"
				return cfg
			}(),
			wantErrMsg: "log level is invalid. valid values are: debug, info, warn, error",
		},
		{
			name: "missing vault address",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Address = ""
				return cfg
			}(),
			wantErrMsg: "vault address is required",
		},
		{
			name: "missing vault auth method",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Auth.Method = ""
				return cfg
			}(),
			wantErrMsg: "vault auth method is required",
		},
		{
			name: "invalid vault auth method",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Auth.Method = "invalid"
				return cfg
			}(),
			wantErrMsg: "vault auth method is invalid. valid values are: kubernetes",
		},
		{
			name: "missing vault auth mount",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Auth.Mount = ""
				return cfg
			}(),
			wantErrMsg: "vault auth mount is required",
		},
		{
			name: "missing vault auth role",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Auth.Role = ""
				return cfg
			}(),
			wantErrMsg: "vault auth role is required",
		},
		{
			name: "missing vault secret mount",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Secret.Mount = ""
				return cfg
			}(),
			wantErrMsg: "vault secret mount is required",
		},
		{
			name: "missing vault secret path",
			config: func() Configuration {
				cfg := defaultConfig
				cfg.Vault.Secret.Path = ""
				return cfg
			}(),
			wantErrMsg: "vault secret path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("unexpected error message: got %v, want %v", err.Error(), tt.wantErrMsg)
			}
		})
	}
}
