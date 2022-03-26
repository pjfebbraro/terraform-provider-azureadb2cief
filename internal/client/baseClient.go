package client

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"strings"
)

type MsGraphClientConfig struct {
	TenantID            string
	ClientID            string
	ClientSecret        string
	EnableAzureCliToken bool
}

type baseClient struct {
	cred    azcore.TokenCredential
	config  MsGraphClientConfig
	baseUrl string
	scopes  []string
	client  *retryablehttp.Client
}

func newBaseClient(config MsGraphClientConfig) (*baseClient, error) {
	var cred azcore.TokenCredential
	var err error
	if config.EnableAzureCliToken && strings.TrimSpace(config.TenantID) != "" {
		cred, err = azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: config.TenantID})
		if err != nil {
			return nil, fmt.Errorf("could not configure AzureCli Authorizer: %s", err)
		}
	} else if strings.TrimSpace(config.TenantID) != "" && strings.TrimSpace(config.ClientID) != "" && strings.TrimSpace(config.ClientSecret) != "" {
		cred, err = azidentity.NewClientSecretCredential(config.TenantID, config.ClientID, config.ClientSecret, &azidentity.ClientSecretCredentialOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientCertificate Authorizer: %s", err)
		}
	}

	client := retryablehttp.NewClient()
	client.RetryMax = 3

	if cred != nil {
		return &baseClient{
			cred:    cred,
			config:  config,
			baseUrl: "https://graph.microsoft.com/beta",
			scopes: []string{
				"https://graph.microsoft.com/.default",
			},
			client: client,
		}, nil
	}

	return nil, fmt.Errorf("no Authorizer could be configured, please check your configuration")
}

func (gc *baseClient) doRequest(ctx context.Context, path string, method string, body io.Reader, contentType *string) (*http.Response, error) {
	token, err := gc.cred.GetToken(ctx, policy.TokenRequestOptions{
		TenantID: gc.config.TenantID,
		Scopes:   gc.scopes,
	})
	if err != nil {
		return nil, err
	}
	url := gc.baseUrl + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	if contentType != nil {
		req.Header.Set("Content-Type", *contentType)
	}

	retryablereq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	return gc.client.Do(retryablereq)
}

func formatHttpErrorResponse(response *http.Response) error {
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unexpected status %d, could not read response body", response.StatusCode)
	}
	if len(respBody) == 0 {
		return fmt.Errorf("unexpected status %d received with no body", response.StatusCode)
	}
	errText := fmt.Sprintf("response: %s", respBody)
	return fmt.Errorf("unexpected status %d with %s", response.StatusCode, errText)
}
