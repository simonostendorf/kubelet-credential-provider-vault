package credentialFetcher

import (
	"context"
	"fmt"

	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/config"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/vault"
	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type VaultCredentialFetcher struct {
	vaultClientBuilder vault.ClientBuilder
	vaultConfig        *config.VaultConfiguration
}

func NewVaultCredentialFetcher(vaultClientBuilder vault.ClientBuilder, vaultConfig *config.VaultConfiguration) CredentialFetcher {
	return &VaultCredentialFetcher{
		vaultClientBuilder: vaultClientBuilder,
		vaultConfig:        vaultConfig,
	}
}

func (f *VaultCredentialFetcher) Fetch(ctx context.Context, request *credentialproviderV1.CredentialProviderRequest) (*credentialproviderV1.AuthConfig, error) {
	// get service account token from request
	serviceAccountToken := request.ServiceAccountToken
	if serviceAccountToken == "" {
		return nil, fmt.Errorf("service account token is required")
	}

	// setup vault client
	vaultClient, err := f.setupVaultClient(ctx, serviceAccountToken)
	if err != nil {
		return nil, fmt.Errorf("failed to setup vault client: %w", err)
	}

	// read auth config from vault
	authConfig, err := f.readAuthConfig(ctx, vaultClient)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth config from vault: %w", err)
	}

	return authConfig, nil
}

func (f *VaultCredentialFetcher) setupVaultClient(ctx context.Context, serviceAccountToken string) (vault.Client, error) {
	// authenticate with vault
	switch f.vaultConfig.Auth.Method {
	case config.VaultAuthMethodKubernetes:
		// kubernetes auth method not possible when service account token is not provided
		if serviceAccountToken == "" {
			return nil, fmt.Errorf("service account token is required for kubernetes auth method")
		}

		// authenticate with kubernetes auth method
		return f.vaultClientBuilder.
			WithAddress(f.vaultConfig.Address).
			InsecureSkipVerify(f.vaultConfig.InsecureSkipVerify).
			WithKubernetesAuth(f.vaultConfig.Auth.Mount, f.vaultConfig.Auth.Role, serviceAccountToken).
			Build(ctx)
	default:
		return nil, fmt.Errorf("unsupported vault auth method: %s", f.vaultConfig.Auth.Method)
	}
}

func (f *VaultCredentialFetcher) readAuthConfig(ctx context.Context, vaultClient vault.Client) (*credentialproviderV1.AuthConfig, error) {
	// read vault secret
	secretData, err := vaultClient.Secrets().KvV2(f.vaultConfig.Secret.Mount, f.vaultConfig.Secret.Path).Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from vault: %w", err)
	}

	username, ok := secretData["username"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to read username from secret data")
	}
	password, ok := secretData["password"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to read password from secret data")
	}

	return &credentialproviderV1.AuthConfig{
		Username: username,
		Password: password,
	}, nil
}
