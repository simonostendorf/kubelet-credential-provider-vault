package vault

import "context"

type ClientBuilder interface {
	WithAddress(address string) ClientBuilder
	InsecureSkipVerify(insecureSkipVerify bool) ClientBuilder
	WithKubernetesAuth(mount string, role string, serviceAccountToken string) ClientBuilder
	validate() error
	Build(ctx context.Context) (Client, error)
}

type Client interface {
	Secrets() SecretsClient
}

type SecretsClient interface {
	KvV2(mount string, path string) SecretKvV2Client
}

type SecretKvV2Client interface {
	Read(ctx context.Context) (map[string]any, error)
}
