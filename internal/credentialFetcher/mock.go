package credentialFetcher

import (
	"context"

	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type MockCredentialFetcher struct {
	authConfig *credentialproviderV1.AuthConfig
}

func NewMockCredentialFetcher(authConfig *credentialproviderV1.AuthConfig) CredentialFetcher {
	return &MockCredentialFetcher{
		authConfig: authConfig,
	}
}

func (f *MockCredentialFetcher) Fetch(_ context.Context, _ *credentialproviderV1.CredentialProviderRequest) (*credentialproviderV1.AuthConfig, error) {
	return f.authConfig, nil
}
