package communicationInterface

import (
	"context"

	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type MockCommunicationInterface struct {
	request      *credentialproviderV1.CredentialProviderRequest
	lastResponse *credentialproviderV1.CredentialProviderResponse
}

func NewMockCommunicationInterface(request *credentialproviderV1.CredentialProviderRequest) CommunicationInterface {
	return &MockCommunicationInterface{
		request:      request,
		lastResponse: nil,
	}
}

func (i *MockCommunicationInterface) ReadRequest(_ context.Context) (*credentialproviderV1.CredentialProviderRequest, error) {
	return i.request, nil
}

func (i *MockCommunicationInterface) WriteResponse(_ context.Context, response *credentialproviderV1.CredentialProviderResponse) error {
	i.lastResponse = response
	return nil
}

func (i *MockCommunicationInterface) LastResponse() *credentialproviderV1.CredentialProviderResponse {
	return i.lastResponse
}
