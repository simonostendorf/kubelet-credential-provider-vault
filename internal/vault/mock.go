package vault

import (
	"context"
)

type MockClientBuilder struct {
	address             *string
	insecureSkipVerify  *bool
	mount               *string
	role                *string
	serviceAccountToken *string

	mockSecretResponse map[string]any
}

func NewMockClientBuilder(mockSecretResponse map[string]any) ClientBuilder {
	return &MockClientBuilder{
		mockSecretResponse: mockSecretResponse,
	}
}

func (b *MockClientBuilder) WithAddress(address string) ClientBuilder {
	b.address = &address
	return b
}

func (b *MockClientBuilder) InsecureSkipVerify(insecureSkipVerify bool) ClientBuilder {
	b.insecureSkipVerify = &insecureSkipVerify
	return b
}

func (b *MockClientBuilder) WithKubernetesAuth(mount string, role string, serviceAccountToken string) ClientBuilder {
	b.mount = &mount
	b.role = &role
	b.serviceAccountToken = &serviceAccountToken
	return b
}

func (b *MockClientBuilder) validate() error {
	return nil
}

func (b *MockClientBuilder) Build(_ context.Context) (Client, error) {
	return newMockClient(b.mockSecretResponse), nil
}

type MockClient struct {
	secretsClient SecretsClient
}

func newMockClient(mockSecretResponse map[string]any) *MockClient {
	return &MockClient{
		secretsClient: newMockSecretsClient(mockSecretResponse),
	}
}

func (c *MockClient) Secrets() SecretsClient {
	return c.secretsClient
}

type MockSecretsClient struct {
	mockSecretResponse map[string]any
}

func newMockSecretsClient(mockSecretResponse map[string]any) *MockSecretsClient {
	return &MockSecretsClient{
		mockSecretResponse: mockSecretResponse,
	}
}

func (c *MockSecretsClient) KvV2(mount string, path string) SecretKvV2Client {
	return &MockSecretKvV2Client{
		mount:              mount,
		path:               path,
		mockSecretResponse: c.mockSecretResponse,
	}
}

type MockSecretKvV2Client struct {
	mount string
	path  string

	mockSecretResponse map[string]any
}

func (c *MockSecretKvV2Client) Read(_ context.Context) (map[string]any, error) {
	return c.mockSecretResponse, nil
}
