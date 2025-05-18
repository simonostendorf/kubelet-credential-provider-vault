package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/helpers"

	hashiVault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type HashiCorpClientAuthMethod string

const (
	HashiCorpClientAuthMethodKubernetes HashiCorpClientAuthMethod = "kubernetes"
)

type HashiCorpClientBuilder struct {
	address             *string
	insecureSkipVerify  *bool
	authMethod          *HashiCorpClientAuthMethod
	mount               *string
	role                *string
	serviceAccountToken *string
}

func NewHashicorpClientBuilder() ClientBuilder {
	return &HashiCorpClientBuilder{}
}

func (b *HashiCorpClientBuilder) WithAddress(address string) ClientBuilder {
	b.address = &address
	return b
}

func (b *HashiCorpClientBuilder) InsecureSkipVerify(insecureSkipVerify bool) ClientBuilder {
	b.insecureSkipVerify = &insecureSkipVerify
	return b
}

func (b *HashiCorpClientBuilder) WithKubernetesAuth(mount string, role string, serviceAccountToken string) ClientBuilder {
	b.authMethod = helpers.Ptr(HashiCorpClientAuthMethodKubernetes)
	b.mount = &mount
	b.role = &role
	b.serviceAccountToken = &serviceAccountToken
	return b
}

func (b *HashiCorpClientBuilder) validate() error {
	var errs []error
	if b.address == nil {
		errs = append(errs, fmt.Errorf("address is required"))
	}
	if b.authMethod == nil {
		errs = append(errs, fmt.Errorf("auth method is required"))
	}
	if b.mount == nil {
		errs = append(errs, fmt.Errorf("mount is required"))
	}
	if b.authMethod != nil && *b.authMethod == HashiCorpClientAuthMethodKubernetes {
		if b.role == nil {
			errs = append(errs, fmt.Errorf("role is required for kubernetes auth method"))
		}
		if b.serviceAccountToken == nil {
			errs = append(errs, fmt.Errorf("service account token is required for kubernetes auth method"))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (b *HashiCorpClientBuilder) Build(ctx context.Context) (Client, error) {
	// validate current builder state
	if err := b.validate(); err != nil {
		return nil, fmt.Errorf("builder validation failed: %w", err)
	}

	// setup tls config
	tlsConfig := hashiVault.TLSConfiguration{}
	if b.insecureSkipVerify != nil && *b.insecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	// build vault client
	client, err := hashiVault.New(
		hashiVault.WithAddress(*b.address),
		hashiVault.WithTLS(tlsConfig),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// authenticate
	switch *b.authMethod {
	case HashiCorpClientAuthMethodKubernetes:
		resp, err := client.Auth.KubernetesLogin(ctx, schema.KubernetesLoginRequest{
			Jwt:  *b.serviceAccountToken,
			Role: *b.role,
		},
			hashiVault.WithMountPath(*b.mount),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate with kubernetes: %w", err)
		}
		if err := client.SetToken(resp.Auth.ClientToken); err != nil {
			return nil, fmt.Errorf("failed to set token on vault client: %w", err)
		}
	}

	return newHashiCorpClient(client), nil
}

type HashiCorpClient struct {
	client *hashiVault.Client

	secretsClient SecretsClient
}

func newHashiCorpClient(client *hashiVault.Client) *HashiCorpClient {
	return &HashiCorpClient{
		client:        client,
		secretsClient: newHashiCorpSecretsClient(client),
	}
}

func (c *HashiCorpClient) Secrets() SecretsClient {
	return c.secretsClient
}

type HashiCorpSecretsClient struct {
	client *hashiVault.Client
}

func newHashiCorpSecretsClient(client *hashiVault.Client) *HashiCorpSecretsClient {
	return &HashiCorpSecretsClient{
		client: client,
	}
}

func (c *HashiCorpSecretsClient) KvV2(mount string, path string) SecretKvV2Client {
	return &HashiCorpSecretKvV2Client{
		client: c.client,
		mount:  mount,
		path:   path,
	}
}

type HashiCorpSecretKvV2Client struct {
	client *hashiVault.Client
	mount  string
	path   string
}

func (c *HashiCorpSecretKvV2Client) Read(ctx context.Context) (map[string]any, error) {
	s, err := c.client.Secrets.KvV2Read(ctx, c.path,
		hashiVault.WithMountPath(c.mount),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}
	return s.Data.Data, nil
}
