package provider

import (
	"reflect"
	"testing"

	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/communicationInterface"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/credentialFetcher"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/logger"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		request    credentialproviderV1.CredentialProviderRequest
		want       *credentialproviderV1.CredentialProviderResponse
		wantErrMsg string
	}{
		{
			name: "valid request",
			request: credentialproviderV1.CredentialProviderRequest{
				TypeMeta: metaV1.TypeMeta{
					APIVersion: credentialproviderV1.SchemeGroupVersion.String(),
					Kind:       "CredentialProviderRequest",
				},
				Image:               "registry.example.com/my-image:latest",
				ServiceAccountToken: "token",
			},
			want: &credentialproviderV1.CredentialProviderResponse{
				TypeMeta: metaV1.TypeMeta{
					APIVersion: credentialproviderV1.SchemeGroupVersion.String(),
					Kind:       "CredentialProviderResponse",
				},
				Auth: map[string]credentialproviderV1.AuthConfig{
					"registry.example.com": {
						Username: "user",
						Password: "password",
					},
				},
				CacheKeyType: credentialproviderV1.RegistryPluginCacheKeyType,
			},
			wantErrMsg: "",
		},
		{
			name: "no registry in image",
			request: credentialproviderV1.CredentialProviderRequest{
				TypeMeta: metaV1.TypeMeta{
					APIVersion: credentialproviderV1.SchemeGroupVersion.String(),
					Kind:       "CredentialProviderRequest",
				},
				Image:               "my-image:latest",
				ServiceAccountToken: "token",
			},
			want:       nil,
			wantErrMsg: "failed to extract registry name: no registry name found in image name: my-image:latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock dependencies
			communicationInterface := communicationInterface.NewMockCommunicationInterface(&tt.request)
			credentialFetcher := credentialFetcher.NewMockCredentialFetcher(&credentialproviderV1.AuthConfig{
				Username: "user",
				Password: "password",
			})
			logger, err := logger.NewFileLogger(false, "", "error")
			if err != nil {
				t.Fatalf("failed to create logger: %v", err)
			}

			// create KubeletCredentialProvider
			provider := NewKubeletCredentialProvider(communicationInterface, credentialFetcher)

			// run the provider
			err = provider.Run(t.Context(), logger)
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("unexpected error: got %v, want %v", err, tt.wantErrMsg)
			}
			if communicationInterface.LastResponse() != nil && !reflect.DeepEqual(communicationInterface.LastResponse(), tt.want) {
				t.Errorf("unexpected response: got %v, want %v", communicationInterface.LastResponse(), tt.want)
			}
		})
	}
}
