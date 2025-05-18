package credentialFetcher

import (
	"reflect"
	"testing"

	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/config"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/vault"
	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

func TestReadAuthConfig(t *testing.T) {
	tests := []struct {
		name            string
		vaultSecretData map[string]any
		want            *credentialproviderV1.AuthConfig
		wantErrMsg      string
	}{
		{
			name: "successful read",
			vaultSecretData: map[string]any{
				"username": "username",
				"password": "password",
			},
			want: &credentialproviderV1.AuthConfig{
				Username: "username",
				Password: "password",
			},
			wantErrMsg: "",
		},
		{
			name: "missing username",
			vaultSecretData: map[string]any{
				"password": "password",
			},
			want:       nil,
			wantErrMsg: "failed to read username from secret data",
		},
		{
			name: "missing password",
			vaultSecretData: map[string]any{
				"username": "username",
			},
			want:       nil,
			wantErrMsg: "failed to read password from secret data",
		},
		{
			name:            "empty vault secret data",
			vaultSecretData: map[string]any{},
			want:            nil,
			wantErrMsg:      "failed to read username from secret data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup vault credential fetcher with dummy data
			fetcher := VaultCredentialFetcher{
				vaultConfig: &config.VaultConfiguration{
					Address:            "http://localhost:8200",
					InsecureSkipVerify: false,
					Auth: config.VaultAuthConfiguration{
						Method: config.VaultAuthMethodKubernetes,
						Mount:  "kubernetes",
						Role:   "example",
					},
					Secret: config.VaultSecretConfiguration{
						Mount: "secret",
						Path:  "example",
					},
				},
				vaultClientBuilder: nil,
			}

			// create vault client mock
			vaultClient, err := vault.NewMockClientBuilder(tt.vaultSecretData).Build(t.Context())
			if err != nil {
				t.Fatalf("failed to create vault client: %v", err)
			}

			got, err := fetcher.readAuthConfig(t.Context(), vaultClient)
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("unexpected error: got %v, want %v", err, tt.wantErrMsg)
			}
			if got != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unexpected auth config: got %v, want %v", got, tt.want)
			}
		})
	}
}
