package provider

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/communicationInterface"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/credentialFetcher"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/logger"
	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type KubeletCredentialProvider struct {
	communicationInterface communicationInterface.CommunicationInterface
	credentialFetcher      credentialFetcher.CredentialFetcher
}

func NewKubeletCredentialProvider(communicationInterface communicationInterface.CommunicationInterface, credentialFetcher credentialFetcher.CredentialFetcher) *KubeletCredentialProvider {
	return &KubeletCredentialProvider{
		communicationInterface: communicationInterface,
		credentialFetcher:      credentialFetcher,
	}
}

func (k *KubeletCredentialProvider) Run(ctx context.Context, log logger.Logger) error {
	// read request
	request, err := k.communicationInterface.ReadRequest(ctx)
	if err != nil {
		return fmt.Errorf("failed to read request: %w", err)
	}
	log.Log(ctx, slog.LevelDebug, "Received request", "request", request)

	// fetch credentials
	authConfig, err := k.credentialFetcher.Fetch(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to fetch credentials: %w", err)
	}
	log.Log(ctx, slog.LevelDebug, "Fetched credentials", "credentials", authConfig)

	// get registry name
	registryName, err := extractRegistryName(request.Image)
	if err != nil {
		return fmt.Errorf("failed to extract registry name: %w", err)
	}
	log.Log(ctx, slog.LevelDebug, "Extracted registry name", "registryName", registryName)

	// create response
	response := &credentialproviderV1.CredentialProviderResponse{
		Auth: map[string]credentialproviderV1.AuthConfig{
			registryName: *authConfig,
		},
		CacheKeyType: credentialproviderV1.RegistryPluginCacheKeyType,
	}
	response.APIVersion = credentialproviderV1.SchemeGroupVersion.String()
	response.Kind = "CredentialProviderResponse"
	log.Log(ctx, slog.LevelDebug, "Created response", "response", response)

	// write response
	err = k.communicationInterface.WriteResponse(ctx, response)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	log.Log(ctx, slog.LevelDebug, "Wrote response", "response", response)
	return nil
}

func extractRegistryName(imageName string) (string, error) {
	parts := strings.Split(imageName, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("no registry name found in image name: %s", imageName)
	}
	return parts[0], nil
}
