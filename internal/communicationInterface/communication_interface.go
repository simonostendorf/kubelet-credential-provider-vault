package communicationInterface

import (
	"context"

	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type CommunicationInterface interface {
	ReadRequest(ctx context.Context) (*credentialproviderV1.CredentialProviderRequest, error)
	WriteResponse(ctx context.Context, response *credentialproviderV1.CredentialProviderResponse) error
	LastResponse() *credentialproviderV1.CredentialProviderResponse // mainly for testing purposes
}
