package communicationInterface

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	credentialproviderV1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type StdIOCommunicationInterface struct {
	lastResponse *credentialproviderV1.CredentialProviderResponse
}

func NewStdIOCommunicationInterface() CommunicationInterface {
	return &StdIOCommunicationInterface{
		lastResponse: nil,
	}
}

func (p *StdIOCommunicationInterface) ReadRequest(ctx context.Context) (*credentialproviderV1.CredentialProviderRequest, error) {
	// basic read from io.ReadAll(os.Stdin) is not possible because this wouldnt be context aware
	// so canceling the context in the main function would not stop the read and the program would hang

	// create pipe to read from stdin
	pipeReader, pipeWriter := io.Pipe()

	// goroutine to copy from stdin to the pipe
	go func() {
		// nolint:errcheck
		defer pipeWriter.Close() //gosec:disable G104
		_, err := io.Copy(pipeWriter, os.Stdin)
		if err != nil {
			pipeWriter.CloseWithError(err)
		}
	}()

	// create a reader from the pipe (this will work at the same time as the write above (because its a gotoutine))
	done := make(chan error, 1)
	var in []byte
	go func() {
		var err error
		in, err = io.ReadAll(pipeReader)
		done <- err
	}()

	// handle context cancellation
	select {
	case <-ctx.Done():
		// stop reading if the context is done
		// nolint:errcheck
		pipeReader.Close() //gosec:disable G104
		return nil, ctx.Err()
	case err := <-done:
		// reading completed
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
	}

	// unmarshal the request
	request := &credentialproviderV1.CredentialProviderRequest{}
	err := json.Unmarshal(in, request)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	// validate the request
	if request.APIVersion != credentialproviderV1.SchemeGroupVersion.String() {
		return nil, fmt.Errorf("invalid request: expected apiVersion %s, got %s", credentialproviderV1.SchemeGroupVersion.String(), request.APIVersion)
	}
	if request.Kind != "CredentialProviderRequest" {
		return nil, fmt.Errorf("invalid request: expected kind %s, got %s", "CredentialProviderRequest", request.Kind)
	}

	return request, nil
}

func (p *StdIOCommunicationInterface) WriteResponse(ctx context.Context, response *credentialproviderV1.CredentialProviderResponse) error {
	// set last response
	p.lastResponse = response

	// marshal the response
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// write to stdout
	fmt.Println(string(data))

	return nil
}

func (p *StdIOCommunicationInterface) LastResponse() *credentialproviderV1.CredentialProviderResponse {
	return p.lastResponse
}
