package credentialFetcher

import (
	"context"

	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type CredentialFetcher interface {
	Fetch(ctx context.Context, request *credentialproviderV1.CredentialProviderRequest) (*credentialproviderV1.AuthConfig, error)
}
